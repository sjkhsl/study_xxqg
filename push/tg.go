package push

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/lib"
	"github.com/johlanse/study_xxqg/lib/state"
	"github.com/johlanse/study_xxqg/model"
	"github.com/johlanse/study_xxqg/utils"
	"github.com/johlanse/study_xxqg/utils/update"
)

var (
	handles sync.Map
	tgPush  func(id string, kind string, message string)
)

func TgInit() {
	defer func() {
		err := recover()
		if err != nil {
			tgPush = func(id string, kind string, message string) {

			}
		}
	}()
	config := conf.GetConfig()
	log.Infoln("已采用tg交互模式")
	telegram := Telegram{
		Token:  config.TG.Token,
		ChatId: config.TG.ChatID,
		Proxy:  config.TG.Proxy,
	}
	tgPush = func(id string, kind string, message string) {
		defer func() {
			err := recover()
			if err != nil {
				log.Errorln("推送tg消息出现错误")
				log.Errorln(err)
			}
		}()
		chatId := telegram.ChatId
		id1, err := strconv.Atoi(id)
		if err == nil {
			chatId = int64(id1)
		} else {
			log.Warningln("转化pushID错误，将发送给默认用户")
		}
		if kind == "flush" {
			telegram.SendMsg(chatId, strings.ReplaceAll(message, "</br>", "\n"))
		} else if kind == "image" {
			bytes, _ := base64.StdEncoding.DecodeString(message)
			telegram.SendPhoto(chatId, bytes)
		} else {
			if log.GetLevel() == log.DebugLevel {
				telegram.SendMsg(chatId, message)
			}
		}
	}

	telegram.Init()
}

// Telegram
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
	handle  func(bot *Telegram, from int64, args []string)
}

func (m Mather) getCommand() string {
	return m.command
}

func (m Mather) execute(bot *Telegram, from int64, args []string) {
	m.handle(bot, from, args)
}

func newPlugin(command string, handle func(bot *Telegram, from int64, args []string)) {
	handles.Store(command, handle)
}

// Init
/**
 * @Description:
 * @receiver t
 * @return func(kind string, message string)
 */
func (t *Telegram) Init() {

	newPlugin("/login", login)
	newPlugin("/get_users", getAllUser)
	newPlugin("/study", study)
	newPlugin("/get_scores", getScores)
	newPlugin("/quit", quit)
	newPlugin("/study_all", studyAll)
	newPlugin("/delete", deleteUser)
	newPlugin("/version", checkVersion)
	newPlugin("/update", botUpdate)
	newPlugin("/restart", botRestart)
	newPlugin("/get_fail_users", getFailUser)
	var err error
	var uri *url.URL
	if t.Proxy != "" {
		uri, err = url.Parse(t.Proxy)
		if err != nil {
			log.Errorln("代理解析失败" + err.Error())
			err = nil
		}
	}
	t.bot, err = tgbotapi.NewBotAPIWithClient(t.Token, conf.GetConfig().TG.CustomApi+"/bot%s/%s", &http.Client{Transport: &http.Transport{
		// 设置代理
		Proxy: func(r *http.Request) (*url.URL, error) {
			if uri != nil {
				return uri, nil
			} else {
				return http.ProxyFromEnvironment(r)
			}
		},
		//TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}})

	if err != nil {
		log.Errorln("telegram token鉴权失败或代理使用失败")
		log.Errorln(err.Error())
	}

	channel := t.bot.GetUpdatesChan(tgbotapi.NewUpdate(1))
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Errorln("处理tg消息时发生异常，请尝试重启程序")
				return
			}
		}()
		for {
			update := <-channel
			if update.Message == nil {
				if update.CallbackQuery != nil {
					update.Message = update.CallbackQuery.Message
					update.Message.Text = update.CallbackQuery.Data
					t.bot.Send(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
				} else {
					data, _ := json.Marshal(update)
					log.Infoln(string(data))
				}
			}

			if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
				update.Message.Text = strings.ReplaceAll(update.Message.Text, "@"+t.bot.Self.UserName, "")
			}
			log.Infoln(fmt.Sprintf("收到tg消息,来自%v,内容 ==》 %v", update.Message.Chat.ID, update.Message.Text))
			if len(conf.GetConfig().TG.WhiteList) > 0 {
				inWhiteList := false
				for _, id := range conf.GetConfig().TG.WhiteList {
					if id == update.Message.Chat.ID {
						inWhiteList = true
						break
					}
				}
				if !inWhiteList {
					log.Warningln("已过滤非白名单的消息,若需允许用户使用，请将user_id添加到配置文件white_list中")
					continue
				}
			}
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
						(value.(func(bot *Telegram, from int64, args []string)))(t, update.Message.Chat.ID, strings.Split(update.Message.Text, " ")[1:])
					}()
				}
				return true
			})
		}
	}()

	_, err = t.bot.Request(tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "login", Description: "登录一个账号"},
		tgbotapi.BotCommand{Command: "get_users", Description: "获取所有cookie有效的用户"},
		tgbotapi.BotCommand{Command: "get_fail_users", Description: "获取所有cookie失效的用户"},
		tgbotapi.BotCommand{Command: "study", Description: "对一个账户进行学习"},
		tgbotapi.BotCommand{Command: "get_scores", Description: "获取用户成绩"},
		tgbotapi.BotCommand{Command: "quit", Description: "退出所有正在学习的实例,或者跟上实例ID退出对应实例"},
		tgbotapi.BotCommand{Command: "study_all", Description: "对当前所有用户进行按顺序学习"},
		tgbotapi.BotCommand{Command: "delete", Description: "删除选中的用户"},
		tgbotapi.BotCommand{Command: "version", Description: "获取程序当前的版本"},
		tgbotapi.BotCommand{Command: "restart", Description: "重启程序！"},
		tgbotapi.BotCommand{Command: "update", Description: "更新程序"},
	))
	if err != nil {
		return
	}
}

func (t *Telegram) SendPhoto(id int64, image []byte) {
	photo := tgbotapi.NewPhoto(id, tgbotapi.FileBytes{
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

func (t *Telegram) SendMsg(id int64, message string) int {
	if id == 0 {
		id = t.ChatId
	}
	msg := tgbotapi.NewMessage(id, message)
	messa, err := t.bot.Send(msg)
	if err != nil {
		return 0
	}
	return messa.MessageID
}

func getFailUser(bot *Telegram, from int64, args []string) {
	user, err := model.QueryFailUser()
	if err != nil {
		bot.SendMsg(from, err.Error())
		return
	}
	msg := "当前过期用户:\n"
	for _, u := range user {
		msg += u.Nick + "\n"
	}
	bot.SendMsg(from, "当前过期用户:\n"+msg)
}

//
//  checkVersion
//  @Description: 检查版本信息
//  @param bot
//  @param from
//  @param args
//
func checkVersion(bot *Telegram, from int64, args []string) {
	about := utils.GetAbout()
	bot.SendMsg(from, about)
}

//
//  botRestart
//  @Description: 重启程序
//  @param bot
//  @param from
//  @param args
//
func botRestart(bot *Telegram, from int64, args []string) {
	if from != conf.GetConfig().TG.ChatID {
		bot.SendMsg(from, "请联系管理员解决！！")
		return
	}
	bot.SendMsg(from, "即将重启程序！！！")
	utils.Restart()
}

//
//  botUpdate
//  @Description: 更新程序
//  @param bot
//  @param from
//  @param args
//
func botUpdate(bot *Telegram, from int64, args []string) {
	if from != conf.GetConfig().TG.ChatID {
		bot.SendMsg(from, "请联系管理员解决！！")
		return
	}
	bot.SendMsg(from, "即将更新程序！！")
	update.SelfUpdate("", conf.GetVersion())
	bot.SendMsg(from, "更新完成，即将重启程序！")
	utils.Restart()
}

func login(bot *Telegram, from int64, args []string) {
	config := conf.GetConfig()
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Errorln(err)
			}
		}()
		core := lib.Core{
			ShowBrowser: config.ShowBrowser,
			Push: func(id string, kind string, message string) {
				switch {
				case kind == "image":
					bytes, _ := base64.StdEncoding.DecodeString(message)
					bot.SendPhoto(from, bytes)
				case kind == "markdown":
					newMessage := tgbotapi.NewMessage(bot.ChatId, message)
					newMessage.ParseMode = tgbotapi.ModeMarkdownV2
					bot.bot.Send(newMessage)
				case kind == "text":
					if log.GetLevel() == log.DebugLevel {
						bot.SendMsg(from, message)
					}
				case kind == "flush":
					bot.SendMsg(from, message)
				}
			},
		}
		core.Init()
		defer core.Quit()
		_, err := core.L(config.Retry.Times, strconv.Itoa(int(from)))
		if err != nil {
			bot.SendMsg(from, err.Error())
			return
		}
		bot.SendMsg(from, "登录成功")
	}()
}

func getAllUser(bot *Telegram, from int64, args []string) {
	users, err := model.Query()
	if err != nil {
		bot.SendMsg(from, "获取用户失败"+err.Error())
		return
	}
	message := fmt.Sprintf("共获取到%v个有效用户信息\n", len(users))
	for i, user := range users {
		message += fmt.Sprintf("%v   %v", i, user.Nick)
		message += "\n"
	}
	bot.SendMsg(from, message)
}

func studyAll(bot *Telegram, from int64, args []string) {
	config := conf.GetConfig()
	users, err := model.Query()
	if err != nil {
		bot.SendMsg(from, err.Error())
		return
	}
	if len(users) == 0 {
		bot.SendMsg(from, "未发现用户信息，请输入/login进行用户登录")
		return
	}
	getAllUser(bot, from, args)
	for _, user := range users {
		s := func() {
			core := lib.Core{
				ShowBrowser: config.ShowBrowser,
				Push: func(id string, kind string, message string) {
					switch {
					case kind == "image":
						bytes, _ := base64.StdEncoding.DecodeString(message)
						bot.SendPhoto(from, bytes)
					case kind == "markdown":
						newMessage := tgbotapi.NewMessage(bot.ChatId, message)
						newMessage.ParseMode = tgbotapi.ModeMarkdownV2
						_, _ = bot.bot.Send(newMessage)

					case kind == "text":
						if log.GetLevel() == log.DebugLevel {
							bot.SendMsg(from, message)
						}
					case kind == "flush":
						bot.SendMsg(from, message)
					}
				},
			}

			timer := time.After(time.Minute * 30)
			c := make(chan int, 1)
			go func() {

				bot.SendMsg(from, "已创建运行实例："+user.Uid)
				state.Add(user.Uid, &core)
				defer state.Delete(user.Uid)
				core.Init()
				defer core.Quit()
				core.LearnArticle(user)
				core.LearnVideo(user)
				core.RespondDaily(user, "daily")
				core.RespondDaily(user, "weekly")
				core.RespondDaily(user, "special")
				c <- 1
			}()

			select {
			case <-timer:
				{
					bot.SendMsg(from, "学习超时，请重新学习或检查日志")
					log.Errorln("学习超时，已自动退出")
					core.Quit()
				}
			case <-c:
				{
				}
			}
			score, _ := lib.GetUserScore(user.ToCookies())
			bot.SendMsg(from, fmt.Sprintf("%v已学习完成\n%v", user.Nick, lib.PrintScore(score)))
		}
		s()
	}
}

func deleteUser(bot *Telegram, from int64, args []string) {
	config := conf.GetConfig()
	if from != config.TG.ChatID {
		bot.SendMsg(from, "请联系管理员删除！")
		return
	}
	users, err := model.Query()
	if err != nil {
		bot.SendMsg(from, err.Error())
		return
	}
	failUser, err := model.QueryFailUser()
	if err != nil {
		bot.SendMsg(from, err.Error())
		return
	}
	if len(args) < 1 {
		msgID := bot.SendMsg(from, "请选择删除的用户")
		markup := tgbotapi.InlineKeyboardMarkup{}
		for i, user := range users {
			markup.InlineKeyboard = append(markup.InlineKeyboard, append([]tgbotapi.InlineKeyboardButton{}, tgbotapi.NewInlineKeyboardButtonData(user.Nick, "/delete "+strconv.Itoa(i))))
		}

		for i, user := range failUser {
			markup.InlineKeyboard = append(markup.InlineKeyboard, append([]tgbotapi.InlineKeyboardButton{}, tgbotapi.NewInlineKeyboardButtonData(user.Nick+"  (已失效)", "/delete "+strconv.Itoa(len(users)+i))))

		}

		replyMarkup := tgbotapi.NewEditMessageReplyMarkup(from, msgID, markup)
		_, err := bot.bot.Send(replyMarkup)
		if err != nil {
			return
		}
		return
	} else {
		users = append(users, failUser...)
		i, err := strconv.Atoi(args[0])
		if err != nil {
			bot.SendMsg(from, err.Error())
			return
		}
		if i >= len(users) {
			bot.SendMsg(from, "错误的序号")
			return
		}
		err = model.DeleteUser(users[i].Uid)
		if err != nil {
			bot.SendMsg(from, err.Error())
			return
		}
		bot.SendMsg(from, "删除用户"+users[i].Nick+"成功")
	}
}

func study(bot *Telegram, from int64, args []string) {
	config := conf.GetConfig()
	users, err := model.Query()
	if err != nil {
		bot.SendMsg(from, err.Error())
		return
	}
	var user *model.User
	switch {
	case len(users) == 1:
		bot.SendMsg(from, "仅存在一名用户信息，自动进行学习")
		user = users[0]
	case len(users) == 0:
		bot.SendMsg(from, "未发现用户信息，请输入/login进行用户登录")
		return
	default:
		if 0 < len(args) {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				bot.SendMsg(from, err.Error())
				return
			}
			user = users[i]
		} else {
			msgID := bot.SendMsg(from, "存在多名用户，未输入用户序号")
			markup := tgbotapi.InlineKeyboardMarkup{}
			for i, user := range users {
				markup.InlineKeyboard = append(markup.InlineKeyboard, append([]tgbotapi.InlineKeyboardButton{}, tgbotapi.NewInlineKeyboardButtonData(user.Nick, "/study "+strconv.Itoa(i))))
			}

			replyMarkup := tgbotapi.NewEditMessageReplyMarkup(from, msgID, markup)
			_, err := bot.bot.Send(replyMarkup)
			if err != nil {
				return
			}
			return
		}
	}
	core := lib.Core{
		ShowBrowser: config.ShowBrowser,
		Push: func(id string, kind string, message string) {
			switch {
			case kind == "image":
				bytes, _ := base64.StdEncoding.DecodeString(message)
				bot.SendPhoto(from, bytes)
			case kind == "markdown":
				newMessage := tgbotapi.NewMessage(bot.ChatId, message)
				newMessage.ParseMode = tgbotapi.ModeMarkdownV2
				_, _ = bot.bot.Send(newMessage)

			case kind == "text":
				if log.GetLevel() == log.DebugLevel {
					bot.SendMsg(from, message)
				}
			case kind == "flush":
				bot.SendMsg(from, message)
			}
		},
	}
	timer := time.After(time.Minute * 30)
	c := make(chan int, 1)
	go func() {
		bot.SendMsg(from, "已创建运行实例："+user.Uid)
		state.Add(user.Uid, &core)
		defer state.Delete(user.Uid)
		core.Init()
		defer core.Quit()

		core.LearnArticle(user)

		core.LearnVideo(user)

		if config.Model == 2 {
			core.RespondDaily(user, "daily")
		} else if config.Model == 3 {
			core.RespondDaily(user, "daily")
			core.RespondDaily(user, "weekly")
			core.RespondDaily(user, "special")
		} else if config.Model == 4 {
			core.RespondDaily(user, "special")
		}

		c <- 1
	}()

	select {
	case <-timer:
		{
			log.Errorln("学习超时，已自动退出")
			bot.SendMsg(from, "学习超时，请重新登录或检查日志")
			core.Quit()
		}
	case <-c:
		{

		}
	}
	score, _ := lib.GetUserScore(user.ToCookies())
	bot.SendMsg(from, fmt.Sprintf("%v已学习完成\n%v", user.Nick, lib.PrintScore(score)))
}

func getScores(bot *Telegram, from int64, args []string) {
	users, err := model.Query()
	if err != nil {
		log.Errorln(err.Error())
		bot.SendMsg(from, "获取用户信息失败"+err.Error())
		return
	}
	message := ""
	for _, user := range users {
		message += user.Nick + "\n"
		score, err := lib.GetUserScore(user.ToCookies())
		if err != nil {
			message += err.Error() + "\n"
			continue
		}
		message += lib.FormatScore(score) + "\n"
		bot.SendMsg(from, message)
		message = ""
	}
}

func quit(bot *Telegram, from int64, args []string) {
	if len(args) < 1 {
		state.Range(func(key, value interface{}) bool {
			bot.SendMsg(from, "已退出运行实例"+key.(string))
			core := value.(*lib.Core)
			core.Quit()
			return true
		})
	} else {
		state.Range(func(key, value interface{}) bool {
			if key.(string) == args[0] {
				core := value.(*lib.Core)
				core.Quit()
			}
			return true
		})
	}
}
