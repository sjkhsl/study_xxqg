package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/lib"
	"github.com/johlanse/study_xxqg/model"
)

type Task struct {
	Core *lib.Core
	User *model.User
	wg   *sync.WaitGroup
}

var (
	pool *ants.PoolWithFunc
)

func run(task *Task) {
	pool.Invoke(task)
}

func inittask() {
	study := func(core2 *lib.Core, u *model.User) {
		defer func() {
			err := recover()
			if err != nil {
				log.Errorln("学习过程异常")
				log.Errorln(err)
			}
		}()
		startTime := time.Now()

		core2.LearnArticle(u)

		core2.LearnVideo(u)

		core2.LearnVideo(u)
		if config.Model == 2 {
			core2.RespondDaily(u, "daily")
		} else if config.Model == 3 {
			core2.RespondDaily(u, "daily")
			core2.RespondDaily(u, "weekly")
			core2.RespondDaily(u, "special")
		}
		endTime := time.Now()
		score, err := lib.GetUserScore(u.ToCookies())
		if err != nil {
			log.Errorln("获取成绩失败")
			log.Debugln(err.Error())
			return
		}

		score, _ = lib.GetUserScore(u.ToCookies())
		message := fmt.Sprintf("%v 学习完成,用时%.1f分钟\n%v", u.Nick, endTime.Sub(startTime).Minutes(), lib.FormatScoreShort(score))
		core2.Push(u.PushId, "flush", message)
	}

	pool1, err := ants.NewPoolWithFunc(config.PoolSize, func(i2 interface{}) {
		task := i2.(*Task)
		log.Infoln("开始执行" + task.User.Nick)
		study(task.Core, task.User)
		defer task.Core.Quit()
		task.wg.Done()
	})
	if err != nil {
		log.Errorln("创建定时任务协程池失败" + err.Error())
	}
	pool = pool1
}
