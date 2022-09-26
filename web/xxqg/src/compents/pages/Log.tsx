import React, {Component} from "react";
import {getLog} from "../../utils/api";
import {TextArea} from "antd-mobile";

class Log extends Component<any, any>{

    constructor(props:any) {
        super(props);
        this.state = {
            data : ""
        }
    }

    reverse = ( str:string ):string=>{
        return str.split("\n").reverse().join("\n").trim()
    };

    timer: any

    componentDidMount() {
        getLog().then(data=>{
            this.setState({
                data:this.reverse(data)
            })
        })
        this.timer = setInterval(()=>{
            getLog().then((data:string)=>{
                this.setState({
                    data:this.reverse(data)
                })
            })
        },30000)
    }

    componentWillUnmount() {
        clearInterval(this.timer)
    }

    render() {
        return <>
            <TextArea style={{margin:10}} autoSize disabled={true} value={this.state.data}/>
        </>
    }
}

export  default  Log
