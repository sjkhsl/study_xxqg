package push

import (
	"fmt"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/utils"
)

func InitPushDeer() func(id, kind, message string) {
	config := conf.GetConfig()

	return func(id, kind, message string) {
		if strings.Contains(message, "login.xuexi.cn") {
			message = fmt.Sprintf("[点击登录](%v)", message)
		}
		if kind == "image" {
			message = fmt.Sprintf("![](%v)", "data:image/png;base64,"+message)
		}
		values := url.Values{}
		values.Add("pushkey", config.PushDeer.Token)
		values.Add("text", strings.ReplaceAll(message, "</br>", "\n"))
		if kind == "flush" {
			_, _ = utils.GetClient().R().SetBody(values.Encode()).
				SetHeader("Content-type", "application/x-www-form-urlencoded").
				Post(config.PushDeer.Api + "/message/push")

		} else {
			if log.GetLevel() == log.DebugLevel {
				_, _ = utils.GetClient().R().SetBody(values.Encode()).
					SetHeader("Content-type", "application/x-www-form-urlencoded").
					Post(config.PushDeer.Api + "/message/push")

			}
		}

	}

}
