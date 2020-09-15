go语言中的io操作主要分布在以下几个包中：

* [io](http://docs.studygolang.com/pkg/io/) 

    属于底层接口定义库，其作用是定义一些基本接口和一些基本常量，并对这些接口的作用给出说明，常见的接口有`Reader`、`Writer`等。
    
* [io/ioutil](http://docs.studygolang.com/pkg/io/ioutil/) 

    ioutil库包含在io目录下，它的主要作用是作为一个工具包，里面有一些比较实用的函数，比如 ReadAll(从某个源读取数据)、ReadFile（读取文件内容）、WriteFile（将数据写入文件）、ReadDir（获取目录）
   
* [os]()
    
    os库主要是跟操作系统打交道，所以文件操作基本都会跟os库挂钩，比如创建文件、打开一个文件等。这个库往往会和ioutil库、bufio库等配合使用
    
* [bufio](http://docs.studygolang.com/pkg/bufio/) 

    理解为在io库上再封装一层，加上了缓存功能。
    
    > 1. `bufio VS ioutil`：两者都提供了对文件的读写功能，唯一的不同就是bufio多了一层缓存的功能，这个优势主要体现读取大文件的时候（ioutil.ReadFile是一次性将内容加载到内存，如果内容过大，很容易爆内存）
    > 2.  `bufio VS bytes.Buffer`：两者都提供一层缓存功能，它们的不同主要在于 bufio 针对的是文件到内存的缓存，而 bytes.Buffer 的针对的是内存到内存的缓存
    >

* [bytes]() 和 [strings]()

    `bytes`和`strings`库都实现了Reader接口，所以它们的不同主要在于针对的对象不同。bytes针对的是字节，strings针对的是字符串。
    另一个区别就是 bytes还带有Buffer的功能，但是 strings没提供。

# 一、io-基本的io接口

## 1. Reader

 接口的定义如下：

```go
type Reader interface {
    Read(p []byte) (n int, err error)  //p中没数据，从数据源将数据写入到p中
}
```

接口方法说明：

> Read 将输入源中的lep(p)个字节读取到p中。它返回读取的字节数 n（0 <= n <= len(p)） 以及任何遇到的错误。

```go
func readFrom(reader io.Reader, num int) {
  p := make([]byte, num)
  n, err := reader.Read(p)  //从输入源读取数据到p中
  if n > 0 {
    return p[:], nil
  }
  return p, err
}
```

`readFrom`从任意地方读取数据，只要数据源实现了`io.Reader`接口。比如，我们可以从标准输入、文件、字符串等读取数据

```go
//标准输入：
data, err := readFrom(os.Stdin, 11)

//普通文件读取,file 是os.File 的实例
data, err := readFrom(file, 9)

//从字符串读取：
data, err := readFrom(strings.NewReader("Hello world"), 12)
```

## 2. Writer

接口的定义如下：

```go
type Writer interface {
    Write(p []byte) (n int, err error)  // p中有数据，将数据写入数据流
}
```

接口文件说明：

> Write 将 len(p) 个字节从 p 中写入到基本数据流中。

fmt的标准库中，有一组函数：Fprint/Fprintf/Fprintln，它们接收一个 io.Wrtier 类型参数（第一个参数），也就是说它们将数据格式化输出到 io.Writer 中。

```go
// 打印到前台
func Println(a... interface{}) {
  return Fprintln(os.stdout, a...)
}
```

## 3. 实现了io.Reader或io.Writer 接口的类型

- os.File 同时实现了 io.Reader 和 io.Writer
- strings.Reader 实现了 io.Reader
- bufio.Reader/Writer 分别实现了 io.Reader 和 io.Writer
- bytes.Buffer 同时实现了 io.Reader 和 io.Writer
- compress/gzip.Reader/Writer 分别实现了 io.Reader 和 io.Writer
- crypto/cipher.StreamReader/StreamWriter 分别实现了 io.Reader 和 io.Writer
- crypto/tls.Conn 同时实现了 io.Reader 和 io.Writer
- encoding/csv.Reader/Writer 分别实现了 io.Reader 和 io.Writer
- mime/multipart.Part 实现了 io.Reader
- net/conn 分别实现了 io.Reader 和 io.Writer(Conn接口定义了Read/Write)

## 4. Copy函数

```go
func Copy(dst Writer, src Reader) (written int64, err error)

io.Copy(os.Stdout, strings.NewReader("Go语言中文网"))
```

# 二、ioutil-方便的IO操作函数集

## ReadAll函数

一次性读取 io.Reader 中的数据

```go
func ReadAll(r io.Reader) ([]byte, error)
```

## ReadDir函数

读取目录并返回排好序的文件和子目录名

```go
func main() {
    dir := os.Args[1]
    listAll(dir,0)
}

func listAll(path string, curHier int){
    fileInfos, err := ioutil.ReadDir(path)
    if err != nil{fmt.Println(err); return}

    for _, info := range fileInfos{
        if info.IsDir(){
            for tmpHier := curHier; tmpHier > 0; tmpHier--{
                fmt.Printf("|\t")
            }
            fmt.Println(info.Name(),"\\")
            listAll(path + "/" + info.Name(),curHier + 1)
        }else{
            for tmpHier := curHier; tmpHier > 0; tmpHier--{
                fmt.Printf("|\t")
            }
            fmt.Println(info.Name())
        }
    }
}
```

## ReadFile/WriteFile函数

ReadFile 读取整个文件的内容，

```go
func ReadFile(filename string) ([]byte, error)
```

WriteFile 将data写入filename文件中，当文件不存在时会根据perm指定的权限进行创建一个,文件存在时会先清空文件内容。对于 perm 参数，我们一般可以指定为：0666，具体含义在os包中

```go
func WriteFile(filename string, data []byte, perm os.FileMode) error
```

# 三、fmt-格式化的io

```go
package main

import "fmt"

type user struct {
	name string
}

func main() {
	u := user{"tang"}
	//Printf 格式化输出
	fmt.Printf("%+v\n", u)     //格式化输出结构 -- {name:tang}
	fmt.Printf("%#v\n", u)       //输出值的 Go 语言表示方法 -- main.user{name:"tang"}
	fmt.Printf("%T\n", u)        //输出值的类型的 Go 语言表示 -- main.user
	fmt.Printf("%t\n", true)     //输出值的 true 或 false -- true
	fmt.Printf("%b\n", 10)     //二进制表示 -- 1010
	fmt.Printf("%c\n", 11111111) //数值对应的 Unicode 编码字符 -- 乱码显示
	fmt.Printf("%d\n", 10)       //十进制表示 -- 10
	fmt.Printf("%o\n", 8)        //八进制表示 -- 10
	fmt.Printf("%q\n", 22)       //转化为十六进制并附上单引号 -- '\x16'
	fmt.Printf("%x\n", 1223)     //十六进制表示，用a-f表示 -- 4c7
	fmt.Printf("%X\n", 1223)     //十六进制表示，用A-F表示 -- 4C7
	fmt.Printf("%U\n", 1233)     //Unicode表示 -- U+04D1
	fmt.Printf("%b\n", 12.34)    //无小数部分，两位指数的科学计数法 -- 6946802425218990p-49
	fmt.Printf("%e\n", 12.345)   //科学计数法，e表示 -- 1.234500e+01
	fmt.Printf("%E\n", 12.34455) //科学计数法，E表示 -- 1.234455E+01
	fmt.Printf("%f\n", 12.3456)  //有小数部分，无指数部分
	fmt.Printf("%g\n", 12.3456)  //根据实际情况采用%e或%f输出
	fmt.Printf("%G\n", 12.3456)  //根据实际情况采用%E或%f输出
	fmt.Printf("%s\n", "wqdew")  //直接输出字符串或者[]byte -- wqdew
	fmt.Printf("%q\n", "dedede") //双引号括起来的字符串 -- "dedede"
	fmt.Printf("%x\n", "abczxc") //每个字节用两字节十六进制表示，a-f表示 -- 6162637a7863
	fmt.Printf("%X\n", "asdzxc") //每个字节用两字节十六进制表示，A-F表示 -- 6173647A7863
	fmt.Printf("%p\n", 0x123)    //0x开头的十六进制数表示 -- %!p(int=291)
}

```

*如果格式化输出某种类型的值，只要它实现了 String() 方法，那么会调用 String() 方法进行处理。*

# 四、bufio - 缓存IO

**bufio.Reader** 结构包装了一个 io.Reader 对象，提供缓存功能，同时实现了 io.Reader 接口。

```go
    type Reader struct {
        buf          []byte        // 缓存
        rd           io.Reader    // 底层的io.Reader
        // r:从buf中读走的字节（偏移）；w:buf中填充内容的偏移；
        // w - r 是buf中可被读的长度（缓存数据的大小），也是Buffered()方法的返回值
        r, w         int
        err          error        // 读过程中遇到的错误
        lastByte     int        // 最后一次读到的字节（ReadByte/UnreadByte)
        lastRuneSize int        // 最后一次读到的Rune的大小 (ReadRune/UnreadRune)
    }
```

两个实例化 bufio.Reader 对象的函数：NewReader 和 NewReaderSize

```go
 func NewReader(rd io.Reader) *Reader {
        // 默认缓存大小：defaultBufSize=4096
        return NewReaderSize(rd, defaultBufSize)
    }


func NewReaderSize(rd io.Reader, size int) *Reader {
    // 已经是bufio.Reader类型，且缓存大小不小于 size，则直接返回
    b, ok := rd.(*Reader)
    if ok && len(b.buf) >= size {
      return b
    }
    // 缓存大小不会小于 minReadBufferSize （16字节）
    if size < minReadBufferSize {
      size = minReadBufferSize
    }
    // 构造一个bufio.Reader实例
    return &Reader{
      buf:          make([]byte, size),
      rd:           rd,
      lastByte:     -1,
      lastRuneSize: -1,
    }
}
```

**bufio.Writer** 结构包装了一个 io.Writer 对象，提供缓存功能，同时实现了 io.Writer 接口。

Writer 结构没有任何导出的字段，结构定义如下：

```go
    type Writer struct {
        err error        // 写过程中遇到的错误
        buf []byte        // 缓存
        n   int            // 当前缓存中的字节数
        wr  io.Writer    // 底层的 io.Writer 对象
    }
```

bufio 包提供了两个实例化 bufio.Writer 对象的函数：NewWriter 和 NewWriterSize。

```go
 func NewWriter(wr io.Writer) *Writer {
        // 默认缓存大小：defaultBufSize=4096
        return NewWriterSize(wr, defaultBufSize)
    }

    func NewWriterSize(wr io.Writer, size int) *Writer {
        // 已经是 bufio.Writer 类型，且缓存大小不小于 size，则直接返回
        b, ok := wr.(*Writer)
        if ok && len(b.buf) >= size {
            return b
        }
        if size <= 0 {
            size = defaultBufSize
        }
        return &Writer{
            buf: make([]byte, size),
            wr:  w,
        }
    }
```

