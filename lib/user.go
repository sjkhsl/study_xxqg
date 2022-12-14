package lib

import (
	"net/http"

	"github.com/guonaihong/gout"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/sjkhsl/study_xxqg/model"
)

// GetUserInfo
/**
 * @Description: 获取用户信息
 * @param cookies
 * @return string
 * @return string
 * @return error
 */
func GetUserInfo(cookies []*http.Cookie) (string, string, error) {
	var resp []byte
	err := gout.GET(userInfoUrl).
		SetCookies(cookies...).
		SetHeader(gout.H{
			"Cache-Control": "no-cache",
		}).BindBody(&resp).Do()
	if err != nil {
		log.Errorln("获取用户信息失败")

		return "", "", err
	}
	log.Debugln("[user] 用户信息：", gjson.GetBytes(resp, "@this|@pretty").String())
	uid := gjson.GetBytes(resp, "data.uid").String()
	nick := gjson.GetBytes(resp, "data.nick").String()

	return uid, nick, err
}

// CheckUserCookie
/**
 * @Description: 获取用户成绩
 * @param user
 * @return bool
 */
func CheckUserCookie(user *model.User) bool {
	_, err := GetUserScore(user.ToCookies())
	if err != nil && err.Error() == "token check failed" {
		return false
	}
	return true
}
