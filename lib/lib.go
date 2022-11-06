package lib

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/model"
)

func Study(core2 *Core, u *model.User) {
	config := conf.GetConfig()
	defer func() {
		err := recover()
		if err != nil {
			logrus.Errorln("学习过程异常")
			logrus.Errorln(err)
		}
	}()
	startTime := time.Now()

	core2.LearnArticle(u)

	core2.LearnVideo(u)

	core2.LearnVideo(u)
	if config.Model == 2 {
		core2.RespondDaily(u, "daily")
	} else if config.Model == 3 {
		core2.RespondDaily(u, "daily")
		core2.RespondDaily(u, "weekly")
		core2.RespondDaily(u, "special")
	}
	endTime := time.Now()
	score, err := GetUserScore(u.ToCookies())
	if err != nil {
		logrus.Errorln("获取成绩失败")
		logrus.Debugln(err.Error())
		return
	}

	score, _ = GetUserScore(u.ToCookies())
	message := fmt.Sprintf("%v 学习完成,用时%.1f分钟\n%v", u.Nick, endTime.Sub(startTime).Minutes(), FormatScoreShort(score))
	core2.Push(u.PushId, "flush", message)
}
