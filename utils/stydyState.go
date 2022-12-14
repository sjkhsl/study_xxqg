package utils

import (
	"errors"
	"sync"

	"github.com/sjkhsl/study_xxqg/lib"
)

// 该文件的方法为保存当前正在学习的用户

var (
	state sync.Map
)

func Add(uid string, core *lib.Core) error {
	_, ok := state.Load(uid)
	if ok {
		return errors.New("the user is studying")
	} else {
		state.Store(uid, core)
		return nil
	}
}

func Delete(uid string) error {
	state.Delete(uid)
	return nil
}

func Item(item func(uid string, core *lib.Core) bool) {
	state.Range(func(key, value interface{}) bool {
		return item(key.(string), value.(*lib.Core))
	})
}
