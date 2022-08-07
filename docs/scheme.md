因为学习强国官方app的scheme是dtxuexi://，但是大部分浏览器并不能识别该scheme,所以可以通过自行搭建跳板进行跳转。

study_xxqg官方搭建的跳板是使用github page进行搭建的，可能访问情况会比较慢，所以可以进自行搭建跳板。

## 搭建方法
+ 在config目录下的创建dist目录
+ 将仓库下docs目录里面的**scheme.html**和**qrcode.js**放入该目录
+ 配置scheme为**http://ip:port/dist/scheme.html?**
+ 重启程序

> 也可以通过其他静态文件服务器搭建，如 nginx等