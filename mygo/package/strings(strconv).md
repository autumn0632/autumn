# strings

strings包提供了很多操作字符串的简单函数

> 字符串常见操作有：
>
> - 字符串长度
> - 求子串
> - 是否存在某个字符或子串
> - 子串出现的次数（字符串匹配）
> - 字符串分割（切分）为[]string
> - 字符串是否有某个前缀或后缀
> - 字符或子串在字符串中首次出现的位置或最后一次出现的位置
> - 通过某个字符串将[]string连接起来
> - 字符串重复几次
> - 字符串中子串替换
> - 大小写转换
> - ...

## 是否存在某个字符或子串

```go
// 子串substr在s中，返回true
func Contains(s, substr string) bool
// chars中任何一个Unicode代码点在s中，返回true
func ContainsAny(s, chars string) bool
// Unicode代码点r在s中，返回true
func ContainsRune(s string, r rune) bool
```

查看这三个函数的源码，发现它们只是调用了相应的Index函数（子串出现的位置），然后和 0 作比较返回true或fale。如，Contains：

```go
func Contains(s, substr string) bool {
	return Index(s, substr) >= 0
}
```



## 子串出现次数(字符串匹配)

在Go中，查找子串出现次数即字符串模式匹配，实现的是Rabin-Karp算法。Count 函数的签名如下：

```
func Count(s, sep string) int
```

## 字符串分割为[]string

该包提供了六个三组分割函数：Fields 和 FieldsFunc、Split 和 SplitAfter、SplitN 和 SplitAfterN。

**Fields和fieldsFunc**

这两个函数签名如下：

```go
func Fields(s string) []string
func FieldsFunc(s string, f func(rune) bool) []string
```

Fields 用一个或多个连续的空格分隔字符串 s，返回子字符串的数组（slice）。

**Split 和 SplitAfter、 SplitN 和 SplitAfterN**

这四个函数都是通过同一个内部函数genSplit 来实现的。它们的函数签名及其实现：

```go
func Split(s, sep string) []string { return genSplit(s, sep, 0, -1) }
func SplitAfter(s, sep string) []string { return genSplit(s, sep, len(sep), -1) }
func SplitN(s, sep string, n int) []string { return genSplit(s, sep, 0, n) }
func SplitAfterN(s, sep string, n int) []string { return genSplit(s, sep, len(sep), n) }
```

Split 和 SplitAfter 区别：Split 会将 s 中的 sep 去掉，而 SplitAfter 会保留 sep。

```go
fmt.Printf("%q\n", strings.Split("foo,bar,baz", ","))
fmt.Printf("%q\n", strings.SplitAfter("foo,bar,baz", ","))

//输出：
["foo" "bar" "baz"]
["foo," "bar," "baz"]
```

 带N 的方法可以通过最后一个参数 n 控制返回的结果中的 slice 中的元素个数，当 n < 0 时，返回所有的子字符串；当 n == 0 时，返回的结果是 nil；当 n > 0 时，表示返回的 slice 中最多只有 n 个元素，其中，最后一个元素不会分割。

```go
fmt.Printf("%q\n", strings.SplitN("foo,bar,baz", ",", 2))

//输出：
["foo" "bar,baz"]
```



## 字符串是否有某个前缀或后缀

```
// s 中是否以 prefix 开始
func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}
// s 中是否以 suffix 结尾
func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
```



## 字符或子串在字符串中出现的位置

```
// 在 s 中查找 sep 的第一次出现，返回第一次出现的索引
func Index(s, sep string) int
// chars中任何一个Unicode代码点在s中首次出现的位置
func IndexAny(s, chars string) int
// 查找字符 c 在 s 中第一次出现的位置，其中 c 满足 f(c) 返回 true
func IndexFunc(s string, f func(rune) bool) int
// Unicode 代码点 r 在 s 中第一次出现的位置
func IndexRune(s string, r rune) int

// 有三个对应的查找最后一次出现的位置
func LastIndex(s, sep string) int
func LastIndexAny(s, chars string) int
func LastIndexFunc(s string, f func(rune) bool) int
```



# strconv

*strconv* 包提供了基本数据类型和字符串之间的转换。



# regexp

regexp包提供了正则表达式的功能，**regexp/syntax**子包进行正则表达式解析



# unicode

go代码使用UTF-8编码（且不能带BOM），同时标识符支持Unicode字符。在标准库 *unicode* 包及其子包 utf8、utf16中，提供了对 Unicode 相关编码、解码的支持，同时提供了测试 Unicode 码点（Unicode code points）属性的功能。