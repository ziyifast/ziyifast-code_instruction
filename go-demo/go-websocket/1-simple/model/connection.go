package model

import (
	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/x/errors"
	"sync"
)

/*
	整体思路：
	1. 维护连接的读写channel
	2. 分别启两个协程for循环，一个用于读，一个用于写
	//中间多了一层Channel，保证了线程安全
	readLoop -> inChannel -> c.ReadMessage拿到data -> c.WriteMessage(data) -> outChannel -> writeLoop
*/

type Connection struct {
	Conn *websocket.Conn
	// 读消息队列
	inChannel chan []byte
	//写消息队列
	outChannel chan []byte
	// 监听Channel是否关闭
	closeChan chan byte
	// 标识
	isClosed bool
	lock     sync.Mutex
}

// InitConnection 初始化封装的conn
func InitConnection(conn *websocket.Conn) (c *Connection, err error) {
	c = &Connection{
		Conn:       conn,
		inChannel:  make(chan []byte, 1000),
		outChannel: make(chan []byte, 1000),
		closeChan:  make(chan byte),
		isClosed:   false,
	}
	//启动协程读取消息
	go c.readLoop()
	go c.writeLoop()
	return c, nil
}

// ReadMessage 读取消息，从inChannel中读取数据（channel保证线程安全，阻塞读取）
func (c *Connection) ReadMessage() (data []byte, err error) {
	//从inChannel读取数据
	for {
		select {
		case data = <-c.inChannel:
			return data, nil
		//监听连接关闭信号，避免一直阻塞读取数据
		case <-c.closeChan:
			return nil, errors.New("conn is closed")
		}
	}
}

// WriteMessage 写消息，将数据写入outChannel（channel保证线程安全,等待write loop从outChannel中获取数据写回连接）
func (c *Connection) WriteMessage(data []byte) (err error) {
	for {
		select {
		case c.outChannel <- data:
			return nil
		case <-c.closeChan:
			return errors.New("conn is closed")
		}
	}
}

// 从连接中不断读取数据写入inChannel
func (c *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = c.Conn.ReadMessage(); err != nil {
			//读取数据失败，关闭连接
			c.Close()
			return
		}
		select {
		//读取到数据写到inChannel
		case c.inChannel <- data:
		case <-c.closeChan:
			c.Close()
		}
	}
}

// 从outChannel中不断读取数据并发送数据写回对端
func (c *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-c.outChannel:
			if err = c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				c.Close()
				return
			}
		case <-c.closeChan:
			c.Close()
		}
	}
}

func (c *Connection) Close() {
	c.Conn.Close()
	c.lock.Lock()
	if !c.isClosed {
		close(c.closeChan)
		c.isClosed = true
	}
	c.lock.Unlock()
}
