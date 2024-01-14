package main

import (
	"github.com/aobco/log"
	"net/http"
	"time"
)

/*
	后端解决跨域问题
*/

func main() {
	mux := http.NewServeMux()
	mux.Handle("/cros/smoke", interceptor(http.HandlerFunc(smoke)))
	http.ListenAndServe(":8080", mux)
}

func smoke(w http.ResponseWriter, r *http.Request) {
	now := time.Now().String()
	_, err := w.Write([]byte(now))
	if err != nil {
		log.Errorf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

//拦截器
func interceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//resolve the cross origin[解决预请求]
		//w3c规范要求，当浏览器判定请求为复杂请求时，会在真实携带数据发送请求前，多一个预处理请求：
		//1. 请求方法不是get head post
		//2. post 的content-type不是application/x-www-form-urlencode,multipart/form-data,text/plain [也就是把content-type设置成"application/json"]
		//3. 请求设置了自定义的header字段: 比如业务需求，传一个字段，方面后端获取，不需要每个接口都传
		if r.Method == "OPTIONS" {
			//handle the preflight request
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept,yi-token")
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Accept,yi-token")
		next.ServeHTTP(w, r)
	})
}
