package push

import (
	"fmt"
	"strings"

	"github.com/guonaihong/gout"
	log "github.com/sirupsen/logrus"
)

type PushPlus struct {
	Token string
}

func (p *PushPlus) Init() func(kind, message string) {
	var datas []string

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
		if kind == "image" {
			message = fmt.Sprintf("![](%v)", "data:image/png;base64,"+message)
			send(message)
		} else if kind == "flush" {
			send(strings.Join(datas, "\n"))
		} else {
			datas = append(datas, message)
			if len(datas) > 10 {
				send(strings.Join(datas, "\n"))
			}
		}
	}
}
