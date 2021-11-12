package lib

import (
	"fmt"

	"github.com/guonaihong/gout"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Score struct {
	TotalScore int
	TodayScore int
	Content    map[string]Data
}

type Data struct {
	CurrentScore int
	MaxScore     int
}

func GetUserScore(cookies []cookie) (Score, error) {
	var score Score
	var resp []byte
	// 获取用户总分
	err := gout.GET(user_totalScore_url).SetCookies(cookieToJar(cookies)...).SetHeader(gout.H{
		"Cache-Control": "no-cache",
	}).BindBody(&resp).Do()
	if err != nil {
		log.Errorln("获取用户总分错误" + err.Error())

		return Score{}, err
	}
	log.Debugln(gjson.GetBytes(resp, "@this|@pretty"))
	score.TotalScore = int(gjson.GetBytes(resp, "data.score").Int())

	// 获取用户今日得分
	err = gout.GET(user_todayTotalScore_url).SetCookies(cookieToJar(cookies)...).SetHeader(gout.H{
		"Cache-Control": "no-cache",
	}).BindBody(&resp).Do()
	if err != nil {
		log.Errorln("获取用户每日总分错误" + err.Error())

		return Score{}, err
	}
	log.Debugln(gjson.GetBytes(resp, "@this|@pretty"))
	score.TodayScore = int(gjson.GetBytes(resp, "data.score").Int())

	err = gout.GET(user_rateScore_url).SetCookies(cookieToJar(cookies)...).SetHeader(gout.H{
		"Cache-Control": "no-cache",
	}).BindBody(&resp).Do()
	if err != nil {
		log.Errorln("获取用户积分出现错误" + err.Error())
		return Score{}, err
	}
	log.Debugln(gjson.GetBytes(resp, "@this|@pretty"))
	datas := gjson.GetBytes(resp, "data.taskProgress").Array()
	m := make(map[string]Data, 7)
	m["article"] = Data{
		CurrentScore: int(datas[0].Get("currentScore").Int()),
		MaxScore:     int(datas[0].Get("dayMaxScore").Int()),
	}
	m["video"] = Data{
		CurrentScore: int(datas[1].Get("currentScore").Int()),
		MaxScore:     int(datas[1].Get("dayMaxScore").Int()),
	}
	m["weekly"] = Data{
		CurrentScore: int(datas[2].Get("currentScore").Int()),
		MaxScore:     int(datas[2].Get("dayMaxScore").Int()),
	}
	m["video_time"] = Data{
		CurrentScore: int(datas[3].Get("currentScore").Int()),
		MaxScore:     int(datas[3].Get("dayMaxScore").Int()),
	}
	m["login"] = Data{
		CurrentScore: int(datas[4].Get("currentScore").Int()),
		MaxScore:     int(datas[4].Get("dayMaxScore").Int()),
	}
	m["special"] = Data{
		CurrentScore: int(datas[5].Get("currentScore").Int()),
		MaxScore:     int(datas[5].Get("dayMaxScore").Int()),
	}
	m["daily"] = Data{
		CurrentScore: int(datas[6].Get("currentScore").Int()),
		MaxScore:     int(datas[6].Get("dayMaxScore").Int()),
	}

	score.Content = m

	return score, err
}

func PrintScore(score Score) {
	log.Infoln(fmt.Sprintf("当前学习总积分：%d  今日得分：%d", score.TodayScore, score.TodayScore))
	for s, data := range score.Content {
		log.Infoln(s, ": ", data.CurrentScore, "/", data.MaxScore)
	}
}
