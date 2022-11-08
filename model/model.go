package model

import (
	"database/sql"
	"os"
	"sync"

	_ "modernc.org/sqlite"
	"xorm.io/xorm"

	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/utils"
)

var (
	engine *xorm.Engine

	db1 *sql.DB
)

func init() {
	en, err := xorm.NewEngine("sqlite", "./config/user.db")
	if err != nil {
		log.Errorln("打开数据库失败！" + err.Error())
		os.Exit(3)
	}
	err = en.Sync2(new(User), new(WechatUser))
	if err != nil {
		log.Errorln("同步数据库结构失败" + err.Error())
		return
	}
	en.SetLogger(&utils.MyLog{})
	en.ShowSQL(true)

	engine = en
}

func initQuestionDb() {
	once := sync.Once{}
	once.Do(func() {
		var err error
		db1, err = sql.Open("sqlite", "./QuestionBank.db")
		if err != nil {
			log.Errorln("题目数据库打开失败，请检查QuestionDB是否存在")
			log.Panicln(err.Error())
		}
		db1.SetMaxOpenConns(1)
	})
}

// User
/**
 * @Description:
 */
type User struct {
	Nick      string `xorm:"TEXT" json:"nick,omitempty"`
	Uid       string `xorm:"TEXT" json:"uid,omitempty"`
	Token     string `xorm:"TEXT" json:"token,omitempty"`
	LoginTime int64  `xorm:"integer" json:"loginTime,omitempty"`
	PushId    string `xorm:"TEXT" json:"pushId,omitempty"`
	Status    int    `xorm:"integer" json:"status,omitempty"`
}

type WechatUser struct {
	OpenId          string `xorm:"TEXT" json:"openId,omitempty"`
	Remark          string `xorm:"TEXT" json:"remark,omitempty"`
	Status          int    `xorm:"INTEGER" json:"status,omitempty"`
	LastRequestTime int64  `xorm:"INTEGER" json:"lastRequestTime,omitempty"`
}
