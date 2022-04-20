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
	return func(kind, message string) {
		if kind == "image" {
			message = fmt.Sprintf("![](%v)", "data:image/png;base64,"+message)
		}
		err := gout.POST("http://www.pushplus.plus/send").SetJSON(gout.H{
			"token":    p.Token,
			"title":    "study_xxqg",
			"content":  message,
			"template": "markdown",
			"channel":  "wechat",
		}).Do()
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}
}
