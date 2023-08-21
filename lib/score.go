package lib

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/sjkhsl/study_xxqg/utils"
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

// 获取用户总分
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
		log.Errorln("获取用户今日得分错误" + err.Error())
		return Score{}, err
	}
	resp = response.Bytes()
	score.TodayScore = int(gjson.GetBytes(resp, "data.score").Int())

	response, err = client.R().SetCookies(cookies...).SetHeaders(header).Get(userRatescoreUrl)
	if err != nil {
		log.Errorln("获取用户详情得分错误" + err.Error())
		return Score{}, err
	}
	resp = response.Bytes()
	datas := gjson.GetBytes(resp, "data.taskProgress").Array()
	m := make(map[string]Data, 4)
	m["article"] = Data{
		CurrentScore: int(datas[0].Get("currentScore").Int()),
		MaxScore:     int(datas[0].Get("dayMaxScore").Int()),
	}
	m["video"] = Data{
		CurrentScore: int(datas[1].Get("currentScore").Int()),
		MaxScore:     int(datas[1].Get("dayMaxScore").Int()),
	}
	m["login"] = Data{
		CurrentScore: int(datas[2].Get("currentScore").Int()),
		MaxScore:     int(datas[2].Get("dayMaxScore").Int()),
	}
	m["daily"] = Data{
		CurrentScore: int(datas[3].Get("currentScore").Int()),
		MaxScore:     int(datas[3].Get("dayMaxScore").Int()),
	}

	score.Content = m

	return score, err
}

// 输出总分
func PrintScore(score Score) string {
	result := ""
	result += fmt.Sprintf("当前学习总积分：%d\n今日得分：%d\n", score.TotalScore, score.TodayScore)
	result += fmt.Sprintf("[%v] [INFO]: 登录：%v/%v\n文章学习：%v/%v\n视频学习：%v/%v\n[%v] [INFO]: 每日答题：%v/%v",
		time.Now().Format("2006-01-02 15:04:05"),
		score.Content["login"].CurrentScore, score.Content["login"].MaxScore,
		score.Content["article"].CurrentScore, score.Content["article"].MaxScore,
		score.Content["video"].CurrentScore, score.Content["video"].MaxScore,
		time.Now().Format("2006-01-02 15:04:05"),
		score.Content["daily"].CurrentScore, score.Content["daily"].MaxScore,
	)
	log.Infoln(result)
	return result
}

// 格式化总分
func FormatScore(score Score) string {
	result := ""
	result += fmt.Sprintf("当前学习总积分：%d\n今日得分：%d\n", score.TotalScore, score.TodayScore)
	result += fmt.Sprintf("登录：%v/%v\n文章学习：%v/%v\n视频学习：%v/%v\n每日答题：%v/%v\n",
		score.Content["login"].CurrentScore, score.Content["login"].MaxScore,
		score.Content["article"].CurrentScore, score.Content["article"].MaxScore,
		score.Content["video"].CurrentScore, score.Content["video"].MaxScore,
		score.Content["daily"].CurrentScore, score.Content["daily"].MaxScore,
	)
	return result
}

// 格式化短格式总分
func FormatScoreShort(score Score) string {
	result := ""
	result += fmt.Sprintf("当前学习总积分：%d\n今日得分：%d\n", score.TotalScore, score.TodayScore)
	result += fmt.Sprintf("登录：%v/%v\n文章学习：%v/%v\n视频学习：%v/%v\n每日答题：%v/%v",
		score.Content["login"].CurrentScore, score.Content["login"].MaxScore,
		score.Content["article"].CurrentScore, score.Content["article"].MaxScore,
		score.Content["video"].CurrentScore, score.Content["video"].MaxScore,
		score.Content["daily"].CurrentScore, score.Content["daily"].MaxScore,
	)
	return result
}
