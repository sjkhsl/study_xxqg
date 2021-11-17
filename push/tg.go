package push

import (
	"fmt"

	"github.com/guonaihong/gout"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Telegram struct {
	Token  string
	ChatId string
}

type TGMsg struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func (t *Telegram) Init() func(kind string, message string) {
	return func(kind string, message string) {
		var resp []byte
		if kind == "markdown" {
			data := TGMsg{
				ChatID:    t.ChatId,
				Text:      message,
				ParseMode: "MarkdownV2",
			}
			log.Infoln(data)
			err := gout.GET(fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", t.Token)).BindBody(&resp).SetQuery(data).SetProxy("http://127.0.0.1:7890").Do()
			if err != nil {
				return
			}
			log.Infoln("向tg推送消息成功")
		} else if kind == "html" {
			data := TGMsg{
				ChatID:    t.ChatId,
				Text:      message,
				ParseMode: "HTML",
			}
			log.Infoln(data)
			err := gout.POST(fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", t.Token)).BindBody(&resp).SetProxy("http://127.0.0.1:7890").SetJSON(gout.H{
				"chat_id":    t.ChatId,
				"text":       message,
				"parse_mode": "HTML",
			}).Do()
			if err != nil {
				return
			}
			log.Infoln("向tg推送消息成功")
		}
		log.Infoln(gjson.GetBytes(resp, "@this|@pretty").String())
	}
}
