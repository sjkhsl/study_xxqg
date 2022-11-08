package conf

import (
	"bytes"
	_ "embed"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config
//  @Description:
//
type Config struct {
	Model       int    `json:"model" yaml:"model" mapstructure:"model"`
	LogLevel    string `json:"log_level" yaml:"log_level" mapstructure:"log_level"`
	ShowBrowser bool   `json:"show_browser" yaml:"show_browser" mapstructure:"show_browser"`
	Scheme      string `json:"scheme" yaml:"scheme" mapstructure:"scheme"`
	Push        struct {
		Ding struct {
			Enable      bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
			AccessToken string `json:"access_token" yaml:"access_token" mapstructure:"access_token"`
			Secret      string `json:"secret" yaml:"secret" mapstructure:"secret"`
		} `json:"ding" yaml:"ding" mapstructure:"ding"`
		PushPlus struct {
			Enable bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
			Token  string `json:"token" yaml:"token" mapstructure:"token"`
			Topic  string `json:"topic" yaml:"topic" mapstructure:"topic"`
		} `json:"push_plus" yaml:"push_plus" mapstructure:"push_plus"`
	} `json:"push" yaml:"push" mapstructure:"push"`
	TG struct {
		Enable    bool    `json:"enable" yaml:"enable" mapstructure:"enable"`
		Token     string  `json:"token" yaml:"token" mapstructure:"token"`
		ChatID    int64   `json:"chat_id" yaml:"chat_id" mapstructure:"chat_id"`
		Proxy     string  `json:"proxy" yaml:"proxy" mapstructure:"proxy"`
		CustomApi string  `json:"custom_api" yaml:"custom_api" mapstructure:"custom_api"`
		WhiteList []int64 `json:"white_list" yaml:"white_list" mapstructure:"white_list"`
	} `json:"tg" yaml:"tg" mapstructure:"tg"`
	Web struct {
		Enable     bool              `json:"enable" yaml:"enable" mapstructure:"enable"`
		Account    string            `json:"account" yaml:"account" mapstructure:"account"`
		Password   string            `json:"password" yaml:"password" mapstructure:"password"`
		Host       string            `json:"host" yaml:"host" mapstructure:"host"`
		Port       int               `json:"port" yaml:"port" mapstructure:"port"`
		CommonUser map[string]string `json:"common_user" yaml:"common_user" mapstructure:"common_user"`
	} `json:"web" yaml:"web" mapstructure:"web"`
	QQ struct {
		Enable      bool    `json:"enable" mapstructure:"enable"`
		PostAddr    string  `json:"post_addr" mapstructure:"post_addr"`
		SuperUser   int64   `json:"super_user" mapstructure:"super_user"`
		WhiteList   []int64 `json:"white_list" mapstructure:"white_list"`
		AccessToken string  `json:"access_token" mapstructure:"access_token"`
	}
	Cron           string `json:"cron" yaml:"cron" mapstructure:"cron"`
	CronRandomWait int    `json:"cron_random_wait" yaml:"cron_random_wait" mapstructure:"cron_random_wait"`
	EdgePath       string `json:"edge_path" yaml:"edge_path" mapstructure:"edge_path"`
	StartWait      int    `json:"start_wait" yaml:"start_wait" mapstructure:"start_wait"`
	// cookie强制过期时间，单位为h
	ForceExpiration int `json:"force_expiration" yaml:"force_expiration" mapstructure:"force_expiration"`
	Retry           struct {
		// 重试次数
		Times int `json:"times" yaml:"times" mapstructure:"times"`
		// 重试时间
		Intervals int `json:"intervals" yaml:"intervals" mapstructure:"intervals"`
	} `json:"retry" yaml:"retry" mapstructure:"retry"`

	Wechat struct {
		Enable        bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
		Token         string `json:"token" yaml:"token" mapstructure:"token"`
		Secret        string `json:"secret" yaml:"secret" mapstructure:"secret"`
		AppID         string `json:"app_id" yaml:"app_id" mapstructure:"app_id"`
		LoginTempID   string `json:"login_temp_id" yaml:"login_temp_id" mapstructure:"login_temp_id"`
		NormalTempID  string `json:"normal_temp_id" yaml:"normal_temp_id" mapstructure:"normal_temp_id"`
		PushLoginWarn bool   `json:"push_login_warn" yaml:"push_login_warn" mapstructure:"push_login_warn"`
		SuperOpenID   string `json:"super_open_id" yaml:"super_open_id" mapstructure:"super_open_id"`
	} `json:"wechat" yaml:"wechat" mapstructure:"wechat"`
	// 专项答题可接受的最小值
	SpecialMinScore int `json:"special_min_score" yaml:"special_min_score" mapstructure:"special_min_score"`

	PushDeer struct {
		Enable bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
		Api    string `json:"api" yaml:"api" mapstructure:"api"`
		Token  string `json:"token" yaml:"token" mapstructure:"token"`
	} `json:"push_deer" yaml:"push_deer" mapstructure:"push_deer"`

	ReverseOrder bool `json:"reverse_order" yaml:"reverse_order" mapstructure:"reverse_order"`

	JiGuangPush struct {
		Enable bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
		Secret string `json:"secret" yaml:"secret" mapstructure:"secret"`
		AppKey string `json:"app_key" yaml:"app_key" mapstructure:"app_key"`
	} `json:"ji_guang_push" yaml:"ji_guang_push" mapstructure:"ji_guang_push"`

	// github的代理地址，用于检查更新或者其他的
	GithubProxy string `json:"github_proxy" yaml:"github_proxy" mapstructure:"github_proxy"`
	// 热重载
	HotReload bool `json:"hot_reload" yaml:"hot_reload" mapstructure:"hot_reload"`

	// 自定义消息推送
	CustomMessage string `json:"custom_message" yaml:"custom_message"  mapstructure:"custom_message"`

	CustomCron string `json:"custom_cron" yaml:"custom_cron" mapstructure:"custom_cron"`

	PoolSize int `json:"pool_size" yaml:"pool_size" mapstructure:"pool_size"`

	version string `mapstructure:"version"`
}

var (
	config = Config{
		Model: 1,
	}

	configPath = "./config/config.yml"
)

//go:embed config_default.yml
var defaultConfig []byte

// SetVersion
/* @Description: 设置应用程序版本号
 * @param string2
 */
func SetVersion(string2 string) {
	config.version = string2
}

// GetVersion
/* @Description: 获取应用程序版本号
 * @return string
 */
func GetVersion() string {
	return config.version
}

func SetConfig(config2 Config) error {
	data, err := yaml.Marshal(&config2)
	if err != nil {
		log.Errorln("不能正确解析配置文件" + err.Error())
		return err
	}
	err = viper.ReadConfig(bytes.NewReader(data))
	if err != nil {
		log.Errorln("viper不能正确解析配置文件" + err.Error())
		return err
	}
	err = viper.WriteConfig()
	if err != nil {
		log.Errorln("保存到文件失败" + err.Error())
		return err
	}
	return err
}

// InitConfig
/* @Description: 初始化配置文件
 * @param path
 */
func InitConfig(path string, restart func()) {
	if path == "" {
		path = "./config/config.yml"
	}
	configPath = path
	pathDir := strings.TrimSuffix(path, "config.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(pathDir)
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warningln("检测到配置文件可能不存在，即将写入默认配置")
			err := os.WriteFile(path, defaultConfig, 0666)
			if err != nil {
				log.Errorln("写入到配置文件出现错误")
				log.Errorln(err.Error())
				return
			}
			log.Infoln("成功写入到配置文件,即将重启应用")
			restart()
		} else {
			log.Panicln("读取配置文件出现未知错误" + err.Error())
		}
	}
	viper.SetDefault("scheme", "https://johlanse.github.io/study_xxqg/scheme.html?")
	viper.SetDefault("special_min_score", 10)
	viper.SetDefault("tg.custom_api", "https://api.telegram.org")
	viper.SetDefault("pool_size", 1)
	viper.AutomaticEnv()
	err := viper.Unmarshal(&config, func(decoderConfig *mapstructure.DecoderConfig) {

	})
	if err != nil {
		log.Panicln("解析配置文件出现错误" + err.Error())
		return
	}
	if viper.GetBool("hot_reload") {
		log.Infoln("程序已开启热重载！")
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			log.Infoln("检测到配置文件变化，即将重启程序")
			restart()
		})
	}
}

// GetConfig
/**
 * @Description: 获取配置信息
 * @return Config
 */
func GetConfig() Config {
	return config
}

// GetConfigFile
/* @Description:
*  @return string
 */
func GetConfigFile() string {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return err.Error()
	}
	return string(file)
}

func SaveConfigFile(data string) error {
	err := os.WriteFile(configPath, []byte(data), 0666)
	if err != nil {
		return err
	}
	return err
}
