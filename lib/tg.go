package lib

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	handles sync.Map
	datas   sync.Map
)

func init() {
	newPlugin("/login", login)
	newPlugin("/get_users", getAllUser)
	newPlugin("/study", study)
	newPlugin("/get_scores", getScores)
	newPlugin("/quit", quit)
	newPlugin("/study_all", studyAll)
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
								log.Errorln(err)
								log.Errorln("handle执行出现了不可挽回的错误")
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
		tgbotapi.BotCommand{Command: "quit", Description: "退出所有正在学习的实例,或者跟上实例ID退出对应实例"},
		tgbotapi.BotCommand{Command: "study_all", Description: "对当前所有用户进行按顺序学习"},
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
	config := GetConfig()
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
			ShowBrowser: config.ShowBrowser,
			Push: func(kind string, message string) {
				switch {
				case kind == "image":
					bytes, _ := base64.StdEncoding.DecodeString(message)
					bot.SendPhoto(bytes)
				case kind == "markdown":
					newMessage := tgbotapi.NewMessage(bot.ChatId, message)
					newMessage.ParseMode = tgbotapi.ModeMarkdownV2
					bot.bot.Send(newMessage)
				default:
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

func studyAll(bot *Telegram, args []string) {
	config := GetConfig()
	users, err := GetUsers()
	if err != nil {
		bot.SendMsg(err.Error())
		return
	}
	if len(users) == 0 {
		bot.SendMsg("未发现用户信息，请输入/login进行用户登录")
		return
	}
	getAllUser(bot, args)
	for _, user := range users {
		s := func() {
			core := Core{
				pw:          nil,
				browser:     nil,
				context:     nil,
				ShowBrowser: config.ShowBrowser,
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

			timer := time.After(time.Minute * 30)
			c := make(chan int, 1)
			go func() {
				u := uuid.New().String()
				bot.SendMsg("已创建运行实例：" + u)
				datas.Store(u, &core)
				defer datas.Delete(u)
				core.Init()
				defer core.Quit()
				core.LearnArticle(user.Cookies)
				core.LearnVideo(user.Cookies)
				core.RespondDaily(user.Cookies, "daily")
				core.RespondDaily(user.Cookies, "weekly")
				core.RespondDaily(user.Cookies, "special")
				c <- 1
			}()

			select {
			case <-timer:
				{
					bot.SendMsg("学习超时，请重新登录或检查日志")
					log.Errorln("学习超时，已自动退出")
					core.Quit()
				}
			case <-c:
				{
				}
			}
			score, _ := GetUserScore(user.Cookies)
			bot.SendMsg(fmt.Sprintf("当前学习总积分：%v,今日得分：%v", score.TotalScore, score.TodayScore))
		}
		s()
	}
}

func study(bot *Telegram, args []string) {
	config := GetConfig()
	users, err := GetUsers()
	if err != nil {
		bot.SendMsg(err.Error())
		return
	}
	var cookies []Cookie
	switch {
	case len(users) == 1:
		bot.SendMsg("仅存在一名用户信息，自动进行学习")
		cookies = users[0].Cookies
	case len(users) == 0:
		bot.SendMsg("未发现用户信息，请输入/login进行用户登录")
		return
	default:
		if 0 <= len(args) {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				bot.SendMsg(err.Error())
				return
			}
			cookies = users[i].Cookies
		} else {
			bot.SendMsg("存在多名用户，未输入用户序号")
			return
		}
	}
	core := Core{
		pw:          nil,
		browser:     nil,
		context:     nil,
		ShowBrowser: config.ShowBrowser,
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
	timer := time.After(time.Minute * 30)
	c := make(chan int, 1)
	go func() {
		u := uuid.New().String()
		bot.SendMsg("已创建运行实例：" + u)
		datas.Store(u, &core)
		defer datas.Delete(u)
		core.Init()
		defer core.Quit()
		go core.LearnArticle(cookies)
		go core.LearnVideo(cookies)
		WaitStudy(&User{Cookies: cookies})
		core.RespondDaily(cookies, "daily")
		core.RespondDaily(cookies, "weekly")
		core.RespondDaily(cookies, "special")
		c <- 1
	}()

	select {
	case <-timer:
		{
			log.Errorln("学习超时，已自动退出")
			bot.SendMsg("学习超时，请重新登录或检查日志")
			core.Quit()
		}
	case <-c:
		{

		}
	}
	score, _ := GetUserScore(cookies)
	bot.SendMsg(fmt.Sprintf("当前学习总积分：%v,今日得分：%v", score.TotalScore, score.TodayScore))
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
		message += foramet_score(score) + "\n"
	}
	bot.SendMsg(message)
}

func quit(bot *Telegram, args []string) {
	if len(args) < 1 {
		datas.Range(func(key, value interface{}) bool {
			bot.SendMsg("已退出运行实例" + key.(string))
			core := value.(*Core)
			core.Quit()
			return true
		})
	} else {
		datas.Range(func(key, value interface{}) bool {
			if key.(string) == args[0] {
				core := value.(*Core)
				core.Quit()
			}
			return true
		})
	}
}
