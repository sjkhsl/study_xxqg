package push

import (
	log "github.com/sirupsen/logrus"

	"github.com/huoxue1/study_xxqg/lib"
)

func GetPush(config lib.Config) func(kind string, message string) {
	if config.Push.Ding.Enable {
		ding := &Ding{
			Secret: config.Push.Ding.Secret,
			Token:  config.Push.Ding.AccessToken,
		}
		log.Infoln("已配置钉钉推送")
		return ding.Send()
	} else if config.Push.PushPlus.Enable {
		log.Infoln("已配置pushplus推送")
		return (&PushPlus{Token: config.Push.PushPlus.Token}).Init()
	}
	return func(kind string, message string) {
		log.Infoln("")
	}
}
