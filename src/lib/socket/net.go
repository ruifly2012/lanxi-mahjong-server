/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2016-01-23 10:24
 * Filename      : net.go
 * Description   : 负通道读写数据，断线处理
 * *******************************************************/
package socket

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type Packet struct {
	proto   uint32
	count   uint32
	content []byte
}

func (this *Packet) SetProto(proto uint32) {
	this.proto = proto
}

func (this *Packet) SetContent(content []byte) {
	this.content = content
}

func (this *Packet) GetProto() uint32 {
	return this.proto
}
func (this *Packet) GetContent() []byte {
	return this.content
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 9 * time.Second
	maxMessageSize = 1024 * 1024 * 30
	//连接建立后5秒内没有收到登陆请求，断开socket
	waitForLogin = time.Second * 5
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512 * 30,
	WriteBufferSize: 512 * 30,
}

func wSHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln(err)
		}
	}()
	if r.Method != "GET" {
		return
	}
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	c := newConnection(socket)
	go c.Reader(c.ReadChan)
	go c.LoginTimeout()
	go c.WritePump()
	c.ReadPump()
}
