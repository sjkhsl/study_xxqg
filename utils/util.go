package utils

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/sjkhsl/study_xxqg/conf"
	"github.com/sjkhsl/study_xxqg/utils/update"
)

// Restart
/* @Description:
 */
func Restart() {
	cmd := exec.Command("./study_xxqg")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	go func() {
		cmd.Start()
		os.Exit(3)
	}()
}

func GetAbout() string {
	msg := "study_xxqg\n程序版本："
	msg += conf.GetVersion()
	msg += "\n" + update.CheckUpdate(conf.GetVersion())
	return msg
}

// CheckUserCookie
/**
 * @Description: 获取用户成绩
 * @param user
 * @return bool
 */
func CheckUserCookie(cookies []*http.Cookie) bool {
	client := req.C().DevMode()
	response, err := client.R().SetCookies(cookies...).Get("https://pc-api.xuexi.cn/open/api/score/get")
	if err != nil {
		logrus.Errorln("获取用户总分错误" + err.Error())
		return false
	}
	if !gjson.GetBytes(response.Bytes(), "ok").Bool() {
		return false
	}
	return true
}
