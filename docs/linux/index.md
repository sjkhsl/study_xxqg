## 可执行文件运行

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

## docker-compose运行

```shell
wget https://raw.githubusercontent.com/johlanse/study_xxqg/main/docker-compose.yml
docker-compose up -d
```

## 二种运行方式的区别

+ #### 可执行文件运行

    可执行文件运行节省存储空间，拥有更低的占用，但是可能会存在浏览器依赖安装的问题，适合拥有一定linux基础的用户使用
    ，如果系统为debian11用户，可以参考DockerFile文件中的依赖安装语句执行即可，centos用户推荐使用docker.
+ #### docker运行
    docker运行不需要解决依赖问题，但是可能面临更高的运行占用，建议使用docker控制内存占用