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
	"github.com/huoxue1/study_xxqg/push"
)

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

	log.SetLevel(level)
}

var (
	config lib.Config
)

func init() {
	_, err := os.Stat(`./config/`)
	if err != nil {
		os.Mkdir("./config/", 0666)
		return
	}
}

func main() {
	switch {
	case config.Cron != "":
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
			do()
		})
		if err != nil {
			log.Errorln(err.Error())
			return
		}
		c.Start()
		select {}
	case config.TG.Enable:
		log.Infoln("已采用tg交互模式")
		telegram := lib.Telegram{
			Token:  config.TG.Token,
			ChatId: config.TG.ChatID,
			Proxy:  config.TG.Proxy,
		}
		telegram.Init()
		select {}
	default:
		log.Infoln("已采用普通学习模式")
		do()
	}
}

func do() {
	log.Infoln(` 刷课模式，默认为1，
 1：只刷文章何视频
 2：只刷文章和视频和每日答题
 3：刷文章和视频和每日答题每周答题和专项答题`)
	log.Infoln("检测到模式", config.Model)

	getPush := push.GetPush(config)
	core := lib.Core{ShowBrowser: config.ShowBrowser, Push: getPush}
	defer core.Quit()
	core.Init()
	var cookies []lib.Cookie
	users, _ := lib.GetUsers()
	switch {
	case len(users) < 1:
		log.Infoln("未检测到有效用户信息，将采用登录模式")
		cookies, _ = core.Login()
	case len(users) == 1:
		log.Infoln("检测到1位有效用户信息，采用默认用户")
		cookies = users[0].Cookies
		log.Infoln("已选择用户: ", users[0].Nick)
	default:
		for i, user := range users {
			log.Infoln("序号：", i+1, "   ===> ", user.Nick)
		}
		log.Infoln("请输入对应序号选择对应账户")
		var i int
		_, _ = fmt.Scanln(&i)
		cookies = users[i-1].Cookies
		log.Infoln("已选择用户: ", users[i-1].Nick)
	}

	core.LearnArticle(cookies)
	core.LearnVideo(cookies)
	if config.Model == 2 {
		core.RespondDaily(cookies, "daily")
	} else if config.Model == 3 {
		core.RespondDaily(cookies, "daily")
		core.RespondDaily(cookies, "weekly")
		core.RespondDaily(cookies, "special")
	}
	score, err := lib.GetUserScore(cookies)
	if err != nil {
		log.Errorln("获取成绩失败")
		log.Debugln(err.Error())
		return
	}
	message := "学习完成：今日得分:" + strconv.Itoa(score.TodayScore)
	core.Push("markdown", message)
}
