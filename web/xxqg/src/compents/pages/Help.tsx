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
            <h2 style={{margin:10}}>项目地址：<a href="https://github.com/johlanse/study_xxqg">https://github.com/johlanse/study_xxqg</a></h2>
            <br/><h2 style={{margin:10}}>{this.state.about}</h2>
        </>
    }
}

export default Help
