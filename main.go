package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	rotates "github.com/lestrrat-go/file-rotatelogs"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"

	"github.com/huoxue1/study_xxqg/lib"
	"github.com/huoxue1/study_xxqg/model"
	"github.com/huoxue1/study_xxqg/push"
	"github.com/huoxue1/study_xxqg/web"
)

var VERSION = "unknown"

func init() {
	config = lib.GetConfig()
	logFormatter := &easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%time%] [%lvl%]: %msg% \n",
	}
	w, err := rotates.New(path.Join("logs", "%Y-%m-%d.log"), rotates.WithRotationTime(time.Hour*24))
	if err != nil {
		log.Errorf("rotates init err: %v", err)
		panic(err)
	}
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
	config lib.Config
)

func init() {
	_, err := os.Stat(`./config/`)
	if err != nil {
		os.Mkdir("./config/", 0666) //nolint:errcheck
		return
	}
}

func main() {
	if config.Web.Enable {
		engine := web.RouterInit()
		go func() {
			err := engine.Run(fmt.Sprintf("%s:%d", config.Web.Host, config.Web.Port))
			if err != nil {
				return
			}
		}()
	}

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
				do("cron")
			})
			if err != nil {
				log.Errorln(err.Error())
				return
			}
			c.Start()
			select {}
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
			select {}
		}()
	}

	if !config.TG.Enable && config.Cron == "" {
		log.Infoln("已采用普通学习模式")
		do("normal")
	} else {
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
	core := lib.Core{ShowBrowser: config.ShowBrowser, Push: getPush}
	defer core.Quit()
	core.Init()
	var user *model.User
	users, _ := model.Query()

	// 用户小于1时自动登录
	if len(users) < 1 {
		log.Infoln("未检测到有效用户信息，将采用登录模式")
		u, err := core.L(config.Retry.Times)
		if err != nil {
			log.Errorln(err.Error())
			return
		}
		user = u
	} else {
		// 如果为定时模式则直接循环所以用户依次运行
		if m == "cron" {
			for _, u := range users {
				go core.LearnArticle(u)
				go core.LearnVideo(u)
				lib.WaitStudy(u, "")
				if config.Model == 2 {
					core.RespondDaily(u, "daily")
				} else if config.Model == 3 {
					core.RespondDaily(u, "daily")
					core.RespondDaily(u, "weekly")
					core.RespondDaily(u, "special")
				}
				score, err := lib.GetUserScore(u.ToCookies())
				if err != nil {
					log.Errorln("获取成绩失败")
					log.Debugln(err.Error())
					return
				}
				message := u.Nick + " 学习完成：今日得分:" + strconv.Itoa(score.TodayScore)
				core.Push("markdown", message)
			}
			return
		}

		for i, user := range users {
			log.Infoln("序号：", i+1, "   ===> ", user.Nick)
		}
		log.Infoln("请输入对应序号选择对应账户，输入0添加用户：")
		var i int
		_, _ = fmt.Scanln(&i)
		if i == 0 {
			u, err := core.L(config.Retry.Times)
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

	// switch {
	// case len(users) < 1:
	//	log.Infoln("未检测到有效用户信息，将采用登录模式")
	//	u, err := core.L(config.Retry.Times)
	//	if err != nil {
	//		log.Errorln(err.Error())
	//		return
	//	}
	//	user = u
	//case len(users) == 1:
	//	log.Infoln("检测到1位有效用户信息，采用默认用户")
	//	user = users[0]
	//	log.Infoln("已选择用户: ", users[0].Nick)
	//default:
	//	if m == "cron" {
	//
	//	}
	//	for i, user := range users {
	//		log.Infoln("序号：", i+1, "   ===> ", user.Nick)
	//	}
	//	log.Infoln("请输入对应序号选择对应账户")
	//	var i int
	//	_, _ = fmt.Scanln(&i)
	//	user = users[i-1]
	//	log.Infoln("已选择用户: ", users[i-1].Nick)
	//}

	go core.LearnArticle(user)
	go core.LearnVideo(user)
	lib.WaitStudy(user, "")
	if config.Model == 2 {
		core.RespondDaily(user, "daily")
	} else if config.Model == 3 {
		core.RespondDaily(user, "daily")
		core.RespondDaily(user, "weekly")
		core.RespondDaily(user, "special")
	}

	score, err := lib.GetUserScore(user.ToCookies())
	if err != nil {
		log.Errorln("获取成绩失败")
		log.Debugln(err.Error())
		return
	}
	message := "学习完成：今日得分:" + strconv.Itoa(score.TodayScore)
	core.Push("markdown", message)
	core.Push("flush", "学习完成")
}
