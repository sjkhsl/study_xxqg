package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/johlanse/wechat/mp"
	"github.com/johlanse/wechat/mp/request"
	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/lib"
	"github.com/johlanse/study_xxqg/model"
	"github.com/johlanse/study_xxqg/utils"
	"github.com/johlanse/study_xxqg/utils/update"
)

func init() {
	initWechat()
}

var (
	wx        *mp.WeiXin
	lastNonce = ""
	datas     sync.Map
)

const (
	login      = "login"
	StartStudy = "start_study"
	getUser    = "stop_study"
	SCORE      = "score"

	checkUpdate = "check_update"
	updateBtn   = "updateBtn"
	restart     = "restart"
)

func initWechat() {
	config := conf.GetConfig()
	if !config.Wechat.Enable {
		return
	}
	log.Infoln(config.Wechat)
	wx = mp.New(config.Wechat.Token, config.Wechat.AppID, config.Wechat.Secret, "123", "123")
	err := wx.CreateMenu(&mp.Menu{Buttons: []mp.MenuButton{
		{
			Name:       "登录",
			Type:       "click",
			Key:        login,
			Url:        "",
			MediaId:    "",
			SubButtons: nil,
		},
		{
			Name:    "学习管理",
			Type:    "click",
			Key:     "study",
			Url:     "",
			MediaId: "",
			SubButtons: []mp.MenuButton{
				{
					Name:       "开始学习",
					Type:       "click",
					Key:        StartStudy,
					Url:        "",
					MediaId:    "",
					SubButtons: nil,
				},
				{
					Name:       "获取用户",
					Type:       "click",
					Key:        getUser,
					Url:        "",
					MediaId:    "",
					SubButtons: nil,
				},
				{
					Name:       "积分查询",
					Type:       "click",
					Key:        SCORE,
					Url:        "",
					MediaId:    "",
					SubButtons: nil,
				},
			},
		},
		{
			Name:    "关于",
			Type:    "click",
			Key:     "",
			Url:     "",
			MediaId: "",
			SubButtons: []mp.MenuButton{
				{
					Name:       "检查更新",
					Type:       "click",
					Key:        checkUpdate,
					Url:        "",
					MediaId:    "",
					SubButtons: nil,
				},
				{
					Name:       "重启程序",
					Type:       "click",
					Key:        restart,
					Url:        "",
					MediaId:    "",
					SubButtons: nil,
				},
				{
					Name:       "更新程序",
					Type:       "click",
					Key:        updateBtn,
					Url:        "",
					MediaId:    "",
					SubButtons: nil,
				},
			},
		},
	}})
	if err != nil {
		log.Errorln("设置自定义菜单出现异常" + err.Error())
		return
	}
	if conf.GetConfig().Wechat.PushLoginWarn {
		model.SetPush(sendMsg)
	}
	wx.HandleFunc("eventCLICK", func(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
		if lastNonce == nonce {
			return
		}
		lastNonce = nonce
		switch r.EventKey {
		case login:
			go handleLogin(r.FromUserName)
		case StartStudy:
			go handleStartStudy(r.FromUserName)
		case getUser:
			go handleGetUser(r.FromUserName)
		case SCORE:
			go handleScore(r.FromUserName)
		case checkUpdate:
			handleCheckUpdate(r.FromUserName)
		case updateBtn:
			handleUpdate(r.FromUserName)
		case restart:
			handleRestart(r.FromUserName)
		}
	})
}

func handleCheckUpdate(id string) {
	about := utils.GetAbout()
	sendMsg(id, about)
}

func handleUpdate(id string) {
	update.SelfUpdate("", conf.GetVersion())
	sendMsg(id, "检查更新已完成，即将重启程序")
	utils.Restart()
}

func handleRestart(id string) {
	sendMsg(id, "即将重启程序")
	utils.Restart()
}

func sendMsg(id, message string) {
	m := map[string]interface{}{
		"data": map[string]string{
			"value": message,
		},
	}
	data, _ := json.Marshal(m)
	_, err := wx.SendTemplateMessage(&mp.TemplateMessage{
		ToUser:      id,
		TemplateId:  conf.GetConfig().Wechat.NormalTempID,
		URL:         "",
		TopColor:    "",
		RawJSONData: data,
	})
	if err != nil {
		return
	}
}

// HandleWechat
/* @Description:
 * @param rep
 * @param req
 */
func HandleWechat(rep http.ResponseWriter, req *http.Request) {
	wx.ServeHTTP(rep, req)
}

func handleLogin(id string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("处理微信事件错误")
			log.Errorln(err)
		}
	}()
	core := &lib.Core{Push: func(kind, message string) {
		if kind == "flush" && strings.HasPrefix(message, "登录链接") {
			l := strings.ReplaceAll(message, "登录链接：\r\n", "")
			_, err := wx.SendTemplateMessage(&mp.TemplateMessage{
				ToUser:      id,
				TemplateId:  conf.GetConfig().Wechat.LoginTempID,
				URL:         l,
				TopColor:    "",
				RawJSONData: nil,
			})
			if err != nil {
				log.Errorln(err.Error())
				return
			}
		}
	}}
	_, err := core.L(0, id)
	if err != nil {
		return
	}
	sendMsg(id, "登录成功")
}

func handleStartStudy(id string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("处理微信事件错误")
			log.Errorln(err)
		}
	}()
	users, err := model.QueryByPushID(id)
	if err != nil {
		return
	}
	if users == nil {
		log.Warningln("还未存在绑定的用户登录")
		sendMsg(id, "你还没有已登陆的用户，请点击下方登录按钮登录！")
		return
	}
	core := &lib.Core{ShowBrowser: conf.GetConfig().ShowBrowser, Push: func(kind, msg string) {
	}}
	core.Init()
	defer core.Quit()
	for i, user := range users {
		_, ok := datas.Load(user.UID)
		if ok {
			log.Warningln("用户" + user.Nick + "已经在学习中了，跳过该用户")
			continue
		} else {
			datas.Store(user.UID, "")
		}
		sendMsg(id, fmt.Sprintf("开始学习第%d个用户，用户名：%v", i+1, user.Nick))
		core.LearnArticle(user)
		core.LearnVideo(user)
		core.RespondDaily(user, "daily")
		core.RespondDaily(user, "weekly")
		core.RespondDaily(user, "special")
		datas.Delete(user.UID)
		score, _ := lib.GetUserScore(user.ToCookies())
		sendMsg(id, fmt.Sprintf("第%d个用户%v学习完成，学习积分\n%v", i+1, user.Nick, lib.FormatScore(score)))
	}
}

func handleGetUser(id string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("处理微信事件错误")
			log.Errorln(err)
		}
	}()
	users, err := model.Query()
	if err != nil {
		return
	}
	if users == nil {
		log.Warningln("还未存在绑定的用户登录")
		sendMsg(id, "你还没有已登陆的用户，请点击下方登录按钮登录！")
		return
	}
	message := ""
	for _, user := range users {
		message += fmt.Sprintf("%v ==>  %v", user.Nick, time.Unix(user.LoginTime, 0).Format("2006-01-02 15:04:05"))
		if user.PushId == id {
			message += "(已绑定)\n"
		}
	}
	sendMsg(id, message)
}

func handleScore(id string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("处理微信事件错误")
			log.Errorln(err)
		}
	}()

	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("处理微信事件错误")
			log.Errorln(err)
		}
	}()
	users, err := model.Query()
	if err != nil {
		return
	}
	if users == nil {
		log.Warningln("还未存在绑定的用户登录")
		sendMsg(id, "你还没有已登陆的用户，请点击下方登录按钮登录！")
		return
	}
	for _, user := range users {
		score, _ := lib.GetUserScore(user.ToCookies())
		sendMsg(id, "用户："+user.Nick+"\n"+lib.FormatScore(score))
	}
}
