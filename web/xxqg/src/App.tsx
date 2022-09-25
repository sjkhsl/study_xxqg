import React from 'react';
import './App.css';
import {List, NavBar, Popup,} from "antd-mobile";
import {UnorderedListOutline} from "antd-mobile-icons";
import {ListItem} from "antd-mobile/es/components/list/list-item";
import {checkToken} from "./utils/api";
import Login from './compents/Login';
import Router from './compents/Router';


class App extends React.Component<any, any> {
    constructor(props: any) {
        super(props);
        this.state = {
            popup_visible: false,
            index: "login",
            is_login: false,
            // 用户等级，1是管理员，2是普通用户
            level: 2
        };
    }

    set_level = (level: number) => {
        this.setState({
            level: level
        })

            this.items.push(
                {
                    "key":"config",
                    "text":"配置管理"
                }
            )
            this.items.push(
                {
                    "key":"log",
                    "text":"日志查看"
                }
            )
        this.items.push(
            {
                "key":"other",
                "text":"其他功能"
            }
        )



        this.items.map((value, index, array)=>{
            if (value.key === "config" || value.key === "log"){
                this.elements.push(

                    <ListItem disabled={this.state.level === 2}  onClick={() => {
                        this.setState({"index": value.key})
                    }}>{value.text}</ListItem>
                );
            }else {
                this.elements.push(

                    <ListItem  onClick={() => {
                        this.setState({"index": value.key})
                    }}>{value.text}</ListItem>
                );
            }

            return true;
        })
    }
    set_login = () => {
        this.setState({
            is_login: true
        })
        this.check_token()
        window.location.reload()
    }

    check_token = () => {
        checkToken().then((t) => {
            console.log(t)
            if (!t) {
                console.log("未登录")
            } else {
                if (t.data === 1) {
                    console.log("管理员登录")
                    this.set_level(1)
                } else {
                    console.log("不是管理员登录")
                    this.set_level(2)
                }
                this.setState({
                    is_login: true
                })
            }
        })
    }

    componentDidMount() {
        this.check_token()


    }

    elements:any = []
    items = [
        {
            "key":"login",
            "text":"添加用户"
        },
        {
            "key":"user_list",
            "text":"用户管理"
        }
    ]


    render() {

        let home = (
            <>
                <NavBar style={{background: "#c0a8c0", margin: 10}} backArrow={false}
                        left={<UnorderedListOutline fontSize={36} onClick={this.back}/>}>
                    {"study_xxqg"}
                </NavBar>
                <Router data={this.state.index} level={this.state.level} set_level={this.set_level}/>
                <Popup
                    bodyStyle={{width: '50vw'}}
                    visible={this.state.popup_visible}
                    position={"left"}
                    onMaskClick={(() => {
                        this.setState({popup_visible: false})
                    })}>
                    <h1 style={{textAlign: "center"}}>XXQG</h1>
                    <List>
                        {this.elements}
                        <ListItem onClick={() => {
                            window.localStorage.removeItem("xxqg_token")
                            this.setState({"index": "help"})
                        }}>帮助</ListItem>
                        <ListItem onClick={() => {
                            window.localStorage.removeItem("xxqg_token")
                            this.setState({
                                is_login: false
                            })
                        }}>退出登录</ListItem>
                    </List>
                </Popup>
            </>
        )
        if (this.state.is_login) {
            return home
        } else {
            return <Login parent={this}/>
        }
    }


    back = () => {
        this.setState({
            popup_visible: true,
        })

    }
}


export default App;
