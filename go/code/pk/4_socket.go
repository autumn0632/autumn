package pk

import (
	"fmt"
	"log"
	"net"
	"os"
)

type UnixSocket struct {
	filename string
	bufsize int
	Handler func(string) string
}

func NewUnixSocket(filename string, size ...int) *UnixSocket {
	size1 := 10480
	if size != nil {
		size1 = size[0]
	}
	us := UnixSocket{filename: filename, bufsize: size1}
	return &us
}

func (this *UnixSocket) createServer() {
	os.Remove(this.filename)
	addr, err := net.ResolveUnixAddr("unix", this.filename)
	if err != nil {
		panic("Cannot resolve unix addr: " + err.Error())
	}
	listener, err := net.ListenUnix("unix", addr)
	defer listener.Close()
	if err != nil {
		panic("Cannot listen to unix domain socket: " + err.Error())
	}
	fmt.Println("Listening on", listener.Addr())
	for {
		c, err := listener.Accept()
		if err != nil {
			panic("Accept: " + err.Error())
		}
		go this.HandleServerConn(c)
	}
}

func (this *UnixSocket) HandleServerConn(c net.Conn) {
	defer c.Close()
	buf := make([]byte, this.bufsize)
	nr, err := c.Read(buf)
	if err != nil {
		log.Printf("%s", err)
	}
	log.Printf("recv data: %s", string(buf[:nr]))
}
func (this *UnixSocket) StartServer() {
	this.createServer()
}

func (this *UnixSocket) ClientSentContext(context string) {
	addr, err := net.ResolveUnixAddr("unix", this.filename)
	if err != nil {
		panic("Cannot resolve unix addr: " + err.Error())
	}
	//拔号
	c, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		panic("DialUnix failed.")
	}
	//发数据
	_, err = c.Write([]byte(context))
	if err != nil {
		panic("Writes failed.")
	}
}