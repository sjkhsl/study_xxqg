import React, {Component} from "react";
import {getAbout} from "../../utils/api";

class Help extends Component<any, any> {

    constructor(props: any) {
        super(props);
        this.state = {
            about: ""
        };
    }

    componentDidMount() {
        getAbout().then((value)=>{
            this.setState({
                about:value.data

            })
        })

    }
    render() {
        return <>
            <h1 style={{color:"red",margin:10}}>该软件为免费软件，若你目前正在付费使用，请速度举报管理员</h1><br/>
            <h2 style={{margin:10}}>项目地址：<a href="https://github.com/johlanse/study_xxqg">https://github.com/johlanse/study_xxqg</a></h2>
            <br/><h2 style={{margin:10}}>{this.state.about}</h2>
        </>
    }
}

export default Help
