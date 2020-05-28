# 第一部分 go与web应用

```go 
package main
import （
	"fmt"
	"net/http"
）

func handler(w http.ResponseWriter, request http.Request) {
    fmt.Fprintf(w, "Hello World, %s", request.URL.Path[:1])
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080",nil)   
}
```



## web应用工作原理

web服务通常会被其它软件调用，而web应用则是为人类提供服务。

一个程序满足以下两个条件，就可以看作是一个**web应用：**

- 这个程序必须要向发送请求的客户端返回HTML，而客户端则会向用户展示渲染后的HTML。
- 这个程序在向客户端传送数据时，必须使用http协议

## http 学习

> **http是一种无状态、由文本构成的请求-响应（request-response）协议，这种协议使用的是客户端-服务器（client-server）计算模型**

http是一种无状态协议，它唯一知道的就是客户端会向服务器发送请求，而服务器则会向客户端返回响应，并且后续发生的请求对之前的请求一无所知。

http是以纯文本方式而不是二进制方式发送和接收协议数据的。

### http请求

**请求方法**

定义了发送请求的客户端想要执行的动作：

- GET：命令服务器返回指定的资源
- HEAD：与GET方法类似，不同之处在于不要求返回报文的主体。
- POST：命令服务器将报文主体中的数据传递给URL指定的资源
- PUT：命令服务器将报文主体中的数据设置为URL指定的资源。如果URL指定的位置上已经有数据存在，那么使用报文主体中的数据去代替已有的数据。
- DELETE：

**请求首部**

http请求的首部记录了与请求本身以及客户端有关的信息。最后已回车（CR）和换行（LF）结尾。

### http响应

### URI

> 统一资源标识符，Uniform Response Identitier，URI
>
>统一资源名称，Uniform Response Name，URN
>
>统一资源定位符，Uniform Response location，URL
>
>URI是一个涵盖性术语，包含了URN和URL

URI的一般格式为：

**<方案名称>:<分层部分>\[ ? <查询参数> ][ # <片段>]**

url里面是不能包含空格的。此外，问好（？）和井号（#）等符号在url中具有特殊含义，所以这些符号是不能用于其它限制的。为了避开这些限制，使用url编码对这些特殊符号进行转换。

## http/2 简介

这一版本对性能非常关注。与http 1.x 不同的是，http/2 是一种二进制协议。使用多路复用的方式，使多个请求和响应可以在同一时间内使用同一个链接。

## web应用的各个组成

任何web应用都包含处理器和模板引擎，处理器负责接收http请求并处理它们。模板引擎负责生成HTML，这些HTML会作为http响应的其中一部分被回传至客户端。

## ChitChat论坛

**使用cookie进行访问控制**

登录流程：

1. 查询用户名是否存在数据库中，如不存在，重定向到登录页面
2. 若存在，进行密码匹配，若匹配失败，重定向到登录页面
3. 密码匹配成功，新建一个session结构，其中 Uuid 字段存储一个随机生成的唯一ID，是实现会话机制的核心。服务器会通过cookie把这个ID存储到前端中。
4. 创建cookie结构，并将cookie返回到前端

**使用模板生成html响应**

首先，函数把每个需要用到的模板文件放到切片里，HTML模板文件都包含了特定的嵌入命令，这些命令被称为**动作（action）**，会被`{{`符号和`}}`符号包围。

**整体流程**

> 1. 客户端向服务端发送请求
> 2. 多路复用器接收到请求，并将其重定向到正确的处理器
> 3. 处理器对请求进行处理
> 4. 在需要访问数据库的情况下，处理器会使用一个或多个数据结构，这些数据结构都是根据数据库中的数据建模而来的。
> 5. 当请求处理完毕时，处理器会调用模板引擎，有时还会向模板引擎传递通过数据模型获取到的数据；
> 6. 模板引擎会对模板文件进行分析并创建相应的模板，而这些模板又会与处理器传递的数据一起合并生成最终的HTML
> 7. 生成的HTML作为响应的一部分回传至客户端。

![](./png/web.jpg)



# 第二部分 Web应用的基本组成部分

## 使用go搭建服务器

使用go创建一个服务器，只要调用**ListenAndServe**，并传入网络地址以及负责处理请求的**处理器（handler）**作为参数就可以

```go
//sample 1
package main
import (
	"net/http"
)
func main() {
    http.ListenAndServe("", nil)
}

//sample 2
package main
import (
	"net/http"
)
func main() {
    server := http.Server{
        Addr: "127.0.0.1:8080",
        Handler:nil,
    }
    server.ListenAndServe()
}
```

## 处理器和处理器函数

前面内容启动了一个web服务器，但这个服务器尚未实现任何功能，所以访问该服务器会获得一个404 HTTP响应代码。原因是尚未对服务器编写任何处理器。

**处理请求**

在go中，一个**处理器**就是一个拥有ServerHTTP方法的接口，这个ServerHTTP方法需要接受两个参数：第一个参数是一个ResponseWriter接口，第二个参数是一个指向Request结构的指针。换句话说，任何接口只要拥有一个ServerHTTP方法，并且该方法带有以下签名，那么它就是一个处理器：

**ServerHTTP(http.ResponseWriter, *http.Request)**

**处理器函数**

处理器函数就是与处理器拥有相同行为的函数：

ServeMux是一个HTTP请求多路复用器，它负责接收HTTP请求并根据请求中的URL将请求重定向到正确的处理器。

**Request结构**



**模板引擎**

模板引擎通过将数据和模板组合在一起生成最终的HTML，而处理器则负责调用模板引擎并将引擎生成的HTML返回给客户端。

**动作**

Go模板的*动作*就是嵌入在模板里面的命令，这些命令在模板中使用两个大括号{{和}}进行包围。

* 条件动作

  ```
  {{ if arg }}
    some content
  {{ end }}
  ```

* 迭代动作

* 设置动作

* 包含动作

  包含动作允许用户在一个模板里面包含另一个模板，从而构建出嵌套的模板

  ```
  {{ template "name" }}
  ```


