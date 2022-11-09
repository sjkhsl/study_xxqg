import React, {Component} from "react";
import {deleteUser, getExpiredUsers, getScore, getUsers, stopStudy, study} from "../../utils/api";
import {Button, Dialog, Divider, List, Modal, Toast} from "antd-mobile";
import {ListItem} from "antd-mobile/es/components/list/list-item";

class Users extends Component<any, any>{

    constructor(props: any) {
        super(props);
        this.state = {
            users:[],
            expired_users:[]
        };
    }

    componentDidMount() {
        getUsers().then(users =>{
            console.log(users)
            this.setState({
                users: users.data
            })
        })

        getExpiredUsers().then(users => {
            console.log(users)
            if (users.data !== null){
                this.setState({
                    expired_users: users.data
                })
            }

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

    delete_user = (uid:string,nick:string)=>{
        Dialog.confirm({content:"你确定要删除用户"+nick+"吗?"}).then((confirm) => {
            if (confirm){
                deleteUser(uid).then((data) => {
                    if (data.success){
                        getUsers().then(users =>{
                            console.log(users)
                            if (users.data != null){
                                this.setState({
                                    users: users.data
                                })
                            }

                        })
                    }else {
                        Dialog.show({content:data.error,closeOnMaskClick:true,closeOnAction:true})
                    }
                })
            }
        })
    }

    render() {
        let elements = []
        for (let i = 0; i < this.state.users.length; i++) {
            console.log(this.props.level)
            elements.push(
                <>
                    <ListItem key={this.state.users[i].uid} style={{border:"blue soild 1px"}}>
                        <h4>姓名：{this.state.users[i].nick}</h4>
                        <h4>UID: {this.state.users[i].uid}</h4>
                        <h4>登录时间：{this.format(this.state.users[i].login_time)}</h4>
                        <Button onClick={this.study.bind(this,this.state.users[i].uid,this.state.users[i].is_study)} color={this.checkStudyColor(this.state.users[i].is_study)} block={true}>
                            {this.checkStudy(this.state.users[i].is_study)}
                        </Button>
                        <br />
                        <Button onClick={this.getScore.bind(this,this.state.users[i].token,this.state.users[i].nick)} color={"success"} block={true}>积分查询</Button>
                        <br />
                        <Button  style={{display: this.props.level !== "1" ? "none" : "inline"}} onClick={this.delete_user.bind(this,this.state.users[i].uid,this.state.users[i].nick)} color={"danger"} block={true}>删除用户</Button>
                    </ListItem>
                    <Divider />
                </>
            )
        }
        if (this.state.users.length === 0){
            elements.push(<>
                <ListItem key={"none"} style={{border:"red soild 1px"}}><h2>未获取到登录用户</h2></ListItem>
            </>)
        }
        for (let i = 0; i < this.state.expired_users.length; i++) {
            console.log(this.state.expired_users[i].uid)
            elements.push(
                <>
                    <ListItem key={this.state.expired_users[i].uid} style={{border:"blue soild 1px",backgroundColor:"#cdced0"}}>
                        <h4>姓名：{this.state.expired_users[i].nick}<span style={{color:"red"}}>（已失效）</span></h4>
                        <h4>UID: {this.state.expired_users[i].uid}</h4>
                        <h4>登录时间：{this.format(this.state.expired_users[i].login_time)}</h4>
                        <Button  style={{display: this.props.level !== "1" ? "none" : "inline"}} onClick={this.delete_user.bind(this,this.state.expired_users[i].uid,this.state.expired_users[i].nick)} color={"danger"} block={true}>删除用户</Button>
                    </ListItem>
                    <Divider />
                </>
            )
        }



        return <List>{elements}</List>;
    }
}

export default Users
