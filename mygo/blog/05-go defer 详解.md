## 一、defer

* **功能**：

  > `defer` 会在当前函数或者方法返回之前执行传入的函数。它会经常被用于关闭文件描述符、关闭数据库连接以及解锁资源。

* **对`defer`的疑问**

  > * `defer` 关键字的调用时机以及多次调用 `defer` 时执行顺序是如何确定的；
  > * `defer` 关键字使用**传值**的方式传递参数时会进行**预计算**，导致不符合预期的结果；

* **预计算参数**

  > ```go 
  > func main() {
  > 	startedAt := time.Now()
  > 	defer fmt.Println(time.Since(startedAt))
  > 	
  > 	time.Sleep(time.Second)
  > }
  > 
  > //$ go run main.go
  > //0s  // 结果返回的是0，而不是1s之后
  > ```
  >
  > 原因分析：
  >
  > 调用 `defer` 关键字会立刻对函数中引用的外部参数进行拷贝，所以 `time.Since(startedAt)` 的结果不是在 `main` 函数退出之前计算的，而是在 `defer` 关键字调用时计算的，最终导致上述代码输出 0s。
  >
  > 解决方法：
  >
  > 向 `defer` 关键字传入匿名函数。
  >
  > 虽然调用 `defer` 关键字时也使用值传递，但是因为拷贝的是函数指针，所以 `time.Since(startedAt)` 会在 `main` 函数执行前被调用并打印出符合预期的结果。
  >
  > ```go
  > func main() {
  > 	startedAt := time.Now()
  > 	defer func() { fmt.Println(time.Since(startedAt)) }()
  > 	
  > 	time.Sleep(time.Second)
  > }
  > 
  > // $ go run main.go
  > // 1s
  > ```
  >
  > 

* **底层数据结构**

  > ```go
  > type _defer struct {
  >   // 1. 参数和结果的内存大小
  > 	siz     int32
  > 	// 2. 
  >   started bool
  >   // 3. 栈指针
  > 	sp      uintptr
  >   // 4. 程序计数器指针
  > 	pc      uintptr
  >   // 5. 传入的函数
  > 	fn      *funcval
  > 	_panic  *_panic
  >   // 链接成链表
  > 	link    *_defer
  > }
  > ```
  >
  > [`runtime._defer`](https://github.com/golang/go/blob/cfe3cd903f018dec3cb5997d53b1744df4e53909/src/runtime/runtime2.go#L853-L878) 结构体是延迟调用链表上的一个元素，所有的结构体都会通过 `link` 字段串联成链表。

* **编译过程**

  > 中间代码生成阶段执行的被 [`cmd/compile/internal/gc.state.stmt`](https://github.com/golang/go/blob/4d5bb9c60905b162da8b767a8a133f6b4edcaa65/src/cmd/compile/internal/gc/ssa.go#L1023-L1502) 函数会处理 `defer` 关键字。
  >
  > 编译器不仅将 `defer` 关键字都转换成 [`runtime.deferproc`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L218-L258) 函数，它还会为所有调用 `defer` 的函数末尾插入 [`runtime.deferreturn`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L526-L571) 的函数调用

* **运行过程**

  > `defer` 关键字的运行时实现分成两个部分：
  >
  > - [`runtime.deferproc`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L218-L258) 函数负责创建新的延迟调用；
  > - [`runtime.deferreturn`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L526-L571) 函数负责在函数调用结束时执行所有的延迟调用；
  >
  > 1. 创建延迟调用：
  >
  >    [`runtime.deferproc`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L218-L258) 会为 `defer` 创建一个新的 [`runtime._defer`](https://github.com/golang/go/blob/cfe3cd903f018dec3cb5997d53b1744df4e53909/src/runtime/runtime2.go#L853-L878) 结构体、设置它的函数指针 `fn`、程序计数器 `pc` 和栈指针 `sp` 并将相关的参数拷贝到相邻的内存空间中。获取 [`runtime._defer`](https://github.com/golang/go/blob/cfe3cd903f018dec3cb5997d53b1744df4e53909/src/runtime/runtime2.go#L853-L878)之后，它都会被追加到所在的 Goroutine `_defer` 链表的最前面。
  >
  > 2. 执行延迟调用：
  >
  >    [`runtime.deferreturn`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L526-L571) 会从 Goroutine 的 `_defer` 链表中取出最前面的 [`runtime._defer`](https://github.com/golang/go/blob/cfe3cd903f018dec3cb5997d53b1744df4e53909/src/runtime/runtime2.go#L853-L878) 结构体并调用 [`runtime.jmpdefer`](https://github.com/golang/go/blob/a38a917aee626a9b9d5ce2b93964f586bf759ea0/src/runtime/asm_386.s#L614-L624) 函数传入需要执行的函数和参数。[`runtime.jmpdefer`](https://github.com/golang/go/blob/a38a917aee626a9b9d5ce2b93964f586bf759ea0/src/runtime/asm_386.s#L614-L624) 是一个用汇编语言实现的运行时函数，它的工作就是跳转 `defer` 所在的代码段并在执行结束之后跳转回 [`runtime.deferreturn`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L526-L571)。

* **小结**

  > `defer` 关键字的实现主要依靠编译器和运行时的协作：
  >
  > - 编译期；
  >   - 将 `defer` 关键字被转换 [`runtime.deferproc`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L218-L258)；
  >   - 在调用 `defer` 关键字的函数返回之前插入 [`runtime.deferreturn`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L526-L571)；
  > - 运行时：
  >   - [`runtime.deferproc`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L218-L258) 会将一个新的 [`runtime._defer`](https://github.com/golang/go/blob/cfe3cd903f018dec3cb5997d53b1744df4e53909/src/runtime/runtime2.go#L853-L878) 结构体追加到当前 Goroutine 的链表头；
  >   - [`runtime.deferreturn`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L526-L571) 会从 Goroutine 的链表中取出 [`runtime._defer`](https://github.com/golang/go/blob/cfe3cd903f018dec3cb5997d53b1744df4e53909/src/runtime/runtime2.go#L853-L878) 结构并依次执行；

* **`defer`疑问解答**

  > 后调用的 `defer` 函数会先执行：
  >
  > - 后调用的 `defer` 函数会被追加到 Goroutine `_defer` 链表的最前面；
  > - 运行 [`runtime._defer`](https://github.com/golang/go/blob/cfe3cd903f018dec3cb5997d53b1744df4e53909/src/runtime/runtime2.go#L853-L878) 时是从前到后依次执行；
  >
  > 函数的参数会被预先计算；
  >
  > - 调用 [`runtime.deferproc`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L218-L258) 函数创建新的延迟调用时就会立刻拷贝函数的参数，函数的参数不会等到真正执行时计算；

* **扩展**

  > - 延迟函数可能操作主函数的具名返回值
  >
  > 当函数有具名返回值值，返回执行顺序为：1. 先给返回值赋值；2. 执行defer语句；包裹函数return返回

## 二、panic 和 recover

* **功能**

  > - `panic` 能够改变程序的控制流，函数调用`panic` 时会立刻停止执行函数的其他代码，并在执行结束后在当前 Goroutine 中递归执行调用方的延迟函数调用 `defer`；
  > - `recover` 可以中止 `panic` 造成的程序崩溃。它是一个只能在 `defer` 中发挥作用的函数，在其他作用域中调用不会发挥任何作用；

* **特点**

  > - `panic` 只会触发当前 Goroutine 的延迟函数调用， 跨协程失效
  > - `recover` 只有在 `defer` 函数中调用才会生效。换句话说，`recover` 只有在发生 `panic` 之后调用才会生效。
  > - `panic` 允许在 `defer` 中嵌套多次调用；

* **数据结构**

  > ```go
  > type _panic struct {
  > // 1. 指向 defer 调用时参数的指针；
  > 	argp      unsafe.Pointer
  > // 2. 指向 defer 调用时参数的指针；
  > 	arg       interface{}
  > // 3. 指向了更早调用的 runtime._panic 结构；
  > 	link      *_panic
  > // 4. 表示当前 runtime._panic 是否被 recover 恢复；
  > 	recovered bool
  > // 5. 表示当前的 panic 是否被强行终止；
  > 	aborted   bool
  > 
  > 	pc        uintptr
  > 	sp        unsafe.Pointer
  > 	goexit    bool
  > }
  > ```
  >
  > 

* **执行流程**

  > 编译器会将关键字 `panic` 转换成 [`runtime.gopanic`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L887-L1062)，该函数的执行过程包含以下几个步骤：
  >
  > * 创建新的 [`runtime._panic`](https://github.com/golang/go/blob/cfe3cd903f018dec3cb5997d53b1744df4e53909/src/runtime/runtime2.go#L891-L900) 结构并添加到所在 Goroutine `_panic` 链表的最前面
  > * 在循环中不断从当前 Goroutine 的 `_defer` 中链表获取 [`runtime._defer`](https://github.com/golang/go/blob/cfe3cd903f018dec3cb5997d53b1744df4e53909/src/runtime/runtime2.go#L853-L878) 并调用 [`runtime.reflectcall`](https://github.com/golang/go/blob/a38a917aee626a9b9d5ce2b93964f586bf759ea0/src/runtime/asm_386.s#L496-L526) 运行延迟调用函数
  > * 调用 [`runtime.fatalpanic`](https://github.com/golang/go/blob/22d28a24c8b0d99f2ad6da5fe680fa3cfa216651/src/runtime/panic.go#L1185-L1220) 中止整个程序

## 三、select

* **概述**

  > 