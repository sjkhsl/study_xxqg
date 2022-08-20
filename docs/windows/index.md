## windows使用教程

+ 浏览器访问[Release](https://github.com/johlanse/study_xxqg/releases)
+ 选择最新版本下载 ```study_xxqg_amd64.zip```
+ 将其解压到合适的位置
+ 进入解压后的文件夹，双击运行```study_xxqg.exe```,第一次打开可能会出现闪退，发现文件夹下生成了config文件夹
+ 打开config目录下的```confif.yml```文件，进行编辑，详情内容见[配置文件](../config.md)
+ 再次进行运行```study_xxqg.exe```
+ 使用浏览器打开```http://127.0.0.1:8080```
+ 推送配置请参考[推送](../push.md)

### 自定义浏览器位置

>windows默认调用系统的edge浏览器，调用目录**C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe**
>
> 若不存在该浏览器会自动尝试下载浏览器到目录下的tools文件夹下，当然也可以自定义配置浏览器位置
> 
> 修改配置文件的**edge_path**选项即可配置，配置为配置可执行文件的路径
> 
> 自定义浏览器支持chromium内核的系列浏览器，但是版本不能太高
> 
> 例如，我的chrome.exe文件在D盘的browser文件夹下，配置为**D:/browser/chrome.exe**或者**D:\\\browser\\\chrome.exe**