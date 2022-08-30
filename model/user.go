// Package model
// @Description:
package model

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/utils"
)

var (
	lock sync.RWMutex
)

func init() {
	go check()
}

var (
	pushFunc func(id, kind, message string)
)

func SetPush(push func(id, kind, message string)) {
	pushFunc = push
}

// User
/**
 * @Description:
 */
type User struct {
	Nick      string `json:"nick"`
	UID       string `json:"uid"`
	Token     string `json:"token"`
	LoginTime int64  `json:"login_time"`
	PushId    string `json:"push_id"`
	Status    int    `json:"status"`
}

// Query
/**
 * @Description: 查询所有未掉线的用户
 * @return []*User
 * @return error
 */
func Query() ([]*User, error) {
	var users []*User
	ping()
	lock.Lock()
	defer lock.Unlock()
	results, err := db.Query("select * from user")
	if err != nil {
		return nil, err
	}
	var failusers []*User
	for results.Next() {
		u := new(User)
		err := results.Scan(&u.Nick, &u.UID, &u.Token, &u.LoginTime, &u.PushId, &u.Status)
		if err != nil {
			_ = results.Close()
			return nil, err
		}
		if u.Status != 0 {

			if ok, _ := utils.CheckUserCookie(u.ToCookies()); ok {
				users = append(users, u)
			} else {
				log.Warningln(u.Nick + "的cookie已失效")
				failusers = append(failusers, u)
				if pushFunc != nil {
					pushFunc(u.PushId, "flush", u.Nick+"的cookie已失效")
				}
			}
		}
	}
	_ = results.Close()
	for _, failuser := range failusers {
		changeStatus(failuser.UID, 0)
	}
	return users, err
}

func changeStatus(uid string, status int) {
	ping()
	_, err := db.Exec("update user set status = ? where uid = ?", status, uid)
	if err != nil {
		log.Errorln("改变status失败" + err.Error())
		return
	}
}

func QueryFailUser() ([]*User, error) {
	var users []*User
	ping()
	lock.Lock()
	defer lock.Unlock()
	results, err := db.Query("select * from user where status = 0")
	if err != nil {
		return nil, err
	}
	for results.Next() {
		u := new(User)
		err := results.Scan(&u.Nick, &u.UID, &u.Token, &u.LoginTime, &u.PushId, &u.Status)
		if err != nil {
			_ = results.Close()
			return nil, err
		}
		users = append(users, u)
	}
	_ = results.Close()
	return users, err
}

// QueryByPushID
/**
 * @Description: 根据推送平台的key查询用户
 * @return []*User
 * @return error
 */
func QueryByPushID(pushID string) ([]*User, error) {
	lock.Lock()
	defer lock.Unlock()
	var users []*User
	ping()
	results, err := db.Query("select * from user where push_id = ?", pushID)
	if err != nil {
		return users, err
	}

	var failusers []*User
	for results.Next() {
		u := new(User)
		err := results.Scan(&u.Nick, &u.UID, &u.Token, &u.LoginTime, &u.PushId, &u.Status)
		if err != nil {
			_ = results.Close()
			return users, err
		}
		if u.Status != 0 {
			if ok, _ := utils.CheckUserCookie(u.ToCookies()); ok {
				users = append(users, u)
			} else {
				log.Warningln(u.Nick + "的cookie已失效")
				failusers = append(failusers, u)
			}
		}
	}
	_ = results.Close()
	for _, failuser := range failusers {
		changeStatus(failuser.UID, 0)
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
	err := db.QueryRow("select * from user where uid=?;", uid).Scan(&u.Nick, &u.UID, &u.Token, &u.LoginTime, &u.PushId, &u.Status)
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
	lock.Lock()

	ping()
	count := UserCount(user.UID)
	if count < 1 {
		_, err := db.Exec("insert into user (nick, uid, token, login_time,push_id) values (?,?,?,?,?)", user.Nick, user.UID, user.Token, user.LoginTime, user.PushId)
		if err != nil {
			log.Errorln("数据库插入失败")
			log.Errorln(err.Error())
			lock.Unlock()
			return err
		}
		lock.Unlock()
		return err
	}
	lock.Unlock()
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
	_, err := db.Exec("update user set token=?,login_time=?,push_id=?,status=1 where uid = ?", user.Token, user.LoginTime, user.PushId, user.UID)
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

// DeleteUser
/* @Description:
 * @param uid
 * @return error
 */
func DeleteUser(uid string) error {
	lock.Lock()
	defer lock.Unlock()
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

func TokenToCookies(token string) []*http.Cookie {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Domain:   "xuexi.cn",
		Expires:  time.Now().Add(time.Hour * 12),
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	return []*http.Cookie{cookie}
}

func (u *User) ToBrowserCookies() []playwright.BrowserContextAddCookiesOptionsCookies {
	cookie := playwright.BrowserContextAddCookiesOptionsCookies{
		Name:     playwright.String("token"),
		Value:    playwright.String(u.Token),
		Path:     playwright.String("/"),
		Domain:   playwright.String(".xuexi.cn"),
		Expires:  playwright.Float(float64(time.Now().Add(time.Hour * 12).Unix())),
		Secure:   playwright.Bool(false),
		HttpOnly: playwright.Bool(false),
		SameSite: playwright.SameSiteAttributeStrict,
	}
	return []playwright.BrowserContextAddCookiesOptionsCookies{cookie}
}

func check() {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorf("%v 出现错误，%v", "auth check", err)
		}
	}()
	c := cron.New()
	cr := "0 */1 * * *"
	if crEnv, ok := os.LookupEnv("CHECK_ENV"); ok {
		cr = crEnv
		log.Infoln("已成功自定义保活cron : " + cr)
	}
	_, err := c.AddFunc(cr, func() {
		log.Infoln("开始执行保活任务")
		users, _ := Query()
		for _, user := range users {
			response, _ := utils.GetClient().R().SetCookies(user.ToCookies()...).Get("https://pc-api.xuexi.cn/open/api/auth/check")
			token := ""
			for _, cookie := range response.Cookies() {
				if cookie.Name == "token" {
					token = cookie.Value
				}
			}
			if token != "" && user.Token != token {
				user.Token = token
				_ = UpdateUser(user)
				log.Infoln("用户" + user.Nick + "的ck已成功保活cookie")
			}
		}
	})
	if err != nil {
		log.Errorln("添加保活任务失败" + err.Error())
		return
	}
	c.Start()
}
