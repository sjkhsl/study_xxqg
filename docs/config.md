## 配置文件
```yaml
# 刷课模式，默认为3，
# 1：只刷文章何视频
# 2：只刷文章和视频和每日答题
# 3：刷文章和视频和每日答题每周答题和专项答题
model: 3

# 日志等级
# panic
# fatal
# error
# warn, warning
# info
# debug
# trace
log_level: "info"

# 是否显示浏览器
show_browser: true


push:
  ding:
    enable: false
    access_token: ""
    secret: ""
  # 目前仅支持pushplus推送二维码，默认建议使用pushplus推送
  # pushplus使用方法见：http://www.pushplus.plus/
  push_plus:
    enable: true
    token: ""

# 通过telegram进行交互模式，当配置tg.enable为true时会自动注册bot命令，
# telegram_bot使用教程https://www.dazhuanlan.com/leemode/topics/927496
tg:
  enable: false
  chat_id: 0
  token: ""
  proxy: ""

# 设置是否定时执行学习程序，格式为cron格式
# "9 19 * * *" 每天19点9分执行一次
# "* 10 * * *” 每天早上十点执行一次
cron: ""
```
