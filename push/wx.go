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

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/lib"
	"github.com/johlanse/study_xxqg/lib/state"
	"github.com/johlanse/study_xxqg/model"
	"github.com/johlanse/study_xxqg/utils"
	"github.com/johlanse/study_xxqg/utils/update"
)

var (
	wx        *mp.WeiXin
	lastNonce = ""
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

	useRequest = "use_request"
)

type WechatHandler func(id string, msg string)

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
	RegisterHandler(useRequest, handleEventUseRequest)

	RegisterHandler("发送", handleTextSendMsg)

	RegisterHandler("申请使用", handleEventUseRequest)
	RegisterHandler("通过", handleTextPass)
	RegisterHandler("拒绝", handleTextReject)
	RegisterHandler("使用用户列表", handleTextUserList)

	// 发送”/remark test即可添加备注信息”
	RegisterHandler("/remark", handleTextRemark)

	wx = mp.New(config.Wechat.Token, config.Wechat.AppID, config.Wechat.Secret, "123", "123")
	err := wx.CreateMenu(mean)
	if err != nil {
		log.Errorln("设置自定义菜单出现异常" + err.Error())
		return
	}
	list, err := wx.GetKFList()
	if err != nil {
		log.Errorln("获取客服列表错误" + err.Error())
		return
	}
	if len(list.KF_list) < 1 {
		err := wx.AddKFAccount("xxqg@xxqg", "xxqg", utils.StrMd5("123"))
		if err != nil {
			log.Errorln("添加客服失败" + err.Error())
			return
		}
	} else {

	}
	wx.HandleFunc("eventCLICK", func(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
		if lastNonce == nonce {
			return
		}
		log.Infoln("收到微信点击事件：" + r.EventKey)
		lastNonce = nonce

		if !checkPermission(r.FromUserName, r.EventKey) {
			log.Infoln("未通过权限检测的用户事件！")
			return
		}

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
			(value.(WechatHandler))(r.FromUserName, "")
		}()
	})

	wx.HandleFunc(mp.GenHttpRouteKey(mp.MsgTypeEvent, mp.EventSubscribe), func(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
		sendMsg(r.FromUserName, "你已关注该公众号，请发送申请使用向管理员申请权限吧！")
	})

	wx.HandleFunc("text", func(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
		log.Infoln(fmt.Sprintf("收到了来自用户%v的文本消息：%v", r.FromUserName, r.Content))
		key := strings.Split(r.Content, " ")[0]

		if !checkPermission(r.FromUserName, key) {
			log.Infoln("未通过权限检测的用户事件！")
			return
		}

		value, ok := handlers.Load(key)
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
			(value.(WechatHandler))(r.FromUserName, r.Content)
		}()

	})
}

func checkPermission(id string, key string) bool {

	// 这三个不检查权限
	keys := []string{"use_request", "get_open_id", "/remark"}

	// 通过管理员所有权限
	if conf.GetConfig().Wechat.SuperOpenID == id {
		return true
	}
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	user, err := model.FindWechatUser(id)
	if err != nil {
		log.Errorln("获取用户出现错误" + err.Error())
		sendMsg(id, "请发送申请使用向管理员申请权限！")
		return false
	}
	if user == nil {
		log.Errorln("不存在该用户！")
		sendMsg(id, "请发送申请使用向管理员申请权限！")
		return false
	}
	if user.Status != 1 {
		return false
	} else {
		return true
	}

}

// handleTextUserList
/* @Description: 获取当前使用用户列表的处理器
*  @param id
*  @param msg
 */
func handleTextUserList(id, msg string) {

	if id != conf.GetConfig().Wechat.SuperOpenID {
		return
	}

	users, err := model.QueryWechatUser()
	if err != nil {
		log.Errorln("获取用户列表出现错误" + err.Error())
		return
	}
	message := ""
	for _, user := range users {
		message += fmt.Sprintf("open_id:%v\n\n备注：%v\n\n状态:%d", user.OpenId, user.Remark, user.Status)
	}
	sendMsg(id, message)
}

func handleTextReject(id, msg string) {

	if id != conf.GetConfig().Wechat.SuperOpenID {
		return
	}

	openID := strings.Split(msg, " ")[1]
	user, err := model.FindWechatUser(openID)
	if err != nil {
		log.Errorln("查询wechat用户出现错误" + err.Error())
		return
	}
	user.Status = -1
	err = model.UpdateWechatUser(user)
	if err != nil {
		log.Errorln("更新用户信息出现错误" + err.Error())
		return
	}
	sendMsg(user.OpenId, "管理员已拒绝了你的使用申请！")
	sendMsg(id, fmt.Sprintf("已拒绝用户(%v)%v使用", user.Remark, user.OpenId))
}

func handleTextPass(id, msg string) {

	if id != conf.GetConfig().Wechat.SuperOpenID {
		return
	}

	openID := strings.Split(msg, " ")[1]
	user, err := model.FindWechatUser(openID)
	if err != nil {
		log.Errorln("查询wechat用户出现错误" + err.Error())
		return
	}
	user.Status = 1
	err = model.UpdateWechatUser(user)
	if err != nil {
		log.Errorln("更新用户信息出现错误" + err.Error())
		return
	}
	sendMsg(user.OpenId, "管理员已通过了你的使用申请！")
	sendMsg(id, fmt.Sprintf("已允许用户(%v)%v使用", user.Remark, user.OpenId))
}

// handleEventUseRequest
/* @Description: 处理申请使用的点击事件
*  @param id
*  @param msg
 */
func handleEventUseRequest(id, msg string) {
	user, err := model.FindWechatUser(id)
	if user.OpenId == "" {

		err := model.AddWechatUser(&model.WechatUser{
			OpenId:          id,
			Remark:          "",
			Status:          0,
			LastRequestTime: time.Now().Unix(),
		})
		if err != nil {
			log.Errorln("添加用户出现错误" + err.Error())
			return
		}
		sendMsg(conf.GetConfig().Wechat.SuperOpenID, fmt.Sprintf("用户%v申请使用测试号，通过则回复信息：\n通过 %v\n\n拒绝则回复:\n拒绝 %v", id, id, id))

	} else {
		if err != nil {
			log.Errorln("查询wechat用户错误" + err.Error())
			return
		}
		if user.Status == 1 {
			sendMsg(id, "你已拥有使用权！")
			return
		} else if user.Status == -1 {
			sendMsg(id, "你已被拉黑，请联系管理员！")
			return
		}

		if (time.Now().Unix()-user.LastRequestTime)/3600 < 1 {
			sendMsg(id, fmt.Sprintf("你已在%v申请过使用权了，请一个小时后再申请！", time.Unix(user.LastRequestTime, 0).Format("2006-01-02 15:04:05")))
			return
		}

		if user.LastRequestTime == 0 {
			user.LastRequestTime = time.Now().Unix()
			err := model.UpdateWechatUser(user)
			if err != nil {
				log.Errorln("更新信息出现错误" + err.Error())
				return
			}
		}

		sendMsg(conf.GetConfig().Wechat.SuperOpenID, fmt.Sprintf("用户(%v)%v申请使用测试号，通过则回复信息：\n通过 %v\n\n拒绝则回复:\n拒绝 %v", user.Remark, id, id, id))
	}
}

// handleTextRemark
/* @Description: 添加备注信息
*  @param id
*  @param msg
 */
func handleTextRemark(id, msg string) {
	data := strings.Split(msg, " ")[1]
	count := model.WechatUserCount(id)
	if count < 1 {
		err := model.AddWechatUser(&model.WechatUser{
			OpenId:          id,
			Remark:          data,
			Status:          0,
			LastRequestTime: 0,
		})
		if err != nil {
			log.Errorln("remark出现错误" + err.Error())
			return
		}
	} else {
		user, err := model.FindWechatUser(id)
		if err != nil {
			log.Errorln("查找用户失败" + err.Error())
			return
		}
		user.Remark = data
		err = model.UpdateWechatUser(user)
		if err != nil {
			log.Errorln("remark出现错误" + err.Error())
			return
		}
	}
	sendMsg(id, "添加备注信息成功！")
}

// handleTextSendMsg
/* @Description: 自定义发送消息
*  @param id
*  @param content
 */
func handleTextSendMsg(id string, content string) {
	if conf.GetConfig().Wechat.SuperOpenID != id {
		return
	}
	msg := strings.SplitN(content, " ", 3)
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
}

func handleGetOpenID(id string, msg string) {
	sendMsg(id, "你的open_id为"+id)
}

//
//  handleCheckUpdate
//  @Description: 检查更新
//  @param id
//
func handleCheckUpdate(id string, msg string) {
	about := utils.GetAbout()
	sendMsg(id, about)
}

//
//  handleUpdate
//  @Description: 开始更新
//  @param id
//
func handleUpdate(id string, msg string) {
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
func handleRestart(id string, msg string) {
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

	if wx == nil {
		initWx()
	}

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

	if id == "" {
		id = conf.GetConfig().Wechat.SuperOpenID
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
	err := wx.PostText(id, message, "xxqg@xxqg")
	if err != nil {
		log.Errorln("发送客服消息错误" + err.Error())
		log.Warningln("开始尝试使用模板消息发送")
		m := map[string]interface{}{
			"data": map[string]string{
				"value": message,
			},
		}
		data, _ := json.Marshal(m)

		_, err = wx.SendTemplateMessage(&mp.TemplateMessage{
			ToUser:      id,
			TemplateId:  conf.GetConfig().Wechat.NormalTempID,
			URL:         "",
			TopColor:    "",
			RawJSONData: data,
		})
		if err != nil {
			return
		}
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
func handleLogin(id string, msg string) {
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
func handleStartStudy(id string, msg string) {
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
		if state.IsStudy(user.Uid) {
			log.Infoln("该用户已经在学习中了，跳过学习")
			continue
		} else {
			state.Add(user.Uid, core)
		}
		sendMsg(id, fmt.Sprintf("开始学习第%d个用户，用户名：%v", i+1, user.Nick))
		core.LearnArticle(user)
		core.LearnVideo(user)
		if conf.GetConfig().Model == 2 {
			core.RespondDaily(user, "daily")
		} else if conf.GetConfig().Model == 3 {
			core.RespondDaily(user, "weekly")
			core.RespondDaily(user, "special")
		}

		state.Delete(user.Uid)
		score, _ := lib.GetUserScore(user.ToCookies())
		sendMsg(id, fmt.Sprintf("第%d个用户%v学习完成，学习积分\n%v", i+1, user.Nick, lib.FormatScore(score)))
	}
}

func handleGetUser(id string, msg string) {
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
	if message == "" {
		log.Warningln("还未存在绑定的用户登录")
		sendMsg(id, "你还没有已登陆的用户，请点击下方登录按钮登录！")
		return
	}
	sendMsg(id, message)
}

func handleScore(id string, msg string) {
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

var (
	mean = &mp.Menu{Buttons: []mp.MenuButton{
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
				{
					Name: "申请使用",
					Type: "click",
					Key:  useRequest,
				},
			},
		},
	}}
)
