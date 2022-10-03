package push

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/utils"
)

type Ding struct {
	Secret string `json:"Secret"`
	Token  string `json:"token"`
}

func (d *Ding) Send() func(id string, kind string, message string) {
	s := TypeSecret{Secret: d.Secret, Webhook: d.Token}
	return func(id string, kind string, message string) {
		if kind == "flush" {

			if strings.Contains(message, "login.xuexi.cn") {
				message = fmt.Sprintf("[点我登录](%v)", message)
			}

			err := s.SendMessage(map[string]interface{}{
				"msgtype": "markdown",
				"markdown": map[string]string{
					"title": "study_xxqg信息推送",
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
	_, err := utils.GetClient().R().SetBodyJsonMarshal(data).Post(t.getURL())
	if err != nil {
		log.Errorln(err.Error())
	}
	return err
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
