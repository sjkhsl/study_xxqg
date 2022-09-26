import React, {Component} from "react";
import Users from "./pages/User";
import Help from "./pages/Help";
import Log from "./pages/Log";
import Config from "./pages/Config";
import Other from "./Other";
import AddUser from "./pages/AddUser";

class Router extends Component<any, any> {

    constructor(props: any) {
        super(props);
        this.state = {
            img: "你还未获取登录链接"
        };
    }


    render() {
        let login = <AddUser/>;

        let userList = <Users data={"12"} level={this.props.level}/>;
        let config = <Config/>
        let help = <Help/>
        let log = <Log/>
        if (this.props.data === "login") {
            return login;
        } else if (this.props.data === "user_list") {
            return userList;
        } else if (this.props.data === "help") {
            return help;
        } else if (this.props.data === "log") {
            return log;
        } else if (this.props.data === "other") {
            return <Other/>
        } else {
            return config;
        }
    }


}

export default Router
