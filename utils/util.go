package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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
	log.Infoln(gjson.GetBytes(response.Bytes(), "@this|@pretty"))
	if !gjson.GetBytes(response.Bytes(), "ok").Bool() &&
		gjson.GetBytes(response.Bytes(), "code").Int() == 401 &&
		gjson.GetBytes(response.Bytes(), "message").String() == "token check failed" {
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

/*时间戳->时间对象*/
func Stamp2Time(stamp int64) time.Time {
	stampStr := Stamp2Str(stamp)
	timer := Str2Time(stampStr)
	return timer
}

/**时间对象->字符串*/
func Time2Str() string {
	const shortForm = "2006-01-01 15:04:05"
	t := time.Now()
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	return str
}

/**字符串->时间对象*/
func Str2Time(formatTimeStr string) time.Time {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, formatTimeStr, loc) //使用模板在对应时区转化为time.time类型

	return theTime

}

/*时间对象->时间戳*/
func Time2Stamp() int64 {
	t := time.Now()
	millisecond := t.UnixNano() / 1e6
	return millisecond
}

/*时间戳->字符串*/
func Stamp2Str(stamp int64) string {
	timeLayout := "2006-01-02 15:04:05"
	str := time.Unix(stamp, 0).Format(timeLayout)
	return str
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
	response, err := http.Get("https://github.com/johlanse/study_xxqg/releases/download/v1.0.37-beta3/QuestionBank.db")
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
