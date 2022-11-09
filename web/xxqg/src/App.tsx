import React, {Component, useEffect} from 'react';
import './App.css';
import {NavBar,} from "antd-mobile";
import {checkToken} from "./utils/api";
import Login from './compents/pages/Login';
import {Route, Routes, useLocation, useNavigate} from 'react-router-dom';
import Home from './compents/Home';
import Log from "./compents/pages/Log";
import Config from "./compents/pages/Config";
import Help from "./compents/pages/Help";


// class Test extends React.Component<any, any> {
//     elements: any = []
//     items = [
//         {
//             "key": "login",
//             "text": "添加用户"
//         },
//         {
//             "key": "user_list",
//             "text": "用户管理"
//         }
//     ]
//
//     constructor(props: any) {
//         super(props);
//         this.state = {
//             popup_visible: false,
//             index: "login",
//             is_login: false,
//             // 用户等级，1是管理员，2是普通用户
//             level: 2
//         };
//     }
//
//     set_level = (level: number) => {
//         this.setState({
//             level: level
//         })
//
//         this.items.push(
//             {
//                 "key": "config",
//                 "text": "配置管理"
//             }
//         )
//         this.items.push(
//             {
//                 "key": "log",
//                 "text": "日志查看"
//             }
//         )
//         this.items.push(
//             {
//                 "key": "other",
//                 "text": "其他功能"
//             }
//         )
//
//
//         this.items.map((value, index, array) => {
//             if (value.key === "config" || value.key === "log") {
//                 this.elements.push(
//                     <ListItem disabled={this.state.level === 2} onClick={() => {
//                         this.setState({"index": value.key})
//                     }}>{value.text}</ListItem>
//                 );
//             } else {
//                 this.elements.push(
//                     <ListItem onClick={() => {
//                         this.setState({"index": value.key})
//                     }}>{value.text}</ListItem>
//                 );
//             }
//
//             return true;
//         })
//     }
//
//     set_login = () => {
//         this.setState({
//             is_login: true
//         })
//         this.check_token()
//         window.location.reload()
//     }
//
//     check_token = () => {
//         checkToken().then((t) => {
//             console.log(t)
//             if (!t) {
//                 console.log("未登录")
//             } else {
//                 if (t.data === 1) {
//                     console.log("管理员登录")
//                     this.set_level(1)
//                 } else {
//                     console.log("不是管理员登录")
//                     this.set_level(2)
//                 }
//                 this.setState({
//                     is_login: true
//                 })
//             }
//         })
//     }
//
//     componentDidMount() {
//         this.check_token()
//
//
//     }
//
//     render() {
//
//         let home = (
//             <>
//                 <NavBar style={{background: "#c0a8c0", margin: 10}} backArrow={false}
//                         left={<UnorderedListOutline fontSize={36} onClick={this.back}/>}>
//                     {"study_xxqg"}
//                 </NavBar>
//                 <Router data={this.state.index} level={this.state.level} set_level={this.set_level}/>
//                 <Popup
//                     bodyStyle={{width: '50vw'}}
//                     visible={this.state.popup_visible}
//                     position={"left"}
//                     onMaskClick={(() => {
//                         this.setState({popup_visible: false})
//                     })}>
//                     <h1 style={{textAlign: "center"}}>XXQG</h1>
//                     <List>
//                         {this.elements}
//                         <ListItem onClick={() => {
//                             window.localStorage.removeItem("xxqg_token")
//                             this.setState({"index": "help"})
//                         }}>帮助</ListItem>
//                         <ListItem onClick={() => {
//                             window.localStorage.removeItem("xxqg_token")
//                             this.setState({
//                                 is_login: false
//                             })
//                         }}>退出登录</ListItem>
//                     </List>
//                 </Popup>
//             </>
//         )
//         if (this.state.is_login) {
//             return home
//         } else {
//             return <Login parent={this}/>
//         }
//     }
//
//
//     back = () => {
//         this.setState({
//             popup_visible: true,
//         })
//
//     }
// }


function App(props: any, states: any) {

    let a = 2;
    let navigate = useNavigate();
    let location = useLocation();


    useEffect(() => {
        checkToken().then((t) => {
            console.log(t)
            if (!t) {
                console.log("未登录")
                navigate("/login")
            } else {
                if (t.data === 1) {
                    console.log("管理员登录")
                    sessionStorage.setItem("level", "1")
                } else {
                    console.log("不是管理员登录")
                    sessionStorage.setItem("level", "2")
                }
                navigate("/home/user_manager")
            }
        })
    }, [a])


    return <>
        <Routes>


            <Route path={"/login"} element={<Login navigate={navigate} location={location}/>}>
            </Route>

            <Route path={"/home/*"} element={<Home navigate={navigate} location={location}/>}>
            </Route>
            <Route path={"*"} element={<OtherPages navigate={navigate} location={location}/>}>

            </Route>

        </Routes>

    </>
}

class OtherPages extends Component<any, any> {
    render() {
        return <>
            <NavBar back='返回' style={{color: "blue", backgroundColor: "#bad7ba"}} onBack={() => {
                window.history.back()
            }}><h3>StudyXXQG</h3></NavBar>
            <Routes>
                <Route path={"/log"} element={<Log navigate={this.props.navigate} location={this.props.location}/>}>

                </Route>

                <Route path={"/config"}
                       element={<Config navigate={this.props.navigate} location={this.props.location}/>}/>

                <Route path={"/help"} element={<Help navigate={this.props.navigate} location={this.props.location}/>}/>


            </Routes>
        </>;
    }
}

export default App;
