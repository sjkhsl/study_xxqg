import React, {Component} from "react";
import {login} from "../utils/api";
import {Button, Dialog, Form, Input} from "antd-mobile";

class Login extends Component<any, any>{
    constructor(props: any) {
        super(props);
        this.state = {
            img : "你还未获取登录链接"
        };
    }

    onFinish = (value:string)=>{
        login(JSON.stringify(value)).then(resp => {
            console.log(resp.message)
            Dialog.alert({content: resp.message,closeOnMaskClick:false})
            if (resp.success){
                window.localStorage.setItem("xxqg_token",resp.data)
                this.props.parent.set_login()
            }

        })
    }

    render() {
        return  <>
            <Form
                onFinish = {this.onFinish}
                footer={
                    <Button block type='submit' color='primary' size='large'>
                        登录
                    </Button>
                }
            >
                <Form.Header><h1>XXQG 登录页</h1></Form.Header>
                <Form.Item name='account' label='账号' rules={[{ required: true }]}>
                    <Input placeholder='请输入账号' />
                </Form.Item>
                <Form.Item name='password' label='密码' rules={[{ required: true }]}>
                    <Input placeholder='请输入密码'  type={"password"}/>
                </Form.Item>
            </Form>
        </>;
    }
}

export default Login
