# 高效的数据压缩编码方式

**Protocol buffers 是一种语言中立，平台无关，可扩展的结构化数据存储格式，可用于通信协议，数据存储等。**由于它是一种二进制的格式，比使用 xml 进行数据交换快许多。可以把它用于分布式应用之间的数据通信或者异构环境下的数据交换。作为一种效率和兼容性都很优秀的二进制数据传输格式，可以用于诸如网络传输、配置文件、数据存储等诸多领域。

重点：

1. Protocol Buffer 利用 varint 原理压缩数据以后，二进制数据非常紧凑，option 也算是压缩体积的一个举措。所以 pb 体积更小，如果选用它作为网络数据传输，势必相同数据，消耗的网络流量更少。但是并没有压缩到极限，float、double 浮点型都没有压缩。
2. Protocol Buffer 比 JSON 和 XML 少了 {、}、: 这些符号，体积也减少一些。再加上 varint 压缩，gzip 压缩以后体积更小！
3. Protocol Buffer 是 Tag - Value (Tag - Length - Value)的编码方式的实现，减少了分隔符的使用，数据存储更加紧凑。
4. Protocol Buffer 另外一个核心价值在于提供了一套工具，一个编译工具，自动化生成 get/set 代码。简化了多语言交互的复杂度，使得编码解码工作有了生产力。
5. Protocol Buffer 不是自我描述的，离开了**数据描述 `.proto` 文件**，就无法理解二进制数据流。这点即是优点，使数据具有一定的“加密性”，也是缺点，数据可读性极差。所以 Protocol Buffer 非常适合内部服务之间 RPC 调用和传递数据。
6. Protocol Buffer 具有向后兼容的特性，更新数据结构以后，老版本依旧可以兼容，这也是 Protocol Buffer 诞生之初被寄予解决的问题。因为编译器对不识别的新增字段会跳过不处理。



## 定义消息类型

假设定义一个“搜索请求”的消息格式，在这个请求里面就，含有一个**查询字符串**、你感兴趣的**查询结果所在的页数**，以及每一页**多少条查询结果**。可以采用如下的方式来定义消息类型的.proto文件

```shell
//指明使用的是proto3语法，如果不指明，protobuf编译器将默认使用proto2语法，
//必须是文件的第一个非空的非注释行。
syntax = "proto3"; 

message SearchRequest {
   string query = 1;
   int32 page_number = 2;
   int32 result_per_page = 3;
}
```

* 第一行指定正在使用`proto3`语法,不这样做，protobuf 编译器将假定您正在使用proto2。这必须是文件的第一个非空的非注释行。

* 指定字段类型，

* 分配标识号，在消息定义中，每个字段都有唯一的一个**数字标识符**。这些标识符是用来在消息的二进制格式中识别各个字段的，一旦开始使用就不能够再改变。

  > [1,15]之内的标识号在编码的时候会占用一个字节。[16,2047]之内的标识号则占用2个字节。所以应该为那些频繁出现的消息元素保留 [1,15]之内的标识号。切记：要为将来有可能添加的、频繁出现的标识号预留一些标识号。

  **指定变量规则**：

  在 proto3 中，可以给变量指定以下两个规则：

  * `singular`：0或者1个，但不能多于1个
  * `repeated`：任意数量（包括0）

  **.proto 文件编译**

  当你使用 [protoc 编译器](https://link.jianshu.com/?t=https://developers.google.com/protocol-buffers/docs/proto3#generating) 编译一个.proto 文件的时候，编译器会根据你选择的语言和你在这个.proto 文件定义的消息类型生成代码，，这些代码的
  功能包括：字段值的 getter，setter，消息序列化并写入到输出流，从输入流接反序列化读取消息等。

  对于go语言来说，针对每一个定义的消息类型编译器会创建一个带类型的.pb.go 文件。

  **引用其他proto文件**


## 定义服务

如果要将消息类型与RPC（远程过程调用）系统一起使用，则可以在`.proto`文件中定义RPC服务接口，protobuf 编译器将使用选择的语言生成服务**接口代码**和**存根**。gRPC允许定义4种类型的service方法。

**简单rpc**

如果要定义RPC服务请求方法为:`SearchRequest`和返回方法为:`SearchResponse`，可以`.proto`按如下方式在文件中定义它。

```
 service SearchService {
    rpc Search（SearchRequest）returns（SearchResponse）;
  }
```

**服务器端（应答）流式rpc**

客户端发送请求到服务器，拿到一个流去读取返回的消息序列。 客户端读取返回的流，直到里面没有任何消息。从例子中可以看出，通过在 响应 类型前插入 `stream` 关键字，可以指定一个服务器端的流方法。

```
rpc ListFeatures(Rectangle) returns (stream Feature) {}
```

**客户端（请求）流式rpc**

 客户端写入一个消息序列并将其发送到服务器，同样也是使用流。一旦客户端完成写入消息，它等待服务器完成读取返回它的响应。通过在 请求 类型前指定 `stream` 关键字来指定一个客户端的流方法。

```
  rpc RecordRoute(stream Point) returns (RouteSummary) {}
```

**双向流式rpc**

一个 双向流式 RPC 是双方使用读写流去发送一个消息序列。两个流独立操作，因此客户端和服务器可以以任意喜欢的顺序读写：比如， 服务器可以在写入响应前等待接收所有的客户端消息，或者可以交替的读取和写入消息，或者其他读写的组合。 每个流中的消息顺序被预留。你可以通过在请求和响应前加 stream 关键字去制定方法的类型。

```
rpc RouteChat(stream RouteNote) returns (stream RouteNote) {}
```



## 生成所需要的类

Protobuf 编译器的调用如下：

> ```
> protoc --proto_path = IMPORT_PATH --cpp_out = DST_DIR --java_out = DST_DIR --python_out = DST_DIR --go_out = DST_DIR --ruby_out = DST_DIR --objc_out = DST_DIR --csharp_out = DST_DIR  path / to / file .proto
>   
> ```