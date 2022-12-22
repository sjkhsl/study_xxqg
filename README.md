# 克隆自原作者johlanse的项目

### 学习强国自动化学习


该项目基于[playwright-go](https://github.com/mxschmitt/playwright-go) 开发，支持*windows*，*linux*,*mac*


### 文档地址: https://johlanse.shhy.xyz
> 请先看文档再提出问题
 

##  申明，该项目仅用于学习。

## 鸣谢

+ ### [imkenf/XueQG](https://github.com/imkenf/XueQG)

## windows使用教程

- 浏览器访问[Release](https://github.com/sjkhsl/study_xxqg/releases)
- 选择最新版本下载 `study_xxqg_amd64.zip`
- 将其解压到合适的位置
- 进入解压后的文件夹，双击运行`study_xxqg.exe`,第一次打开可能会出现闪退，发现文件夹下生成了config文件夹
- 打开config目录下的`confif.yml`文件，进行编辑，详情内容见<u>配置文件</u>
- 再次进行运行`study_xxqg.exe`
- 使用浏览器打开`http://127.0.0.1:8080`
- 推送配置请参考<u>推送</u>

### 自定义浏览器位置

> windows默认调用系统的edge浏览器，调用目录**C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe**
> 
> 若不存在该浏览器会自动尝试下载浏览器到目录下的tools文件夹下，当然也可以自定义配置浏览器位置
> 
> 修改配置文件的**edge_path**选项即可配置，配置为配置可执行文件的路径
> 
> 自定义浏览器支持chromium内核的系列浏览器，但是版本不能太高
> 
> 例如，我的chrome.exe文件在D盘的browser文件夹下，配置为**D:/browser/chrome.exe**或者**D:\\browser\\chrome.exe**

## 可执行文件运行

- 本地访问[Releases](https://johlanse.shhy.xyz/[Release](https://github.com/johlanse/study_xxqg/releases)) ,查找对应版本并复制链接
  
- 使用wget下载对应版本压缩包
  
- > tar -xzvf study_xxqg_linux_amd64.tar.gz
  
- 运行 `./study_xxqg --init`,首次运行会生成默认配置文件
  
- 使用vim修改对应配置文件，linux建议使用tg模式运行，详情配置参考[配置](https://johlanse.shhy.xyz/#/../config),推送方式查看[push](https://johlanse.shhy.xyz/#/../push)
  
- 再次运行即可
  

### 一键安装脚本  废弃

```
wget  https://raw.githubusercontent.com/johlanse/study_xxqg/main/docs/study_xxqg_install.py && python3 study_xxqg_install.py  废弃
```

## docker运行

```
docker run --name study_xxqg -d -p 8080:8080 -v /etc/study_xxqg/:/opt/config/  sjkhsl/study_xxqg:latest
```

各个参数的含义：

- **--name study_xxqg** 运行的容器的名称，可以根据自己实际情况进行修改
- **-p 8080:8080** 将容器内部的8080端口映射到容器外面，前面是宿主机的端口，就是网页上访问的端口，后面是容器内部需要运行的端口，对应配置文件内web配置的端口就好
- **-v /etc/study_xxqg/:/opt/config/** 将容器内的/opt/config/目录映射到宿主机的/etc/study_xxqg/目录，可根据实际情况修改前面宿主机路径，映射后对应的config.yml配置文件位置就在该目录下
- **jolanse/study_xxqg:latest**镜像名称和镜像的版本，latest代表开发中的最新版本

## docker-compose运行

```
wget https://raw.githubusercontent.com/sjkhsl/study_xxqg/main/docker-compose.yml
docker-compose up -d
```

## 二种运行方式的区别

- #### 可执行文件运行
  
  可执行文件运行节省存储空间，拥有更低的占用，但是可能会存在浏览器依赖安装的问题，适合拥有一定linux基础的用户使用 ，如果系统为debian11用户，可以参考DockerFile文件中的依赖安装语句执行即可，centos用户推荐使用docker.
  
- #### docker运行
  
  docker运行不需要解决依赖问题，但是可能面临更高的运行占用，建议使用docker控制内存占用
  

## 源码运行

### 安装golang环境

- 去golang[官网](https://studygolang.com/dl) 下载对应系统的安装包，建议安装golang 1.7+
- 配置环境变量
- 具体可百度搜索golang环境安装
- 验证，任意终端中输入`go version`,显示版本信息即安装完成

### 运行项目

- 再任意终端输入一下命令
  
  ```
  cd study_xxqg
  go mod tidy
  go build ./
  ./study_xxqg
  ```
  

## 推送配置

*一共有以下五种推送方式*

- 微信公众号测试号推送
- 网页推送
- telegram推送
- 微信pushPlus推送
- 钉钉推送

> 其中pushPlus和钉钉推送相互冲突，因为两种推送模式都只能单方面配合定时运行功能使用，只能接收消息，不能发送消息；

> 在公众号测试号和tg推送以及定时三种只要配置了任意一种，程序将自动卡住等待用户指令。

> 若您想打开程序就运行，请关闭这三项配置；当前程序默认开启cron定时,所以新版若不想程序一直等待则关闭cron即可。

> 微信公众号和网页需要公网ip,若没有建议更换其他推送方式，或者自行配置内网穿透，tg推送需要配置代理或者自己反代tg的api，钉钉和pushPlus仅支持单向推送，一般配合定时使用

### 定时配置

定时任务和一下所有推送均可配合一起使用，cron的语法遵循linux标准cron语法，详情可百度自行查询

因为一些不知名的bug,观看视频时可能卡住不加分，所以建议一天运行三次左右定时，同时多次定时之间间隔不要太短

为防止定时任务每天在同一时间触发，可以配置**cron_random_wait**,等待随机时间再运行任务

除此之外，还支持以下语法

- @yearly：也可以写作@annually，表示每年第一天的 0 点。等价于0 0 1 1 *；
- @monthly：表示每月第一天的 0 点。等价于0 0 1 * *；
- @weekly：表示每周第一天的 0 点，注意第一天为周日，即周六结束，周日开始的那个 0 点。等价于0 0 * * 0；
- @daily：也可以写作@midnight，表示每天 0 点。等价于0 0 * * *；
- @hourly：表示每小时的开始。等价于0 * * * *。
- @every duration: duration为任意时间端，例如 1h,1s,1s，1h30m2s，代表间隔时间段就指向一次

### 微信公众号推送

配置config.yml的如下部分

```
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
  # 微信管理员的openid,可点击关于按钮获得，配置后请重启程序
  super_open_id: ""
```

- 前往微信[公众号开发者平台](http://mp.weixin.qq.com/debug/cgi-bin/sandbox?t=sandbox/login)，手机微信扫码登录
- 配置url为**[http://ip:port/wx](http://ip:port/wx)**,ip为你运行机器的公网ip,若使用docker运行，端口则为宿主机中映射出来的端口，ip和端口的配置和web使用同一个配置
- 设置token,需和配置项中一样
- 分别添加登录模板消息和普通模板消息，添加要求:

- 在配置文件中配置所有内容，启动程序
- 运行程序后，在浏览器中访问配置的url,页面会返回`No Signature!`,然后提交配置，若成功则关注公众号尝试运行
- docker运行方式参考<u>linux运行</u>
- 配置成功后可点击关于按钮获取open_id，然后填写到配置项的super_open_id中，然后重启容器生效

### web推送

> 适用于部署在服务器上或者家里有公网IP的设备上

配置config.yml的如下部分

```
web:
  # 启用web
  enable: true
  # 监听的ip,若只需要本机访问则设置为127.0.0.1，监听本机所有ip为0.0.0.0
  host: 0.0.0.0
  # 监听的端口号 0-65535可选
  port: 8081
  # web端登录管理员的账号
  account：admin
  # web端登录管理员的密码
  password: admin
  # web端登录普通用户的账号密码，支持多个用户,普通用户只能看到自己的信息
  common_user:
    # 代表账号为user,密码为123的普通用户，可添加多个，继续在下面写就好了
    user: 123

    # user1: 123
    # user2: 123
```

- 开启后通过浏览器访问 *[http://ip:port](http://ip:port/)*或者*[http://ip:port/new](http://ip:port/new)*即可打开网址 ,若为docker运行，则ip为宿主机公网ip,端口为docker映射到宿主机的端口
- 若无法访问，首先检查程序运行日志，查看有无报错，其次查看docker的运行情况，端口是否映射正常，然后可以通过curl命令检测在宿主机中能否访问，然后检查防火墙之类的
- 若点击登录之后出现一个小框然后无反应，则说明账户密码错误，请重新配置程序账户密码并重启程序

> 登录的账号密码是在配置文件中配置，不是学习强国的登录账号，管理员登录支持删除用户，同时能看到所有人的用户信息，普通用户就是`common_user`下面配置的用户，支持多个用户，键是账号，值是密码

### 钉钉推送

配置config.yml的如下部分,具体使用教程详情参考[钉钉](https://developers.dingtalk.com/document/robots/custom-robot-access?spm=ding_open_doc.document.0.0.7f875e5903iVpC#topic-2026027)

```
ding:
    enable: true
    access_token: ""
    secret: ""
```

- 在电脑端钉钉中创建群聊，在聊天设置中选择只能群助手，选择添加机器人，机器人类别选择webhook自定义机器人
- 机器人名字任意，机器人安全设置勾选加签，复制加签的密钥，作为secret配置项填入配置文件中
- 勾选协议，确认添加，会出现一个webhook地址，形如这样：`https://oapi.dingtalk.com/robot/send?access_token=aaabbbbcccc`
- 将上述地址中的后半段，就是access_token=之后的内容作为access_token配置项填入配置文件中，例如上述网址，则填入aaabbbccc到access_token中
- 设置定时cron,启动程序，程序会在定时时间运行脚本

### pushplus推送

配置config.yml的如下部分，具体使用教程参考[pushplus](https://www.pushplus.plus/)

```
  push_plus:
    enable: true
    token: ""
```

### telegram推送

## Telegram Bot

配置 config.yml的如下部分

```
tg:
  enable: false
  chat_id: 0
  token: ""
  # telegram的代理，不配置默认走系统代理
  proxy: ""
  # 自定义tg的api,可通过cloudflare搭建，需自备域名
  custom_api: "https://api.telegram.org"
  # 白名单id,包括群id或者用户id,若为空，则允许所有群所有用户使用，若仅用于单人，直接配置上面的chat_id就可以
  white_list:
    - 123
```

### 配置

1. 在 Tg 中搜索[`@BotFather`](https://t.me/BotFather) ，发送指令`/newbot`创建一个 bot
2. 获取你创建好的 API Token 格式为`123456789:AAaaaa-Uuuuuuuuuuu` ,要完整复制**全部内容**
3. 在 Tg 中搜索[`@userinfobot`](https://t.me/userinfobot) ，点击`START`，它就会给你发送你的信息，记住 Id 即可，是一串数字。
4. 跟你创建的 bot 会话，点击`START`，或者发送`/start`
5. 将第 2 步获取的 token 放在`tokenn`中，第 3 步获取的 Id 放到`chat_id`中，`enable`设置为 true。
6. 因为众所周知的原因，telegram推送需要进行配置代理，例如clash的代理配置为`http://127.0.0.1:7890`即可，若通过cf反代的api,,则填写到**custom_api**配置项
7. 若不配置代理的情况下会默认走系统代理，white_list建议填写自己的chat_id,为可以使用机器人的白名单，若需要在群组中使用，请相应进行配置

增加 telegram bot 指令支持

`/login` 添加一个用户

`/get_users` 获取所有cookie有效的用户

`/study 张三` 指定账号学习,若只存在一个用户则自动选择学习

`/get_scores` 获取账户积分

`/quit` 退出正在学习的实例，当长时间无响应时建议退出并查看日志然后提交issue

`/study_all` 按顺序对cookie有效的所有用户进行学习

### PushDeer推送配置

pishDeer也仅支持单向推送

配置：

```
push_deer:
  enable: true
  api: "https://api2.pushdeer.com"
  token: ""
```

自行注册pushDeer后获取token,配置token到配置文件即可，api默认为官方api,若为自建，则配置对应接口即可

## 配置文件

> 配置文件修改后需要重启程序才能生效

```
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
show_browser: false

# 跳转学习强国的scheme,默认使用本仓库的action自建scheme,若需自行修改，可直接复制仓库下/docs/scheme.html到任意静态文件服务器
scheme: "https://johlanse.github.io/study_xxqg/scheme.html?"


push:
  ding:
    enable: false
    access_token: ""
    secret: ""
  # 目前仅支持push-plus推送二维码，默认建议使用push-plus推送
  # push-plus使用方法见：http://www.pushplus.plus/
  push_plus:
    enable: false
    token: ""

# telegram交互模式配置
tg:
  enable: false
  chat_id: 0
  token: ""
  # telegram的代理，不配置默认走系统代理
  proxy: ""
  # 自定义tg的api,可通过cloudflare搭建，需自备域名
  custom_api: "https://api.telegram.org"
  # 白名单id,包括群id或者用户id,若为空，则允许所有群所有用户使用，若仅用于单人，直接配置上面的chat_id就可以
  white_list:
    - 123
# 网页端配置
web:
  # 是否启用网页
  enable: true
  #
  host: 0.0.0.0
  port: 8080
  # 网页端登录账号
  account: admin
  # 网页端登录密码
  password: admin

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
  # 微信管理员的openid,可点击关于按钮获得，配置后请重启程序
  super_open_id: ""


# pushDeer推送配置,详情参考psuhDeer官网：http://www.pushdeer.com/official.html
push_deer:
  enable: false
  api: "https://api2.pushdeer.com"
  token: ""

# 登录重试配置
retry:
  # 重试次数
  times: 0

  # 重试之间的时间间隔，单位为分钟
  intervals: 5


# 设置是否定时执行学习程序，格式为cron格式
# "9 19 * * *" 每天19点9分执行一次
# "* 10 * * *” 每天早上十点执行一次
cron: "0 0 * * *"

# 定时任务随机等待时间，单位：分钟
cron_random_wait: 0

#windows环境自定义浏览器路径，仅支持chromium系列
edge_path: ""

# 是否推送二维码
qr_code: false

# 启动时等待时间，为了防止代理启动比软件慢而报错，默认不等待，单位秒
start_wait: 0

# 专项答题可接受的最小分值，因一天重复运行的时候，若专项答题未能答满会继续答新的一套题，会浪费题
special_min_score: 10

# 题目搜索的顺序，为true则从2018年最开始搜题，否则从现在最新开始搜题
reverse_order: false
```

## 跳板搭建

因为学习强国官方app的scheme是dtxuexi://，但是大部分浏览器并不能识别该scheme,所以可以通过自行搭建跳板进行跳转。

study_xxqg官方搭建的跳板是使用github page进行搭建的，可能访问情况会比较慢，所以可以进自行搭建跳板。

## 搭建方法

- 在config目录下的创建dist目录
- 将仓库下docs目录里面的**scheme.html**和**qrcode.js**放入该目录
- 配置scheme为**[http://ip:port/dist/scheme.html](http://ip:port/dist/scheme.html)?**
- 重启程序

> 也可以通过其他静态文件服务器搭建，如 nginx等

## 常见问题

- ### 遇到问题的常用解决办法
  

```
首先将日志项中的日志等级调整为debug

复现出现的错误，在issue中查找错误日志的关键字

通过搜索引擎查找问题

在群聊中查找聊天记录，查找置顶信息

若无解决方案，可附上关键日志和相关配置文件，在群聊中提出问题或者在github提出issue
```

- ### windows打开**study_xxqg.exe**出现直接闪退
  
  ```
  在文件路径栏输入**cmd**,然后再黑色命令窗口中输入```./study_xxqg.exe```,
  然后查看报错内容截图并在[github](https://github.com/johlanse/study_xxqg/issues) 提交issue
  ```
  
- ### web端账号密码
  
  ```
  web端账号密码默认都是admin，不是你学习强国的手机号，需要修改可自行修改配置文件
  ```
  
- ### 关于cookie的时间问题
  
  ```
  原理是是通过带上当前cookie访问一个api即可，在1.0.35版本之后我通过cron定时执行保活，默认的cron是 0 */1 * * *
  ```
  

目前暂不知道能够续期的次数

如果你想让访问间隔时间更短或者更长，可以通过添加环境变量 CHECK_ENV 为cron值

````
+ ### windows下出现找不到浏览器的问题

```yaml

自行安装chromium内核的浏览器，包括chrome，edge浏览器之类，然后在配置文件中配置 edge_path 配置项，配置时将路径中的 \ 换成 / 或者 \\ 
````

- ### arm设备报错`could not download driver: could not check if driver is up2date: could not run driver: exit status 127`
  
  ```
  因为playwright官方对arm设备支持会出现一些问题，所以需要手动安装node.js和chromium，并创建软连接
  
  apt-get install nodejs
  
  apt-get install chromium
  
  ln -s /usr/bin/chromium ./tools/browser/chromium-978106/chrome-linux/chrome
  
  ln -s /usr/bin/node ./tools/driver/ms-playwright-go/1.20.0-beta-1647057403000/node
  ```
  

- ### linux退出终端后脚本停止运行
  
  ```
  可以使用screen或者nohup等命令后台运行，具体命令自行百度
  
  nohup参考命令
  
  nohup ./study_xxqg > xxqg.log 2>&1
  
  退出程序可以通过**cat pid.pid**查看程序pid,然后kill对应pid进行退出
  ```
  

````
+ ### linux上退出后台正在执行的进程

```yaml
study_xxqg进程会在运行的时候将pid输出到目录下的pid.pid文件，使用kil -9 命令即可退出后台进程
````

- ### 刷文章或者视频无法加分
  
  ```
  偶尔出现视频和文章无法加分的bug,可以进行等待一段时间后重启程序再次测试，目前尚不清楚造成原因
  ```
  
- ### Host system is missing
  
  报错大概为这样：
  
  ```
  [ERROR]: [core]  初始化chrome失败 
  [2022-05-13 13:43:47] [ERROR]: [core]  could not send message: could not send message to server: Host system is missing dependencies!
  
  Missing libraries are:
      libgtk-3.so.0
      libgdk-3.so.0
  ```
  

````
```shell
sudo ./tools/driver/ms-playwright-go/1.20.0-beta-1647057403000/playwright.sh install-deps
````

> 若运行后显示未找到apt-get，可百度对应系统安装apt-get的方法

- ### 为什么运行了就卡住了
  

当开启了cron定时配置，微信公众号测试号配置，telegram配置这三项的任意一项后，

程序就会等待用户的指令从而卡住，所以只需要修改配置文件就可以解决

study_xxqg作为一个开源程序，欢迎大家尽自己的一份力做出贡献

## 贡献要求

- ###### 热爱祖国
  
- 愿意参与开源贡献
  

## 贡献文档

> 文档采用docsify框架加上github page进行自动部署，
> 
> 你只需要在docs目录进行marikdown编写，提交pr后会
> 
> 自动生成，本地可通过docsify进行查看运行结果

## 贡献代码

### 技术需求

项目采用go语言编写  
web框架采用gin框架  
爬虫框架采用req库  
浏览器自动化框架采用playwright-go库

### 目录结构

- **\.github\** github的相关配置目录，主要存储了action自动化脚本
- **\conf\** 程序的配置文件解析，默认配置文件存放目录
- **\lib\** 程序的主要代码目录，包含了核心的各项功能，包括答题，看文章，看视频，telegram
- **\model\** 程序的用户数据存储，封装了对sqlite的操作
- **\push\** 程序的推送配置，主要包括钉钉推送和pushplus推送
- **\web\** 程序的web端操作和微信公众的操作
- **\utils\** 程序的一些工具类封装

运行时生成的目录

- **\config\** 程序的配置文件config.yml的存放目录
- **\config\logs\** 程序的日志文件存放目录
- **\dist\** 留给用户存放自定义静态文件的目录，映射路径为 \dist，需手动添加并重启程序
