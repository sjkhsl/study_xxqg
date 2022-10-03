import React, {Component} from "react";
import {Route, Routes} from "react-router-dom";
import AddUser from "./pages/AddUser";
import Users from "./pages/User";
import Other from "./Other";
import {NavBar, TabBar} from "antd-mobile";
import {MoreOutline, UserAddOutline, UserOutline} from "antd-mobile-icons";

let tableItems = [
    {
        key: '/add_user',
        title: '添加',
        icon: <UserAddOutline />,
    },
    {
        key: '/user_manager',
        title: '用户',
        icon: <UserOutline />
    },
    {
        key: '/other',
        title: '其他',
        icon: <MoreOutline />
    },
]

class Home extends Component<any, any>{
    render() {
        return <>
            <NavBar backArrow={false} style={{color:"blue",backgroundColor:"#bad7ba"}} onBack={()=>{
                window.history.back()
            }} ><h3>StudyXXQG</h3></NavBar>
            <Routes>

                <Route path={"add_user"} element={<AddUser navigate={this.props.navigate} location={this.props.location}/>}>
                </Route>

                <Route path={"user_manager"} element={<Users level={sessionStorage.getItem("level")} navigate={this.props.navigate} location={this.props.location}/>}>

                </Route>
                <Route path={"other"} element={<Other navigate={this.props.navigate} location={this.props.location}/>}>

                </Route>



            </Routes>
            <div style={{position:"fixed","height":"60px",width:"100%",bottom:0,zIndex:9,color: "#f0f", backgroundColor: "#5f6d6e"}}>

                <TabBar activeKey={this.props.location.pathname} onChange={value => this.props.navigate("/home"+value)}>
                    {tableItems.map(item => (
                        <TabBar.Item key={item.key} icon={item.icon} title={item.title} />
                    ))}
                </TabBar>
            </div>

        </>;
    }
}




export default Home
