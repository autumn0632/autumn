## new() 与make()的区别

**值类型：**int，string，数组，结构体，当不指定变量的默认值时，默认值时它们的零值，可以直接使用

**引用类型：**chan，slice，函数类型，map，引用类型的零值是nil。

对于引用类型的变量，我们**不光要声明它，还要为它分配内容空间**，否则我们的值放在哪，值类型的声明不需要，是因为已经默认帮我们分配好了。

```go
func main() {
 var i *int
 i=new(int)
 *i=10
 fmt.Println(*i)
}

//接受一个参数，这个参数是一个类型，分配好内存后，返回一个指向该类型内存地址的指针。同时请注意它同时把分配的内存置为零，也就是类型的零值。
```



make也是用于内存分配的，但是和new不同，它只用于chan、map以及切片的内存创建，而且它返回的类型就是这三个类型本身，而不是他们的指针类型，因为这三种类型就是引用类型，所以就没有必要返回他们的指针了。



## 反射

**1. 什么是反射**

反射就是用来**检测存储在接口变量内部(值value；类型concrete type) pair对的一种机制。**

反射就是建立在类型之上的，主要与Golang的interface类型相关（它的type是concrete type），只有interface类型才有反射一说。

> 变量包括（type, value）两部分
>
> type 包括 static type和concrete type. 简单来说 static type是你在编码是看见的类型(如int、string)，concrete type是runtime系统看见的类型
>
> 类型断言能否成功，取决于变量的concrete type，而不是static type. 因此，一个 reader变量如果它的concrete type也实现了write方法的话，它也可以被类型断言为writer.

在Golang的实现中，每个interface变量都有一个对应pair，pair中记录了实际变量的值和类型:**(value, type)**

value是实际变量值，type是实际变量的类型。一个interface{}类型的变量包含了2个指针，一个指针指向值的类型【对应concrete type】，另外一个指针指向实际的值【对应value】。



**2. reflect包**

Golang语言实现了反射，反射机制就是在运行时动态的调用对象的方法和属性，官方自带的reflect包就是反射相关的，只要包含这个包就可以使用。

 reflect包提供了两种类型（或者说两个方法）让我们可以很容易的访问接口变量内容，分别是**reflect.ValueOf()** 和 **reflect.TypeOf()**，

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var num float64 = 1.2345

	fmt.Println("type: ", reflect.TypeOf(num))
	fmt.Println("value: ", reflect.ValueOf(num))
}

//运行结果:
//type:  float64
//value:  1.2345
```



## iota

iota是golang语言的常量计数器,只能在常量的表达式中使用。iota在const关键字出现时将被重置为0(const内部的第一行之前)，const中每新增一行常量声明将使iota计数一次



## 类型断言

接口类型向普通类型的转换称为类型断言(运行期确定)

```go
//方式一：
if instance,ok := 接口对象.(实际类型1);ok{

}else if instance,ok := 接口对象.(实际类型2);ok{

}else if instance,ok := 接口对象.(实际类型3);ok{

}
//方式二:
switch instance := 接口对象.(type) {
case 实际类型1：
case 实际类型2：
case 实际类型3：
}
```



## 结构体中的反引号

```go
type Account struct {
	// 把struct编码成json字符串时，common.Address字段的key是address
	Address common.Address `json:"address"` // Ethereum account address derived from the key
	// 把struct编码成json字符串时，URL字段的key是url
	URL     URL            `json:"url"`     // Optional resource locator within a backend
}
```

