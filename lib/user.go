package lib

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/sjkhsl/study_xxqg/model"
	"github.com/sjkhsl/study_xxqg/utils"
)

// 获取用户信息
func GetUserInfo(cookies []*http.Cookie) (string, string, error) {
	var resp []byte
	response, err := utils.GetClient().R().SetCookies(cookies...).SetHeader("Cache-Control", "no-cache").Get(userInfoUrl)
	if err != nil {
		log.Errorln("获取用户信息失败" + err.Error())
		return "", "", err
	}
	resp = response.Bytes()
	log.Debugln("[user] 用户信息：", gjson.GetBytes(resp, "@this|@pretty").String())
	uid := gjson.GetBytes(resp, "data.uid").String()
	nick := gjson.GetBytes(resp, "data.nick").String()

	return uid, nick, err
}

// 获取用户成绩
func CheckUserCookie(user *model.User) bool {
	_, err := GetUserScore(user.ToCookies())
	if err != nil && err.Error() == "token check failed" {
		return false
	}
	return true
}
