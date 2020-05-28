# 一、为什么选择微服务

**传统服务**：从立项到开发上线，随着时间和需求的不断激增，会越来越复杂，变成一个大项目，如果前期项目架构没设计的不好，代码会越来越臃肿，难以维护，后期的每次产品迭代上线都会牵一发而动全身。

**微服务**：项目微服务化，松耦合模块间的关系，

比较重要的问题：

* 服务间数据传输的效率和安全性。
* 服务的动态扩充，也就是服务的注册和发现，服务集群化。
* 微服务功能的可订制化，因为并不是所有的功能都会很符合你的需求，难免需要根据自己的需要二次开发一些功能。



# 二、go-micro架构

**go-micro**：o是go语言下的一个的rpc微服务框架，功能很完善，且高度可定制化，几个重要的问题也解决的很好：

* 服务间传输格式为protobuf，效率非常高，也很安全。
* go-micro的服务注册和发现是多种多样的。
* 主要的功能都有相应的接口，只要实现相应的接口，就可以根据自己的需要订制插件。

 go-micro之所以可以高度订制和他的框架结构是分不开的，go-micro由8个关键的interface组成，每一个interface都可以根据自己的需求重新实现，这8个主要的inteface也构成了go-micro的框架结构。 

![](./../png/go-micro.png)

## 2.1 Transort

​    服务之间通信的接口。也就是服务发送和接收的最终实现方式，是由这些接口定制的。

http传输是go-micro默认的同步通信机制。在go-plugins里面实现的的方式有：grpc,nats,tcp,udp,rabbitmq,nats

```go 
type Socket interface {
    Recv(*Message) error
    Send(*Message) error
    Close() error
}

//Client 封装Socket，实现发送和接收通信的消息
type Client interface {
    Socket  
}

type Listener interface {
    Addr() string
    Close() error
    Accept(func(Socket)) error
}

type Transport interface {
    //客户端进行调用，连接服务端
    Dial(addr string, opts ...DialOption) (Client, error)
    //srever端进行调用，监听一个端口，等待客户端连接
    Listen(addr string, opts ...ListenOption) (Listener, error) 
    String() string
}
```

## 2.2 Codec

编码方式。默认的实现方式是protobuf，go-plugins里面还有实现其它方式

```go
type Codec interface {
    //解码过程
    ReadHeader(*Message, MessageType) error
    ReadBody(interface{}) error
    //编码过程
    Write(*Message, interface{}) error
    Close() error
    String() string
}

type Message struct {
    Id     uint64
    Type   MessageType
    Target string
    Method string
    Error  string
    Header map[string]string
}
```

## 2.3 Registry

服务的注册和发现，目前实现的consul,mdns, etcd,etcdv3,zookeeper,kubernetes.等。

简单来说，就是Service 进行Register，来进行注册，Client 使用watch方法进行监控，当有服务加入或者删除时这个方法会被触发，以提醒客户端更新Service信息。

```go
type Registry interface {
    Register(*Service, ...RegisterOption) error
    Deregister(*Service) error
    GetService(string) ([]*Service, error)
    ListServices() ([]*Service, error)
    Watch(...WatchOption) (Watcher, error)
    String() string
    Options() Options
}
```

## 2.4 Selector

   以Registry为基础，Selector 是客户端级别的负载均衡，当有客户端向服务发送请求时， selector根据不同的算法从Registery中的主机列表，得到可用的Service节点，进行通信。目前实现的有循环算法和随机算法，默认的是随机算法。

​     默认的是实现是本地缓存，当前实现的有blacklist,label,named等方式。

```go
type Selector interface {
    Init(opts ...Option) error
    Options() Options
    // Select returns a function which should return the next node
    Select(service string, opts ...SelectOption) (Next, error)
    // Mark sets the success/error against a node
    Mark(service string, node *registry.Node, err error)
    // Reset returns state back to zero for a service
    Reset(service string)
    // Close renders the selector unusable
    Close() error
    // Name of the selector
    String() string
}
```

## 2.5 Broker

Broker是发布和订阅的接口。很简单的一个例子，因为服务的节点是不固定的，如果有需要修改所有服务行为的需求，可以使服务订阅某个主题，当有信息发布时，所有的监听服务都会收到信息，根据你的需要做相应的行为。

Broker默认的实现方式是http方式，但是这种方式不要在生产环境用。go-plugins里有很多成熟的消息队列实现方式，有kafka、nsq、rabbitmq、redis，等等。

```go
type Broker interface {
    Options() Options
    Address() string
    Connect() error
    Disconnect() error
    Init(...Option) error
    Publish(string, *Message, ...PublishOption) error
    Subscribe(string, Handler, ...SubscribeOption) (Subscriber, error)
    String() string
}
```

## 2.6 Client

​    Client是请求服务的接口，他封装Transport和Codec进行rpc调用，也封装了Brocker进行信息的发布。

```go
type Client interface {
    Init(...Option) error
    Options() Options
    NewMessage(topic string, msg interface{}, opts ...MessageOption) Message
    NewRequest(service, method string, req interface{}, reqOpts ...RequestOption) Request
    Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error
    Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error)
    Publish(ctx context.Context, msg Message, opts ...PublishOption) error
    String() string
}
```

## 2.7 Server

监听等待rpc请求。监听broker的订阅消息，等待消息队列的推送等。

```go
type Server interface {
    Options() Options
    Init(...Option) error
    Handle(Handler) error
    NewHandler(interface{}, ...HandlerOption) Handler
    NewSubscriber(string, interface{}, ...SubscriberOption) Subscriber
    Subscribe(Subscriber) error
    Register() error
    Deregister() error
    Start() error
    Stop() error
    String() string
}
```

## 2.8 Service

​     Service是Client和Server的封装，他包含了一系列的方法使用初始值去初始化Service和Client，使我们可以很简单的创建一个rpc服务。

```go
type Service interface {
    Init(...Option)
    Options() Options
    Client() client.Client
    Server() server.Server
    Run() error
    String() string
}
```

