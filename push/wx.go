package push

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

	"github.com/sjkhsl/study_xxqg/conf"
	"github.com/sjkhsl/study_xxqg/lib"
	"github.com/sjkhsl/study_xxqg/model"
	"github.com/sjkhsl/study_xxqg/utils"
	"github.com/sjkhsl/study_xxqg/utils/update"
)

var (
	wx        *mp.WeiXin
	lastNonce = ""
	datas1    sync.Map
	wxPush    func(id, kind, message string)
)

const (
	loginBtn   = "loginBtn"
	StartStudy = "start_study"
	getUser    = "get_user"
	SCORE      = "score"

	checkUpdate = "check_update"
	updateBtn   = "updateBtn"
	restart     = "restart"
	getOpenID   = "get_open_id"
)

type WechatHandler func(id string)

var (
	handlers sync.Map
)

func RegisterHandler(key string, action WechatHandler) {
	handlers.Store(key, action)
}

func initWx() {
	once := sync.Once{}
	once.Do(initWechat)
}

func initWechat() {
	config := conf.GetConfig()
	if !config.Wechat.Enable {
		return
	}

	if config.Wechat.SuperOpenID == "" {
		log.Warningln("你还未配置super_open_id选项")
	}

	// 注册插件
	RegisterHandler(loginBtn, handleLogin)
	RegisterHandler(StartStudy, handleStartStudy)
	RegisterHandler(getUser, handleGetUser)
	RegisterHandler(SCORE, handleScore)
	RegisterHandler(checkUpdate, handleCheckUpdate)
	RegisterHandler(updateBtn, handleUpdate)
	RegisterHandler(restart, handleRestart)
	RegisterHandler(getOpenID, handleGetOpenID)

	wx = mp.New(config.Wechat.Token, config.Wechat.AppID, config.Wechat.Secret, "123", "123")
	err := wx.CreateMenu(&mp.Menu{Buttons: []mp.MenuButton{
		{
			Name:       "登录",
			Type:       "click",
			Key:        loginBtn,
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
					Name: "开始学习",
					Type: "click",
					Key:  StartStudy,
				},
				{
					Name: "获取用户",
					Type: "click",
					Key:  getUser,
				},
				{
					Name: "积分查询",
					Type: "click",
					Key:  SCORE,
				},
			},
		},
		{
			Name: "关于",
			Type: "click",
			SubButtons: []mp.MenuButton{
				{
					Name: "检查更新",
					Type: "click",
					Key:  checkUpdate,
				},
				{
					Name: "重启程序",
					Type: "click",
					Key:  restart,
				},
				{
					Name: "更新程序",
					Type: "click",
					Key:  updateBtn,
				},
				{
					Name: "获取open_id",
					Type: "click",
					Key:  getOpenID,
				},
			},
		},
	}})
	if err != nil {
		log.Errorln("设置自定义菜单出现异常" + err.Error())
		return
	}
	wx.HandleFunc("eventCLICK", func(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
		if lastNonce == nonce {
			return
		}
		lastNonce = nonce
		value, ok := handlers.Load(r.EventKey)
		if !ok {
			log.Warningln("未注册key为" + r.EventKey + "的调用方法")
			return
		}
		go func() {
			defer func() {
				err := recover()
				if err != nil {
					log.Errorln("处理微信事件错误")
					log.Errorln(err)
				}
			}()
			(value.(WechatHandler))(r.FromUserName)
		}()
	})
	wx.HandleFunc("text", func(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {

		if r.FromUserName != conf.GetConfig().Wechat.SuperOpenID {
			log.Infoln("收到了微信文本消息，但不是管理员")
			return
		}

		msg := strings.SplitN(r.Content, " ", 3)
		if len(msg) < 3 {
			return
		} else {
			if msg[0] == "发送" {
				if msg[1] == "all" || msg[1] == "所有" {
					sendMsg("all", msg[2])
				} else {
					sendMsg(msg[1], msg[2])
				}
			}
		}
	})
}

func handleGetOpenID(id string) {
	sendMsg(id, "你的open_id为"+id)
}

//
//  handleCheckUpdate
//  @Description: 检查更新
//  @param id
//
func handleCheckUpdate(id string) {
	about := utils.GetAbout()
	sendMsg(id, about)
}

//
//  handleUpdate
//  @Description: 开始更新
//  @param id
//
func handleUpdate(id string) {
	if conf.GetConfig().Wechat.SuperOpenID != id {
		sendMsg(id, "请联系管理员处理！")
		return
	}
	update.SelfUpdate("", conf.GetVersion())
	sendMsg(id, "检查更新已完成，即将重启程序")
	utils.Restart()
}

//
//  handleRestart
//  @Description: 重启程序
//  @param id
//
func handleRestart(id string) {
	if conf.GetConfig().Wechat.SuperOpenID != id {
		sendMsg(id, "请联系管理员处理！")
		return
	}
	sendMsg(id, "即将重启程序")
	utils.Restart()
}

//
//  sendMsg
//  @Description: 发送消息
//  @param id
//  @param message
//
func sendMsg(id, message string) {

	if id == "all" {
		userList, err := wx.GetUserList("")
		if err != nil {
			log.Errorln("获取关注列表错误")
			return
		}
		url := ""
		color := ""
		if strings.Contains(message, "$$$") {
			splits := strings.Split(message, "$$$")
			message = splits[0]
			url = splits[1]
			if len(splits) == 3 {
				color = splits[2]
			}
		}
		for _, user := range userList.Data.OpenId {
			m := map[string]interface{}{
				"data": map[string]string{
					"value": message,
				},
			}
			data, _ := json.Marshal(m)

			_, err = wx.SendTemplateMessage(&mp.TemplateMessage{
				ToUser:      user,
				TemplateId:  conf.GetConfig().Wechat.NormalTempID,
				URL:         url,
				TopColor:    color,
				RawJSONData: data,
			})
			if err != nil {
				log.Errorln("向用户" + user + "推送消息错误")
				continue
			}
		}
	}

	// 登录消息单独采用模板发送
	if strings.Contains(message, "login.xuexi.cn") {
		_, err := wx.SendTemplateMessage(&mp.TemplateMessage{
			ToUser:      id,
			TemplateId:  conf.GetConfig().Wechat.LoginTempID,
			URL:         message,
			TopColor:    "",
			RawJSONData: nil,
		})
		if err != nil {
			log.Errorln(err.Error())
			return
		}
		return
	}

	m := map[string]interface{}{
		"data": map[string]string{
			"value": message,
		},
	}
	data, _ := json.Marshal(m)
	if wx == nil {
		initWx()
	}
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
/* @Description:处理wechat的请求接口
 * @param rep
 * @param req
 */
func HandleWechat(rep http.ResponseWriter, req *http.Request) {
	if wx == nil {
		initWx()
	}
	wx.ServeHTTP(rep, req)
}

//
//  handleLogin
//  @Description: 用户登录
//  @param id
//
func handleLogin(id string) {
	core := &lib.Core{Push: func(id1 string, kind, message string) {
		if kind == "flush" && strings.Contains(message, "login.xuexi.cn") {
			_, err := wx.SendTemplateMessage(&mp.TemplateMessage{
				ToUser:      id,
				TemplateId:  conf.GetConfig().Wechat.LoginTempID,
				URL:         message,
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

//
//  handleStartStudy
//  @Description: 开始学习
//  @param id
//
func handleStartStudy(id string) {
	users, err := model.QueryByPushID(id)
	if err != nil {
		return
	}
	if users == nil {
		log.Warningln("还未存在绑定的用户登录")
		sendMsg(id, "你还没有已登陆的用户，请点击下方登录按钮登录！")
		return
	}
	core := &lib.Core{ShowBrowser: conf.GetConfig().ShowBrowser, Push: func(id1 string, kind, msg string) {
	}}
	core.Init()
	defer core.Quit()
	for i, user := range users {
		_, ok := datas1.Load(user.UID)
		if ok {
			log.Warningln("用户" + user.Nick + "已经在学习中了，跳过该用户")
			continue
		} else {
			datas1.Store(user.UID, "")
		}
		sendMsg(id, fmt.Sprintf("开始学习第%d个用户，用户名：%v", i+1, user.Nick))
		core.LearnArticle(user)
		core.LearnVideo(user)
		core.RespondDaily(user, "daily")
		core.RespondDaily(user, "weekly")
		core.RespondDaily(user, "special")
		datas1.Delete(user.UID)
		score, _ := lib.GetUserScore(user.ToCookies())
		sendMsg(id, fmt.Sprintf("第%d个用户%v学习完成，学习积分\n%v", i+1, user.Nick, lib.FormatScore(score)))
	}
}

func handleGetUser(id string) {
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
	config := conf.GetConfig()
	for _, user := range users {
		if config.Wechat.SuperOpenID == id {
			message += fmt.Sprintf("%v ==>  %v", user.Nick, time.Unix(user.LoginTime, 0).Format("2006-01-02"))
			if user.PushId == id {
				message += "(已绑定)\r\n"
			}
		} else {
			if user.PushId == id {
				message += fmt.Sprintf("%v ==>  %v", user.Nick, time.Unix(user.LoginTime, 0).Format("2006-01-02"))
			}
		}

	}
	sendMsg(id, message)
}

func handleScore(id string) {
	users, err := model.Query()
	if err != nil {
		return
	}
	if users == nil {
		log.Warningln("还未存在绑定的用户登录")
		sendMsg(id, "你还没有已登陆的用户，请点击下方登录按钮登录！")
		return
	}
	config := conf.GetConfig()
	for _, user := range users {
		score, _ := lib.GetUserScore(user.ToCookies())
		if config.Wechat.SuperOpenID == id {
			sendMsg(id, "用户："+user.Nick+"\n"+lib.FormatScore(score))
		} else {
			if user.PushId == id {
				sendMsg(id, "用户："+user.Nick+"\n"+lib.FormatScore(score))
			}
		}

	}
}
