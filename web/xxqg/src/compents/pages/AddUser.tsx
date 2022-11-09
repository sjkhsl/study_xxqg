import React, {Component} from "react";
import {Button, Dialog, Toast} from "antd-mobile";
import {checkQrCode, getLink, getToken} from "../../utils/api";

class AddUser extends Component<any, any>{


    constructor(props: any) {
        super(props);
        this.state = {
            img : "你还未获取登录链接",
            link: "",
        };
    }

    componentWillUnmount() {
        if (this.state.check !== undefined){
            clearInterval(this.state.check)
        }

    }
   isMobile = ()=> {
       return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(window.navigator.userAgent);
    }

    isWechat = ()=> {
        if (/MicroMessenger/i.test(window.navigator.userAgent)){
            return "inline"
        }else {
            return "none"
        }
    }

    click = async () => {
        console.log(this.isMobile())
        if (!this.isMobile()){
            Dialog.show({
                title:"提醒",
                content:"网页端不再生成二维码，请在移动端访问该网页，会自动跳转到xxqg",
                closeOnAction:true,
                closeOnMaskClick:true,
            })
            return
        }
        let data = await getLink()

        this.setState({
            img: data.url,
            link: data.code
        })
        let check = setInterval(async ()=>{
            let resp = await checkQrCode(data.code);
            if (resp.success){
                clearInterval(check)
                console.log("登录成功")
                console.log(resp.data)

                let token = await getToken(resp.data.split("=")[1],data.sign)
                console.log(token)
                if (token.success){
                    Toast.show("登录成功\n该软件为免费软件，若你正在付费使用，请速度举报管理员")
                    this.setState({
                        link : ""
                    })
                }

            }
        },5000)
        this.setState({
            check: check
        })
        setTimeout(()=>{
            clearInterval(check)
        },1000*300)

        let element = document.createElement("a");
        element.href = "dtxuexi://appclient/page/study_feeds?url="+escape(data.url)
        element.click()

    }

    render() {
        return <div style={{width:"100%",height:"50%"}}>
            <h2 style={{margin:10,color:"red",display:this.isWechat()}}>当前环境为微信环境，请点击右上角在浏览器中打开</h2>
            <Button onClick={this.click} color={"primary"} style={{marginRight:10,marginTop:10,marginBottom:10}} block>生成链接</Button>
            <span>{this.state.link}</span>
            {/*<QRCode*/}
            {/*    id="qrCode"*/}
            {/*    value={this.state.img}*/}
            {/*    size={400} // 二维码的大小*/}
            {/*    fgColor="#000000" // 二维码的颜色*/}
            {/*    style={{ margin: 'auto' ,display:this.state.img === "你还未获取登录链接" ? "none" : "block"}}*/}
            {/*    imageSettings={{ // 二维码中间的logo图片*/}
            {/*        src: qr,*/}
            {/*        height: 100,*/}
            {/*        width: 100,*/}
            {/*        excavate: true, // 中间图片所在的位置是否镂空*/}
            {/*    }}*/}
            {/*/>*/}
        </div>;
    }
}

export default AddUser
