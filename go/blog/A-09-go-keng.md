#go 坑记录

## for-range

**遍历取不到所有元素的指针**

查看[go编译源码](https://github.com/golang/gofrontend/blob/e387439bfd24d5e142874b8e68e7039f74c744d7/go/statements.cc#L5501)可以了解到, for-range其实是语法糖，内部调用还是for循环，初始化会拷贝带遍历的列表（如array，slice，map），然后每次遍历的`v`都是对同一个元素的遍历赋值。也就是说如果直接对`v`取地址，最终只会拿到一个地址，而对应的值就是最后遍历的那个元素所附给`v`的值。所以输出都是3。

```go

func TestForRange(t *testing.T) {

	slice := []int{0, 1, 2, 3}
	myMap := make(map[int]*int)

	for index, value := range slice {
		myMap[index] = &value
	}
	fmt.Println("=====new map=====")
	for key, value := range myMap {
		fmt.Printf("map[%v]=%v\n", key, *value)
	}
}
/*
输出：
=====new map=====
map[0]=3
map[1]=3
map[2]=3
map[3]=3
*/
```

正解：

* 使用局部变量

  ```go
  for index, value := range slice {
      v := value
  	myMap[index] = &v
  }
  ```

* 直接索引获取原来的元素

  ```
  for index, value := range slice {
  	myMap[index] = &slice[index]
  }
  ```

  

