# 一、简介

在计算机性能调试领域里，profiling 就是对应用的画像，这里画像就是应用使用 CPU 和内存的情况。也就是说应用使用了多少 CPU 资源？都是哪些部分在使用？每个函数使用的比例是多少？有哪些函数在等待 CPU 资源？知道了这些，我们就能对应用进行规划，也能快速定位性能瓶颈。

golang 是一个对性能特别看重的语言，因此语言中自带了 profiling 的库。

在 go 语言中，主要关注的应用运行情况主要包括以下几种：

- CPU profile：报告程序的 CPU 使用情况，按照一定频率去采集应用程序在 CPU 和寄存器上面的数据
- Memory Profile（Heap Profile）：报告程序的内存使用情况
- Block Profiling：报告 goroutines 不在运行状态的情况，可以用来分析和查找死锁等性能瓶颈
- Goroutine Profiling：报告 goroutines 的使用情况，有哪些 goroutine，它们的调用关系是怎样的

# 二、 使用

## 1. 两种性能搜集方法

### 工具型应用

### 服务型应用-`net/http/pprof`

如果你的应用是一直运行的，比如 web 应用，那么可以使用 `net/http/pprof` 库，它能够在提供 HTTP 服务进行分析。

* 代码

```go
import _ "net/http/pprof"

func main() {
  ...
  
  ...
  
  ...
  
  http.ListenAndServe("0.0.0.0:8000", nil)
  
}

```

* 调试：

  1. web访问：http://http://62.234.81.226/:8080//debug/pprof/profile

  2. go tool：`go tool pprof [binary] [source]`（go tool pprof ./sniff-go-pprof http://62.234.81.226:8080//debug/pprof/profile）

     profile：查看cpu

     heap：查看内存

     命令执行后会进入交互模式：输入一下命令，

     * topN：列出最耗时的N个函数
     * list handleName：查看匹配函数的代码以及每行代码的耗时
     * web：生成函数调用图，在浏览器中显示（需安装Graphviz）

  3. 火焰图：

     下载to-touch工具：go get github.com/uber/go-torch

     执行：go-torch -u  http://62.234.81.226:8080//debug/pprof 生成svg文件

     

