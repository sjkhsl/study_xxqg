package lib

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/model"
	"github.com/johlanse/study_xxqg/utils"
)

// Core
// @Description:
type Core struct {
	pw          *playwright.Playwright
	browser     playwright.Browser
	ShowBrowser bool
	Push        func(id string, kind string, message string)
}

// Cookie
// @Description:
type Cookie struct {
	Name     string `json:"name" yaml:"name"`
	Value    string `json:"value" yaml:"value"`
	Domain   string `json:"domain" yaml:"domain"`
	Path     string `json:"path" yaml:"path"`
	Expires  int    `json:"expires" yaml:"expires"`
	HTTPOnly bool   `json:"httpOnly" yaml:"http_only"`
	Secure   bool   `json:"secure" yaml:"secure"`
	SameSite string `json:"same_site" yaml:"same_site"`
}

type signResp struct {
	Data struct {
		Sign string `json:"sign"`
	} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Error   interface{} `json:"error"`
	Ok      bool        `json:"ok"`
}
type gennerateResp struct {
	Success   bool        `json:"success"`
	ErrorCode interface{} `json:"errorCode"`
	ErrorMsg  interface{} `json:"errorMsg"`
	Result    string      `json:"result"`
	Arguments interface{} `json:"arguments"`
}
type checkQrCodeResp struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// Init
/**
 * @Description:
 * @receiver c
 */
func (c *Core) Init() {
	if runtime.GOOS == "windows" {
		c.initWindows()
	} else {
		c.initNotWindows()
	}
}

func GetToken(code, sign, pushId string) (bool, error) {
	client := utils.GetClient()
	response, err := client.R().SetQueryParams(map[string]string{
		"code":  code,
		"state": sign + uuid.New().String(),
	}).Get("https://pc-api.xuexi.cn/login/secure_check")
	if err != nil {
		log.Errorln(err.Error())
		return false, err
	}
	uid, nick, err := GetUserInfo(response.Cookies())
	if err != nil {
		log.Errorln(err.Error())
		return false, err
	}
	var token string

	for _, cookie := range response.Cookies() {
		if cookie.Name == "token" {
			token = cookie.Value
		}
	}
	user := &model.User{
		Nick:      nick,
		Uid:       uid,
		Token:     token,
		LoginTime: time.Now().Unix(),
		PushId:    pushId,
	}
	err = model.AddUser(user)
	if err != nil {
		log.Errorln(err.Error())
	}
	log.Infoln("添加数据库成功")
	return true, err
}

// GenerateCode
/* @Description: 生成二维码
 * @receiver c
 * @return string 二维码连接
 * @return string 二维码回调查询的code
 * @return error
 */
func (c *Core) GenerateCode(pushID string) (string, string, error) {
	client := utils.GetClient()
	g := new(gennerateResp)
	_, err := client.R().SetResult(g).Get("https://login.xuexi.cn/user/qrcode/generate")
	if err != nil {
		log.Errorln(err.Error())
		return "", "", err
	}
	log.Infoln(g.Result)
	codeURL := fmt.Sprintf("https://login.xuexi.cn/login/qrcommit?showmenu=false&code=%v&appId=dingoankubyrfkttorhpou", g.Result)
	log.Infoln("登录链接： " + conf.GetConfig().Scheme + url.QueryEscape(codeURL))
	c.Push(pushID, "flush", conf.GetConfig().Scheme+url.QueryEscape(codeURL))
	c.Push(pushID, "flush", "请在一分钟内点击链接登录")
	return codeURL, g.Result, err
}

func (c *Core) CheckQrCode(code, pushID string) (*model.User, bool, error) {
	client := utils.GetClient()
	checkQrCode := func() (bool, string) {
		res := new(checkQrCodeResp)
		_, err := client.R().SetResult(res).SetFormData(map[string]string{
			"qrCode":   code,
			"goto":     "https://oa.xuexi.cn",
			"pdmToken": ""}).SetHeader("content-type", "application/x-www-form-urlencoded;charset=UTF-8").Post("https://login.xuexi.cn/login/login_with_qr")
		if err != nil {
			return false, ""
		}
		if res.Success {
			return true, res.Data
		}
		return false, ""
	}
	qrCode, s := checkQrCode()
	if !qrCode {
		return nil, false, nil
	} else {
		sign := new(signResp)
		_, err := client.R().SetResult(s).Get("https://pc-api.xuexi.cn/open/api/sns/sign")
		if err != nil {
			log.Errorln(err.Error())
			return nil, false, err
		}
		s2 := strings.Split(s, "=")[1]
		response, err := client.R().SetQueryParams(map[string]string{
			"code":  s2,
			"state": sign.Data.Sign + uuid.New().String(),
		}).Get("https://pc-api.xuexi.cn/login/secure_check")
		if err != nil {
			return nil, false, err
		}

		uid, nick, err := GetUserInfo(response.Cookies())
		if err != nil {
			return nil, false, err
		}
		user := &model.User{
			Nick:      nick,
			Uid:       uid,
			Token:     response.Cookies()[0].Value,
			LoginTime: time.Now().Unix(),
			PushId:    pushID,
		}
		err = model.AddUser(user)
		if err != nil {
			return nil, false, err
		}
		c.Push(pushID, "text", "登录成功，用户名："+nick)
		return user, true, err
	}
}

// L
/**
 * @Description:
 * @receiver c
 * @return *model.User
 * @return error
 */
func (c *Core) L(retryTimes int, pushID string) (*model.User, error) {
	_, codeData, err := c.GenerateCode(pushID)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 30; i++ {
		user, b, err := c.CheckQrCode(codeData, pushID)
		if b && err == nil {
			return user, err
		}
	}
	if retryTimes == 0 {
		return nil, errors.New("time out")
	}
	// 等待几分钟后重新执行
	time.Sleep(time.Duration(conf.GetConfig().Retry.Intervals) * time.Minute)
	c.Push(pushID, "flush", fmt.Sprintf("登录超时，将进行第%d重新次登录", retryTimes))
	return c.L(retryTimes-1, pushID)
}

func (c *Core) initWindows() {
	_, err := os.Stat("C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe")
	if err != nil {
		if os.IsNotExist(err) && conf.GetConfig().EdgePath == "" {
			log.Warningln("检测到edge浏览器不存在并且未配置edge_path，将再次运行时自动下载chrome浏览器")
			c.initNotWindows()
			return
		}
		err = nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	pwo := &playwright.RunOptions{
		DriverDirectory:     dir + "/tools/driver/",
		SkipInstallBrowsers: true,
		Browsers:            []string{"msedge"},
	}

	err = playwright.Install(pwo)
	if err != nil {
		log.Errorln("[core]", "安装playwright失败")
		log.Errorln("[core] ", err.Error())

		return
	}

	pwt, err := playwright.Run(pwo)
	if err != nil {
		log.Errorln("[core]", "初始化playwright失败")
		log.Errorln("[core] ", err.Error())

		return
	}
	c.pw = pwt
	path := conf.GetConfig().EdgePath
	if path == "" {
		path = "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
	}
	browser, err := pwt.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Args: []string{
			"--disable-extensions",
			"--disable-gpu",
			"--start-maximized",
			"--no-sandbox",
			"--window-size=500,450",
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
		DownloadsPath:   playwright.String("./tools/temp/"),
		ExecutablePath:  playwright.String(path),
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
}

func (c *Core) initNotWindows() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	_, b := os.LookupEnv("PLAYWRIGHT_BROWSERS_PATH")
	if !b {
		err = os.Setenv("PLAYWRIGHT_BROWSERS_PATH", dir+"/tools/browser/")
		if err != nil {
			log.Errorln("设置环境变量PLAYWRIGHT_BROWSERS_PATH失败" + err.Error())
			err = nil
		}
	}

	pwo := &playwright.RunOptions{
		DriverDirectory:     dir + "/tools/driver/",
		SkipInstallBrowsers: false,
		Browsers:            []string{"chromium"},
	}

	err = playwright.Install(pwo)
	if err != nil {
		log.Errorln("[core]", "安装playwright失败")
		log.Errorln("[core] ", err.Error())

		return
	}

	pwt, err := playwright.Run(pwo)
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
			"--start-maximized",
			"--no-sandbox",
			"--window-size=500,450",
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
}

func (c *Core) Quit() {
	err := c.browser.Close()
	if err != nil {
		log.Errorln("关闭浏览器失败" + err.Error())
		return
	}
	err = c.pw.Stop()
	if err != nil {
		return
	}
}

func (c *Core) IsQuit() bool {
	return !c.browser.IsConnected()
}
