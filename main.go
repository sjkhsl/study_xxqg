package main

import (
	"os"
	"path"
	"time"

	rotates "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"

	"github.com/huoxue1/study_xxqg/lib"
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
	log.SetOutput(w)
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
	log.Infoln(`// 刷课模式，默认为1，
 1：只刷文章何视频
 2：只刷文章和视频和每日答题
 3：刷文章和视频和每日答题每周答题和专项答题`)
	log.Infoln("检测到模式", config.Model)
	core := lib.Core{ShowBrowser: config.ShowBrowser}
	defer core.Quit()
	core.Init()
	login, err := core.Login()
	if err != nil {
		return
	}

	core.LearnArticle(login)
	core.LearnVideo(login)
	if config.Model == 2 {
		core.RespondDaily(login, "daily")
	} else if config.Model == 3 {
		core.RespondDaily(login, "daily")
		core.RespondDaily(login, "weekly")
		core.RespondDaily(login, "special")
	}
}
