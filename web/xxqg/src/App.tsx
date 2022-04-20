import React, {Component} from 'react';
import './App.css';
import {Button, Divider, List, Modal, NavBar, Popup, TextArea, Toast,} from "antd-mobile";
import {UnorderedListOutline} from "antd-mobile-icons";
import {ListItem} from "antd-mobile/es/components/list/list-item";
import {checkQrCode, getLog, getScore, getToken, getUsers, login, stopStudy, study} from "./utils/api";
import QrCode from 'qrcode.react';


class App extends React.Component<any, any> {
  constructor(props: any) {
    super(props);
    this.state = {
      popup_visible: false,
      index: "login"
    };
  }


  render() {


    return <><>
      <NavBar style={{background: "#c0a8c0", margin: 10}} backArrow={false}
              left={<UnorderedListOutline fontSize={36} onClick={this.back}/>}>
        {"study_xxqg"}
      </NavBar>
      <Router data={this.state.index}/>
      <Popup
          bodyStyle={{width: '50vw'}}
          visible={this.state.popup_visible}
          position={"left"}
          onMaskClick={(() => {
            this.setState({popup_visible: false})
          })}>
        <h1 style={{textAlign:"center"}}>XXQG</h1>
        <List>
          <ListItem onClick={()=>{this.setState({"index":"login"})}}>添加用户</ListItem>
          <ListItem onClick={()=>{this.setState({"index":"user_list"})}}>用户管理</ListItem>
          <ListItem onClick={()=>{this.setState({"index":"config"})}}>配置管理</ListItem>
          <ListItem onClick={()=>{this.setState({"index":"log"})}}>日志查看</ListItem>
          <ListItem onClick={()=>{this.setState({"index":"help"})}}>帮助</ListItem>
        </List>
      </Popup>
    </>
    </>;
  }



  back = () => {
    this.setState({
      popup_visible: true,
    })

  }
}


class Router extends Component<any, any>{

  constructor(props: any) {
    super(props);
    this.state = {
      img : "你还未获取登录链接"
    };
  }

  render() {
    let login =  <>
      <Button onClick={this.click} color={"primary"} style={{margin:10,marginRight:10}} block>生成链接</Button>
      <QrCode style={{margin:10}} fgColor={"#000000"} size={200} value={this.state.img} />
    </>;

    let userList = <Users data={"12"}/>;
    let config = <h1>配置管理</h1>
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
    } else {
      return config;
    }
  }

  click = async () => {
    let data = await login()

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

    setTimeout(()=>{
      clearInterval(check)
    },1000*300)

    let element = document.createElement("a");
    element.href = "dtxuexi://appclient/page/study_feeds?url="+escape(data.url)
    element.click()
  }
}

class Log extends Component<any, any>{

  constructor(props:any) {
    super(props);
    this.state = {
      data : ""
    }
  }

  timer: any

  componentDidMount() {
    getLog().then(data=>{
      this.setState({
        data:data
      })
    })
    this.timer = setInterval(()=>{
      getLog().then((data)=>{
        console.log(data)
        this.setState({
          data:data
        })
      })
    },30000)
  }

  componentWillUnmount() {
    clearInterval(this.timer)
  }

  render() {
    console.log(this.state.data)
    return <>
    <TextArea autoSize disabled={true} value={this.state.data}/>
    </>
  }
}

class Help extends Component<any, any> {
  render() {
    return <>
        <h2>项目地址：<a href="https://github.com/johlanse/study_xxqg">https://github.com/johlanse/study_xxqg</a></h2>
    </>
  }
}

class Users extends Component<any, any>{

  constructor(props: any) {
    super(props);
    this.state = {
      users:[]
    };
  }

  componentDidMount() {
    getUsers().then(users =>{
      console.log(users)
      this.setState({
        users: users.data
      })
    })

  }

  format = (value:any)=> {
    const date = new Date(value*1000);
    let y = date.getFullYear(),
        m = date.getMonth() + 1,
        d = date.getDate(),
        h = date.getHours(),
        i = date.getMinutes(),
        s = date.getSeconds();
    if (m < 10) { m = parseInt('0') + m; }
    if (d < 10) { d = parseInt('0') + d; }
    if (h < 10) { h = parseInt('0') + h; }
    if (i < 10) { i = parseInt('0') + i; }
    if (s < 10) { s = parseInt('0') + s; }
    return y + '-' + m + '-' + d + ' ' + h + ':' + i + ':' + s;
  }

  getScore = (token:string,nick:string)=>{
      getScore(token).then((data)=>{
        console.log(data)
        Modal.alert({
            title: nick,
            content: data.data,
          closeOnMaskClick: true,
        })
      })
  }

  checkStudy = (is_study:boolean)=>{
    if (is_study){
      return "停止学习"
    }else {
      return "开始学习"
    }
  }

  checkStudyColor = (is_study:boolean)=>{
    if (is_study){
      return "danger"
    }else {
      return "primary"
    }
  }

  study = (uid:string,is_study:boolean) =>{
    if (!is_study){
        study(uid).then(()=>{
          Toast.show("开始学习成功")
          getUsers().then(users =>{
            console.log(users)
            this.setState({
              users: users.data
            })
          })
        })
    }else {
      stopStudy(uid).then(()=>{
        Toast.show("已停止学习")
        getUsers().then(users =>{
          console.log(users)
          this.setState({
            users: users.data
          })
        })
      })
    }
  }


  render() {
    let elements = []
    for (let i = 0; i < this.state.users.length; i++) {
      elements.push(
          <>
         <ListItem key={this.state.users[i].uid} style={{border:"blue soild 1px"}}>
           <h3>姓名：{this.state.users[i].nick}</h3>
           <h3>UID: {this.state.users[i].uid}</h3>
           <h3>登录时间：{this.format(this.state.users[i].login_time)}</h3>
           <Button onClick={this.study.bind(this,this.state.users[i].uid,this.state.users[i].is_study)} color={this.checkStudyColor(this.state.users[i].is_study)} block={true}>
             {this.checkStudy(this.state.users[i].is_study)}
           </Button>
           <br />
           <Button onClick={this.getScore.bind(this,this.state.users[i].token,this.state.users[i].nick)} color={"success"} block={true}>积分查询</Button>
         </ListItem>
          <Divider />
          </>
      )
    }

    return <List>{elements}</List>;
  }
}

export default App;
