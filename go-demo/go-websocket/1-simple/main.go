package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"myTest/demo_home/go-demo/go-websocket/model"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	conn       *websocket.Conn
	connection *model.Connection
	err        error
	data       []byte
	wsHandler  = func(w http.ResponseWriter, r *http.Request) {
		if conn, err = upgrader.Upgrade(w, r, nil); err != nil {
			return
		}
		//初始化连接
		if connection, err = model.InitConnection(conn); err != nil {
			return
		}
		go func() {
			for {
				if data, err = connection.ReadMessage(); err != nil {
					return
				}
				if err = connection.WriteMessage(data); err != nil {
					return
				}
			}
		}()
		//心跳检测
		go func() {
			for {
				if err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("heart beat...%s", time.Now().String()))); err != nil {
					return
				}
				time.Sleep(time.Second * 10)
			}
		}()
	}
)

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":7777", nil)
}
