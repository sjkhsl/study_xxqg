import React, {Component} from "react";
import {getConfig, saveConfig} from "../../utils/api";
import MonacoEditor from 'react-monaco-editor';
import {Button, Dialog, Toast} from "antd-mobile";
import { setDiagnosticsOptions } from 'monaco-yaml';
import {editor, Uri} from "monaco-editor";
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
        try {
            this.initYaml()
        }catch (e) {

        }
    }
    editorDidMount = (editor:any, monaco:any) => {
        editor.focus();
    }
    onChange = (newValue:any, e:any)=> {
        console.log(this.monaco.current)
    }

    onSave = ()=> {
        // @ts-ignore
        let data = this.monaco.current.editor?.getModel().getValue()
        saveConfig(data).then(resp => {
            if (resp.code === 200){
                Toast.show("保存成功")
            }else {
                console.log(resp)
                Dialog.show({content:"配置提交失败"+resp.error,closeOnMaskClick:true,closeOnAction:true})
            }
        })
    }
    render() {
        const options = {
            selectOnLineNumbers: true,
            minimap: {
                enabled: false
            }
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

    initYaml = ()=>{
        const modelUri = Uri.parse('a://b/foo.yaml');

        setDiagnosticsOptions({
            enableSchemaRequest: true,
            hover: true,
            completion: true,
            validate: true,
            format: true,
            schemas: [
                {
                    // Id of the first schema
                    uri: 'http://myserver/foo-schema.json',
                    // Associate with our model
                    fileMatch: [String(modelUri)],
                    schema: {
                        type: 'object',
                        properties: {
                            p1: {
                                enum: ['v1', 'v2'],
                            },
                            p2: {
                                // Reference the second schema
                                $ref: 'http://myserver/bar-schema.json',
                            },
                        },
                    },
                },
                {
                    // Id of the first schema
                    uri: 'http://myserver/bar-schema.json',
                    fileMatch:[],
                    schema: {
                        type: 'object',
                        properties: {
                            q1: {
                                enum: ['x1', 'x2'],
                            },
                        },
                    },
                },
            ],
        });
        editor.create(document.createElement('editor'), {
            // Monaco-yaml features should just work if the editor language is set to 'yaml'.
            language: 'yaml',
            model: editor.createModel('p1: \n', 'yaml', modelUri),
        });
    }
}


export default Config
