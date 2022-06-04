package lib

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/mxschmitt/playwright-go"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
	goqrcode "github.com/skip2/go-qrcode"
	"golang.org/x/image/bmp"

	"github.com/huoxue1/study_xxqg/model"
)

// Core
// @Description:
//
type Core struct {
	pw          *playwright.Playwright
	browser     playwright.Browser
	context     *playwright.BrowserContext
	ShowBrowser bool
	Push        func(kind string, message string)
}

// Cookie
// @Description:
//
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

func GetToken(code, sign string) (bool, error) {
	client := req.C()
	client.SetCommonHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
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
		UID:       uid,
		Token:     token,
		LoginTime: time.Now().Unix(),
	}
	err = model.AddUser(user)
	if err != nil {
		log.Errorln(err.Error())
	}
	log.Infoln("添加数据库成功")
	return true, err
}

// L
/**
 * @Description:
 * @receiver c
 * @return *model.User
 * @return error
 */
func (c *Core) L(retryTimes int) (*model.User, error) {
	client := req.C()
	client.OnAfterResponse(func(client *req.Client, response *req.Response) error {
		return nil
	})
	client.SetCommonHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")

	s := new(signResp)
	_, err := client.R().SetResult(s).Get("https://pc-api.xuexi.cn/open/api/sns/sign")
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}
	log.Debugln("获取sign成功==》" + s.Data.Sign)
	g := new(gennerateResp)
	_, err = client.R().SetResult(g).Get("https://login.xuexi.cn/user/qrcode/generate")
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}
	log.Infoln(g.Result)
	codeURL := fmt.Sprintf("https://login.xuexi.cn/login/qrcommit?showmenu=false&code=%v&appId=dingoankubyrfkttorhpou", g.Result)

	err = goqrcode.WriteFile(codeURL, goqrcode.Medium, 128, "qrcode.png")
	if err != nil {
		log.Errorln("图片生成错误" + err.Error())
		err = nil
	} else {
		log.Infoln("二维码已生成到目录下的qrcode.png")
	}
	if GetConfig().QrCOde {
		data, _ := os.ReadFile("qrcode.png")
		c.Push("image", base64.StdEncoding.EncodeToString(data))
	}

	qrCodeString := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlack, qrcodeTerminal.ConsoleColors.BrightWhite, qrcodeTerminal.QRCodeRecoveryLevels.Low).Get(codeURL)
	qrCodeString.Print()
	c.Push("flush", "登录链接：\r\n"+config.Scheme+url.QueryEscape(codeURL))
	checkQrCode := func() (bool, string) {
		res := new(checkQrCodeResp)
		_, err := client.R().SetResult(res).SetFormData(map[string]string{
			"qrCode":   g.Result,
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
	for i := 0; i < 150; i++ {
		code, data := checkQrCode()
		if code {
			s2 := strings.Split(data, "=")[1]
			response, err := client.R().SetQueryParams(map[string]string{
				"code":  s2,
				"state": s.Data.Sign + uuid.New().String(),
			}).Get("https://pc-api.xuexi.cn/login/secure_check")
			if err != nil {
				return nil, err
			}

			uid, nick, err := GetUserInfo(response.Cookies())
			if err != nil {
				return nil, err
			}
			user := &model.User{
				Nick:      nick,
				UID:       uid,
				Token:     response.Cookies()[0].Value,
				LoginTime: time.Now().Unix(),
			}
			err = model.AddUser(user)
			if err != nil {
				return nil, err
			}
			c.Push("text", "登录成功，用户名："+nick)
			return user, err
		}
	}
	if retryTimes == 0 {
		return nil, errors.New("time out")
	} else {
		// 等待几分钟后重新执行
		time.Sleep(time.Duration(GetConfig().Retry.Intervals) * time.Minute)
		c.Push("text", fmt.Sprintf("登录超时，将进行第%d重新次登录", retryTimes))
		return c.L(retryTimes - 1)
	}
}

func (c *Core) initWindows() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	pwt, err := playwright.Run(&playwright.RunOptions{
		DriverDirectory:     dir + "/tools/driver/",
		SkipInstallBrowsers: true,
		Browsers:            []string{"msedge"},
	})
	if err != nil {
		log.Errorln("[core]", "初始化playwright失败")
		log.Errorln("[core] ", err.Error())

		return
	}
	c.pw = pwt
	path := GetConfig().EdgePath
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

	context, err := c.browser.NewContext()
	_ = context.AddInitScript(playwright.BrowserContextAddInitScriptOptions{
		Script: playwright.String("Object.defineProperties(navigator, {webdriver:{get:()=>undefined}});")})
	if err != nil {
		return
	}
	c.context = &context
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

	pwt, err := playwright.Run(&playwright.RunOptions{
		DriverDirectory:     dir + "/tools/driver/",
		SkipInstallBrowsers: false,
		Browsers:            []string{"msedge", "chromium"},
	})
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
	context, err := c.browser.NewContext()
	_ = context.AddInitScript(playwright.BrowserContextAddInitScriptOptions{
		Script: playwright.String("Object.defineProperties(navigator, {webdriver:{get:()=>undefined}});")})
	if err != nil {
		return
	}
	c.context = &context
}

func (c *Core) Quit() {
	err := (*c.context).Close()
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

func (c *Core) IsQuit() bool {
	return !c.browser.IsConnected()
}

// func (c *Core) Login() ([]Cookie, error) {
//	defer func() {
//		i := recover()
//		if i != nil {
//			log.Errorln("登录模块出现无法挽救的错误")
//			log.Errorln(i)
//		}
//	}()
//	c.Push("text", "开始添加用户")
//	page, err := (*c.context).NewPage()
//
//	if err != nil {
//		return nil, err
//	}
//	_, err = page.Goto("https://pc.xuexi.cn/points/login.html", playwright.PageGotoOptions{
//		Referer:   nil,
//		Timeout:   playwright.Float(30000),
//		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
//	})
//	if err != nil {
//		log.Errorln("[core] ", "打开登录页面失败")
//		log.Errorln("[core] ", err.Error())
//
//		return nil, err
//	}
//	log.Infoln("[core] ", "正在等待二维码加载")
//	c.Push("text", "正在加载二维码")
//	if runtime.GOOS == "windows" {
//		time.Sleep(3 * time.Second)
//	} else {
//		_, _ = page.WaitForSelector(`#app > div > div.login_content > div > div.login_qrcode `, playwright.PageWaitForSelectorOptions{
//			State: playwright.WaitForSelectorStateVisible,
//		})
//	}
//
//	_, err = page.Evaluate(`let h = document.body.scrollWidth/2;document.documentElement.scrollTop=h;`)
//
//	if err != nil {
//		fmt.Println(err.Error())
//
//		return nil, err
//	}
//
//	log.Infoln("[core] ", "加载验证码中，请耐心等待")
//
//	//frame := page.Frame(playwright.PageFrameOptions{
//	//	Name: playwright.String(`ddlogin-iframe`),
//	//	URL:  nil,
//	//})
//	//if frame == nil {
//	//	log.Errorln("获取frame失败")
//	//}
//
//	removeNode(page)
//
//	screen, _ := page.Screenshot()
//
//	var result []byte
//	buffer := bytes.NewBuffer(result)
//	_ = Clip(bytes.NewReader(screen), buffer, 0, 0, 529, 70, 748, 284, 0)
//
//	c.Push("markdown", fmt.Sprintf("![screenshot](%v) \n>点开查看登录二维码\n>请在五分钟内完成扫码", "data:image/png;base64,"+base64.StdEncoding.EncodeToString(buffer.Bytes())))
//	c.Push("image", base64.StdEncoding.EncodeToString(buffer.Bytes()))
//	matrix := GetPaymentStr(bytes.NewReader(buffer.Bytes()))
//	log.Debugln("已获取到二维码内容：" + matrix.GetText())
//
//	c.Push("text", GetConfig().Scheme+url.QueryEscape(matrix.GetText()))
//
//	qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlack, qrcodeTerminal.ConsoleColors.BrightWhite, qrcodeTerminal.QRCodeRecoveryLevels.Low).Get(matrix.GetText()).Print()
//	_, err = page.WaitForNavigation(playwright.PageWaitForNavigationOptions{
//		Timeout:   playwright.Float(30 * 1000 * 5),
//		URL:       nil,
//		WaitUntil: nil,
//	})
//	if err != nil {
//		log.Errorln(err.Error())
//
//		return nil, err
//	}
//	cookies, err := (*c.context).Cookies() //nolint:wsl
//	if err != nil {
//		log.Errorln("[core] ", "获取cookie失败")
//		return nil, err
//	}
//
//	var (
//		cos []Cookie
//	)
//
//	for _, c := range cookies {
//		co := Cookie{}
//		co.Name = c.Name
//		co.Path = c.Path
//		co.Value = c.Value
//		co.Domain = c.Domain
//		co.Expires = c.Expires
//		co.HTTPOnly = c.HttpOnly
//		co.SameSite = c.SameSite
//		co.Secure = c.Secure
//		cos = append(cos, co)
//	}
//	info, nick, err := GetUserInfo(cos)
//	if err != nil {
//		return nil, err
//	}
//	c.Push("text", "登录成功，用户名："+nick)
//	err = SaveUser(User{
//		Cookies: cos,
//		Nick:    nick,
//		Uid:     info,
//		Time:    time.Now().Add(time.Hour * 24).Unix(),
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return cos, err
//}

func removeNode(page playwright.Page) {
	page.Evaluate(`document.getElementsByClassName("layout-header")[0].remove()`) //nolint:errcheck
	page.Evaluate(`document.getElementsByClassName("layout-footer")[0].remove()`) //nolint:errcheck
	page.Evaluate(`document.getElementsByClassName("redflag-2")[0].remove()`)     //nolint:errcheck
	page.Evaluate(`document.getElementsByClassName("ddlogintext")[0].remove()`)   //nolint:errcheck
	page.Evaluate(`document.getElementsByClassName("oath")[0].remove()`)
}

// Clip
// *  图片裁剪
// * 入参:图片输入、输出、缩略图宽、缩略图高、Rectangle{Pt(x0, y0), Pt(x1, y1)}，精度
// * 规则:如果精度为0则精度保持不变
//*
// * 返回:error
// */
func Clip(in io.Reader, out io.Writer, wi, hi, x0, y0, x1, y1, quality int) (err error) {
	err = errors.New("unknow error")
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	var origin image.Image
	var fm string
	origin, fm, err = image.Decode(in)
	if err != nil {
		log.Println(err)
		return err
	}

	if wi == 0 || hi == 0 {
		wi = origin.Bounds().Max.X
		hi = origin.Bounds().Max.Y
	}
	var canvas image.Image
	if wi != origin.Bounds().Max.X {
		// 先缩略
		canvas = resize.Thumbnail(uint(wi), uint(hi), origin, resize.Lanczos3)
	} else {
		canvas = origin
	}

	switch fm {
	case "jpeg":
		img := canvas.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		return jpeg.Encode(out, subImg, &jpeg.Options{quality})
	case "png":
		switch canvas.(type) {
		case *image.NRGBA:
			img := canvas.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			return png.Encode(out, subImg)
		case *image.RGBA:
			img := canvas.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			return png.Encode(out, subImg)
		}
	case "gif":
		img := canvas.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		return gif.Encode(out, subImg, &gif.Options{})
	case "bmp":
		img := canvas.(*image.RGBA)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
		return bmp.Encode(out, subImg)
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}

func WaitStudy(user *model.User, id string) {
	i := 0
	for i <= 180 {
		score, err := GetUserScore(user.ToCookies())
		if err != nil {
			return
		}
		if (score.Content["video"].CurrentScore >= score.Content["video"].MaxScore && score.Content["video_time"].CurrentScore >= score.Content["video_time"].MaxScore) &&
			score.Content["article"].CurrentScore >= score.Content["article"].MaxScore {
			return
		}

		time.Sleep(10 * time.Second)
		i++
	}
}

func GetPaymentStr(fi io.Reader) (paymentCodeUrl *gozxing.Result) {
	img, _, err := image.Decode(fi)
	if err != nil {
		fmt.Println(err)
	}
	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		fmt.Println(err)
	}

	return result
}
