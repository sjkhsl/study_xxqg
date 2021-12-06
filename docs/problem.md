## 常见问题

+ ### windows打开**study_xxqg.exe**出现直接闪退
> 在文件路径栏输入**cmd**,然后再黑色命令窗口中输入```./study_xxqg.exe```,
> 然后查看报错内容截图并在[github](https://github.com/johlanse/study_xxqg/issues) 提交issue


+ ### arm设备报错```could not download driver: could not check if driver is up2date: could not run driver: exit status 127```

>因为playwright官方对arm设备支持会出现一些问题，所以需要手动安装node.js，并创建软连接
> 
> > apt-get install nodejs
> 
> > ln -s /usr/bin/node ~/.cache/ms-playwright-go/1.14.0/node

+ ### linux退出终端后脚本停止运行
> 可以使用screen或者nohup等命令后台运行，具体命令自行百度
> 
> nohup参考命令
> > nohup ./study_xxqg > xxqg.log 2>&1 & echo $! >pid.pid
>
> 退出程序可以通过**cat pid.pid**查看程序pid,然后kill对应pid进行退出

+ ### 刷文章或者视频无法加分
> 偶尔出现视频和文章无法加分的bug,可以进行等待一段时间后重启程序再次测试，目前尚不清楚造成原因