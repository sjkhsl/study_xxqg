package main

import (
	"flag"
	"fmt"
	"github.com/johlanse/study_xxqg/utils/stop"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	rotates "github.com/lestrrat-go/file-rotatelogs"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	nested "github.com/Lyrics-you/sail-logrus-formatter/sailor"
	"github.com/huoxue1/xdaemon"

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/lib/state"
	"github.com/johlanse/study_xxqg/utils"
	// "github.com/johlanse/study_xxqg/gui"
	"github.com/johlanse/study_xxqg/lib"
	"github.com/johlanse/study_xxqg/model"
	"github.com/johlanse/study_xxqg/push"
	"github.com/johlanse/study_xxqg/utils/update"
	"github.com/johlanse/study_xxqg/web"
)

var (
	u          bool
	i          bool
	now        bool
	configPath = ""
)

var VERSION = "unknown"

func init() {
	if runtime.GOOS != "windows" {
		runBack()
	}
	fmt.Printf("\033[1;31;40m%s\033[0m\n\n", "******************************************************************")

	fmt.Printf("\033[1;31;40m%s\033[0m\n\n", "软件仅可用户学习和个人使用，禁止用于任何商业活动！！！！")

	fmt.Printf("\033[1;31;40m%s\033[0m\n\n", "软件仅可用户学习和个人使用，禁止用于任何商业活动！！！！")

	fmt.Printf("\033[1;31;40m%s\033[0m\n\n", "软件仅可用户学习和个人使用，禁止用于任何商业活动！！！！")

	fmt.Printf("\033[1;31;40m%s\033[0m\n\n", "******************************************************************")
	time.Sleep(3 * time.Second)

	flag.BoolVar(&u, "u", false, "更新应用")
	flag.BoolVar(&i, "init", false, "init the app")
	flag.BoolVar(&now, "now", false, "run cron now")
	flag.StringVar(&configPath, "config", "./config/config.yml", "设置配置文件路径")
	flag.Parse()
	// 初始化配置文件
	conf.InitConfig(configPath, utils.Restart)
	config = conf.GetConfig()
	w, err := rotates.New(path.Join("config", "logs", "%Y-%m-%d.log"), rotates.WithRotationTime(time.Hour*24))
	if err != nil {
		log.Errorf("rotates init err: %v", err)
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(w, &utils.LogWriter{})
	log.SetOutput(io.MultiWriter(w, os.Stdout))

	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.SetLevel(log.DebugLevel)
	}
	log.SetLevel(level)

	showPosition := false

	if level == log.DebugLevel {
		showPosition = true
	}
	log.SetFormatter(&nested.Formatter{
		FieldsOrder:           nil,
		TimeStampFormat:       "2006-01-02 15:04:05",
		CharStampFormat:       "",
		HideKeys:              false,
		Position:              showPosition,
		Colors:                true,
		FieldsColors:          true,
		FieldsSpace:           true,
		ShowFullLevel:         false,
		LowerCaseLevel:        true,
		TrimMessages:          true,
		CallerFirst:           false,
		CustomCallerFormatter: nil,
	})
	if !utils.CheckQuestionDB() {
		go utils.DownloadDbFile()
		//log.Errorln("题库文件不存在或已损坏，请手动前往 https://github.com/johlanse/study_xxqg/blob/main/conf/QuestionBank.db 下载并放入程序根目录")
	}
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
		_ = os.Mkdir("./config/", 0666)
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
		if config.QQ.Enable {
			h.Handle("/qq", push.InitQQ())
			log.Infoln(fmt.Sprintf("已开启qq配置,监听地址： ==》 %v:%v", config.Web.Host, config.Web.Port))
		}
		if config.Web.Enable {
			log.Infoln(fmt.Sprintf("已开启web配置，web监听地址 ==> %v:%v", config.Web.Host, config.Web.Port))
			h.Handle("/", engine)
		}
		if config.Wechat.Enable {
			log.Infoln(fmt.Sprintf("已开启wechat公众号配置,监听地址： ==》 %v:%v", config.Web.Host, config.Web.Port))
			h.HandleFunc("/wx", push.HandleWechat)
		}
		if config.Web.Enable || config.Wechat.Enable || config.QQ.Enable {
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

			log.Infoln("已采用定时执行模式")
			c := cron.New()

			_, err := c.AddFunc(config.Cron, func() {
				log.Infoln("即将开始执行定时任务")
				// 检测是否开启了随机等待
				if config.CronRandomWait > 0 {
					rand.Seed(time.Now().UnixNano())
					r := rand.Intn(config.CronRandomWait)
					log.Infoln(fmt.Sprintf("随机延迟%d分钟", r))
					time.Sleep(time.Duration(r) * time.Minute)
				}
				do()
			})
			if err != nil {
				log.Errorln(err.Error())
				return
			}
			c.Start()
		}()
	}

	if config.TG.Enable {
		push.TgInit()
	}
	getPush := push.GetPush(config)
	getPush("", "flush", "学习强国助手已上线")

	if config.CustomCron != "" {
		c2 := cron.New()
		_, err := c2.AddFunc(config.CustomCron, func() {
			getPush("all", "flush", config.CustomMessage)
		})
		if err != nil {
			log.Errorln("添加自定义定时消息推送错误" + err.Error())
			return
		}
		c2.Run()
	}
	initTask()
	model.SetPush(getPush)
	if now {
		do()
	}
	stop.Stop()

}

func do() {

	log.Infoln("检测到模式", config.Model)

	getPush := push.GetPush(config)
	getPush("", "flush", "学习强国助手已上线")
	users, _ := model.Query()

	failUser, _ := model.QueryFailUser()
	for _, user := range failUser {
		go func(user2 *model.User) {
			c := &lib.Core{Push: getPush, ShowBrowser: config.ShowBrowser}
			getPush(user2.PushId, "flush", user2.Nick+"的cookie已过期")
			newUser, err := c.L(config.Retry.Times, user2.PushId)
			if err != nil {
				c.Push(user2.PushId, "flush", "用户"+user2.Nick+"登录超时！")
				return
			}
			c.Init()
			defer c.Quit()
			lib.Study(c, newUser)
		}(user)
	}

	s := &sync.WaitGroup{}
	// 如果为定时模式则直接循环所以用户依次运行

	if config.PoolSize == 1 {
		for _, user := range users {
			if state.IsStudy(user.Uid) {
				log.Infoln("检测到该用户已在学习中！")
				continue
			} else {
				core := &lib.Core{ShowBrowser: config.ShowBrowser, Push: getPush}
				core.Init()
				state.Add(user.Uid, core)
				lib.Study(core, user)
				core.Quit()
				state.Delete(user.Uid)
			}

		}
	} else {
		for _, u := range users {
			if state.IsStudy(u.Uid) {
				log.Infoln("检测到该用户已在学习中！")
				continue
			} else {
				core := &lib.Core{ShowBrowser: config.ShowBrowser, Push: getPush}
				core.Init()
				t := &Task{
					Core: core,
					User: u,
					wg:   s,
				}

				run(t)
				s.Add(1)
			}
		}
		s.Wait()
	}
	log.Infoln("定时任务执行完成")
	return

}

func runBack() {
	cmd, err := xdaemon.Background(os.Stdout, false)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if xdaemon.IsParent() {
		go onKill(cmd)
		for true {
			_ = cmd.Wait()
			if cmd.ProcessState.Exited() {
				log.Infoln(cmd.ProcessState)
				if cmd.ProcessState.ExitCode() != 201 {
					break
				} else {
					log.Infoln("检测到重启,开始重启程序")
				}
			}
			cmd, err = xdaemon.Background(os.Stdout, false)
			if err != nil {
				return
			}
		}
		os.Exit(0)
	}
}

func onKill(cmd *exec.Cmd) {
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	if cmd.Process != nil {
		cmd.Process.Kill()
	}
	os.Exit(1)
}
