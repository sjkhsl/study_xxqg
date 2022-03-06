package utils

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/guonaihong/gout"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func SelfUpdate(github string, baseVersion string) {
	if github == "" {
		github = "https://github.com"
	}

	log.Infof("正在检查更新.")
	latest, err := GetLastVersion()
	if err != nil {
		log.Warnf("获取最新版本失败: %v", err)
		wait()
	}
	url := fmt.Sprintf("%v/johlanse/study_xxqg/releases/download/%v/%v", github, latest, binaryName())
	if baseVersion == latest {
		log.Info("当前版本已经是最新版本!")
		wait()
	}
	log.Info("当前最新版本为 ", latest)
	log.Warn("是否更新(y/N): ")
	r := strings.TrimSpace(readLine())
	if r != "y" && r != "Y" {
		log.Warn("已取消更新！")
		wait()
	}
	log.Info("正在更新,请稍等...")

	err = update(url)
	if err != nil {
		log.Error("更新失败: ", err)
	} else {
		log.Info("更新成功!")
	}
	wait()
}

func readLine() (str string) {
	console := bufio.NewReader(os.Stdin)
	str, _ = console.ReadString('\n')
	str = strings.TrimSpace(str)
	return
}

func wait() {
	log.Info("按 Enter 继续....")
	readLine()
	os.Exit(0)
}

func GetLastVersion() (string, error) {
	var data string
	err := gout.GET("https://api.github.com/repos/johlanse/study_xxqg/releases/latest").BindBody(&data).Do()
	if err != nil {
		return "", err
	}
	return gjson.Get(data, "tag_name").Str, err
}

func CheckVersion(oldVersion, newVersion string) bool {
	if oldVersion == "UnKnow" {
		log.Infoln("使用action版本或者自编译版本")
		return false
	}
	oldVersion = strings.ReplaceAll(oldVersion, ".", "")
	oldVersion = strings.ReplaceAll(oldVersion, "v", "")
	newVersion = strings.ReplaceAll(newVersion, ".", "")
	newVersion = strings.ReplaceAll(newVersion, "v", "")
	old, err := strconv.Atoi(oldVersion)
	newV, err := strconv.Atoi(newVersion)
	if err != nil {
		return false
	}

	return newV > old
}

func binaryName() string {
	goarch := runtime.GOARCH
	if goarch == "arm" {
		goarch += "v7"
	}
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "exe"
	}
	return fmt.Sprintf("study_xxqg_%v_%v.%v", runtime.GOOS, goarch, ext)
}
