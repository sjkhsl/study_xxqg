## linux基本使用

+ 本地访问[Releases]([Release](https://github.com/johlanse/study_xxqg/releases)) ,查找对应版本并复制链接
+ 使用wget下载对应版本压缩包
+ > tag -xzvf study_xxqg_linux_amd64.tag.gz
+ 运行 ```./study_xxqg```,首次运行会生成默认配置文件
+ 使用vim修改对应配置文件，linux建议使用tg模式运行，详情配置参考[配置](../config.md),推送方式查看[push](../push.md)
+ 再次运行即可


## docker运行

```
docker run --name study_xxqg -d -p 8080:8080 -v /etc/study_xxqg/:/opt/config/  jolanse/study_xxqg
```
