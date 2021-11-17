package lib

import (
	_ "embed"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Model       int    `json:"model" yaml:"model"`
	LogLevel    string `json:"log_level" yaml:"log_level"`
	ShowBrowser bool   `json:"show_browser" yaml:"show_browser"`
	Push        struct {
		Ding struct {
			Enable      bool   `json:"enable" yaml:"enable"`
			AccessToken string `json:"access_token" yaml:"access_token"`
			Secret      string `json:"secret" yaml:"secret"`
		} `json:"ding" yaml:"ding"`
		TG struct {
			Enable bool   `json:"enable" yaml:"enable"`
			Token  string `json:"token" yaml:"token"`
			ChatID string `json:"chat_id" yaml:"chat_id"`
		} `json:"tg" yaml:"tg"`
	} `json:"push" yaml:"push"`
	Cron string `json:"cron"`
}

var (
	config = Config{
		Model: 1,
	}
)

//go:embed config_default.yml
var defaultConfig []byte

func GetConfig() Config {
	file, err := os.ReadFile("./config/config.yml")
	if err != nil {
		log.Warningln("检测到配置文件可能不存在")
		err := os.WriteFile("./config/config.yml", defaultConfig, 0666)
		if err != nil {
			log.Errorln("写入到配置文件出现错误")
			return Config{}
		}
		log.Infoln("成功写入到配置文件,请重启应用")
		os.Exit(3)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Errorln(err.Error())
		return Config{}
	}
	return config
}
