package lib

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	rand2 "math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/sjkhsl/study_xxqg/conf"
	"github.com/sjkhsl/study_xxqg/model"
	"github.com/sjkhsl/study_xxqg/utils"
)

const (
	MyPointsUri = "https://pc.xuexi.cn/points/my-points.html"

	DailyBUTTON   = `#app > div > div.layout-body > div > div.my-points-section > div.my-points-content > div:nth-child(4) > div.my-points-card-footer > div.buttonbox > div`
	WEEKEND       = `#app > div > div.layout-body > div > div.my-points-section > div.my-points-content > div:nth-child(7) > div.my-points-card-footer > div.buttonbox > div`
	SPECIALBUTTON = `#app > div > div.layout-body > div > div.my-points-section > div.my-points-content > div:nth-child(6) > div.my-points-card-footer > div.buttonbox > div`
)

// 每日答题
func (c *Core) RespondDaily(user *model.User, model string) {

	var title string
	retryTimes := 0
	var id int

	// 捕获所有异常，防止程序崩溃
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("答题模块异常结束或答题已完成")
			c.Push(user.PushId, "text", "答题模块异常退出或答题已完成")
			log.Errorln(err)
		}
	}()
	// 判断浏览器是否被退出
	if c.IsQuit() {
		return
	}
	// 获取用户成绩
	score, err := GetUserScore(user.ToCookies())
	if err != nil {
		log.Errorln("获取分数失败，停止每日答题", err.Error())

		return
	}
	// 创建浏览器上下文对象
	context, err := c.browser.NewContext()
	// 添加一个script,防止被检测
	_ = context.AddInitScript(playwright.BrowserContextAddInitScriptOptions{
		Script: playwright.String("Object.defineProperties(navigator, {webdriver:{get:()=>undefined}});")})
	if err != nil {
		log.Errorln("创建实例对象错误" + err.Error())
		return
	}
	// 在退出方法时关闭对象
	defer func(context playwright.BrowserContext) {
		err := context.Close()
		if err != nil {
			log.Errorln("错误的关闭了实例对象" + err.Error())
		}
	}(context)
	// 创建一个页面
	page, err := context.NewPage()
	if err != nil {
		log.Errorln("创建页面失败" + err.Error())
		return
	}
	// 退出时关闭页面
	defer func() {
		page.Close()
	}()
	// 添加用户的cookie
	err = context.AddCookies(user.ToBrowserCookies()...)
	if err != nil {
		log.Errorln("添加cookie失败" + err.Error())
		return
	}
	// 跳转到积分页面
	_, err = page.Goto(MyPointsUri, playwright.PageGotoOptions{
		Referer:   playwright.String(MyPointsUri),
		Timeout:   playwright.Float(10000),
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	if err != nil {
		log.Errorln("跳转页面失败" + err.Error())
		return
	}
	log.Infoln("已加载答题模块")
	// 判断答题类型，然后相应处理
	switch model {
	case "daily":
		{
			// 检测是否已经完成
			if score.Content["daily"].CurrentScore >= score.Content["daily"].MaxScore {
				log.Infoln("检测到每日答题已经完成，即将退出答题")

				return
			}
			// 点击每日答题的按钮
			err = page.Click(DailyBUTTON)
			if err != nil {
				log.Errorln("跳转到积分页面错误")

				return
			}
			c.Push(user.PushId, "text", "已加载每日答题模块")
		}
	case "weekly":
		{
			// 检测是否已经完成
			if score.Content["weekly"].CurrentScore >= score.Content["weekly"].MaxScore {
				log.Infoln("检测到每周答题已经完成，即将退出答题")

				return
			}
			// err = page.Click(WEEKEND)
			// if err != nil {
			//	log.Errorln("跳转到积分页面错误")
			//	return
			//}

			// 获取每周答题的ID
			id, err = getweekID(user.ToCookies())
			if err != nil {
				return
			}
			// 跳转到每周答题界面
			_, err = page.Goto(fmt.Sprintf("https://pc.xuexi.cn/points/exam-weekly-detail.html?id=%d", id), playwright.PageGotoOptions{
				Referer:   playwright.String(MyPointsUri),
				Timeout:   playwright.Float(10000),
				WaitUntil: playwright.WaitUntilStateDomcontentloaded,
			})
			if err != nil {
				log.Errorln("跳转到答题页面错误" + err.Error())
				return
			}
			c.Push(user.PushId, "text", "已加载每周答题模块")
		}
	case "special":
		{
			//获取专项答题ID
			id, err = getSpecialID(user.ToCookies())
			if err != nil {
				return
			}
			// id = 77
			// 跳转到专项答题界面
			_, err = page.Goto(fmt.Sprintf("https://pc.xuexi.cn/points/exam-paper-detail.html?id=%d", id), playwright.PageGotoOptions{
				Referer:   playwright.String(MyPointsUri),
				Timeout:   playwright.Float(10000),
				WaitUntil: playwright.WaitUntilStateDomcontentloaded,
			})
			if err != nil {
				log.Errorln("跳转到答题页面错误" + err.Error())
				return
			}
			c.Push(user.PushId, "text", "已加载专项答题模块")
		}
	}
	time.Sleep(5 * time.Second)
	// 跳转到答题页面，若返回true则说明已答完
	// if getAnswerPage(page, model) {
	//	return
	//}

	tryCount := 0
	for {
	label:
		tryCount++
		if tryCount >= 30 {
			log.Panicln("多次循环尝试答题，已超出30次，自动退出")
		}
		if c.IsQuit() {
			return
		}
		// 查看是否存在答题按钮，若按钮可用则重新提交答题
		btn, err := page.QuerySelector(`#app > div > div.layout-body > div > div.detail-body > div.action-row > button`)
		if err != nil {
			log.Infoln("获取提交按钮失败，本次答题结束" + err.Error())
			return
		}
		if btn != nil {
			enabled, err := btn.IsEnabled()
			if err != nil {
				log.Errorln(err.Error())
				continue
			}
			if enabled {
				err := btn.Click()
				if err != nil {
					log.Errorln("提交答案失败")
				}
			}
		}
		// 该元素存在则说明出现了滑块
		handle, _ := page.QuerySelector("#nc_mask > div")
		if handle != nil {
			log.Infoln(handle)
			en, err := handle.IsVisible()
			if err != nil {
				return
			}
			if en {
				page.Mouse().Move(496, 422)
				time.Sleep(1 * time.Second)
				page.Mouse().Down()

				page.Mouse().Move(772, 416, playwright.MouseMoveOptions{})
				page.Mouse().Up()
				time.Sleep(10 * time.Second)
				log.Infoln("可能存在滑块")
				c.Push(user.PushId, "text", "答题过程出现滑块，正在尝试滑动")
				en, err = handle.IsVisible()
				if err != nil {
					return
				}
				if en {
					page.Evaluate("__nc.reset()")
					goto label
				}
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
				if score.TodayScore >= 34 {
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
		if category != nil {
			_ = category.WaitForElementState(`visible`)
			time.Sleep(1 * time.Second)

			// 获取题目
			question, err := page.QuerySelector(
				`#app > div > div.layout-body > div > div.detail-body > div.question > div.q-body > div`)
			if err != nil {
				log.Errorln("未找到题目问题元素")

				return
			}
			// 获取题目类型
			categoryText, err := category.TextContent()
			if err != nil {
				log.Errorln("获取题目元素失败" + err.Error())

				return
			}
			log.Infoln("## 题目类型：" + categoryText)

			// 获取题目的问题
			questionText, err := question.TextContent()
			if err != nil {
				log.Errorln("获取题目问题失败" + err.Error())
				return
			}

			log.Infoln("## 题目：" + questionText)
			if title == questionText {
				log.Warningln("可能已经卡住，正在重试，重试次数+1")
				retryTimes++
			} else {
				retryTimes = 0
			}
			title = questionText

			// 获取答题帮助
			openTips, err := page.QuerySelector(
				`#app > div > div.layout-body > div > div.detail-body > div.question > div.q-footer > span`)
			if err != nil || openTips == nil {
				log.Errorln("未获取到题目提示信息")

				goto label
			}
			log.Debugln("开始尝试获取打开提示信息按钮")
			// 点击提示的按钮
			err = openTips.Click()
			if err != nil {
				log.Errorln("点击打开提示信息按钮失败" + err.Error())
				goto label
			}
			log.Debugln("已打开提示信息")
			// 获取页面内容
			content, err := page.Content()
			if err != nil {
				log.Errorln("获取网页全体内容失败" + err.Error())
				goto label
			}
			time.Sleep(time.Second * time.Duration(rand2.Intn(3)))
			log.Debugln("以获取网页内容")
			// 关闭提示信息
			err = openTips.Click()
			if err != nil {
				log.Errorln("点击打开提示信息按钮失败" + err.Error())

				goto label
			}
			log.Debugln("已关闭提示信息")
			// 从整个页面内容获取提示信息
			tips := getTips(content)
			log.Infoln("[提示信息]：", tips)

			if retryTimes > 4 {
				log.Warningln("重试次数太多，即将退出答题")
				options, _ := getOptions(page)
				c.Push(user.PushId, "flush", fmt.Sprintf(
					"答题过程出现异常！！</br>答题渠道：%v</br>题目ID:%v</br>题目类型：%v</br>题目：%v</br>题目选项：%v</br>提示信息：%v</br>", model, id, categoryText, questionText, strings.Join(options, " "), strings.Join(tips, " ")))
				return
			}

			// 填空题
			switch {
			case strings.Contains(categoryText, "填空题"):

				// 填充填空题
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
						if strings.Contains(strings.ReplaceAll(option, " ", ""), strings.ReplaceAll(tip, " ", "")) {
							answer = append(answer, option)
						}
					}
				}

				answer = RemoveRepByLoop(answer)

				if len(answer) < 1 {
					answer = append(answer, options...)
					log.Infoln("无法判断答案，自动选择ABCD")
				}
				log.Infoln("根据提示分别选择了", answer)
				// 多选题选择
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

				var answer []string

				if len(tips) > 1 {
					log.Warningln("检测到单选题出现多个提示信息，即将对提示信息进行合并")
					tip := strings.Join(tips, "")
					tips = []string{tip}
				}

				for _, option := range options {
					for _, tip := range tips {
						if strings.Contains(option, tip) {
							answer = append(answer, option)
						}
					}
				}
				if len(answer) < 1 {
					answer = append(answer, options[0])
					log.Infoln("无法判断答案，自动选择A")
				}

				log.Infoln("根据提示分别选择了", answer)
				err = radioCheck(page, answer)
				if err != nil {
					return
				}
			}
		}
		score, _ = GetUserScore(user.ToCookies())
	}
}

func GetAnswerPage(page playwright.Page, model string) bool {
	selectPages, err := page.QuerySelectorAll(`#app .ant-pagination .ant-pagination-item`)
	if err != nil {
		log.Errorln("获取到页码失败")

		return false
	}
	log.Infoln("共获取到", len(selectPages), "页")
	modelName := ""
	modelSlector := ""
	switch model {
	case "daily":
		return false
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
		time.Sleep(2 * time.Second)
		datas, err := page.QuerySelectorAll(modelSlector)
		if err != nil {
			log.Errorln("获取页面内容失败")
			continue
		}
		for _, data := range datas {
			content, err := data.TextContent()
			if err != nil {
				log.Errorln("获取按钮文本失败" + err.Error())
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
					return false
				}
				if enabled {
					log.Infoln("按钮可用")
				}
				data.WaitForElementState("stable", playwright.ElementHandleWaitForElementStateOptions{Timeout: playwright.Float(10000)})
				time.Sleep(5 * time.Second)
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
				return false
			}
		}
	}
	log.Infoln("检测到每周答题已经完成")
	return true
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
			log.Errorln("获取选项答案文本出现错误" + err.Error())
			return err
		}
		for _, s := range answer {
			if textContent == s {
				err := radio.Click()
				if err != nil {
					log.Errorln("点击选项出现错误" + err.Error())
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

// 获取选项
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

// 获取问题提示
func getTips(data string) []string {
	data = strings.ReplaceAll(data, " ", "")
	data = strings.ReplaceAll(data, "\n", "")
	compile := regexp.MustCompile(`<fontcolor="red">(.*?)</font>`)
	match := compile.FindAllStringSubmatch(data, -1)
	var tips []string
	for _, i := range match {
		// 新增判断提示信息为空的逻辑
		if i[1] != "" {
			tips = append(tips, i[1])
		}
	}
	return tips
}

// 填空题
func FillBlank(page playwright.Page, tips []string) error {
	video := false
	var answer []string
	if len(tips) < 1 {
		log.Warningln("检测到未获取到提示信息")
		video = true
	}
	if video {
		data1, err := page.QuerySelector("#app > div > div.layout-body > div > div.detail-body > div.question > div.q-body > div > span:nth-child(1)")
		if err != nil {
			log.Errorln("获取题目前半段失败" + err.Error())
			return err
		}
		data1Text, _ := data1.TextContent()
		log.Infoln("题目前半段：=》" + data1Text)
		searchAnswer := model.SearchAnswer(data1Text)
		if searchAnswer != "" {
			answer = append(answer, searchAnswer)
		} else {
			answer = append(answer, "不知道")
		}
	} else {
		answer = tips
	}
	inouts, err := page.QuerySelectorAll(`div.q-body > div > input`)
	if err != nil {
		log.Errorln("获取输入框错误" + err.Error())
		return err
	}
	log.Infoln("获取到", len(inouts), "个填空")
	if len(inouts) == 1 && len(tips) > 1 {
		temp := ""
		for _, tip := range tips {
			temp += tip
		}
		answer = strings.Split(temp, ",")
		log.Infoln("答案已合并处理")
	}
	var ans string
	for i := 0; i < len(inouts); i++ {
		if len(answer) < i+1 {
			ans = "不知道"
		} else {
			ans = answer[i]
		}

		err := inouts[i].Fill(ans)
		if err != nil {
			log.Errorln("填充答案失败" + err.Error())
			continue
		}
		r := rand2.Intn(4) + 1
		time.Sleep(time.Duration(r) * time.Second)
	}
	r := rand2.Intn(2)
	time.Sleep(time.Duration(r) * time.Second)
	checkNextBotton(page)
	return nil
}

// 检查下一题按钮
func checkNextBotton(page playwright.Page) {
	btns, err := page.QuerySelectorAll(`#app .action-row > button`)
	if err != nil {
		log.Errorln("未检测到按钮" + err.Error())

		return
	}
	if len(btns) <= 1 {
		err := btns[0].Click()
		if err != nil {
			log.Errorln("点击下一题按钮失败")

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

// RemoveRepByLoop 通过两重循环过滤重复元素
func RemoveRepByLoop(slc []string) []string {
	var result []string // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// 获取专项答题ID
func getSpecialID(cookies []*http.Cookie) (int, error) {
	c := req.C()
	c.SetCommonCookies(cookies...)
	// 获取专项答题列表
	repo, err := c.R().SetQueryParams(map[string]string{"pageSize": "1000", "pageNo": "1"}).Get(querySpecialList)
	if err != nil {
		log.Errorln("获取专项答题列表错误" + err.Error())
		return 0, err
	}
	dataB64, err := repo.ToString()
	if err != nil {
		log.Errorln("获取专项答题列表获取string错误" + err.Error())
		return 0, err
	}
	// 因为返回内容使用base64编码，所以需要对内容进行转码
	data, err := base64.StdEncoding.DecodeString(gjson.Get(dataB64, "data_str").String())
	if err != nil {
		log.Errorln("获取专项答题列表转换b64错误" + err.Error())
		return 0, err
	}
	// 创建实例对象
	list := new(SpecialList)
	// json序列号
	err = json.Unmarshal(data, list)
	if err != nil {
		log.Errorln("获取专项答题列表转换json错误" + err.Error())
		return 0, err
	}
	log.Infoln(fmt.Sprintf("共获取到专项答题%d个", list.TotalCount))

	// 判断是否配置选题顺序，若ReverseOrder为true则从后面选题
	if conf.GetConfig().ReverseOrder {
		for i := len(list.List) - 1; i >= 0; i-- {
			if list.List[i].TipScore == 0 {
				log.Infoln(fmt.Sprintf("获取到未答专项答题: %v，id: %v", list.List[i].Name, list.List[i].Id))
				return list.List[i].Id, nil
			}
		}
	} else {
		for _, s := range list.List {
			if s.TipScore == 0 {
				log.Infoln(fmt.Sprintf("获取到未答专项答题: %v，id: %v", s.Name, s.Id))
				return s.Id, nil
			}
		}
	}
	log.Warningln("你已不存在未答的专项答题了")
	return 0, errors.New("未找到专项答题")
}

// 获取每周答题ID
func getweekID(cookies []*http.Cookie) (int, error) {
	c := req.C()
	c.SetCommonCookies(cookies...)
	repo, err := c.R().SetQueryParams(map[string]string{"pageSize": "500", "pageNo": "1"}).Get(queryWeekList)
	if err != nil {
		log.Errorln("获取每周答题列表错误" + err.Error())
		return 0, err
	}
	dataB64, err := repo.ToString()
	if err != nil {
		log.Errorln("获取每周答题列表获取string错误" + err.Error())
		return 0, err
	}
	data, err := base64.StdEncoding.DecodeString(gjson.Get(dataB64, "data_str").String())
	if err != nil {
		log.Errorln("获取每周答题列表转换b64错误" + err.Error())
		return 0, err
	}
	list := new(WeekList)
	err = json.Unmarshal(data, list)
	if err != nil {
		log.Errorln("获取每周答题列表转换json错误" + err.Error())
		return 0, err
	}
	log.Infoln(fmt.Sprintf("共获取到每周答题%d个", list.TotalCount))

	if conf.GetConfig().ReverseOrder {
		for i := len(list.List) - 1; i >= 0; i-- {
			for _, practice := range list.List[i].Practices {
				if practice.TipScore == 0 {
					log.Infoln(fmt.Sprintf("获取到未答每周答题: %v，id: %v", practice.Name, practice.Id))
					return practice.Id, nil
				}
			}
		}
	} else {
		for _, s := range list.List {
			for _, practice := range s.Practices {
				if practice.TipScore == 0 {
					log.Infoln(fmt.Sprintf("获取到未答每周答题: %v，id: %v", practice.Name, practice.Id))
					return practice.Id, nil
				}
			}
		}
	}
	log.Warningln("你已不存在未答的每周答题了")
	return 0, errors.New("未找到每周答题")
}

func GetSpecialContent(cookies []*http.Cookie, id int) *SpecialContent {
	response, err := utils.GetClient().R().SetCookies(cookies...).SetQueryParams(map[string]string{
		"type":   "2",
		"id":     strconv.Itoa(id),
		"forced": "true",
	}).Get("https://pc-proxy-api.xuexi.cn/api/exam/service/detail/queryV3")
	if err != nil {
		return nil
	}
	data, _ := base64.StdEncoding.DecodeString(gjson.GetBytes(response.Bytes(), "data_str").String())
	log.Println(string(data))
	content := new(SpecialContent)
	_ = json.Unmarshal(data, content)
	return content
}

// 获取每周答题ID列表
func GetweekIDs(cookies []*http.Cookie) []int {
	c := req.C()
	c.SetCommonCookies(cookies...)
	repo, err := c.R().SetQueryParams(map[string]string{"pageSize": "500", "pageNo": "1"}).Get(queryWeekList)
	if err != nil {
		log.Errorln("获取每周答题列表错误" + err.Error())
		return nil
	}
	dataB64, err := repo.ToString()
	if err != nil {
		log.Errorln("获取每周答题列表获取string错误" + err.Error())
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(gjson.Get(dataB64, "data_str").String())
	if err != nil {
		log.Errorln("获取每周答题列表转换b64错误" + err.Error())
		return nil
	}
	list := new(WeekList)
	err = json.Unmarshal(data, list)
	if err != nil {
		log.Errorln("获取每周答题列表转换json错误" + err.Error())
		return nil
	}
	log.Infoln(fmt.Sprintf("共获取到每周答题%d个", list.TotalCount))
	var ids []int
	for _, l := range list.List {
		for _, practice := range l.Practices {
			ids = append(ids, practice.Id)
		}
	}
	return ids
}

// 获取专项答题ID列表
func GetSpecialIDs(cookies []*http.Cookie) []int {
	c := req.C()

	c.SetCommonCookies(cookies...)
	// 获取专项答题列表
	repo, err := c.R().SetQueryParams(map[string]string{"pageSize": "1000", "pageNo": "1"}).Get(querySpecialList)
	if err != nil {
		log.Errorln("获取专项答题列表错误" + err.Error())
		return nil
	}
	dataB64, err := repo.ToString()
	if err != nil {
		log.Errorln("获取专项答题列表获取string错误" + err.Error())
		return nil
	}
	// 因为返回内容使用base64编码，所以需要对内容进行转码
	data, err := base64.StdEncoding.DecodeString(gjson.Get(dataB64, "data_str").String())
	if err != nil {
		log.Errorln("获取专项答题列表转换b64错误" + err.Error())
		return nil
	}
	// 创建实例对象
	list := new(SpecialList)
	// json序列号
	err = json.Unmarshal(data, list)
	if err != nil {
		log.Errorln("获取专项答题列表转换json错误" + err.Error())
		return nil
	}
	log.Infoln(fmt.Sprintf("共获取到专项答题%d个", list.TotalCount))
	var ids []int
	for _, l := range list.List {
		ids = append(ids, l.Id)
	}
	return ids
}

type SpecialContent struct {
	Perfect   bool `json:"perfect"`
	TotalTime int  `json:"totalTime"`
	Questions []struct {
		HasDescribe bool `json:"hasDescribe"`
		// 提示信息
		QuestionDesc string `json:"questionDesc"`
		QuestionId   int    `json:"questionId"`
		Origin       string `json:"origin"`
		// 答案
		Answers []struct {
			AnswerId int    `json:"answerId"`
			Label    string `json:"label"`
			Content  string `json:"content"`
		} `json:"answers"`
		QuestionScore int `json:"questionScore"`
		// 题目呢偶然
		Body               string `json:"body"`
		OriginTitle        string `json:"originTitle"`
		AllCorrect         bool   `json:"allCorrect"`
		Supplier           string `json:"supplier"`
		QuestionDescOrigin string `json:"questionDescOrigin"`
		QuestionDisplay    int    `json:"questionDisplay"`
		Recommender        string `json:"recommender"`
	} `json:"questions"`
	Type               int    `json:"type"`
	TotalScore         int    `json:"totalScore"`
	PassScore          int    `json:"passScore"`
	FinishedNum        int    `json:"finishedNum"`
	UsedTime           int    `json:"usedTime"`
	Name               string `json:"name"`
	QuestionNum        int    `json:"questionNum"`
	Id                 int    `json:"id"`
	UniqueId           string `json:"uniqueId"`
	TipScoreReasonType int    `json:"tipScoreReasonType"`
}
