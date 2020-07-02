# 概述

Go Micro是可插拔的微服务开发框架。它提供分布式系统开发的核心库，包含RPC与事件驱动的通信机制。

go-micro里面包含了开发所用到的常用功能模块，主要模块之间关系如下：

![](./png/go-micro.png)



| 模块      | 说明                                                         |
| --------- | ------------------------------------------------------------ |
| services  | 微服务，提供了对微服务功能开发的封装，通过它可以快速创建一个微服务 |
| client    | RPC客户端，提供了诸如 服务发现/负载均衡/RPC代理和调用。以及失败的重试/超时/上下文等功能。 |
| server    | RPC服务端，提供了如何实现RPC请求的方法，功能逻辑主要使用这个实现 |
| codec     | 数据编码模块，提供将程序调用数据转换成RPC调用数据的功能，目前支持：json/protobuf |
| broker    | pub/sub模块，提供事件发布/订阅功能，目前支持：nats，rabbitmq，http |
| transport | 数据传输模块，通过抽象实现对传输协议的无缝替换。目前支持：http，rabbitmq，nats |
| registry  | 服务发现模块，提供集群的服务发现功能，目前支持：consul，etcd，memory，kubernetes |
| selector  | 负载均衡模块，当client发出请求是，负责在多个满足条件的服务器列表中决定使用哪个服务器，目前支持循环，哈希，黑名单 |



# 特性

Go Micro把分布式系统的各种细节抽象出来。下面是它的主要特性。

* **服务发现（Service Discovery）** - 自动服务注册与名称解析。服务发现是微服务开发中的核心。当服务A要与服务B协作时，它得知道B在哪里。默认的服务发现系统是Consul，而multicast DNS (mdns，组播)机制作为本地解决方案，或者零依赖的P2P网络中的SWIM协议（gossip）。

* **负载均衡（Load Balancing）** - 在服务发现之上构建了负载均衡机制。当我们得到一个服务的任意多个的实例节点时，我们要一个机制去决定要路由到哪一个节点。我们使用随机处理过的哈希负载均衡机制来保证对服务请求颁发的均匀分布，并且在发生问题时进行重试。
* **消息编码（Message Encoding）** - 支持基于内容类型（content-type）动态编码消息。客户端和服务端会一起使用content-type的格式来对Go进行无缝编/解码。各种各样的消息被编码会发送到不同的客户端，客户端服服务端默认会处理这些消息。content-type默认包含proto-rpc和json-rpc。
* **Request/Response** - RPC通信基于支持双向流的请求/响应方式，我们提供有抽象的同步通信机制。请求发送到服务时，会自动解析、负载均衡、拨号、转成字节流，默认的传输协议是http/1.1，而tls下使用http2协议。
* **异步消息（Async Messaging）** - 发布订阅（PubSub）头等功能内置在异步通信与事件驱动架构中。事件通知在微服务开发中处于核心位置。默认的消息传送使用点到点http/1.1，激活tls时则使用http2。
* **可插拔接口（Pluggable Interfaces）** - Go Micro为每个分布式系统抽象出接口。因此，Go Micro的接口都是可插拔的，允许其在运行时不可知的情况下仍可支持。所以只要实现接口，可以在内部使用任何的技术。

# 安装protof

代码生成依赖Protobuf，

python编译选型：

```
python -m grpc_tools.protoc --proto_path=./  --python_out=./gen_py --grpc_python_out=./gen_py ./data.proto
```

go 编译选项

```
 protoc -I ../routeguide --go_out=plugins=grpc:../routeguide ../routeguide/route_guide.proto
```

# 服务发现

服务发现用于解析服务名和地址

