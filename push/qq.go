package push

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/lib"
	"github.com/johlanse/study_xxqg/lib/state"
	"github.com/johlanse/study_xxqg/model"
	"github.com/johlanse/study_xxqg/utils"
)

var (
	qq *QQ
)

type qqPlugin func(event *Event, args []string)

type QQ struct {
	postAdd string
	plugins map[string]qqPlugin
}

func InitQQ() *QQ {
	config := conf.GetConfig()
	q := new(QQ)
	qq = q
	q.postAdd = config.QQ.PostAddr
	q.plugins = make(map[string]qqPlugin, 1)
	q.newPlugin("user", qqGetUser)
	q.newPlugin("score", qqGetScore)
	q.newPlugin("fail", qqGetFailUser)
	q.newPlugin("study", qqStudy)
	q.newPlugin("login", qqLogin)
	q.newPlugin("help", qqHelp)
	return q
}

func (q *QQ) newPlugin(text string, plugin qqPlugin) {
	q.plugins[text] = plugin
}

func (q *QQ) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	data, _ := io.ReadAll(request.Body)
	go q.handle(data)
	writer.WriteHeader(204)
}

func (q *QQ) handle(data []byte) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("qq消息处理错误")
		}
	}()
	config := conf.GetConfig()
	e := new(Event)
	err := json.Unmarshal(data, e)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	e.qq = q
	if e.PostType == "message" {
		log.Infoln("收到qq消息 ==》 " + e.Message)
		// 遍历白名单列表
		for _, id := range config.QQ.WhiteList {
			if e.GroupId == id || e.UserId == id {
				for text, plugin := range q.plugins {
					messages := strings.Split(e.Message, " ")
					if messages[0] == "."+text {
						plugin(e, messages[1:])
						return
					}
				}
			}
		}
		log.Infoln("消息来源不在白名单中！")

	}
}

func qqHelp(event *Event, _ []string) {
	help := ".user  查询用户\n.fail  查询过期用户\n.study  对一个用户进行学习\n.score  查询用户分数\n.login 登录一个用户"
	event.Send(help)
}

func qqLogin(event *Event, _ []string) {
	core := &lib.Core{ShowBrowser: conf.GetConfig().ShowBrowser, Push: func(id string, kind string, message string) {
		if kind == "flush" {
			event.Send(message)
		} else {
			if conf.GetConfig().LogLevel == "debug" {
				event.Send(message)
			}
		}
	}}
	_, err := core.L(conf.GetConfig().Retry.Times, "")
	if err != nil {
		event.Send(err.Error())
	}
}

func qqStudy(event *Event, args []string) {
	users, err := model.Query()
	if err != nil {
		event.Send(err.Error())
		return
	}
	var user *model.User
	if len(users) == 1 {
		user = users[0]
	} else {
		if len(args) < 0 {
			event.Send("缺少序号参数，请输入 .study 序号")
			return
		} else {
			index, err := strconv.Atoi(args[0])
			if err != nil {
				event.Send(err.Error())
				return
			}
			user = users[index]
		}
	}
	core := &lib.Core{ShowBrowser: conf.GetConfig().ShowBrowser, Push: func(id string, kind string, message string) {
		if kind == "flush" {
			event.Send(message)
		} else {
			if conf.GetConfig().LogLevel == "debug" {
				event.Send(message)
			}
		}
	}}
	core.Init()
	state.Add(user.Uid, core)
	defer state.Delete(user.Uid)
	defer core.Quit()
	lib.Study(core, user)

}

func qqGetScore(event *Event, _ []string) {
	users, err := model.Query()
	if err != nil {
		event.Send(err.Error())
		return
	}
	for _, user := range users {
		score, err := lib.GetUserScore(user.ToCookies())
		if err != nil {
			event.Send(err.Error())
			continue
		}
		event.Send(user.Nick + "\n" + lib.FormatScore(score))
	}
}

func qqGetFailUser(event *Event, _ []string) {
	user, err := model.QueryFailUser()
	if err != nil {
		event.Send(err.Error())
		return
	}
	result := ""
	for i, user := range user {
		result += fmt.Sprintf("%d  %v  %v\n", i, user.Nick, utils.Stamp2Str(user.LoginTime))
	}
	event.Send(result)
}

func qqGetUser(event *Event, _ []string) {
	users, err := model.Query()
	if err != nil {
		event.Send(err.Error())
		return
	}
	result := ""
	for i, user := range users {
		result += fmt.Sprintf("%d  %v  %v\n", i, user.Nick, utils.Stamp2Str(user.LoginTime))
	}
	event.Send(result)
}

type (
	anonymous struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Flag string `json:"flag"`
	}

	Files struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Size    int64  `json:"size"`
		Busid   int64  `json:"busid"`
		FileUrl string `json:"url"`
	}

	Status struct {
		AppEnabled     bool        `json:"app_enabled"`
		AppGood        bool        `json:"app_good"`
		AppInitialized bool        `json:"app_initialized"`
		Good           bool        `json:"good"`
		Online         bool        `json:"online"`
		PluginsGood    interface{} `json:"plugins_good"`
		Stat           struct {
			PacketReceived  int `json:"packet_received"`
			PacketSent      int `json:"packet_sent"`
			PacketLost      int `json:"packet_lost"`
			MessageReceived int `json:"message_received"`
			MessageSent     int `json:"message_sent"`
			DisconnectTimes int `json:"disconnect_times"`
			LostTimes       int `json:"lost_times"`
			LastMessageTime int `json:"last_message_time"`
		} `json:"stat"`
	}

	MessageIds struct {
		MessageID int32 `json:"message_id"`
	}

	Senders struct {
		Age      int    `json:"age"`
		Area     string `json:"area"`
		Card     string `json:"card"`
		Level    string `json:"level"`
		NickName string `json:"nickname"`
		Role     string `json:"role"`
		Sex      string `json:"sex"`
		Title    string `json:"title"`
		UserId   int    `json:"user_id"`
	}

	// Event
	/*
	 * 事件
	 *
	 */
	Event struct {
		qq            *QQ
		Anonymous     anonymous `json:"anonymous"`
		Font          int       `json:"font"`
		GroupId       int64     `json:"group_id"`
		Message       string    `json:"message"`
		MessageType   string    `json:"message_type"`
		PostType      string    `json:"post_type"`
		RawMessage    string    `json:"raw_message"`
		SelfId        int64     `json:"self_id"`
		Sender        Senders   `json:"sender"`
		SubType       string    `json:"sub_type"`
		UserId        int64     `json:"user_id"`
		Time          int       `json:"time"`
		NoticeType    string    `json:"notice_type"`
		RequestType   string    `json:"request_type"`
		Comment       string    `json:"comment"`
		Flag          string    `json:"flag"`
		OperatorID    int       `json:"operator_id"`
		File          Files     `json:"file"`
		Duration      int64     `json:"duration"`
		TargetId      int64     `json:"target_id"` // 运气王id
		HonorType     string    `json:"honor_type"`
		MetaEventType string    `json:"meta_event_type"`
		Status        Status    `json:"status"`
		Interval      int       `json:"interval"`
		CardNew       string    `json:"card_new"` // 新名片
		CardOld       string    `json:"card_old"` // 旧名片
		MessageIds

		GuildID   int64 `json:"guild_id"`
		ChannelID int64 `json:"channel_id"`
	}
)

func (e *Event) sendGroupMsg(groupId int64, message any) int {
	if _, ok := message.(string); ok {
		message = map[string]any{
			"type": "text",
			"data": map[string]any{
				"text": message,
			},
		}
	}
	response, err := utils.GetClient().R().SetHeader("Authorization", conf.GetConfig().QQ.AccessToken).SetBodyJsonMarshal(map[string]any{
		"action": "send_group_msg",
		"params": map[string]any{
			"group_id": groupId,
			"message":  message,
		},
	}).Post(e.qq.postAdd)
	if err != nil {
		return 0
	}
	return int(gjson.GetBytes(response.Bytes(), "data").Int())
}

func (e *Event) Send(message any) {
	if e.MessageType == "group" {
		e.sendGroupMsg(e.GroupId, message)
	} else {
		e.sendPrivateMsg(e.UserId, message)
	}
}

func (e *Event) sendPrivateMsg(userId int64, message any) int {
	if _, ok := message.(string); ok {
		message = map[string]any{
			"type": "text",
			"data": map[string]any{
				"text": message,
			},
		}
	}
	response, err := utils.GetClient().R().SetHeader("Authorization", conf.GetConfig().QQ.AccessToken).SetBodyJsonMarshal(map[string]any{
		"action": "send_private_msg",
		"params": map[string]any{
			"user_id": userId,
			"message": message,
		},
	}).Post(e.qq.postAdd)
	if err != nil {
		return 0
	}
	return int(gjson.GetBytes(response.Bytes(), "data").Int())
}
