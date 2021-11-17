package lib

import (
	rand2 "math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/mxschmitt/playwright-go"
	log "github.com/sirupsen/logrus"
)

const (
	MyPointsUri = "https://pc.xuexi.cn/points/my-points.html"

	DailyBUTTON = `#app > div > div.layout-body > div >
div.my-points-section > div.my-points-content > div:nth-child(5) > div.my-points-card-footer > div.buttonbox > div`
	WEEKEND = `#app > div > div.layout-body > 
div > div.my-points-section > div.my-points-content > div:nth-child(6) > div.my-points-card-footer > div.buttonbox > div`
	SPECIALBUTTON = `#app > div > div.layout-body > 
div > div.my-points-section > div.my-points-content > div:nth-child(7) > div.my-points-card-footer > div.buttonbox > div`
)

func (c *Core) RespondDaily(cookies []Cookie, model string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("答题模块异常结束")
		}
	}()

	// 获取用户成绩
	score, err := GetUserScore(cookies)
	if err != nil {
		log.Errorln("获取分数失败，停止每日答题", err.Error())

		return
	}

	page, err := (*c.context).NewPage()
	if err != nil {
		log.Errorln("创建页面失败" + err.Error())

		return
	}
	err = (*c.context).AddCookies(cookieToParam(cookies)...)
	if err != nil {
		log.Errorln("添加cookie信息失败，已退出答题")

		return
	}
	log.Infoln("已加载每日答题模块")

	_, err = page.Goto(MyPointsUri, playwright.PageGotoOptions{
		Referer:   playwright.String(MyPointsUri),
		Timeout:   playwright.Float(1000),
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	if err != nil {
		log.Errorln("跳转页面失败")

		return
	}
	switch model {
	case "daily":
		{
			// 检测是否已经完成
			if score.Content["daily"].CurrentScore >= score.Content["daily"].MaxScore {
				log.Infoln("检测到每日答题已经完成，即将退出答题")

				return
			}
			err = page.Click(DailyBUTTON)
			if err != nil {
				log.Errorln("跳转到积分页面错误")

				return
			}
		}
	case "weekly":
		{
			// 检测是否已经完成
			if score.Content["weekly"].CurrentScore >= score.Content["weekly"].MaxScore {
				log.Infoln("检测到每周答题已经完成，即将退出答题")

				return
			}
			err = page.Click(WEEKEND)
			if err != nil {
				log.Errorln("跳转到积分页面错误")

				return
			}
		}
	case "special":
		{
			// 检测是否已经完成
			if score.Content["special"].CurrentScore >= score.Content["special"].MaxScore {
				log.Infoln("检测到特殊答题已经完成，即将退出答题")

				return
			}
			err = page.Click(SPECIALBUTTON)
			if err != nil {
				log.Errorln("跳转到积分页面错误")

				return
			}
		}
	}
	time.Sleep(5 * time.Second)
	getAnswerPage(page, model)

	for true {
		// 查看是否存在答题按钮，若按钮可用则重新提交答题
		btn, err := page.QuerySelector(`#app > div > div.layout-body > div > div.detail-body > div.action-row > button`)
		if err != nil {
			log.Infoln("获取提交按钮失败，本次答题结束" + err.Error())
			return
		}

		enabled, err := btn.IsEnabled()
		if err != nil {
			log.Errorln(err.Error())
			continue
		}
		if enabled {
			log.Infoln("检测到有答案未提交，将重新提交")
			err := btn.Click()
			if err != nil {
				log.Errorln("提交答案失败")
			}
		}
		switch model {
		case "daily":
			{
				// 检测是否已经完成
				if score.Content["daily"].CurrentScore >= score.Content["daily"].MaxScore {
					log.Infoln("检测到每日答题已经完成，即将退出答题")

					return
				}
			}
		case "weekly":
			{
				// 检测是否已经完成
				if score.Content["weekly"].CurrentScore >= score.Content["weekly"].MaxScore {
					log.Infoln("检测到每周答题已经完成，即将退出答题")

					return
				}
			}
		case "special":
			{
				// 检测是否已经完成
				if score.Content["special"].CurrentScore >= score.Content["special"].MaxScore {
					log.Infoln("检测到特殊答题已经完成，即将退出答题")

					return
				}
			}
		}

		// 获取题目类型
		category, err := page.QuerySelector(
			`#app > div > div.layout-body > div > div.detail-body > div.question > div.q-header`)
		if err != nil {
			log.Errorln("没有找到题目元素" + err.Error())

			return
		}
		time.Sleep(1 * time.Second)

		// 获取题目
		question, err := page.QuerySelector(
			`#app > div > div.layout-body > div > div.detail-body > div.question > div.q-body > div`)
		if err != nil {
			log.Errorln("未找到题目问题元素")

			return
		}

		categoryText, err := category.TextContent()
		if err != nil {
			log.Errorln("获取题目元素失败" + err.Error())

			return
		}
		log.Infoln("## 题目类型：" + categoryText)

		questionText, err := question.TextContent()
		if err != nil {
			log.Errorln("获取题目问题失败" + err.Error())
			return
		}
		log.Infoln("## 题目：" + questionText)

		// 获取答题帮助
		openTips, err := page.QuerySelector(
			`#app > div > div.layout-body > div > div.detail-body > div.question > div.q-footer > span`)
		if err != nil {
			log.Errorln("为获取到题目提示信息" + err.Error())

			return
		}
		err = openTips.Click()
		if err != nil {
			log.Errorln("点击打开提示信息按钮失败" + err.Error())

			return
		}

		content, err := page.Content()
		if err != nil {
			log.Errorln("获取网页全体内容失败" + err.Error())

			return
		}
		err = openTips.Click()
		if err != nil {
			log.Errorln("点击打开提示信息按钮失败" + err.Error())

			return
		}

		tips := getTips(content)
		log.Infoln("[提示信息]：", tips)
		// 填空题
		switch {
		case strings.Contains(categoryText, "填空题"):
			err := FillBlank(page, tips)
			if err != nil {
				log.Errorln("填空题答题失败" + err.Error())

				return
			}
		case strings.Contains(categoryText, "多选题"):
			log.Infoln("读取到多选题")
			options, err := getOptions(page)
			if err != nil {
				log.Errorln("获取选项失败" + err.Error())

				return
			}
			log.Infoln("获取到选项答案：", options)
			log.Infoln("[多选题选项]：", options)
			var answer []string
			for _, option := range options {
				for _, tip := range tips {
					if strings.Contains(option, tip) {
						answer = append(answer, option)
					}
				}
			}
			err = radioCheck(page, answer)
			if err != nil {
				return
			}
		case strings.Contains(categoryText, "单选题"):
			log.Infoln("读取到单选题")
			options, err := getOptions(page)
			if err != nil {
				log.Errorln("获取选项失败" + err.Error())
				return
			}
			log.Infoln("获取到选项答案：", options)
			log.Infoln("[多选题选项]：", options)
			var answer []string
			for _, option := range options {
				for _, tip := range tips {
					if strings.Contains(option, tip) {
						answer = append(answer, option)
					}
				}
			}
			err = radioCheck(page, answer)
			if err != nil {
				return
			}
		}
		score, _ = GetUserScore(cookies)
	}
}

func getAnswerPage(page playwright.Page, model string) {
	selectPages, err := page.QuerySelectorAll(`#app .ant-pagination .ant-pagination-item`)
	if err != nil {
		log.Errorln("获取到页码失败")

		return
	}
	log.Infoln("共获取到", len(selectPages), "页")
	modelName := ""
	modelSlector := ""
	switch model {
	case "daily":
		return
	case "weekly":
		modelName = "每周答题"
		modelSlector = "button.ant-btn-primary"
	case "special":
		modelName = "专项答题"
		modelSlector = "#app .items .item button"
	}
	for i := 1; i <= len(selectPages); i++ {
		log.Infoln("获取到"+modelName, "第", i, "页")
		err1 := selectPages[i-1].Click()
		if err1 != nil {
			log.Errorln("点击页码失败")
		}
		datas, err := page.QuerySelectorAll(modelSlector)
		if err != nil {
			log.Errorln("获取页面内容失败")
			continue
		}
		for _, data := range datas {
			time.Sleep(3 * time.Second)
			content, err := data.TextContent()
			if err != nil {
				continue
			}
			if strings.Contains(content, "重新") || strings.Contains(content, "满分") {
				continue
			} else {
				if strings.Contains(content, "电影试题") {
					log.Infoln("发现有未答题的电影试题")
					continue
				}
				enabled, err := data.IsEnabled()
				if err != nil {
					return
				}
				if enabled {
					log.Infoln("按钮可用")
				}
				err = data.Click(playwright.ElementHandleClickOptions{
					Button:      nil,
					ClickCount:  playwright.Int(2),
					Delay:       nil,
					Force:       nil,
					Modifiers:   nil,
					NoWaitAfter: nil,
					Position:    nil,
					Timeout:     playwright.Float(100000),
				})
				if err != nil {
					log.Errorln("点击按钮失败" + err.Error())
					time.Sleep(2 * time.Second)
					continue
				}
				time.Sleep(3 * time.Second)
				return
			}
		}
	}
}

func radioCheck(page playwright.Page, answer []string) error {
	radios, err := page.QuerySelectorAll(`.q-answer.choosable`)
	if err != nil {
		log.Errorln("获取选项失败")

		return err
	}
	log.Debugln("获取到", len(radios), "个按钮")
	for _, radio := range radios {
		textContent, err := radio.TextContent()
		if err != nil {
			return err
		}
		for _, s := range answer {
			if textContent == s {
				err := radio.Click()
				if err != nil {
					return err
				}
				r := rand2.Intn(2)
				time.Sleep(time.Duration(r) * time.Second)
			}
		}
	}
	r := rand2.Intn(5)
	time.Sleep(time.Duration(r) * time.Second)
	checkNextBotton(page)
	return nil
}

func getOptions(page playwright.Page) ([]string, error) {
	handles, err := page.QuerySelectorAll(`.q-answer.choosable`)
	if err != nil {
		log.Errorln("获取选项信息失败")
		return nil, err
	}
	var options []string
	for _, handle := range handles {
		content, err := handle.TextContent()
		if err != nil {
			return nil, err
		}
		options = append(options, content)
	}
	return options, err
}

func getTips(data string) []string {
	data = strings.ReplaceAll(data, " ", "")
	data = strings.ReplaceAll(data, "\n", "")
	compile := regexp.MustCompile(`<fontcolor="red">(.*?)</font>`)
	match := compile.FindAllStringSubmatch(data, -1)
	var tips []string
	for _, i := range match {
		tips = append(tips, i[1])
	}
	return tips
}

func FillBlank(page playwright.Page, tips []string) error {
	video := false
	var answer []string
	for _, tip := range tips {
		if tip == "请观看视频" {
			video = true
		}
	}
	if video {
		answer = append(answer, "不知道")
	} else {
		answer = tips
	}
	inouts, err := page.QuerySelectorAll(`div.q-body > div > input`)
	if err != nil {
		log.Errorln("获取输入框错误" + err.Error())
		return err
	}
	log.Debugln("获取到", len(inouts), "个填空")
	if len(inouts) == 1 && len(tips) > 1 {
		temp := ""
		for _, tip := range tips {
			temp += tip
		}
		answer = strings.Split(temp, ",")
		log.Infoln("答案已合并处理" + err.Error())
	}
	for i := 0; i < len(inouts); i++ {
		err := inouts[i].Fill(answer[i])
		if err != nil {
			log.Errorln("填充答案失败" + err.Error())
			continue
		}
		r := rand2.Intn(5)
		time.Sleep(time.Duration(r) * time.Second)
	}
	r := rand2.Intn(2)
	time.Sleep(time.Duration(r) * time.Second)
	checkNextBotton(page)
	return nil
}

func checkNextBotton(page playwright.Page) {
	btns, err := page.QuerySelectorAll(`#app .action-row > button`)
	if err != nil {
		log.Errorln("未检测到按钮" + err.Error())

		return
	}
	if len(btns) <= 1 {
		err := btns[0].Check()
		if err != nil {
			log.Errorln("点击吓一题按钮失败")

			return
		}
		time.Sleep(2 * time.Second)
		_, err = btns[0].GetAttribute("disabled")
		if err != nil {
			log.Infoln("未检测到禁言属性")

			return
		}
	} else {
		err := btns[1].Click()
		if err != nil {
			log.Errorln("提交试卷失败")

			return
		}
		log.Infoln("已成功提交试卷")
	}
}
