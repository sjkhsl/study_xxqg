import axios, {AxiosInstance, AxiosRequestConfig} from "axios";
import {useNavigate} from "react-router-dom";

type TAxiosOption = {
    baseURL: string;
    timeout: number;
}

// const config = {
//     baseURL: '/',
//     timeout: 120000
// }

class Http {
    service: AxiosInstance;
    constructor(config:TAxiosOption) {
        this.service = axios.create(config);
        this.service.defaults.withCredentials = true
        this.service.interceptors.request.use(
            (value)=>{
               if (value.headers !== null){
                   // @ts-ignore
                   value.headers.Authorization = "Bearer "+localStorage.getItem("xxqg_token")
               }
               return value
        },()=>{
                console.log("请求异常")
        })
        this.service.interceptors.response.use((value)=>{
            console.log(value.data)
            return value
        },(error)=>{
            console.log(error.message)
            if (error.message === "Request failed with status code 401"){
                window.location.hash = "/login"
            }
        })
    }

    get<T>(url: string, params?: object, _object = {}): Promise<IResponseData<T>> {
        return this.service.get(url, { params, ..._object })
    }
    post<T>(url: string, data?: object, _object:AxiosRequestConfig = {}): Promise<IResponseData<T>> {
        return this.service.post(url, data, _object)
    }
    put<T>(url: string, params?: object, _object = {}): Promise<IResponseData<T>> {
        return this.service.put(url, params, _object)
    }
    delete<T>(url: string, params?: any, _object = {}): Promise<IResponseData<T>> {
        return this.service.delete(url, { params, ..._object })
    }
}

export default Http

export interface IResponseData<T> {
    success: boolean;
    message?:string;
    data:T;
    code: string;
    error?:string
}
