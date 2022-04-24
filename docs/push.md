## 推送配置

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
```

开启后通过浏览器访问 *http://ip:port*即可打开网址 

### 钉钉推送
配置config.yml的如下部分,具体使用教程详情参考[钉钉](https://developers.dingtalk.com/document/robots/custom-robot-access?spm=ding_open_doc.document.0.0.7f875e5903iVpC#topic-2026027)
```yaml
ding:
    enable: true
    access_token: ""
    secret: ""
```

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



