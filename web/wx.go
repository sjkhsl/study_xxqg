package web

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/yujinliang/wechat/mp"
	"github.com/yujinliang/wechat/mp/request"

	"github.com/huoxue1/study_xxqg/conf"
)

func init() {
	//InitWechat()
}

var (
	wx *mp.WeiXin
)

func InitWechat() {
	config := conf.GetConfig()
	log.Infoln(config.Wechat)
	wx = mp.New(config.Wechat.Token, config.Wechat.AppID, config.Wechat.Secret, "123", "123")
	wx.CreateMenu(&mp.Menu{Buttons: []mp.MenuButton{
		{
			Name:       "登录",
			Type:       "click",
			Key:        "login",
			Url:        "",
			MediaId:    "",
			SubButtons: nil,
		},
	}})
	wx.HandleFunc("eventCLICK", func(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
		log.Infoln(r.EventKey)
	})
}

func HandleWechat(rep http.ResponseWriter, req *http.Request) {
	wx.ServeHTTP(rep, req)
}
