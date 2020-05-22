* **功能**

> reflect 实现了运行时的反射能力，能够让程序操作不同类型的对象。反射包中有两对非常重要的函数和类型，[`reflect.TypeOf`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/type.go#L1365-L1368) 能获取类型信息，[`reflect.ValueOf`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/value.go#L2316-L2328) 能获取数据的运行时表示，另外两个类型是 `Type` 和 `Value`，它们与函数是一一对应的关系。我们通过 [`reflect.TypeOf`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/type.go#L1365-L1368)、[`reflect.ValueOf`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/value.go#L2316-L2328) 可以将一个普通的变量转换成『反射』包中提供的 `Type` 和 `Value`，随后就可以使用反射包中的方法对它们进行复杂的操作。

* **反射三大法则**

> 1. **从 `interface{}` 变量可以反射出反射对象**
>
>    当执行 `reflect.ValueOf(1)` 时，虽然看起来是获取了基本类型 `int` 对应的反射类型，但是由于 [`reflect.TypeOf`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/type.go#L1365-L1368)、[`reflect.ValueOf`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/value.go#L2316-L2328) 两个方法的入参都是 `interface{}` 类型，所以在方法执行的过程中发生了类型转换。
>
>    有了变量的类型之后，我们可以通过 `Method` 方法获得类型实现的方法，通过 `Field` 获取类型包含的全部字段。对于不同的类型，我们也可以调用不同的方法获取相关信息：
>
>    - 结构体：获取字段的数量并通过下标和字段名获取字段 `StructField`；
>    - 哈希表：获取哈希表的 `Key` 类型；
>    - 函数或方法：获取入参和返回值的类型；
>    - ...
>
>    总而言之，使用 [`reflect.TypeOf`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/type.go#L1365-L1368) 和 [`reflect.ValueOf`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/value.go#L2316-L2328) 能够获取 Go 语言中的变量对应的反射对象。一旦获取了反射对象，我们就能得到跟当前类型相关数据和操作，并可以使用这些运行时获取的结构执行方法。
>
>    ```go
>    package main
>    
>    import (
>    	"fmt"
>    	"reflect"
>    )
>    
>    func main() {
>    	author := "draven"
>    	fmt.Println("TypeOf author:", reflect.TypeOf(author))
>    	fmt.Println("ValueOf author:", reflect.ValueOf(author))
>    }
>    
>    // $ go run main.go
>    // TypeOf author: string
>    // ValueOf author: draven
>    ```
>
>    
>
> 2. 从反射对象可以获取 `interface{}` 变量；
>
>    既然能够将接口类型的变量转换成反射对象，那么一定需要其他方法将反射对象还原成接口类型的变量，[`reflect`](https://golang.org/pkg/reflect/) 中的 [`reflect.Value.Interface`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/value.go#L992-L994) 方法就能完成这项工作。
>
>    调用 [`reflect.Value.Interface`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/value.go#L992-L994) 方法只能获得 `interface{}` 类型的变量，如果想要将其还原成最原始的状态还需要经过如下所示的显式类型转换：
>
>    ```go
>    v := reflect.ValueOf(1)
>    v.Interface().(int)
>    ```
>
>    从反射对象到接口值的过程就是从接口值到反射对象的镜面过程，两个过程都需要经历两次转换：
>
>    - 从接口值到反射对象：
>      - 从基本类型到接口类型的类型转换；
>      - 从接口类型到反射对象的转换；
>    - 从反射对象到接口值：
>      - 反射对象转换成接口类型；
>      - 通过显式类型转换变成原始类型；
>
> 3. 要修改反射对象，其值必须可设置
>
>    如果我们想要更新一个 `reflect.Value`，那么它持有的值一定是可以被更新的，假设我们有以下代码：
>
>    ```
>    func main() {
>    	i := 1
>    	v := reflect.ValueOf(i)
>    	v.SetInt(10)
>    	fmt.Println(i)
>    }
>    
>    $ go run reflect.go
>    panic: reflect: reflect.flag.mustBeAssignable using unaddressable value
>    
>    goroutine 1 [running]:
>    reflect.flag.mustBeAssignableSlow(0x82, 0x1014c0)
>    	/usr/local/go/src/reflect/value.go:247 +0x180
>    reflect.flag.mustBeAssignable(...)
>    	/usr/local/go/src/reflect/value.go:234
>    reflect.Value.SetInt(0x100dc0, 0x414020, 0x82, 0x1840, 0xa, 0x0)
>    	/usr/local/go/src/reflect/value.go:1606 +0x40
>    main.main()
>    	/tmp/sandbox590309925/prog.go:11 +0xe0
>    ```
>
>    出错原因：Go 语言的[函数调用](http://draveness.me/golang-function-call)都是传值的，所以我们得到的反射对象跟最开始的变量没有任何关系，所以直接对它修改会导致崩溃。
>
>    可以通过以下方法修改：
>
>    * 调用 [`reflect.ValueOf`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/value.go#L2316-L2328) 函数获取变量指针
>    * 调用 [`reflect.Value.Elem`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/value.go#L788-L821) 方法获取指针指向的变量
>    * 调用 [`reflect.Value.SetInt`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/value.go#L1600-L1616) 方法更新变量的值
>
>    ```go
>    func main() {
>    	i := 1
>    	v := reflect.ValueOf(&i)
>    	v.Elem().SetInt(10)
>    	fmt.Println(i)
>    }
>    
>    // $ go run reflect.go
>    // 10
>    ```
>
>    

* 实现过程

  > **当我们想要将一个变量转换成反射对象时，Go 语言会在编译期间完成类型转换的工作，将变量的类型和值转换成了 `interface{}` 并等待运行期间使用 `reflect`包获取接口中存储的信息。**
  >
  > [`reflect.TypeOf`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/type.go#L1365-L1368) 函数的实现原理将一个 `interface{}` 变量转换成了内部的 `emptyInterface` 表示，然后从中获取相应的类型信息。
  >
  > ```go
  > type emptyInterface struct {
  > 	typ  *rtype // 变量类型
  > 	word unsafe.Pointer // 指向内部封装的数据
  > }
  > 
  > 
  > func TypeOf(i interface{}) Type {
  > 	eface := *(*emptyInterface)(unsafe.Pointer(&i))
  > 	return toType(eface.typ)
  > }
  > 
  > func toType(t *rtype) Type {
  > 	if t == nil {
  > 		return nil
  > 	}
  > 	return t
  > }
  > 
  > func (t *rtype) String() string {
  > 	s := t.nameOff(t.str).name()
  > 	if t.tflag&tflagExtraStar != 0 {
  > 		return s[1:]
  > 	}
  > 	return s
  > }
  > ```
  >
  > [`reflect.ValueOf`](https://github.com/golang/go/blob/52c4488471ed52085a29e173226b3cbd2bf22b20/src/reflect/value.go#L2316-L2328) 实现中先调用了 [`reflect.escapes`](https://github.com/golang/go/blob/4e8d27068df52eb372dc2ba7e929e47850934805/src/reflect/value.go#L2779-L2783) 函数保证当前值逃逸到堆上，然后通过 [`reflect.unpackEface`](https://github.com/golang/go/blob/4e8d27068df52eb372dc2ba7e929e47850934805/src/reflect/value.go#L140-L152) 方法从接口中获取 `Value` 结构体。
  >
  > ```
  > func ValueOf(i interface{}) Value {
  > 	if i == nil {
  > 		return Value{}
  > 	}
  > 
  > 	escapes(i)
  > 
  > 	return unpackEface(i)
  > }
  > 
  > func unpackEface(i interface{}) Value {
  > 	e := (*emptyInterface)(unsafe.Pointer(&i))
  > 	t := e.typ
  > 	if t == nil {
  > 		return Value{}
  > 	}
  > 	f := flag(t.Kind())
  > 	if ifaceIndir(t) {
  > 		f |= flagIndir
  > 	}
  > 	return Value{t, e.word, f}
  > }
  > ```

* **小结**

  > Go 语言的 [`reflect`](https://golang.org/pkg/reflect/) 包为我们提供的多种能力，包括如何使用反射来动态修改变量、判断类型是否实现了某些接口以及动态调用方法等功能。