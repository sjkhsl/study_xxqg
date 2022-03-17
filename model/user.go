// Package model
// @Description:
package model

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/guonaihong/gout"
	"github.com/mxschmitt/playwright-go"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// User
/**
 * @Description:
 */
type User struct {
	Nick      string `json:"nick"`
	UID       string `json:"uid"`
	Token     string `json:"token"`
	LoginTime int64  `json:"login_time"`
}

// Query
/**
 * @Description:
 * @return []*User
 * @return error
 */
func Query() ([]*User, error) {
	var users []*User
	ping()
	results, err := db.Query("select * from user")
	if err != nil {
		return nil, err
	}
	defer func(results *sql.Rows) {
		err := results.Close()
		if err != nil {
			log.Errorln("关闭results失败" + err.Error())
		}
	}(results)
	for results.Next() {
		u := new(User)
		err := results.Scan(&u.Nick, &u.UID, &u.Token, &u.LoginTime)
		if err != nil {
			return nil, err
		}
		if CheckUserCookie(u) {
			users = append(users, u)
		} else {
			log.Infoln("用户" + u.Nick + "cookie已失效")
		}
	}
	return users, err
}

// Find
/**
 * @Description:
 * @param uid
 * @return *User
 */
func Find(uid string) *User {
	u := new(User)
	err := db.QueryRow("select * from user where uid=?;", uid).Scan(&u.Nick, &u.UID, &u.Token, &u.LoginTime)
	if err != nil {
		return nil
	}
	return u
}

// AddUser
/**
 * @Description:
 * @param user
 * @return error
 */
func AddUser(user *User) error {
	ping()
	count := UserCount(user.UID)
	if count < 1 {
		_, err := db.Exec("insert into user (nick, uid, token, login_time) values (?,?,?,?)", user.Nick, user.UID, user.Token, user.LoginTime)
		if err != nil {
			log.Errorln("数据库插入失败")
			log.Errorln(err.Error())
			return err
		}
		return err
	}
	err := UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser
/**
 * @Description:
 * @param user
 * @return error
 */
func UpdateUser(user *User) error {
	ping()
	_, err := db.Exec("update user set token=? where uid = ?", user.Token, user.UID)
	if err != nil {
		log.Errorln("更新数据失败")
		log.Errorln(err.Error())
		return err
	}
	return err
}

// UserCount
/**
 * @Description:
 * @param uid
 * @return int
 */
func UserCount(uid string) int {
	ping()
	var count int
	err := db.QueryRow("select count(*) from user where uid = ?", uid).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func DeleteUser(uid string) error {
	ping()
	_, err := db.Exec("delete from user where uid = ?;", uid)
	if err != nil {
		return err
	}
	return err
}

// ToCookies
/**
 * @Description: 获取user的cookie
 * @receiver u
 * @return []*http.Cookie
 */
func (u *User) ToCookies() []*http.Cookie {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    u.Token,
		Path:     "/",
		Domain:   "xuexi.cn",
		Expires:  time.Now().Add(time.Hour * 12),
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	return []*http.Cookie{cookie}
}

func (u *User) ToBrowserCookies() []playwright.SetNetworkCookieParam {
	cookie := playwright.SetNetworkCookieParam{
		Name:     "token",
		Value:    u.Token,
		Path:     playwright.String("/"),
		Domain:   playwright.String(".xuexi.cn"),
		Expires:  playwright.Int(int(time.Now().Add(time.Hour * 12).UnixNano())),
		Secure:   playwright.Bool(false),
		HttpOnly: playwright.Bool(false),
		SameSite: playwright.String("Strict"),
	}
	return []playwright.SetNetworkCookieParam{cookie}
}

// CheckUserCookie
/**
 * @Description: 获取用户成绩
 * @param user
 * @return bool
 */
func CheckUserCookie(user *User) bool {
	var resp []byte
	err := gout.GET("https://pc-api.xuexi.cn/open/api/score/get").SetCookies(user.ToCookies()...).SetHeader(gout.H{
		"Cache-Control": "no-cache",
	}).BindBody(&resp).Do()
	if err != nil {
		log.Errorln("获取用户总分错误" + err.Error())

		return false
	}
	if !gjson.GetBytes(resp, "ok").Bool() {
		return false
	}
	return true
}
