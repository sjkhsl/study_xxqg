package push

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/conf"
)

func GetPush(config conf.Config) func(id string, kind string, message string) {
	var pushs []func(id, kind, message string)
	if config.Push.Ding.Enable {
		ding := &Ding{
			Secret: config.Push.Ding.Secret,
			Token:  config.Push.Ding.AccessToken,
		}
		log.Infoln("已配置钉钉推送")
		pushs = append(pushs, ding.Send())
	}
	if config.Push.PushPlus.Enable {
		log.Infoln("已配置pushplus推送")
		pushs = append(pushs, (&PushPlus{Token: config.Push.PushPlus.Token}).Init())
	}
	if config.Wechat.Enable {
		log.Infoln("已配置wechat推送")
		pushs = append(pushs, func(id, kind, message string) {
			defer func() {
				err := recover()
				if err != nil {
					log.Errorln("推送微信消息出现错误")
					log.Errorln(err)
				}
			}()
			if kind == "flush" {
				sendMsg(id, message)
			} else {
				if log.GetLevel() == log.DebugLevel {
					sendMsg(id, message)
				}
			}
		})
	}
	if config.TG.Enable {
		log.Infoln("已配置tg推送")
		pushs = append(pushs, tgPush)
	}
	if config.PushDeer.Enable {
		log.Infoln("已配置pushDeer推送")
		pushs = append(pushs, InitPushDeer())
	}
	if config.QQ.Enable {
		log.Infoln("已配置qq推送")
		pushs = append(pushs, func(id, kind, message string) {
			e := &Event{qq: qq}
			if kind == "flush" {
				e.sendPrivateMsg(conf.GetConfig().QQ.SuperUser, message)
			} else {
				if log.GetLevel() == log.DebugLevel {
					e.sendPrivateMsg(conf.GetConfig().QQ.SuperUser, message)
				}
			}

		})
	}
	pushs = append(pushs, func(id, kind, message string) {
		log.Debugln(fmt.Sprintf("消息id: %v，消息类型：%v,消息内容：%v", id, kind, message))
	})
	return multiPush(pushs...)
}

func multiPush(pushs ...func(id, kind, message string)) func(id, kind, message string) {
	return func(id, kind, message string) {
		for _, push := range pushs {
			push(id, kind, message)
		}
	}
}
