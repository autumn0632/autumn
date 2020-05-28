package main

import (
	"fmt"
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


func main() {
	var res int
	client, err := jsonrpc.Dial("tcp", URL)
	if err != nil {
		fmt.Println(err)
	}

	err = client.Call("Arith.Multiply", &Args{10, 20}, &res)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("rpc response:%d\n", res)

}
