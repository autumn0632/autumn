# blademaster （HTTP）架构简介

## 背景
在微服务这种分布式框架中，经常会有一些需求需要调用多个服务，但是还需要确保`服务的安全性`、`统一化每次的请求日志`或者`追踪用户完整的行为`等等。这其中可能会有以下问题：

* 很难让每一个服务都实现上述功能。因为对于开发者而言，他们应当注重的是实现功能。很多项目的开发者经常在一些日常开发中遗漏了这些关键点，经常有人会忘记去打日志或者去记录调用链。对于一些大流量的互联网服务而言，一个线上服务一旦发生故障时，即使故障时间很小，其影响面会非常大。一旦有人在关键路径上忘记路记录日志，那么故障的排除成本会非常高，
* 事实上实现之前叙述的这些功能的成本也非常高。比如说对于鉴权（Identify）这个功能，你要是去一个服务一个服务地去实现，那样的成本也是非常高的。

这是可能需要一个框架来帮助实现这些功能，在一些关键路径的请求上配置必要的鉴权或超时策略。那样服务间的调用会被多层中间件所过滤并检查，确保整体服务的稳定性。



## 设计目标

- 性能优异，不应该掺杂太多业务逻辑的成分
- 方便开发使用，开发对接的成本应该尽可能地小
- 后续鉴权、认证等业务逻辑的模块应该可以通过业务模块的开发接入该框架内
- 默认配置已经是 production ready 的配置，减少开发与线上环境的差异性

## 概述

- 参考`gin`设计整套HTTP框架，去除`gin`中不需要的部分逻辑
- 内置一些必要的中间件，便于业务方可以直接上手使用



## blademaster 架构

![](./png/bm1.png)

`blademaster`由几个非常精简的内部模块组成。其中`Router`用于根据请求的路径分发请求，`Context`包含了一个完整的请求信息，`Handler`则负责处理传入的`Context`，`Handlers`为一个列表，一个串一个地执行。
所有的`middlerware`均以`Handler`的形式存在，这样可以保证`blademaster`自身足够精简且扩展性足够强。

![](./png/bm2.png)

正常情况下每个`Handler`按照顺序一个一个串行地执行下去，但是`Handler`中也可以中断整个处理流程，直接输出`Response`。这种模式常被用于校验登陆的`middleware`中：一旦发现请求不合法，直接响应拒绝。

请求处理的流程中也可以使用`Render`来辅助渲染`Response`，比如对于不同的请求需要响应不同的数据格式`JSON`、`XML`，此时可以使用不同的`Render`来简化逻辑。

 一般而言，业务逻辑作为最后一个`handler`。