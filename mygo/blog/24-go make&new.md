* **概述**

  > - `make` 的作用是初始化内置的数据结构，即切片、哈希表和 Channel
  >
  >   ```go
  >   slice := make([]int, 0, 100)
  >   hash := make(map[int]bool, 10)
  >   ch := make(chan int, 5)
  >   ```
  >
  > - `new` 的作用是根据传入的类型在堆上分配一片内存空间并返回指向这片内存空间的指针；
  >
  >   ```go
  >   i := new(int)
  >   
  >   var v int
  >   i := &v
  >   ```
  >
  >   

  

