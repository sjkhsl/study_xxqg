package push

import (
	"fmt"

	"github.com/guonaihong/gout"
	log "github.com/sirupsen/logrus"
)

type PushPlus struct {
	Token string
}

func (p *PushPlus) Init() func(kind, message string) {
	send := func(data string) {
		err := gout.POST("http://www.pushplus.plus/send").SetJSON(gout.H{
			"token":    p.Token,
			"title":    "study_xxqg",
			"content":  data,
			"template": "markdown",
			"channel":  "wechat",
		}).Do()
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}

	return func(kind, message string) {
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
