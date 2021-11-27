package lib

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/mxschmitt/playwright-go"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/bmp"
)

type Core struct {
	pw          *playwright.Playwright
	browser     playwright.Browser
	context     *playwright.BrowserContext
	ShowBrowser bool
	Push        func(kind string, message string)
}

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
			"--window-size=500,450",
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

func (c *Core) Login() ([]Cookie, error) {
	defer func() {
		i := recover()
		if i != nil {
			log.Errorln("登录模块出现无法挽救的错误")
			log.Errorln(i)
		}
	}()
	page, err := (*c.context).NewPage()

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
	log.Infoln("[core] ", "正在等待二维码加载")

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
	removeNode(page)

	screen, _ := page.Screenshot()
	var result []byte
	buffer := bytes.NewBuffer(result)
	_ = Clip(bytes.NewReader(screen), buffer, 0, 0, 525, 35, 755, 255, 0)
	c.Push("markdown", fmt.Sprintf("![screenshot](%v) \n>点开查看登录二维码\n>请在五分钟内完成扫码", "data:image/png;base64,"+base64.StdEncoding.EncodeToString(buffer.Bytes())))
	c.Push("image", base64.StdEncoding.EncodeToString(buffer.Bytes()))
	os.WriteFile("screen.png", buffer.Bytes(), 0666)
	matrix := GetPaymentStr(bytes.NewReader(buffer.Bytes()))

	qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlack, qrcodeTerminal.ConsoleColors.BrightWhite, qrcodeTerminal.QRCodeRecoveryLevels.Low).Get(matrix.GetRawBytes()).Print()
	_, err = page.WaitForNavigation(playwright.PageWaitForNavigationOptions{
		Timeout:   playwright.Float(30 * 1000 * 5),
		URL:       nil,
		WaitUntil: nil,
	})
	if err != nil {
		log.Errorln(err.Error())

		return nil, err
	}
	cookies, err := (*c.context).Cookies() //nolint:wsl
	if err != nil {
		log.Errorln("[core] ", "获取cookie失败")
		return nil, err
	}

	var (
		cos []Cookie
	)

	for _, c := range cookies {
		co := Cookie{}
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
	c.Push("text", "登录成功，用户名："+nick)
	err = SaveUser(User{
		Cookies: cos,
		Nick:    nick,
		Uid:     info,
		Time:    time.Now().Add(time.Hour * 24).Unix(),
	})
	if err != nil {
		return nil, err
	}

	return cos, err
}

func removeNode(page playwright.Page) {
	page.Evaluate(`document.getElementsByClassName("layout-header")[0].remove()`) //nolint:errcheck
	page.Evaluate(`document.getElementsByClassName("layout-footer")[0].remove()`) //nolint:errcheck
	page.Evaluate(`document.getElementsByClassName("redflag-2")[0].remove()`)     //nolint:errcheck
	page.Evaluate(`document.getElementsByClassName("ddlogintext")[0].remove()`)   //nolint:errcheck
}

// Clip
//*  图片裁剪
//* 入参:图片输入、输出、缩略图宽、缩略图高、Rectangle{Pt(x0, y0), Pt(x1, y1)}，精度
//* 规则:如果精度为0则精度保持不变
//*
//* 返回:error
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
		//先缩略
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
