package lib

import (
	"errors"
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

func GetUserScore(cookies []Cookie) (Score, error) {
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
	data := string(resp)
	log.Infoln(data)
	if !gjson.GetBytes(resp, "ok").Bool() {
		return Score{}, errors.New("token check failed")
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

func PrintScore(score Score) string {
	result := ""
	result += fmt.Sprintf("当前学习总积分：%d  今日得分：%d\n", score.TodayScore, score.TodayScore)
	result += fmt.Sprintf("登录：%v/%v  文章学习：%v/%v  视频学习：%v/%v  视频时长：%v/%v\n每日答题：%v/%v  每周答题：%v/%v   专项答题：%v/%v",
		score.Content["login"].CurrentScore, score.Content["login"].MaxScore,
		score.Content["article"].CurrentScore, score.Content["article"].MaxScore,
		score.Content["video"].CurrentScore, score.Content["video"].MaxScore,
		score.Content["video_time"].CurrentScore, score.Content["video_time"].MaxScore,
		score.Content["daily"].CurrentScore, score.Content["daily"].MaxScore,
		score.Content["weekly"].CurrentScore, score.Content["weekly"].MaxScore,
		score.Content["special"].CurrentScore, score.Content["special"].MaxScore,
	)
	log.Infoln(result)
	return result
}
