## 配置文件<!-- {docsify-ignore} -->
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

# 跳转学习强国的scheme,默认使用本仓库的action自建scheme,若需自行修改，可直接复制仓库下/docs/scheme.html到任意静态文件服务器
scheme: "https://johlanse.github.io/study_xxqg/scheme.html?"


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

# telegram交互模式配置
tg:
  enable: false
  chat_id: 0
  token: ""
  proxy: ""

# 网页端配置
web:
  # 是否启用网页
  enable: true
  #
  host: 0.0.0.0
  port: 80
  # 网页端登录账号
  account: admin
  # 网页端登录密码
  password: admin

# 登录重试配置
retry:
  # 重试次数
  times: 0

  # 重试之间的时间间隔，单位为分钟
  intervals: 5


# 设置是否定时执行学习程序，格式为cron格式
# "9 19 * * *" 每天19点9分执行一次
# "* 10 * * *” 每天早上十点执行一次
cron: ""

#windows环境自定义浏览器路径，仅支持chromium系列
edge_path: ""

# 是否推送二维码
qr_code: false

# 启动时等待时间，为了防止代理启动比软件慢而报错，默认不等待，单位秒
start_wait: 0
```
