package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/utils/update"
)

// Restart
/* @Description:
 */
func Restart() {
	//once := sync.Once{}
	//once.Do(func() {
	//	log.Infoln("程序启动命令： " + strings.Join(os.Args, " "))
	//	cmd := exec.Command(strings.Join(os.Args, " "))
	//	cmd.Stdout = os.Stdout
	//	cmd.Stderr = os.Stderr
	//	cmd.Stdin = os.Stdin
	//	cmd.Start()
	//	os.Exit(3)
	//})
	os.Exit(201)

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
func CheckUserCookie(cookies []*http.Cookie) (bool, error) {
	client := GetClient()
	response, err := client.R().SetCookies(cookies...).Get("https://pc-api.xuexi.cn/open/api/score/get")
	if err != nil {
		log.Errorln("获取用户总分错误" + err.Error())
		return true, err
	}
	if !gjson.GetBytes(response.Bytes(), "ok").Bool() {
		return false, err
	}
	return true, err
}

var (
	dbSum = "d6e455f03b419af108cced07ea1d17f8268400ad1b6d80cb75d58e952a5609bf"
)

// CheckQuestionDB
/**
 * @Description: 检查数据库文件完整性
 * @param user
 * @return bool
 */
func CheckQuestionDB() bool {

	if !FileIsExist("./QuestionBank.db") {
		return false
	}
	f, err := os.Open("./QuestionBank.db")
	if err != nil {
		log.Errorln(err.Error())
		return false
	}

	defer f.Close()
	h := sha256.New()
	//h := sha1.New()
	//h := sha512.New()

	if _, err := io.Copy(h, f); err != nil {
		log.Errorln(err.Error())
		return false
	}

	// 格式化为16进制字符串
	sha := fmt.Sprintf("%x", h.Sum(nil))
	log.Infoln("db_sha: " + sha)
	if sha != dbSum {
		return false
	}
	return true
}

func DownloadDbFile() {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("下载题库文件意外错误")
			log.Errorln(err)
		}
	}()
	log.Infoln("正在从github下载题库文件！")
	response, err := http.Get("https://github.com/johlanse/study_xxqg/releases/download/v1.0.36/QuestionBank.db")
	if err != nil {
		log.Errorln("下载db文件错误" + err.Error())
		return
	}
	data, _ := io.ReadAll(response.Body)
	err = os.WriteFile("./QuestionBank.db", data, 0666)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
}
