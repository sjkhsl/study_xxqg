package lib

import (
	"encoding/json"
	"os"
	"time"

	"github.com/guonaihong/gout"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func init() {
	_, err := os.Stat(`./config/user.json`)
	if err != nil {
		err := os.WriteFile(user_path, []byte("[]"), 0666)
		if err != nil {
			return
		}
		return
	}
}

const (
	user_path = "./config/user.json"
)

type User struct {
	Cookies []Cookie `json:"cookies"`
	Nick    string   `json:"nick"`
	Uid     string   `json:"uid"`
	Time    int64    `json:"time"`
}

func GetUsers() ([]User, error) {
	file, err := os.ReadFile(user_path)
	if err != nil {
		return nil, err
	}
	var users []User
	err = json.Unmarshal(file, &users)
	if err != nil {
		return nil, err
	}
	var newUsers []User
	for i := 0; i < len(users); i++ {
		if CheckUserCookie(users[i]) {
			newUsers = append(newUsers, users[i])
			continue
		}
		log.Infoln("用户" + users[i].Nick + "cookie已失效")
	}
	return users, err
}

func SaveUser(user User) error {
	users, err := GetUsers()
	if err != nil {
		log.Errorln("获取用户信息错误")
		return err
	}
	a := false
	for _, u := range users {
		if u.Uid == user.Uid {
			u.Cookies = user.Cookies
			a = true
		}
	}
	if !a {
		users = append(users, user)
	}

	data, err := json.Marshal(&users)
	if err != nil {
		log.Errorln("序列化用户失败")
		return err
	}
	err = os.WriteFile(user_path, data, 0666)
	if err != nil {
		log.Errorln("写入用户信息到文件错误")

		return err
	}
	return err
}

func GetUserInfo(cookies []Cookie) (string, string, error) {
	var resp []byte
	err := gout.GET(user_Info_url).
		SetCookies(cookieToJar(cookies)...).
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

func CheckUserCookie(user User) bool {
	if time.Now().Unix() <= user.Time {
		return true
	}

	return false
}
