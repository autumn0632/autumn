1. 编译64位的Linux系统可执行程序：

   ```
   GOOS=linux GOARCH=amd64 go build hello.go
   
   1. GOOS：目标操作系统
   2. GOARCH：目标操作系统的架构
   
   ```

| OS             | ARCH              |
| -------------- | ----------------- |
| linux          | 386 / amd64 / arm |
| darwin//mac os | 386 / amd64       |
| freebsd        | 386 / amd64       |
| windows        | 386 / amd64       |

