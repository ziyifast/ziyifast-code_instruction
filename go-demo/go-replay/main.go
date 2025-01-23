package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/getInfo", getUserInfo)
	http.HandleFunc("/addUser", addUser)
	fmt.Println("run server2.....")
	http.ListenAndServe(":9999", nil)
}

func getUserInfo(rsp http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		rsp.Write([]byte(`{"code":-1, "msg":"method not allowed"}`))
		rsp.Header().Set("Content-Type", "application/json")
		rsp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	rsp.Write([]byte(`{"name":"zhangsan", "age":15}`))
	fmt.Printf("get user info\n")
	return
}

func addUser(rsp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		rsp.Write([]byte(`{"code":-1, "msg":"method not allowed"}`))
		rsp.Header().Set("Content-Type", "application/json")
		rsp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	username := req.FormValue("username")
	age := req.FormValue("age")
	rsp.Write([]byte(`{"message":"success", "code":0}`))
	rsp.Header().Set("Content-Type", "application/json")
	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Printf("add user[%s] [%d]\n", username, age)
	return
}
