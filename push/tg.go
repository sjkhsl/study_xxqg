package push

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

//Telegram
// @Description:
//
type Telegram struct {
	Token  string
	ChatId string
}

//TGMsg
// @Description:
//
type TGMsg struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

//Init
/**
 * @Description:
 * @receiver t
 * @return func(kind string, message string)
 */
func (t *Telegram) Init() func(kind string, message string) {
	uri, err := url.Parse("http://127.0.0.1:7890")
	bot, err := tgbotapi.NewBotAPIWithClient(t.Token, tgbotapi.APIEndpoint, &http.Client{Transport: &http.Transport{
		// 设置代理
		Proxy: http.ProxyURL(uri),
	}})

	if err != nil {
		log.Errorln("telegram token鉴权失败")
		return func(kind string, message string) {}
	}
	chatId, err := strconv.ParseInt(t.ChatId, 10, 64)
	if err != nil {
		return func(kind string, message string) {}
	}
	return func(kind string, message string) {
		if kind == "image" {
			bytes, _ := base64.StdEncoding.DecodeString(message)
			photo := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{
				Name:  "123",
				Bytes: bytes,
			})
			_, err := bot.Send(photo)
			if err != nil {
				log.Errorln("发送图片信息失败")
				log.Errorln(err.Error())
				return
			}
		}

		mess := tgbotapi.NewMessage(chatId, message)
		mess.ParseMode = tgbotapi.ModeMarkdownV2
		_, err := bot.Send(mess)
		if err != nil {
			log.Errorln("发送消息失败")
			log.Errorln(err.Error())
			return
		}
	}
}
