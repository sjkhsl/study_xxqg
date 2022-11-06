package state

import (
	"sync"

	"github.com/johlanse/study_xxqg/lib"
)

var (
	state = sync.Map{}
)

func Add(uid string, core *lib.Core) {
	state.Store(uid, core)
}

func IsStudy(uid string) bool {
	_, ok := state.Load(uid)
	return ok
}

func Delete(uid string) {
	state.Delete(uid)
}

func Get(uid string) *lib.Core {
	value, _ := state.Load(uid)
	return value.(*lib.Core)
}

func Range(fun func(key, value interface{}) bool) {
	state.Range(fun)
}
