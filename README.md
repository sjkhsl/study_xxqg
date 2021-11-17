### 学习强国自动化学习


该项目基于[playwright-go](https://github.com/mxschmitt/playwright-go) 开发，支持*windows*，*linux*,*mac*
若在使用过程中出现什么不可预料的问题欢迎提交*issue*


## 使用

+ 从release下载对应版本压缩包
+ windows一般下载 **study_xxqg_windows_amd64.zip**
+ 首次打开会在 ```config\config.yml```生成默认配置文件
+ 生成配置文件后第一次打开会自动安装无头浏览器，可能需要耐心等待
+ 再次打开即可运行

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

# 是否显示浏览器,linux环境请关闭该选项，否则会引发错误
show_browser: true

# 推送信息配置，建议使用钉钉
push:
  ding:
    enable: false
    access_token: ""
    secret: ""
  tg:
    enable: false
    chat_id: ""
    token: ""
# 是否启用定时功能
cron: ""
```

##  申明，该项目仅用于学习。

## 鸣谢

+ ### [imkenf/XueQG](https://github.com/imkenf/XueQG)