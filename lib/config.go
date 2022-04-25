package lib

import (
	_ "embed"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Config
//  @Description:
//
type Config struct {
	Model       int    `json:"model" yaml:"model"`
	LogLevel    string `json:"log_level" yaml:"log_level"`
	ShowBrowser bool   `json:"show_browser" yaml:"show_browser"`
	Scheme      string `json:"scheme" yaml:"scheme"`
	Push        struct {
		Ding struct {
			Enable      bool   `json:"enable" yaml:"enable"`
			AccessToken string `json:"access_token" yaml:"access_token"`
			Secret      string `json:"secret" yaml:"secret"`
		} `json:"ding" yaml:"ding"`
		PushPlus struct {
			Enable bool   `json:"enable" yaml:"enable"`
			Token  string `json:"token" yaml:"token"`
		} `json:"push_plus" yaml:"push_plus"`
	} `json:"push" yaml:"push"`
	TG struct {
		Enable bool   `json:"enable" yaml:"enable"`
		Token  string `json:"token" yaml:"token"`
		ChatID int64  `json:"chat_id" yaml:"chat_id"`
		Proxy  string `json:"proxy" yaml:"proxy"`
	} `json:"tg" yaml:"tg"`
	QQ struct {
	}
	Web struct {
		Enable   bool   `json:"enable" yaml:"enable"`
		Account  string `json:"account" yaml:"account"`
		Password string `json:"password" yaml:"password"`
		Host     string `json:"host" yaml:"host"`
		Port     int    `json:"port" yaml:"port"`
	} `json:"web"`
	Cron      string `json:"cron" yaml:"cron"`
	EdgePath  string `json:"edge_path" yaml:"edge_path"`
	QrCOde    bool   `json:"qr_code" yaml:"qr_code"`
	StartWait int    `json:"start_wait" yaml:"start_wait"`
}

var (
	config = Config{
		Model: 1,
	}
)

//go:embed config_default.yml
var defaultConfig []byte

// GetConfig
/**
 * @Description:
 * @return Config
 */
func GetConfig() Config {
	file, err := os.ReadFile("./config/config.yml")
	if err != nil {
		log.Warningln("检测到配置文件可能不存在")
		err := os.WriteFile("./config/config.yml", defaultConfig, 0666)
		if err != nil {
			log.Errorln("写入到配置文件出现错误")
			log.Errorln(err.Error())
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
	if config.Scheme == "" {
		config.Scheme = "https://johlanse.github.io/study_xxqg/scheme.html?"
	}
	return config
}
