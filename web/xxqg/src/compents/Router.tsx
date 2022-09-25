import React, {Component} from "react";
import {Button, Toast} from "antd-mobile";
import QrCode from "qrcode.react";
import Users from "./User";
import Help from "./Help";
import Log from "./Log";
import {checkQrCode, getLink, getToken} from "../utils/api";
import Config from "./Config";
import Other from "./Other";

class Router extends Component<any, any>{

    constructor(props: any) {
        super(props);
        this.state = {
            img : "你还未获取登录链接"
        };
    }

    isWechat = ()=> {
        if (/MicroMessenger/i.test(window.navigator.userAgent)){
            return "inline"
        }else {
            return "none"
        }
    }

    render() {
        let login =  <>
            <h2 style={{margin:10,color:"red",display:this.isWechat()}}>当前环境为微信环境，请点击右上角在浏览器中打开</h2>
            <Button onClick={this.click} color={"primary"} style={{margin:10,marginRight:10}} block>生成链接</Button>
            <QrCode style={{margin:10}} fgColor={"#000000"} size={200} value={this.state.img} />
        </>;

        let userList = <Users data={"12"} level={this.props.level}/>;
        let config = <Config />
        let help = <Help />
        let log = <Log />
        if (this.props.data === "login"){
            return login;
        }else if (this.props.data === "user_list"){
            return userList;
        }else if (this.props.data === "help"){
            return help;
        } else if (this.props.data === "log"){
            return log;
        }else if (this.props.data === "other") {
            return <Other />
        }
        else {
            return config;
        }
    }

    componentWillUnmount() {
        if (this.state.check !== undefined){
            clearInterval(this.state.check)
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
}

export default Router
