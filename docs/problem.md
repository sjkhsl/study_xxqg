## 常见问题<!-- {docsify-ignore} -->


+ ### 遇到问题的常用解决办法

```yaml
首先将日志项中的日志等级调整为debug
    
复现出现的错误，在issue中查找错误日志的关键字

通过搜索引擎查找问题
    
在群聊中查找聊天记录，查找置顶信息

若无解决方案，可附上关键日志和相关配置文件，在群聊中提出问题或者在github提出issue
```

+ ### windows打开**study_xxqg.exe**出现直接闪退
```yaml
  在文件路径栏输入**cmd**,然后再黑色命令窗口中输入```./study_xxqg.exe```,
  然后查看报错内容截图并在[github](https://github.com/johlanse/study_xxqg/issues) 提交issue
```

+ ### web端账号密码
```yaml
  web端账号密码默认都是admin，不是你学习强国的手机号，需要修改可自行修改配置文件
```

+ ### 关于cookie的时间问题
```yaml
原理是是通过带上当前cookie访问一个api即可，在1.0.35版本之后我通过cron定时执行保活，默认的cron是 0 */1 * * *

目前暂不知道能够续期的次数

如果你想让访问间隔时间更短或者更长，可以通过添加环境变量 CHECK_ENV 为cron值
```


+ ### windows下出现找不到浏览器的问题

```yaml

自行安装chromium内核的浏览器，包括chrome，edge浏览器之类，然后在配置文件中配置 edge_path 配置项，配置时将路径中的 \ 换成 / 或者 \\ 
```


+ ### arm设备报错```could not download driver: could not check if driver is up2date: could not run driver: exit status 127```
 ```yaml
因为playwright官方对arm设备支持会出现一些问题，所以需要手动安装node.js和chromium，并创建软连接

  apt-get install nodejs
  
  apt-get install chromium
 
  ln -s /usr/bin/chromium ./tools/browser/chromium-978106/chrome-linux/chrome
 
  ln -s /usr/bin/node ./tools/driver/ms-playwright-go/1.20.0-beta-1647057403000/node
```



+ ### linux退出终端后脚本停止运行
```yaml
 可以使用screen或者nohup等命令后台运行，具体命令自行百度

 nohup参考命令
 
 nohup ./study_xxqg > xxqg.log 2>&1

 退出程序可以通过**cat pid.pid**查看程序pid,然后kill对应pid进行退出


```

+ ### linux上退出后台正在执行的进程

```yaml
study_xxqg进程会在运行的时候将pid输出到目录下的pid.pid文件，使用kil -9 命令即可退出后台进程
```

+ ### 刷文章或者视频无法加分
```yaml
偶尔出现视频和文章无法加分的bug,可以进行等待一段时间后重启程序再次测试，目前尚不清楚造成原因
```

+ ### Host system is missing
报错大概为这样：
```
[ERROR]: [core]  初始化chrome失败 
[2022-05-13 13:43:47] [ERROR]: [core]  could not send message: could not send message to server: Host system is missing dependencies!

  Missing libraries are:
      libgtk-3.so.0
      libgdk-3.so.0

```
```shell
sudo ./tools/driver/ms-playwright-go/1.20.0-beta-1647057403000/playwright.sh install-deps

```

> 若运行后显示未找到apt-get，可百度对应系统安装apt-get的方法




+ ### 为什么运行了就卡住了

当开启了cron定时配置，微信公众号测试号配置，telegram配置这三项的任意一项后，

程序就会等待用户的指令从而卡住，所以只需要修改配置文件就可以解决