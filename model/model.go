package model

import (
	"database/sql"

	_ "modernc.org/sqlite"

	log "github.com/sirupsen/logrus"
)

var (
	db *sql.DB
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
    login_time integer not null
);
`)
}

func ping() {
	err := db.Ping()
	if err != nil {
		log.Errorln("数据库断开了连接")
	}
}
