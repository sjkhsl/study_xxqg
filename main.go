package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	rotates "github.com/lestrrat-go/file-rotatelogs"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"

	"github.com/johlanse/study_xxqg/conf"
	// "github.com/johlanse/study_xxqg/gui"
	"github.com/johlanse/study_xxqg/lib"
	"github.com/johlanse/study_xxqg/model"
	"github.com/johlanse/study_xxqg/push"
	"github.com/johlanse/study_xxqg/utils/update"
	"github.com/johlanse/study_xxqg/web"
)

var (
	u bool
	i bool
)

var VERSION = "unknown"

func init() {
	flag.BoolVar(&u, "u", false, "update the study_xxqg")
	flag.BoolVar(&i, "init", false, "init the app")
	flag.Parse()

	config = conf.GetConfig()
	logFormatter := &easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%time%] [%lvl%]: %msg% \n",
	}
	w, err := rotates.New(path.Join("config", "logs", "%Y-%m-%d.log"), rotates.WithRotationTime(time.Hour*24))
	if err != nil {
		log.Errorf("rotates init err: %v", err)
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(w, os.Stdout)
	log.SetOutput(io.MultiWriter(w, os.Stdout))
	log.SetFormatter(logFormatter)
	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.SetLevel(log.DebugLevel)
	}
	log.SetLevel(level)
}

func init() {
	pid := os.Getpid()
	pi := strconv.Itoa(pid)
	err := os.WriteFile("pid.pid", []byte(pi), 0666)
	if err != nil {
		log.Errorln("pid写入失败")
		return
	}
}

var (
	config conf.Config
)

func init() {
	_, err := os.Stat(`./config/`)
	if err != nil {
		os.Mkdir("./config/", 0666) //nolint:errcheck
		return
	}
}

func main() {
	conf.SetVersion(VERSION)
	log.Infoln("当前程序运行版本： " + VERSION)
	if i {
		core := &lib.Core{}
		core.Init()
		core.Quit()
		return
	}

	go update.CheckUpdate(VERSION)

	if u {
		update.SelfUpdate("", VERSION)
		log.Infoln("请重启应用")
		os.Exit(1)
	}

	engine := web.RouterInit()
	go func() {
		h := http.NewServeMux()
		if config.Web.Enable {
			log.Infoln(fmt.Sprintf("已开启web配置，web监听地址 ==> %v:%v", config.Web.Host, config.Web.Port))
			h.Handle("/", engine)
		}
		if config.Wechat.Enable {
			log.Infoln(fmt.Sprintf("已开启wechat公众号配置,监听地址： ==》 %v:%v", config.Web.Host, config.Web.Port))
			h.HandleFunc("/wx", web.HandleWechat)
		}
		if config.Web.Enable || config.Wechat.Enable {
			err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Web.Host, config.Web.Port), h)
			if err != nil {
				return
			}
		}
	}()

	if config.StartWait > 0 {
		log.Infoln(fmt.Sprintf("将等待%d秒后启动程序", config.StartWait))
		time.Sleep(time.Second * time.Duration(config.StartWait))
	}

	if config.Cron != "" {
		go func() {
			defer func() {
				err := recover()
				if err != nil {
					log.Errorln("定时任务执行出现问题")
					log.Errorln(err)
				}
			}()
			log.Infoln("已采用定时执行模式")
			c := cron.New()

			_, err := c.AddFunc(config.Cron, func() {
				defer func() {
					i := recover()
					if i != nil {
						log.Errorln(i)
						log.Errorln("执行定时任务出现异常")
					}
				}()
				log.Infoln("即将开始执行定时任务")
				// 检测是否开启了随机等待
				if config.CronRandomWait > 0 {
					rand.Seed(time.Now().UnixNano())
					r := rand.Intn(config.CronRandomWait)
					log.Infoln(fmt.Sprintf("随机延迟%d分钟", r))
					time.Sleep(time.Duration(r) * time.Minute)
				}
				do("cron")
			})
			if err != nil {
				log.Errorln(err.Error())
				return
			}
			c.Start()
		}()
	}

	if config.TG.Enable {
		go func() {
			defer func() {
				err := recover()
				if err != nil {
					log.Errorln("TG模式执行出现问题")
					log.Errorln(err)
				}
			}()
			log.Infoln("已采用tg交互模式")
			telegram := lib.Telegram{
				Token:  config.TG.Token,
				ChatId: config.TG.ChatID,
				Proxy:  config.TG.Proxy,
			}
			telegram.Init()
		}()
	}

	if !config.TG.Enable && config.Cron == "" && !config.Wechat.Enable {
		log.Infoln("已采用普通学习模式")
		do("normal")
	} else {
		// gui.InitWindow()
		select {}
	}
}

func do(m string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("do 方法执行错误")
			log.Errorln(err)
		}
	}()

	log.Infoln(` 刷课模式，默认为1，
 1：只刷文章何视频
 2：只刷文章和视频和每日答题
 3：刷文章和视频和每日答题每周答题和专项答题`)
	log.Infoln("检测到模式", config.Model)

	getPush := push.GetPush(config)
	core := &lib.Core{ShowBrowser: config.ShowBrowser, Push: getPush}
	defer core.Quit()
	core.Init()
	var user *model.User
	users, _ := model.Query()
	study := func(core2 *lib.Core, u *model.User) {
		defer func() {
			err := recover()
			if err != nil {
				log.Errorln("学习过程异常")
				log.Errorln(err)
			}
		}()

		core2.LearnArticle(u)
		core2.LearnVideo(u)
		if config.Model == 2 {
			core2.RespondDaily(u, "daily")
		} else if config.Model == 3 {
			core2.RespondDaily(u, "daily")
			core2.RespondDaily(u, "weekly")
			core2.RespondDaily(u, "special")
		}
		score, err := lib.GetUserScore(u.ToCookies())
		if err != nil {
			log.Errorln("获取成绩失败")
			log.Debugln(err.Error())
			return
		}
		message := u.Nick + " 学习完成：今日得分:" + strconv.Itoa(score.TodayScore)
		score, _ = lib.GetUserScore(user.ToCookies())
		content := lib.FormatScore(score)
		err = push.PushMessage(user.Nick+"学习情况", user.Nick+"学习情况"+content, "score", user.PushId)
		if err != nil {
			log.Errorln(err.Error())
			err = nil
		}
		core2.Push("markdown", message)
		core2.Push("flush", "")
	}

	// 用户小于1时自动登录
	if len(users) < 1 {
		log.Infoln("未检测到有效用户信息，将采用登录模式")
		u, err := core.L(config.Retry.Times, "")
		if err != nil {
			log.Errorln(err.Error())
			return
		}
		user = u
	} else {
		// 如果为定时模式则直接循环所以用户依次运行
		if m == "cron" {
			for _, u := range users {
				study(core, u)
			}
			if len(users) < 1 {
				user, err := core.L(config.Retry.Times, "")
				if err != nil {
					core.Push("msg", "登录超时")
					return
				}
				study(core, user)
			}
			return
		}

		for i, user := range users {
			log.Infoln("序号：", i+1, "   ===> ", user.Nick)
		}
		log.Infoln("请输入对应序号选择对应账户，输入0添加用户：")

		inputChan := make(chan int, 1)
		go func(c chan int) {
			var i int
			_, _ = fmt.Scanln(&i)
			c <- i
		}(inputChan)

		var i int
		select {
		case i = <-inputChan:
			log.Infoln("已获取到输入")
		case <-time.After(time.Minute):
			log.Errorln("获取输入超时，默认选择第一个用户")
			if len(users) < 1 {
				return
			} else {
				i = 1
			}
		}

		if i == 0 {
			u, err := core.L(config.Retry.Times, "")
			if err != nil {
				log.Errorln(err.Error())
				return
			}
			user = u
		} else {
			user = users[i-1]
			log.Infoln("已选择用户: ", users[i-1].Nick)
		}
	}

	study(core, user)
	core.Push("flush", "")
}
