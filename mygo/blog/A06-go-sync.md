# sync

Go 语言在 `sync` 包中提供了用于同步的一些基本原语，包括常见的`sync.Mutex`、`sync.RWMutex`、`sync.WaitGroup`、`sync.Once`和 `sync.Cond`。这些基本原语提高了较为基础的同步功能，但是它们是一种相对原始的同步机制，在多数情况下，我们都应该使用抽象层级的更高的 Channel 实现同步。

## Mutex

互斥锁由`state`和`sema`组成，其中 `state` 表示当前互斥锁的状态，而 `sema` 是用于控制锁状态的信号量。

```go
type Mutex struct {
	state int32
	sema  uint32
}
```

### 状态

在默认情况下，互斥锁的所有状态位都是 `0`，`int32` 中的不同位分别表示了不同的状态，最低三位分别表示 `mutexLocked`、`mutexWoken` 和 `mutexStarving`，剩下的位置用来表示当前有多少个 Goroutine 等待互斥锁的释放：

- `mutexLocked` — 表示互斥锁的锁定状态；
- `mutexWoken` — 表示从正常模式被从唤醒；
- `mutexStarving` — 当前的互斥锁进入饥饿状态；
- `waitersCount` — 当前互斥锁上等待的 Goroutine 个数；

**饥饿模式与正常模式**

相比于饥饿模式，正常模式下的互斥锁能够提供更好地性能，饥饿模式的能避免 Goroutine 由于陷入等待无法获取锁而造成的高尾延时。

### **加锁和解锁**

* 互斥锁的加锁过程比较复杂，它涉及自旋、信号量以及调度等概念：
  - 如果互斥锁处于初始化状态，就会直接通过置位 `mutexLocked` 加锁；
  - 如果互斥锁处于 `mutexLocked` 并且在普通模式下工作，就会进入自旋，执行 30 次 `PAUSE` 指令消耗 CPU 时间等待锁的释放；
  - 如果当前 Goroutine 等待锁的时间超过了 1ms，互斥锁就会切换到饥饿模式；
  - 互斥锁在正常情况下会通过 `sync.runtime_SemacquireMutex`函数将尝试获取锁的 Goroutine 切换至休眠状态，等待锁的持有者唤醒当前 Goroutine；
  - 如果当前 Goroutine 是互斥锁上的最后一个等待的协程或者等待的时间小于 1ms，当前 Goroutine 会将互斥锁切换回正常模式；

* 互斥锁的解锁过程与之相比就比较简单
  - 当互斥锁已经被解锁时，那么调用 `sync.Mutex.Unlock`会直接抛出异常；
  - 当互斥锁处于饥饿模式时，会直接将锁的所有权交给队列中的下一个等待者，等待者会负责设置 `mutexLocked` 标志位；
  - 当互斥锁处于普通模式时，如果没有 Goroutine 等待锁的释放或者已经有被唤醒的 Goroutine 获得了锁，就会直接返回；在其他情况下会通过 `sync.runtime_Semrelease` 唤醒对应的 Goroutine

## RWMutex

读写互斥锁 `sync.RWMutex`是细粒度的互斥锁，它不限制资源的并发读，但是读写、写写操作无法并行执行。

读写互斥锁在互斥锁之上提供了额外的更细粒度的控制，能够在**读操作远远多于写操作时**提升性能。



## WaitGroup

`sync.WaitGroup` 用于等待一组 Goroutine 的返回。将原本顺序执行的代码在多个 Goroutine 中并发执行，加快程序处理的速度。

`sync.WaitGroup` 对外暴露了三个方法 — `sync.WaitGroup.Add`、`sync.WaitGroup.Wait` 和 `sync.WaitGroup.Done`。



## Once

 `sync.Once` 可以保证在 Go 程序运行期间的某段代码只会执行一次。

```go
func main() {
    o := &sync.Once{}
    for i := 0; i < 10; i++ {
        o.Do(func() {
            fmt.Println("only once")
        })
    }
}

$ go run main.go
only once
```

## Cond

`sync.Cond` 是一个条件变量，它可以让一系列的 Goroutine 都在满足特定条件时被唤醒。每一个 `sync.Cond` 结构体在初始化时都需要传入一个互斥锁，

```go
func main() {
	c := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 10; i++ {
		go listen(c)
	}
	time.Sleep(1*time.Second)
	go broadcast(c)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

func broadcast(c *sync.Cond) {
	c.L.Lock()
	c.Broadcast()
	c.L.Unlock()
}

func listen(c *sync.Cond) {
	c.L.Lock()
	c.Wait()
	fmt.Println("listen")
	c.L.Unlock()
}
```

上述代码同时运行了 11 个 Goroutine，这 11 个 Goroutine 分别做了不同事情：

- 10 个 Goroutine 通过 `sync.Cond.Wait` 等待特定条件的满足；
- 1 个 Goroutine 会调用 `sync.Cond.Broadcast` 方法通知所有陷入等待的 Goroutine；

调用 `sync.Cond.Broadcast` 方法后，上述代码会打印出 10 次 “listen” 并结束调用。

`sync.Cond`暴露的方法：

- `sync.Cond.Wait` 方法在调用之前一定要使用获取互斥锁，否则会触发程序崩溃；
- `sync.Cond.Signal` 方法唤醒的 Goroutine 都是队列最前面、等待最久的 Goroutine；
- `sync.Cond.Broadcast` 会按照一定顺序广播通知等待的全部 Goroutine；

## Pool

**保存和复用临时对象，减少内存分配，降低GC压力。**对象越多GC越慢，因为Golang进行三色标记回收的时候，要标记的也越多，自然就慢。

**Pool用于存储那些被分配了但是没有被使用，而未来可能会使用的值，以减小垃圾回收的压力。**

使用：fmt包，fmt包总是需要使用一些[]byte之类的对象，golang建立了一个临时对象池，存放着这些对象，如果需要使用一个[]byte，就去Pool里面拿，如果拿不到就分配一份。

