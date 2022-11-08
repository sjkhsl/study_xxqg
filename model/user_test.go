package model

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestInsert(t *testing.T) {
	err := AddUser(&User{
		Nick:      "123",
		Uid:       "123",
		Token:     "123444444444444444444444444",
		LoginTime: 1031312,
	})
	if err != nil {
		log.Errorln(err.Error())
	}
}

func TestQuery(t *testing.T) {
	users, err := Query()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	log.Infoln(users[0].Uid)
}

func TestFind(t *testing.T) {
	user := Find("123")
	log.Infoln(user.Nick)
}
