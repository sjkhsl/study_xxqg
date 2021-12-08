package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/guonaihong/gout"
	"github.com/mxschmitt/playwright-go"
	log "github.com/sirupsen/logrus"
)

var (
	article_url_list = []string{
		"https://www.xuexi.cn/lgdata/35il6fpn0ohq.json",
		"https://www.xuexi.cn/lgdata/45a3hac2bf1j.json",
		"https://www.xuexi.cn/lgdata/1ajhkle8l72.json",
		"https://www.xuexi.cn/lgdata/1ahjpjgb4n3.json",
		"https://www.xuexi.cn/lgdata/1je1objnh73.json",
		"https://www.xuexi.cn/lgdata/1kvrj9vvv73.json",
		"https://www.xuexi.cn/lgdata/17qonfb74n3.json",
		"https://www.xuexi.cn/lgdata/1i30sdhg0n3.json"}

	video_url_list = []string{
		"https://www.xuexi.cn/lgdata/3j2u3cttsii9.json",
		"https://www.xuexi.cn/lgdata/1novbsbi47k.json",
		"https://www.xuexi.cn/lgdata/31c9ca1tgfqb.json",
		"https://www.xuexi.cn/lgdata/1oajo2vt47l.json",
		"https://www.xuexi.cn/lgdata/18rkaul9h7l.json",
		"https://www.xuexi.cn/lgdata/2qfjjjrprmdh.json",
		"https://www.xuexi.cn/lgdata/3o3ufqgl8rsn.json",
		"https://www.xuexi.cn/lgdata/525pi8vcj24p.json",
		"https://www.xuexi.cn/lgdata/1742g60067k.json"}
)

type Link struct {
	Editor       string   `json:"editor"`
	PublishTime  string   `json:"publishTime"`
	ItemType     string   `json:"itemType"`
	Author       string   `json:"author"`
	CrossTime    int      `json:"crossTime"`
	Source       string   `json:"source"`
	NameB        string   `json:"nameB"`
	Title        string   `json:"title"`
	Type         string   `json:"type"`
	Url          string   `json:"url"`
	ShowSource   string   `json:"showSource"`
	ItemId       string   `json:"itemId"`
	ThumbImage   string   `json:"thumbImage"`
	AuditTime    string   `json:"auditTime"`
	ChannelNames []string `json:"channelNames"`
	Producer     string   `json:"producer"`
	ChannelIds   []string `json:"channelIds"`
	DataValid    bool     `json:"dataValid"`
}

func getLinks(model string) ([]Link, error) {
	UID := rand.Intn(20000000) + 10000000
	learnUrl := ""
	if model == "article" {
		learnUrl = article_url_list[rand.Intn(7)]
	} else if model == "video" {
		learnUrl = video_url_list[rand.Intn(7)]
	} else {
		return nil, errors.New("model选择出现错误")
	}
	var (
		resp []byte
	)

	err := gout.GET(learnUrl + "?_st=" + strconv.Itoa(UID)).BindBody(&resp).Do()
	if err != nil {
		log.Errorln("请求连接列表出现错误" + err.Error())
		return nil, err
	}

	var links []Link
	err = json.Unmarshal(resp, &links)
	if err != nil {
		log.Errorln("解析列表出现错误" + err.Error())
		return nil, err
	}
	return links, err
}

func (c *Core) LearnArticle(cookies []Cookie) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("文章学习模块异常结束")
			log.Errorln(err)
		}
	}()
	if c.IsQuit() {
		return
	}

	score, err := GetUserScore(cookies)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	links, _ := getLinks("article")
	if score.Content["article"].CurrentScore < score.Content["article"].MaxScore {
		log.Infoln("开始加载文章学习模块")
		page, err := (*c.context).NewPage()
		if err != nil {
			return
		}

		err = (*c.context).AddCookies(cookieToParam(cookies)...)
		if err != nil {
			log.Errorln("添加cookie失败" + err.Error())
			return
		}
		tryCount := 0

		for {
			if tryCount < 20 {
				PrintScore(score)
				n := rand.Intn(len(links))
				_, err := page.Goto(links[n].Url, playwright.PageGotoOptions{
					Referer:   playwright.String(links[rand.Intn(len(links))].Url),
					Timeout:   playwright.Float(10000),
					WaitUntil: playwright.WaitUntilStateDomcontentloaded,
				})
				if err != nil {
					log.Errorln("页面跳转失败")
				}
				log.Infoln("正在学习文章：" + links[n].Title)
				c.Push("text", "正在学习文章："+links[n].Title)
				log.Infoln("文章发布时间：" + links[n].PublishTime)
				log.Infoln("文章学习链接：" + links[n].Url)
				learnTime := 70 + rand.Intn(30) + 10
				for i := 0; i < learnTime; i++ {
					if c.IsQuit() {
						return
					}
					fmt.Printf("\r[%v] [INFO]: 正在进行阅读学习中，剩余%d篇，本篇剩余时间%d秒", time.Now().Format("2006-01-02 15:04:05"), score.Content["article"].MaxScore-score.Content["article"].CurrentScore, learnTime-i)

					if rand.Float32() > 0.5 {
						go func() {
							_, err = page.Evaluate(fmt.Sprintf(`let h = document.body.scrollHeight/120*%d;document.documentElement.scrollTop=h;`, i))
							if err != nil {
								log.Errorln("文章滑动失败")
							}
						}()
					}
					time.Sleep(1 * time.Second)
				}
				fmt.Println()

				if score.Content["article"].CurrentScore >= score.Content["article"].MaxScore {
					log.Infoln("检测到本次阅读学习分数已满，退出学习")
					break
				}
				score, _ = GetUserScore(cookies)
				tryCount++
			} else {
				log.Errorln("阅读学习出现异常，稍后可重新学习")
			}
		}
	} else {
		log.Infoln("检测到文章学习已经完成")
	}
}

func (c *Core) LearnVideo(cookies []Cookie) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("视频学习模块异常结束")
			log.Errorln(err)
		}
	}()
	if c.IsQuit() {
		return
	}
	score, err := GetUserScore(cookies)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	links, _ := getLinks("video")
	if !(score.Content["video"].CurrentScore >= score.Content["video"].MaxScore && score.Content["video_time"].CurrentScore >= score.Content["video_time"].MaxScore) {
		log.Infoln("开始加载视频学习模块")
		// core := Core{}
		//core.Init()

		page, err := (*c.context).NewPage()
		if err != nil {
			return
		}

		var resp string
		err = gout.GET("http://1.15.144.22/stealth.min.js").BindBody(&resp).Do()
		if err != nil {
			return
		}
		err = page.AddInitScript(playwright.PageAddInitScriptOptions{
			Script: playwright.String(resp),
			Path:   nil,
		})
		if err != nil {
			return
		}
		err = (*c.context).AddCookies(cookieToParam(cookies)...)
		if err != nil {
			log.Errorln("添加cookie失败" + err.Error())
			return
		}
		tryCount := 0

		for {
			if tryCount < 20 {
				PrintScore(score)
				n := rand.Intn(len(links))
				_, err := page.Goto(links[n].Url, playwright.PageGotoOptions{
					Referer:   playwright.String(links[rand.Intn(len(links))].Url),
					Timeout:   playwright.Float(10000),
					WaitUntil: playwright.WaitUntilStateDomcontentloaded,
				})
				if err != nil {
					log.Errorln("页面跳转失败")
				}
				log.Infoln("正在观看视频：" + links[n].Title)
				c.Push("text", "正在观看视频："+links[n].Title)
				log.Infoln("视频发布时间：" + links[n].PublishTime)
				log.Infoln("视频学习链接：" + links[n].Url)
				learnTime := 70 + rand.Intn(30) + 10
				for i := 0; i < learnTime; i++ {
					if c.IsQuit() {
						return
					}
					fmt.Printf("\r[%v] [INFO]: 正在进行视频学习中，剩余%d个，当前剩余时间%d秒", time.Now().Format("2006-01-02 15:04:05"), score.Content["video"].MaxScore-score.Content["video"].CurrentScore, learnTime-i)

					if rand.Float32() > 0.5 {
						go func() {
							_, err := page.Evaluate(fmt.Sprintf(`let h = document.body.scrollHeight/120*%d;document.documentElement.scrollTop=h;`, i))
							if err != nil {
								log.Errorln("视频滑动失败")
							}
						}()
					}
					time.Sleep(1 * time.Second)
				}
				fmt.Println()
				if score.Content["video"].CurrentScore >= score.Content["video"].MaxScore && score.Content["video_time"].CurrentScore >= score.Content["video_time"].MaxScore {
					log.Infoln("检测到本次视频学习分数已满，退出学习")
					break
				}
				score, _ = GetUserScore(cookies)
				tryCount++
			} else {
				log.Errorln("视频学习出现异常，稍后可重新学习")
			}
		}
	} else {
		log.Infoln("检测到视频学习已经完成")
	}
}
