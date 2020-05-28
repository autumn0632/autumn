

Go 语言中，为了方便开发者使用，将 IO 操作封装在了如下几个包中：

- **io**：为 IO 原语（I/O primitives）提供基本的接口
- **io/ioutil**：封装一些实用的 I/O 函数
- **fmt**：实现格式化 I/O，类似 C 语言中的 printf 和 scanf
- **bufio**：实现带缓冲I/O

# io

io 包为 I/O 原语提供了基本的接口。在 io 包中最重要的是两个接口：**Reader** 和**Writer** 接口。所提到的各种 IO 包，都跟这两个接口有关，也就是说，只要满足这两个接口，它就可以使用 IO 包的功能。

**Reader 接口**

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

```

Read 将 len(p) 个字节**读取到 p 中**。它返回读取的字节数 n（0 <= n <= len(p)） 以及任何遇到的错误。Reader 接口的方法集只包含一个 Read 方法，因此，所有实现了 Read 方法的类型都满足 io.Reader 接口，也就是说，在所有需要 io.Reader 的地方，可以传递实现了 Read() 方法的类型的实例。

my：如果一个数据类型实现了Reader接口，即可以从这个数据类型中读取数据，读到参数p中，并返回字节数以及错误。

```go
//示例:从reader中读取数据并返回
func ReadFrom(reader io.Reader, num int) ([]byte, error) {
    p := make([]byte, num)
    
    n, err := reader.Read(p)
    if n > 0{
        return p[:n],nil
    }
    return p, err    
}

//ReadFrom 可以从任意的地方读取数据，只要来源实现了 io.Reader 接口。比如，我们可以从标准输入、文件、字符串等读取数据，

//从标准输入读
data, err := ReadFrom(os.Stdin, 11)
//从普通文件读，file为os.File实例
data, err := ReadFrom(file, 9)
//从字符串读取
data, err := ReadFrom(string.NewReader("hello"), 6)
```



**Writer 接口**

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

Write 将 len(p) 个字节从 p 中**写入到基本数据流中。**

my：如果一个数据类型实现了Writer接口，即可以往这个数据类型中写数据。



```
实现了Reader和Writer接口的类型列表：
os.File 同时实现了 io.Reader 和 io.Writer
strings.Reader 实现了 io.Reader
bufio.Reader/Writer 分别实现了 io.Reader 和 io.Writer
bytes.Buffer 同时实现了 io.Reader 和 io.Writer
bytes.Reader 实现了 io.Reader
compress/gzip.Reader/Writer 分别实现了 io.Reader 和 io.Writer
crypto/cipher.StreamReader/StreamWriter 分别实现了 io.Reader 和 io.Writer
crypto/tls.Conn 同时实现了 io.Reader 和 io.Writer
encoding/csv.Reader/Writer 分别实现了 io.Reader 和 io.Writer
mime/multipart.Part 实现了 io.Reader
net/conn 分别实现了 io.Reader 和 io.Writer(Conn接口定义了Read/Write)

```

# ioutil-方便的IO操作函数集

虽然 io 包提供了不少类型、方法和函数，但有时候使用起来不是那么方便。比如**读取一个文件中的所有内容**。为此，标准库中提供了一些常用、方便的IO操作函数。

## ReadAll函数

用来从io.Reader中一次读取所有数据。

```
func ReadAll(r io.Reader) ([]byte, error)
```

该函数源码实现是通过bytes.Buffer中的ReadFrom来实现读取所有数据的。该函数成功调用后会返回err == nil而不是err == EOF。

## ReadDir函数

读取目录并返回排好序的文件和子目录名（[]os.FileInfo）

```go
//笔试题：输出某目录下的所有文件
func main() {
    dir := os.Args[1]
    listAll(dir, 0)
}

func listAll(path string, curHier int) {
    fileInfos, err := ioutil.ReadDir(path)
	if err != nil{ 
        fmt.Println(err)
        return
    }
    
    for _, info := range fileInfos {
        if info.IsDir() {
            for tmpHier := curHier; tmpHier > 0; tmpHier--{
				fmt.Printf("|\t")
			}
			fmt.Println(info.Name(),"\\")
			listAll(path + "/" + info.Name(),curHier + 1)
        }else {
            for tmpHier := curHier; tmpHier > 0; tmpHier--{
				fmt.Printf("|\t")
			}
			fmt.Println(info.Name())
        }
    }
}
```

## ReadFile 和 WriteFile函数

ReadFile 读取整个文件的内容。ReadFile的实现和ReadAll类似，不过，ReadFile会先判断文件的大小，给bytes.Buffer一个预定义容量，避免额外分配内存。函数签名如下：

```
func ReadFile(filename string) ([]byte, error)
```

成功的调用返回的err为nil而非EOF。因为本函数定义为读取整个文件，它不会将读取返回的EOF视为应报告的错误。

WriteFile 函数的签名如下：

```
func WriteFile(filename string, data []byte, perm os.FileMode) error
```

WriteFile 将data写入filename文件中，当文件不存在时会根据perm指定的权限进行创建一个,文件存在时会先清空文件内容。对于perm参数，我们一般可以指定为：0666