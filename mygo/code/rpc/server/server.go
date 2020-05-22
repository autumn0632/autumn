// 示例代码
// 1. 定义传入参数和返回参数的数据结构
package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const (
	URL = "127.0.0.1:5001"
)

// 由客户端传入的参数类型
type Args struct {
	A, B int
}

// 返回给客户端的结果
type Quotient struct {
	Quo, Rem int
}

// 2. 定义服务对象。对象可以很简单，比如是int，或者interface{}，重要的是输出的方法
// 可以注册多个不同类型的对象，不允许注册一个类型的不同对象
type Arith int

// 3. 实现类型的两个方法
func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}
func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
func main() {
	// Arith类型只能注册一个对象-arith
	arith := new(Arith)
	rpc.Register(arith)

	tcpAddr, err := net.ResolveTCPAddr("tcp", URL)
	if err != nil {
		fmt.Println(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	fmt.Println("listen ", URL)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go jsonrpc.ServeConn(conn)
	}

}