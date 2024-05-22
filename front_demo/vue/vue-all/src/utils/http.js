import axios from "axios";

//封装请求
const http = axios.create({
    baseURL: '/api',
    timeout: 5000,
    headers: {
        'Content-Type': 'application/json;charset=UTF-8'
    },
    withCredentials:true, //允许跨域
});


function get(url, params = {}) {
    return new Promise((resolve, reject) => {
        http.get(url, {
            params: params
        }).then(res => {
            resolve(res.data);
        }).catch(err => {
            reject(err.data)
        })
    })
}

function post(url, data = {}) {
    return new Promise((resolve, reject) => {
        http.post(url, data)
            .then(response => {
                resolve(response.data);
            })
            .catch(err => {
                reject(err)
            })
    })
}

export default http;