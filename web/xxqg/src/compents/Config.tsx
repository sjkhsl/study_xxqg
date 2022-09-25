import React, {Component} from "react";
import {getConfig, saveConfig} from "../utils/api";
import MonacoEditor from 'react-monaco-editor';
import {Button, Dialog, Toast} from "antd-mobile";
class Config  extends Component<any, any>{
    private monaco: React.RefObject<MonacoEditor>;
    constructor(props: any) {
        super(props);
        this.monaco = React.createRef<MonacoEditor>()
        this.state = {
            config: ""
        };
    }

    componentDidMount() {
        getConfig().then((value)=>{
            this.setState({
                config:value.data
            })
        })

    }
    editorDidMount = (editor:any, monaco:any) => {
        console.log('editorDidMount', editor);
        editor.focus();
    }
    onChange = (newValue:any, e:any)=> {

    }

    onSave = ()=> {
        // @ts-ignore
        let data = this.monaco.current.editor?.getModel().getValue()
        saveConfig(data).then(resp => {
            if (resp.code === 200){
                Toast.show("保存成功")
            }else {
                Dialog.show({content:resp.err})
            }
        })
    }
    render() {
        const options = {
            selectOnLineNumbers: true
        };
        return <>
            <Button style={{margin:10,marginRight:30}} onClick={this.onSave} color={"primary"} block={true}>保存配置</Button><br/>
        <MonacoEditor
            ref={this.monaco}
            width={window.innerWidth}
            height={window.innerHeight}
            language="yaml"
            theme="vs"
            value={this.state.config}
            options={options}
            onChange={this.onChange}
            editorDidMount={this.editorDidMount}
            />
        </>
    }
}


export default Config
