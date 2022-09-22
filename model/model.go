package model

import (
	"database/sql"
	"sync"

	_ "modernc.org/sqlite"

	log "github.com/sirupsen/logrus"
)

var (
	db *sql.DB

	db1 *sql.DB
)

func init() {
	var err error
	db, err = sql.Open("sqlite", "./config/user.db")
	if err != nil {
		log.Errorln("用户数据库打开失败，请检查config目录权限")
		log.Panicln(err.Error())
	}

	_, _ = db.Exec(`create table user
(
    nick       TEXT,
    uid        TEXT    not null
        constraint user_pk
            primary key,
    token      TEXT    not null,
    login_time integer not null,
    push_id TEXT
);
`)

	_, _ = db.Exec(`
		create table wechat_user(
		    open_id TEXT not null constraint user_pk primary key,
		    remark TEXT default '',
		    status INTEGER default 0,
		    last_request_time INTEGER not null 
		)
	`)
	_, _ = db.Exec(`alter table user
    add status integer default 1;
`)
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
	})
}

func ping() {
	err := db.Ping()
	if err != nil {
		log.Errorln("数据库断开了连接")
	}
}
