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
		return ding.Send()
	} else if config.Push.TG.Enable {
		t := &Telegram{
			Token:  config.Push.TG.Token,
			ChatId: config.Push.TG.ChatID,
		}
		return t.Init()
	}
	return func(kind string, message string) {
		log.Infoln("")
	}
}
