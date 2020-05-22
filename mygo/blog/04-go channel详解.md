>  通道：
>
> **通道是Go中的一种一等公民类型**， 它是Go的招牌特性之一。 和另一个招牌特性[协程](https://gfw.go101.org/article/control-flows-more.html#goroutine)一起，这两个招牌特性使得使用Go进行并发编程（concurrent programming）变得十分方便和有趣，并且大大降低了并发编程的难度。 csp模型

# 	一、通道介绍

## 0. 设计模式

**不要通过共享内存通信，而是通过通信来共享内存**

Go 语言的 Channel 在运行时使用 [`runtime.hchan`](https://github.com/golang/go/blob/e35876ec6591768edace6c6f3b12646899fd1b11/src/runtime/chan.go#L32) 结构体表示。我们在 Go 语言中创建新的 Channel 时，实际上创建的都是如下所示的结构体：

```go
type hchan struct {
  // 1. channel 中元素个数
	qcount   uint
  // 2. channel 中循环队列的长度
	dataqsiz uint
  // 3. channel 中缓冲区数据指针
	buf      unsafe.Pointer
  // 4. channel 能够收发的数据元素的大小
	elemsize uint16
  // 5. 关闭标志
	closed   uint32
  // 6. channel 元素类型
	elemtype *_type
  // 7. channel 的发送操作处理到的位置
	sendx    uint  
  // 8. channel 的接收操作处理到的位置
	recvx    uint
  // 9. 缓存空间为空时，用来保存接收数据的goroutinue
	recvq    waitq
  // 10. 缓存空间写满时，用来存放写数据的goroutinue
	sendq    waitq
	
  // 锁
	lock mutex
}
```



## 1. 创建通道

编译器会将 `make(chan int, 10)` 表达式被转换成 `OMAKE` 类型的节点，并在类型检查阶段将 `OMAKE` 类型的节点转换成 `OMAKECHAN` 类型：

```go
func typecheck1(n *Node, top int) (res *Node) {
	switch n.Op {
	case OMAKE:
		...
		switch t.Etype {
		case TCHAN:
			l = nil
			if i < len(args) { // 带缓冲区的异步 Channel
				...
				n.Left = l
			} else { // 不带缓冲区的同步 Channel
				n.Left = nodintconst(0)
			}
			n.Op = OMAKECHAN
		}
	}
}
```

`OMAKECHAN` 类型的节点最终都会在 SSA 中间代码生成阶段之前被转换成调用 [`runtime.makechan`](https://github.com/golang/go/blob/e35876ec6591768edace6c6f3b12646899fd1b11/src/runtime/chan.go#L71) 或者 [`runtime.makechan64`](https://github.com/golang/go/blob/e35876ec6591768edace6c6f3b12646899fd1b11/src/runtime/chan.go#L63) 的函数，其中后者用于处理缓冲区大小大于 2 的 32 次方的情况，我们重点关注 [`runtime.makechan`](https://github.com/golang/go/blob/e35876ec6591768edace6c6f3b12646899fd1b11/src/runtime/chan.go#L71) 函数：

```go
func makechan(t *chantype, size int) *hchan {
	elem := t.elem
	mem, _ := math.MulUintptr(elem.size, uintptr(size))

	var c *hchan
	switch {
  // 如果当前 Channel 中不存在缓冲区，那么就只会为 runtime.hchan 分配一段内存空间
	case mem == 0:
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		c.buf = c.raceaddr()
  // 如果当前 Channel 中存储的类型不是指针类型，就会为当前的 Channel 和底层的数组分配一块连续的内存空间；
	case elem.kind&kindNoPointers != 0:
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	// 在默认情况下会单独为 runtime.hchan 和缓冲区分配内存；
  default:
		c = new(hchan)
		c.buf = mallocgc(mem, elem, true)
	}
	c.elemsize = uint16(elem.size)
	c.elemtype = elem
	c.dataqsiz = uint(size)
	return c
}

```



### 2. 通道的类型和值

**类型**： 和数组、切片以及映射类型一样，每个通道类型也有一个元素类型。 一个通道只能传送它的（通道类型的）元素类型的值。 

 一个通道值可以被看作是先入先出（first-in-first-out，FIFO）队列。一个通道值可能是可读可写的、只读的（receive-only）或者只写的（send-only）。 

*  一个可读可写的通道值也称为一个双向通道。 一个双向通道类型的底层类型可以被表示为`chan T`。 
*  只能向一个只写的通道值发送数据，而不能从其中接收数据。 只写通道类型的底层类型可以被表示为`chan<- T`。  
* 只能从一个只读的通道值接收数据，而不能向其发送数据。 只读通道类型的底层类型可以被表示为`<-chan T`。 
*  双向通道`chan T`的值可以被隐式转换为单向通道类型`chan<- T`和`<-chan T`，但反之不行（即使显式也不行）。 类型`chan<- T`和`<-chan T`的值也不能相互转换。 

**值**： 通道类型的零值也使用预声明的`nil`来表示。 一个非零通道值必须通过内置的`make`函数来创建。 

### 3. 通道的操作

#### 3.1 基本操作

* **关闭通道**：

  ```go
  close(ch)
  ```

* 向通道发送值

  ```go
  ch <- v
  ```

   `v`必须能够赋值给通道`ch`的元素类型。 `ch`不能为单向接收通道。 `<-`称为数据发送操作符。 

* 从通道接收一个值

  ```go 
  <- ch
  ```

   在大多数场合下，一个数据接收操作可以被认为是一个**单值表达式**。

* 查询一个通道容量

  ```go
  cap(ch)
  ```

* 查询一个通道长度

  ```go
  len(ch)
  ```

 除了并发地关闭一个通道和向此通道发送数据这种情形，上面这些所有列出的操作都已经同步过了，因此它们可以在并发协程中安全运行而无需其它同步操作。 我们在编程中应该**避免并发地关闭一个通道和向此通道发送数据这种情形，**

**通道操作详解**

通道可分为以下三类：

* 零值（nil）通道
* 非零值但已关闭的通道
* 非零值且尚未关闭的通道

 三种通道操作施加到三类通道的结果如下：

|   操作   | 一个零值nil通道 | 一个非零值但已关闭的通道 | 一个非零值且尚未关闭的通道 |
| :------: | :-------------: | :----------------------: | :------------------------: |
|   关闭   |    产生恐慌     |         产生恐慌         |        成功关闭(C)         |
| 发送数据 |    永久阻塞     |         产生恐慌         |    阻塞或者成功发送(B)     |
| 接收数据 |    永久阻塞     |       永不阻塞(D)        |    阻塞或者成功接收(A)     |



**通道的元素值的传递都是复制过程**， 对于官方标准编译器，最大支持的通道的元素类型的尺寸为`65535`。  一般说来，为了在数据传递过程中避免过大的复制成本，我们不应该使用尺寸很大的通道元素类型。 如果欲传送的值的尺寸较大，应该改用指针类型做为通道的元素类型。 

# 发送channel

1. 发送数据时先判断channel类型，如果有缓冲区，判断channel是否还有空间，
2. 然后从等待channel中获取等待channel中的接受者，如果取到接收者，则将对象直接传递给接受者，然后将接受者所在的go放入P所在的可运行G队列,发送过程完成，
3. 如果未取到接收者，则将发送者enqueue到发送channel，发送者进入阻塞状态，
4. 有缓冲的channel需要先判断channel缓冲是否还有空间，如果缓冲空间已满，则将发送者enqueue到发送channel，发送者进入阻塞状态如果缓冲空间未满，则将元素copy到缓冲中，这时发送者就不会进入阻塞状态，
5. 最后尝试唤醒等待队列中的一个接受者。

# 接收channel

1. 首先判断channel的类型，然后如果是有缓冲的channel就判断缓冲中是否有元素，
2. 接着从channel中获取接受者，如果取到，则直接从接收者获取元素，并唤醒发送者，本次接收过程完成，
3. 如果没有取到接收者，阻塞当前的goroutine并等待发送者唤醒，
4. 如果是拥有缓冲的channel需要先判断缓冲中是否有元素，缓冲为空时，阻塞当前goroutine并等待发送者唤醒，缓冲如果不为空，则取出缓冲中的第一个元素，然后尝试唤醒channel中的一个发送者