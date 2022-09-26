import React, {Component} from "react";
import {Button, Toast} from "antd-mobile";
import QrCode from "qrcode.react";
import {checkQrCode, getLink, getToken} from "../../utils/api";

class AddUser extends Component<any, any>{


    constructor(props: any) {
        super(props);
        this.state = {
            img : "你还未获取登录链接"
        };
    }

    componentWillUnmount() {
        if (this.state.check !== undefined){
            clearInterval(this.state.check)
        }

    }

    isWechat = ()=> {
        if (/MicroMessenger/i.test(window.navigator.userAgent)){
            return "inline"
        }else {
            return "none"
        }
    }

    click = async () => {
        let data = await getLink()

        this.setState({
            img: data.url
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
                    Toast.show("登录成功")
                    this.setState({
                        img : "你还未获取登录链接"
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
            <QrCode style={{marginLeft:10,display:this.state.img === "你还未获取登录链接" ? "none":"block"}} fgColor={"#000000"} size={200} value={this.state.img} />
        </div>;
    }
}

export default AddUser
