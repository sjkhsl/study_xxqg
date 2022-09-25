import React, {Component} from "react";
import {Button, Toast} from "antd-mobile";
import {restart} from "../utils/api";

class Other extends Component<any, any>{

    onrestart = ()=>{
    restart().then(r => {

    });
    Toast.show("重启完成")
}
    render() {
        return <>
            <Button style={{margin:10,marginRight:30}} onClick={this.onrestart} color={"primary"} block={true}>重启程序</Button><br/>
        </>;
    }
}

export default Other
