package push

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/guonaihong/gout"
	log "github.com/sirupsen/logrus"
)

type Ding struct {
	Secret string `json:"Secret"`
	Token  string `json:"token"`
}

func (d *Ding) Send() func(kind string, message string) {
	s := TypeSecret{Secret: d.Secret, Webhook: d.Token}
	return func(kind string, message string) {
		if kind == "flush" {
			err := s.SendMessage(map[string]interface{}{
				"msgtype": "markdown",
				"markdown": map[string]string{
					"title": "学习强国登录",
					"text":  message,
				},
			})
			if err != nil {
				return
			}
		} else {
			if log.GetLevel() == log.DebugLevel {
				err := s.SendMessage(Text(message))
				if err != nil {
					return
				}
			}
		}
	}
}

type TypeSecret struct {
	Webhook string
	Secret  string
}

func Text(text string, ats ...string) map[string]interface{} {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": text,
		},
		"at": map[string]interface{}{
			"atMobiles": ats,
			"isAtAll":   false,
		},
	}
	return msg
}

func MarkDown(title, text string, ats ...string) map[string]interface{} {
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
		"at": map[string]interface{}{
			"atMobiles": ats,
			"isAtAll":   false,
		},
	}

	return msg
}

// SendMessage Function to send message
//goland:noinspection GoUnhandledErrorResult
func (t *TypeSecret) SendMessage(data map[string]interface{}) error {
	gout.POST(t.getURL()).SetJSON(data).Do()
	return nil
}

func (t *TypeSecret) hmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (t *TypeSecret) getURL() string {
	wh := "https://oapi.dingtalk.com/robot/send?access_token=" + t.Webhook
	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, t.Secret)
	sign := t.hmacSha256(stringToSign, t.Secret)
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", wh, timestamp, sign)
	return url
}
