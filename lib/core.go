package lib

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/mxschmitt/playwright-go"
	log "github.com/sirupsen/logrus"
	"github.com/tuotoo/qrcode"
)

type Core struct {
	pw          *playwright.Playwright
	browser     playwright.Browser
	context     playwright.BrowserContext
	ShowBrowser bool
}

type cookie struct {
	Name     string `json:"name" yaml:"name"`
	Value    string `json:"value" yaml:"value"`
	Domain   string `json:"domain" yaml:"domain"`
	Path     string `json:"path" yaml:"path"`
	Expires  int    `json:"expires" yaml:"expires"`
	HTTPOnly bool   `json:"httpOnly" yaml:"http_only"`
	Secure   bool   `json:"secure" yaml:"secure"`
	SameSite string `json:"same_site" yaml:"same_site"`
}

func (c *Core) Init() {
	pwt, err := playwright.Run()
	if err != nil {
		log.Errorln("[core]", "初始化playwright失败")
		log.Errorln("[core] ", err.Error())

		return
	}
	c.pw = pwt
	browser, err := pwt.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Args: []string{
			"--disable-extensions",
			"--disable-gpu",
			"--no-sandbox",
			"--window-size=540,400",
			"--start-maximized",
			"--mute-audio",
			"--window-position=0,0",
			"--ignore-certificate-errors",
			"--ignore-ssl-errors",
			"--disable-features=RendererCodeIntegrity",
			"--disable-blink-features",
			"--disable-blink-features=AutomationControlled",
		},
		Channel:         nil,
		ChromiumSandbox: nil,
		Devtools:        nil,
		DownloadsPath:   nil,
		ExecutablePath:  nil,
		HandleSIGHUP:    nil,
		HandleSIGINT:    nil,
		HandleSIGTERM:   nil,
		Headless:        playwright.Bool(!c.ShowBrowser),
		Proxy:           nil,
		SlowMo:          nil,
		Timeout:         nil,
	})
	if err != nil {
		log.Errorln("[core] ", "初始化chrome失败")
		log.Errorln("[core] ", err.Error())
		return
	}
	c.browser = browser
	context, err := c.browser.NewContext()
	if err != nil {
		return
	}
	c.context = context
}

func (c *Core) Quit() {
	err := c.context.Close()
	if err != nil {
		return
	}
	err = c.browser.Close()
	if err != nil {
		return
	}
	err = c.pw.Stop()
	if err != nil {
		return
	}
}

func (c *Core) Login() ([]cookie, error) {
	page, err := c.context.NewPage()

	if err != nil {
		return nil, err
	}
	_, err = page.Goto("https://pc.xuexi.cn/points/login.html", playwright.PageGotoOptions{
		Referer:   nil,
		Timeout:   playwright.Float(30000),
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	if err != nil {
		log.Errorln("[core] ", "打开登录页面失败")
		log.Errorln("[core] ", err.Error())

		return nil, err
	}
	log.Infoln("[core] ", "正在等待二维码扫描")

	_, _ = page.WaitForSelector(`#app > div > div.login_content > div > div.login_qrcode `)

	_, err = page.Evaluate(`let h = document.body.scrollWidth/2;document.documentElement.scrollTop=h;`)

	if err != nil {
		fmt.Println(err.Error())

		return nil, err
	}

	log.Infoln("[core] ", "加载验证码中，请耐心等待")

	frame := page.Frame(playwright.PageFrameOptions{
		Name: playwright.String(`ddlogin-iframe`),
		URL:  nil,
	})
	if frame == nil {
		log.Errorln("获取frame失败")
	}

	selector, err := frame.QuerySelector(`img`)

	if err != nil {
		log.Errorln(err.Error())

		return nil, err
	}

	img, err := selector.GetAttribute(`src`)
	if err != nil {
		log.Errorln(err.Error())

		return nil, err
	}
	img = strings.ReplaceAll(img, "data:image/png;base64,", "")
	//go sendToQQ(img)
	data, err := base64.StdEncoding.DecodeString(img)
	if err != nil {
		return nil, err
	}
	os.WriteFile("qrcode.png", data, 0666)
	matrix, err := qrcode.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	qrcodeTerminal.New().Get(matrix.Content).Print()

	_, err = page.WaitForNavigation(playwright.PageWaitForNavigationOptions{
		Timeout:   playwright.Float(30 * 1000 * 5),
		URL:       nil,
		WaitUntil: nil,
	})
	if err != nil {
		log.Errorln(err.Error())

		return nil, err
	}
	cookies, err := c.context.Cookies() //nolint:wsl
	if err != nil {
		log.Errorln("[core] ", "获取cookie失败")
		return nil, err
	}

	var (
		cos []cookie
	)

	for _, c := range cookies {
		co := cookie{}
		co.Name = c.Name
		co.Path = c.Path
		co.Value = c.Value
		co.Domain = c.Domain
		co.Expires = c.Expires
		co.HTTPOnly = c.HttpOnly
		co.SameSite = c.SameSite
		co.Secure = c.Secure
		cos = append(cos, co)
	}
	info, nick, err := GetUserInfo(cos)
	if err != nil {
		return nil, err
	}
	err = SaveUser(User{
		Cookies: cos,
		Nick:    nick,
		Uid:     info,
	})
	if err != nil {
		return nil, err
	}

	return cos, err
}

//func sendToQQ(img string) {
//	err := gout.POST("http://127.0.0.1:5700/send_private_msg").SetJSON(map[string]interface{}{
//		"user_id": int64(3343780376),
//		"message": map[string]interface{}{
//			"type": "image",
//			"data": map[string]interface{}{
//				"file": "base64://" + img,
//			},
//		},
//	}).Do()
//	if err != nil {
//		return
//	}
//}
