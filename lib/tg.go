package lib

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

var (
	handles sync.Map
)

func init() {
	newPlugin("/login", login)
	newPlugin("/get_users", getAllUser)
	newPlugin("/study", study)
	newPlugin("/get_scores", getScores)
}

//Telegram
// @Description:
//
type Telegram struct {
	Token  string
	ChatId int64
	bot    *tgbotapi.BotAPI
	Proxy  string
}

type Handle interface {
	getCommand() string
	execute(bot *Telegram, args []string)
}

type Mather struct {
	command string
	handle  func(bot *Telegram, args []string)
}

func (m Mather) getCommand() string {
	return m.command
}

func (m Mather) execute(bot *Telegram, args []string) {
	m.handle(bot, args)
}

func newPlugin(command string, handle func(bot *Telegram, args []string)) {
	handles.Store(command, handle)
}

//Init
/**
 * @Description:
 * @receiver t
 * @return func(kind string, message string)
 */
func (t *Telegram) Init() {
	uri, err := url.Parse(t.Proxy)
	t.bot, err = tgbotapi.NewBotAPIWithClient(t.Token, tgbotapi.APIEndpoint, &http.Client{Transport: &http.Transport{
		// 设置代理
		Proxy: http.ProxyURL(uri),
	}})

	if err != nil {
		log.Errorln("telegram token鉴权失败或代理使用失败")
		log.Errorln(err.Error())
	}

	channel := t.bot.GetUpdatesChan(tgbotapi.NewUpdate(1))

	go func() {
		for {
			update := <-channel
			handles.Range(func(key, value interface{}) bool {
				if strings.Split(update.Message.Text, " ")[0] == key.(string) {
					go func() {
						defer func() {
							err := recover()
							if err != nil {

							}
						}()
						(value.(func(bot *Telegram, args []string)))(t, strings.Split(update.Message.Text, " ")[1:])
					}()
				}
				return true
			})
		}
	}()

	_, err = t.bot.Request(tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "login", Description: "登录一个账号"},
		tgbotapi.BotCommand{Command: "get_users", Description: "获取所有cookie有效的用户"},
		tgbotapi.BotCommand{Command: "study", Description: "对一个账户进行学习"},
		tgbotapi.BotCommand{Command: "get_scores", Description: "获取用户成绩"},
	))
	if err != nil {
		return
	}
}

func (t *Telegram) SendPhoto(image []byte) {
	photo := tgbotapi.NewPhoto(t.ChatId, tgbotapi.FileBytes{
		Name:  "login code",
		Bytes: image,
	})
	_, err := t.bot.Send(photo)
	if err != nil {
		log.Errorln("发送图片信息失败")
		log.Errorln(err.Error())
		return
	}
}

func (t *Telegram) SendMsg(message string) {
	msg := tgbotapi.NewMessage(t.ChatId, message)
	t.bot.Send(msg)
}

func login(bot *Telegram, args []string) {
	log.Infoln(args)
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Errorln(err)
			}
		}()
		core := Core{
			pw:          nil,
			browser:     nil,
			context:     nil,
			ShowBrowser: false,
			Push: func(kind string, message string) {
				if kind == "image" {
					bytes, _ := base64.StdEncoding.DecodeString(message)
					bot.SendPhoto(bytes)
				} else if kind == "markdown" {
					newMessage := tgbotapi.NewMessage(bot.ChatId, message)
					newMessage.ParseMode = tgbotapi.ModeMarkdownV2
					bot.bot.Send(newMessage)
				} else {
					bot.SendMsg(message)
				}
			},
		}
		core.Init()
		defer core.Quit()
		_, err := core.Login()
		if err != nil {
			bot.SendMsg(err.Error())
			return
		}
		bot.SendMsg("登录成功")
	}()
}

func getAllUser(bot *Telegram, args []string) {
	users, err := GetUsers()
	if err != nil {
		bot.SendMsg("获取用户失败" + err.Error())
		return
	}
	message := fmt.Sprintf("共获取到%v个有效用户信息\n", len(users))
	for i, user := range users {
		message += fmt.Sprintf("%v   %v", i, user.Nick)
		message += "\n"
	}
	bot.SendMsg(message)
}

func study(bot *Telegram, args []string) {
	users, err := GetUsers()
	if err != nil {
		bot.SendMsg(err.Error())
		return
	}
	var cookies []Cookie
	if len(users) == 1 {
		bot.SendMsg("仅存在一名用户信息，自动进行学习")
		cookies = users[0].Cookies
	} else if len(users) == 0 {
		bot.SendMsg("未发现用户信息，请输入/login进行用户登录")
		return
	} else {
		if len(args) < 0 {
			bot.SendMsg("存在多名用户，未输入用户序号")
			return
		} else {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				bot.SendMsg(err.Error())
				return
			}
			cookies = users[i].Cookies
		}
	}
	core := Core{
		pw:          nil,
		browser:     nil,
		context:     nil,
		ShowBrowser: true,
		Push: func(kind string, message string) {
			switch {
			case kind == "image":
				bytes, _ := base64.StdEncoding.DecodeString(message)
				bot.SendPhoto(bytes)
			case kind == "markdown":
				newMessage := tgbotapi.NewMessage(bot.ChatId, message)
				newMessage.ParseMode = tgbotapi.ModeMarkdownV2
				_, _ = bot.bot.Send(newMessage)

			default:
				bot.SendMsg(message)
			}
		},
	}
	core.Init()
	defer core.Quit()
	core.LearnArticle(cookies)
	core.LearnVideo(cookies)
	core.RespondDaily(cookies, "daily")
	core.RespondDaily(cookies, "daily")
	core.RespondDaily(cookies, "weekly")
	core.RespondDaily(cookies, "special")
}

func getScores(bot *Telegram, args []string) {
	users, err := GetUsers()
	if err != nil {
		log.Errorln(err.Error())
		bot.SendMsg("获取用户信息失败" + err.Error())
		return
	}
	message := fmt.Sprintf("共获取到%v个有效用户信息\n", len(users))
	for _, user := range users {
		message += user.Nick + "\n"
		score, err := GetUserScore(user.Cookies)
		if err != nil {
			message += err.Error() + "\n"
		}
		message += PrintScore(score) + "\n"
	}
	bot.SendMsg(message)
}
