## 推送配置<!-- {docsify-ignore} -->

*一共有以下五种推送方式*
+ 微信公众号测试号推送
+ 网页推送
+ telegram推送
+ 微信pushPlus推送
+ 钉钉推送

> 其中pushPlus和钉钉推送相互冲突，因为两种推送模式都只能单方面配合定时运行功能使用，只能接收消息，不能发送消息；

> 在公众号测试号和tg推送以及定时三种只要配置了任意一种，程序将自动卡住等待用户指令。
> 
> 若您想打开程序就运行，请关闭这三项配置；当前程序默认开启cron定时,所以新版若不想程序一直等待则关闭cron即可。

### 微信公众号推送
配置config.yml的如下部分
```yaml
# 微信公众号测试号配置
wechat:
  # 是否启用
  enable: false
  # 开发者平台设置的token
  token: ""
  # 开发者平台的secret
  secret: ""
  # 开发者平台的appId
  app_id: ""
  # 发送登录消息需要使用的消息模板
  # 模板标题，随意  模板内容：  点我登录，然后在浏览器中打开！！
  login_temp_id: ""
  # 发送普通消息需要使用的消息模板
  # 模板标题：随意 模板内容： {{data.DATA}}
  normal_temp_id: ""
  # xxqg会每隔两小时左右检查所有用户的ck有效性，若开启该选项，会在检查失败时推送提醒消息
  push_login_warn: false
```

+ 前往微信[公众号开发者平台](http://mp.weixin.qq.com/debug/cgi-bin/sandbox?t=sandbox/login)，手机微信扫码登录
+ 配置url为**http://ip:port/wx**,ip为你运行机器的公网ip,若使用docker运行，端口则为宿主机中映射出来的端口，ip和端口的配置和web使用同一个配置
+ 设置token,需和配置项中一样
+ 分别添加登录模板消息和普通模板消息，添加要求:![](./img/wx_temp.jpg)
+ 在配置文件中配置所有内容，启动程序
+ 运行程序后，在浏览器中访问配置的url,页面会返回``No Signature!``,然后提交配置，若成功则关注公众号尝试运行
+ docker运行方式参考[linux运行](./linux/index.md)

### web推送
> 适用于部署在服务器上或者家里有公网IP的设备上

配置config.yml的如下部分
```yaml
web:
  # 启用web
  enable: true
  # 监听的ip,若只需要本机访问则设置为127.0.0.1，监听本机所有ip为0.0.0.0
  host: 0.0.0.0
  # 监听的端口号 0-65535可选
  port: 8081
  # web端登录得账号
  account： admin
  # web端登录的密码
  password: admin
```

+ 开启后通过浏览器访问 *http://ip:port*即可打开网址 ,若为docker运行，则ip为宿主机公网ip,端口为docker映射到宿主机的端口
+ 若无法访问，首先检查程序运行日志，查看有无报错，其次查看docker的运行情况，端口是否映射正常，然后可以通过curl命令检测在宿主机中能否访问，然后检查防火墙之类的
+ 若点击登录之后出现一个小框然后无反应，则说明账户密码错误，请重新配置程序账户密码并重启程序

### 钉钉推送
配置config.yml的如下部分,具体使用教程详情参考[钉钉](https://developers.dingtalk.com/document/robots/custom-robot-access?spm=ding_open_doc.document.0.0.7f875e5903iVpC#topic-2026027)
```yaml
ding:
    enable: true
    access_token: ""
    secret: ""
```
+ 在电脑端钉钉中创建群聊，在聊天设置中选择只能群助手，选择添加机器人，机器人类别选择webhook自定义机器人
+ 机器人名字任意，机器人安全设置勾选加签，复制加签的密钥，作为secret配置项填入配置文件中
+ 勾选协议，确认添加，会出现一个webhook地址，形如这样：```https://oapi.dingtalk.com/robot/send?access_token=aaabbbbcccc```
+ 将上述地址中的后半段，就是access_token=之后的内容作为access_token配置项填入配置文件中，例如上述网址，则填入aaabbbccc到access_token中
+ 设置定时cron,启动程序，程序会在定时时间运行脚本

### pushplus推送
配置config.yml的如下部分，具体使用教程参考[pushplus](https://www.pushplus.plus/)
```yaml
  push_plus:
    enable: true
    token: ""
```
### telegram推送 
## Telegram Bot
配置 config.yml的如下部分
```yaml
tg:
  enable: false
  chat_id: 0
  token: ""
  proxy: ""
```

### 配置

1. 在 Tg 中搜索[`@BotFather`](https://t.me/BotFather) ，发送指令`/newbot`创建一个 bot
2. 获取你创建好的 API Token 格式为`123456789:AAaaaa-Uuuuuuuuuuu` ,要完整复制**全部内容**
3. 在 Tg 中搜索[`@userinfobot`](https://t.me/userinfobot) ，点击`START`，它就会给你发送你的信息，记住 Id 即可，是一串数字。
4. 跟你创建的 bot 会话，点击`START`，或者发送`/start`
5. 将第 2 步获取的 token 放在`tokenn`中，第 3 步获取的 Id 放到`chat_id`中，`enable`设置为 true。
6. 因为众所周知的原因，telegram推送需要进行配置代理，例如clash的代理配置为```http://127.0.0.1:7890```即可

增加 telegram bot 指令支持

`/login` 添加一个用户

`/get_users` 获取所有cookie有效的用户

`/study 张三` 指定账号学习,若只存在一个用户则自动选择学习

`/get_scores` 获取账户积分

`/quit` 退出正在学习的实例，当长时间无响应时建议退出并查看日志然后提交issue

`/study_all` 按顺序对cookie有效的所有用户进行学习



