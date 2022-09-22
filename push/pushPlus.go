package push

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/utils"
)

type PushPlus struct {
	Token string
}

func (p *PushPlus) Init() func(id string, kind, message string) {
	send := func(data string) {
		_, err := utils.GetClient().R().SetBodyJsonMarshal(map[string]string{
			"token":    p.Token,
			"title":    "study_xxqg",
			"content":  data,
			"template": "markdown",
			"channel":  "wechat",
		}).Post("http://www.pushplus.plus/send")
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}

	return func(id string, kind, message string) {
		message = strings.ReplaceAll(message, "\n", "<br/>")
		switch {
		case kind == "image":
			message = fmt.Sprintf("![](%v)", "data:image/png;base64,"+message)
			send(message)
		case kind == "flush":
			if message != "" {
				send(message)
			}
		default:
			if log.GetLevel() == log.DebugLevel {
				send(message)
			}
		}
	}
}
