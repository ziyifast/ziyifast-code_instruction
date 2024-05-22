import axios from "axios";

const http = axios.create({
    baseURL: 'http://43.139.239.29/',
    timeout: 5000,
    headers: {
        'Content-Type': 'application/json;charset=UTF-8'
    }
})
//拦截器
http.interceptors.request.use(function (config) {
    //拦截请求
    console.log("拦截到了请求")
    return config
})

http.interceptors.response.use(function (config) {
    //拦截响应
    console.log("拦截到了响应")
    return config
})
export default http;


