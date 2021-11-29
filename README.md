### 学习强国自动化学习


该项目基于[playwright-go](https://github.com/mxschmitt/playwright-go) 开发，支持*windows*，*linux*,*mac*
若在使用过程中出现什么不可预料的问题欢迎提交*issue*


## 使用

+ 从release下载对应版本压缩包
+ windows一般下载 **study_xxqg_windows_amd64.zip**
+ 首次打开会在 ```config\config.yml```生成默认配置文件
+ 生成配置文件后第一次打开会自动安装无头浏览器，可能需要耐心等待
+ 再次打开即可运行
+ windows环境推荐直接打开浏览器扫码或者在控制台出现二维码后打开当前目录```screen.png```进行扫码
+ 其他无头浏览环境请配置[pushplus](http://www.pushplus.plus/) 推送进行扫码


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


>telegram推送使用参考
![](https://raw.githubusercontent.com/johlanse/study_xxqg/main/config/tg.jpg)
##  申明，该项目仅用于学习。

## 鸣谢

+ ### [imkenf/XueQG](https://github.com/imkenf/XueQG)