# 流程概述

使用go-micro架构创建grpc通信的微服务过程：

1. 包导入

   ```go
   import (
   	hello "github.com/micro/examples/greeter/srv/proto/hello"
   	"github.com/micro/go-grpc"
   	"github.com/micro/go-micro"
   
   	"context"
   )
   ```

2. 新建一个服务

   ```go
   service := grpc.NewService(
   		micro.Name("go.micro.srv.greeter"),
   		micro.RegisterTTL(time.Second*30),
   		micro.RegisterInterval(time.Second*10),
   	)
   ```

3. 服务初始化

   ```
   service.Init()
   ```

4. 注册处理器

   ```
   hello.RegisterSayHandler(service.Server(), new(Say))
   ```

5. 启动服务

   ```
   service.Run()
   ```

# 流程分析

## 创建服务

```go
service := grpc.NewService(
		micro.Name("go.micro.srv.greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
```

### 数据结构

`NewService`返回接口类型

```go
type Service interface {
	Init(...Option)         //Init参数为函数类型切片
	Options() Options       // 
	Client() client.Client
	Server() server.Server
	Run() error
	String() string
}

//1. 
type Option func(*Options) // 函数类型

//2. 
type Options struct {
	Broker    broker.Broker // Broker is an interface used for asynchronous messaging.
	Cmd       cmd.Cmd//接口类型
	Client    client.Client //接口类型
	Server    server.Server//接口类型
	Registry  registry.Registry  //服务注册方式，默认为mdns //接口类型
	Transport transport.Transport // 服务间通信方式 //接口类型

	// Before and After funcs
	BeforeStart []func() error //函数类型
	BeforeStop  []func() error //函数类型
	AfterStart  []func() error //函数类型
	AfterStop   []func() error //函数类型

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}
//3.
type Client interface {
	Init(...Option) error
	Options() Options
	NewMessage(topic string, msg interface{}, opts ...MessageOption) Message
	NewRequest(service, endpoint string, req interface{}, reqOpts ...RequestOption) Request
	Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error
	Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error)
	Publish(ctx context.Context, msg Message, opts ...PublishOption) error
	String() string
}
//4. 
type Server interface {
	Options() Options
	Init(...Option) error
	Handle(Handler) error
	NewHandler(interface{}, ...HandlerOption) Handler
	NewSubscriber(string, interface{}, ...SubscriberOption) Subscriber
	Subscribe(Subscriber) error
	Start() error
	Stop() error
	String() string
}
```

### 方法

```go
// 1. NewService returns a grpc service compatible with go-micro.Service
func NewService(opts ...micro.Option) micro.Service {
    //1. 客户端初始化
    c := client.NewClient()
    //2. 服务端初始化
    s := server.NewServer()
    //3. 代理初始化（异步通信）
    b := broker.NewBroker()
    
    //4. 
    options := []micro.Option{
		micro.Client(c),
		micro.Server(s),
		micro.Broker(b),
	}
    
    //5. 
    options = append(options, opts...)
    
    //6.
    // generate and return a service
	return micro.NewService(options...)
    
}

func newService(opts ...Option) Service {
	options := newOptions(opts...) //对各种option进行处理，option是函数类型，也就是执行各种函数

	options.Client = &clientWrapper{
		options.Client,
		metadata.Metadata{
			HeaderPrefix + "From-Service": options.Server.Options().Name,
		},
	}

	return &service{
		opts: options,
	}
}
```