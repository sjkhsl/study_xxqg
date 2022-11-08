package main

import (
	"sync"

	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/lib"
	"github.com/johlanse/study_xxqg/lib/state"
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

func initTask() {

	pool1, err := ants.NewPoolWithFunc(config.PoolSize, func(i2 interface{}) {
		task := i2.(*Task)
		log.Infoln("开始执行" + task.User.Nick)
		state.Add(task.User.Uid, task.Core)
		lib.Study(task.Core, task.User)
		defer task.Core.Quit()
		defer state.Delete(task.User.Uid)
		task.wg.Done()
	})
	if err != nil {
		log.Errorln("创建定时任务协程池失败" + err.Error())
	}
	pool = pool1
}
