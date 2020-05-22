* **什么是接口**

  > * Go 语言中的接口是一种内置的类型，它定义了一组方法的签名。
  > * **接口的本质就是引入一个新的中间层，调用方可以通过接口与具体实现分离，解除上下游的耦合，上层的模块不再需要依赖下层的具体模块，只需要依赖一个约定好的接口。**
  > * 接口还能够帮助我们隐藏底层实现，减少关注点。
  > * 

* 接口的隐式实现

  > 一个类型只要实现了接口中的方法，就说这个类型实现了此接口，不需要显示的声明

* **两种接口**

  > 一种是带有一组方法的接口，另一种是不带任何方法的 `interface{}`。Go 语言使用 `iface` 结构体表示第一种接口，使用 `eface` 结构体表示第二种空接口。由于后者在 Go 语言中非常常见，所以在实现时使用了特殊的类型。

* **结构体、指针和接口**

  > * 当我们使用指针实现接口时，只有指针类型的变量才会实现该接口；当我们使用结构体实现接口时，指针类型和结构体类型都会实现该接口（隐式转换）。
  >
  >   ```go
  >   //作为指针的 &Cat{} 变量能够隐式地获取到指向的结构体，
  >   //所以能在结构体上调用 Walk 和 Quack 方法。
  >   type Cat struct{}
  >   
  >   func (c Cat) Quack() {
  >   	fmt.Println("meow")
  >   }
  >   
  >   func main() {
  >   	var c Duck = &Cat{}
  >   	c.Quack()
  >   }
  >   
  >   // 以下方法编译不过，编译器会提醒我们：Cat 类型没有实现 Duck 接口，Quack 方法的接受者是指针。
  >   type Duck interface {
  >   	Quack()
  >   }
  >   
  >   type Cat struct{}
  >   
  >   func (c *Cat) Quack() {
  >   	fmt.Println("meow")
  >   }
  >   
  >   func main() {
  >   	var c Duck = Cat{}
  >   	c.Quack()
  >   }
  >   
  >   ```
  >
  >   

* **nil != nil**

  > ```go
  > package main
  > 
  > type TestStruct struct{}
  > 
  > func NilOrNot(v interface{}) bool {
  > 	return v == nil
  > }
  > 
  > func main() {
  > 	var s *TestStruct
  > 	fmt.Println(s == nil)      // #=> true
  > 	fmt.Println(NilOrNot(s))   // #=> false
  > }
  > 
  > // $ go run main.go
  > // true
  > //false
  > ```
  >
  > 上述代码的执行结果：
  >
  > - 将上述变量与 `nil` 比较会返回 `true`；
  > - 将上述变量传入 `NilOrNot` 方法并与 `nil` 比较会返回 `false`；
  >
  > 原因：
  >
  > 调用 `NilOrNot` 函数时发生了**隐式的类型转换**，除了向方法传入参数之外，变量的赋值也会触发隐式类型转换。在类型转换时，`*TestStruct` 类型会转换成 `interface{}` 类型，转换后的变量不仅包含转换前的变量，还包含变量的类型信息 `TestStruct`，所以转换后的变量与 `nil` 不相等。
  
* **数据结构**

  > ```go
  > // 空接口类型
  > type eface struct { // 16 bytes
  > 	_type *_type
  > 	data  unsafe.Pointer
  > }
  > 
  > // 包含一组方法的接口类型
  > type iface struct { // 16 bytes
  > 	tab  *itab
  > 	data unsafe.Pointer
  > }
  > 
  > ```
  >
  > `_type`和`itab`类型
  >
  > ```go
  > // _type 是 Go 语言类型的运行时表示。包含了很多元信息，例如：类型的大小、哈希、对齐以及种类等。
  > type _type struct {
  > 	size       uintptr
  > 	ptrdata    uintptr
  > 	hash       uint32
  > 	tflag      tflag
  > 	align      uint8
  > 	fieldAlign uint8
  > 	kind       uint8
  > 	equal      func(unsafe.Pointer, unsafe.Pointer) bool
  > 	gcdata     *byte
  > 	str        nameOff
  > 	ptrToThis  typeOff
  > }
  > 
  > // itab 结构体是接口类型的核心组成部分，每一个 itab 都占 32 字节的空间，我们可以将其看成接口类型和具体类型的组合，它们分别用 inter 和 _type 两个字段表示：
  > type itab struct { // 32 bytes
  > 	inter *interfacetype
  > 	_type *_type
  >   //hash 是对 _type.hash 的拷贝，当我们想将 interface 类型转换成具体类型时，可以使用该字段快速判断目标类型和具体类型 _type 是否一致；
  > 	hash  uint32
  > 	_     [4]byte
  > 	fun   [1]uintptr
  > }
  > ```

* **类型转换**

  > `指针类型`：
  >
  > 1. 结构体 `Cat` 的初始化；
  > 2. 赋值触发的类型转换过程；
  > 3. 调用接口的方法 `Quack()`；
  >
  > ```go
  > // 指针类型
  > ackage main
  > 
  > type Duck interface {
  > 	Quack()
  > }
  > 
  > type Cat struct {
  > 	Name string
  > }
  > 
  > //go:noinline
  > func (c *Cat) Quack() {
  > 	println(c.Name + " meow")
  > }
  > 
  > func main() {
  > 	var c Duck = &Cat{Name: "grooming"}
  > 	c.Quack()
  > }
  > ```

* **类型断言**

  > 类型断言：将接口类型转换成实际类型
  >
  > ```go
  > //方式一：
  > if instance,ok := 接口对象.(实际类型1);ok{
  > 
  > }else if instance,ok := 接口对象.(实际类型2);ok{
  > 
  > }else if instance,ok := 接口对象.(实际类型3);ok{
  > 
  > }
  > //方式二:
  > switch instance := 接口对象.(type) {
  > case 实际类型1：
  > case 实际类型2：
  > case 实际类型3：
  > }
  > ```
  >
  > 