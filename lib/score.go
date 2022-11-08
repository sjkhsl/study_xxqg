package lib

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/johlanse/study_xxqg/utils"
)

type Score struct {
	TotalScore int             `json:"total_score"`
	TodayScore int             `json:"today_score"`
	Content    map[string]Data `json:"content"`
}

type Data struct {
	CurrentScore int `json:"current_score"`
	MaxScore     int `json:"max_score"`
}

func GetUserScore(cookies []*http.Cookie) (Score, error) {
	var score Score
	var resp []byte

	header := map[string]string{
		"Cache-Control": "no-cache",
	}

	client := utils.GetClient()
	response, err := client.R().SetCookies(cookies...).SetHeaders(header).Get(userTotalscoreUrl)
	if err != nil {
		log.Errorln("获取用户总分错误" + err.Error())
		return Score{}, err
	}
	resp = response.Bytes()
	score.TotalScore = int(gjson.GetBytes(resp, "data.score").Int())

	response, err = client.R().SetCookies(cookies...).SetHeaders(header).Get(userTodaytotalscoreUrl)
	if err != nil {
		log.Errorln("获取用户总分错误" + err.Error())
		return Score{}, err
	}
	resp = response.Bytes()
	score.TodayScore = int(gjson.GetBytes(resp, "data.score").Int())

	response, err = client.R().SetCookies(cookies...).SetHeaders(header).Get(userRatescoreUrl)
	if err != nil {
		log.Errorln("获取用户总分错误" + err.Error())
		return Score{}, err
	}
	resp = response.Bytes()
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
	result += fmt.Sprintf("当前学习总积分：%d\n今日得分：%d\n", score.TotalScore, score.TodayScore)
	result += fmt.Sprintf("[%v] [INFO]: 登录：%v/%v\n文章学习：%v/%v\n视频学习：%v/%v\n视频时长：%v/%v\n[%v] [INFO]: 每日答题：%v/%v\n每周答题：%v/%v\n专项答题：%v/%v",
		time.Now().Format("2006-01-02 15:04:05"),
		score.Content["login"].CurrentScore, score.Content["login"].MaxScore,
		score.Content["article"].CurrentScore, score.Content["article"].MaxScore,
		score.Content["video"].CurrentScore, score.Content["video"].MaxScore,
		score.Content["video_time"].CurrentScore, score.Content["video_time"].MaxScore,
		time.Now().Format("2006-01-02 15:04:05"),
		score.Content["daily"].CurrentScore, score.Content["daily"].MaxScore,
		score.Content["weekly"].CurrentScore, score.Content["weekly"].MaxScore,
		score.Content["special"].CurrentScore, score.Content["special"].MaxScore,
	)
	log.Infoln(result)
	return result
}

func FormatScore(score Score) string {
	result := ""
	result += fmt.Sprintf("当前学习总积分：%d\n今日得分：%d\n", score.TotalScore, score.TodayScore)
	result += fmt.Sprintf("登录：%v/%v\n文章学习：%v/%v\n视频学习：%v/%v\n视频时长：%v/%v\n每日答题：%v/%v\n每周答题：%v/%v\n专项答题：%v/%v",
		score.Content["login"].CurrentScore, score.Content["login"].MaxScore,
		score.Content["article"].CurrentScore, score.Content["article"].MaxScore,
		score.Content["video"].CurrentScore, score.Content["video"].MaxScore,
		score.Content["video_time"].CurrentScore, score.Content["video_time"].MaxScore,
		score.Content["daily"].CurrentScore, score.Content["daily"].MaxScore,
		score.Content["weekly"].CurrentScore, score.Content["weekly"].MaxScore,
		score.Content["special"].CurrentScore, score.Content["special"].MaxScore,
	)
	return result
}

func FormatScoreShort(score Score) string {
	result := ""
	result += fmt.Sprintf("当前学习总积分：%d\n今日得分：%d\n", score.TotalScore, score.TodayScore)
	result += fmt.Sprintf("登录：%v/%v\n文章学习：%v/%v\n视频学习：%v/%v\n视频时长：%v/%v\n每日答题：%v/%v\n每周答题：%v/%v\n专项答题：%v/%v",
		score.Content["login"].CurrentScore, score.Content["login"].MaxScore,
		score.Content["article"].CurrentScore, score.Content["article"].MaxScore,
		score.Content["video"].CurrentScore, score.Content["video"].MaxScore,
		score.Content["video_time"].CurrentScore, score.Content["video_time"].MaxScore,
		score.Content["daily"].CurrentScore, score.Content["daily"].MaxScore,
		score.Content["weekly"].CurrentScore, score.Content["weekly"].MaxScore,
		score.Content["special"].CurrentScore, score.Content["special"].MaxScore,
	)
	return result
}
