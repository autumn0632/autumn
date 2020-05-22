# 一、什么是go module

## 1. 为什么需要依赖管理

最早的时候，Go所依赖的所有第三方库都放在`GOPATH`目录下面。这就导致了同一个库只能保存一个版本的代码。如果不同的项目依赖同一个第三方库的不同版本，应该如何解决？

## 2. 依赖管理解决工具

### 2.1 godep

Go语言从V1.5开始引入`vendor`模式。如果项目目录下面有`vendor`目录，那么go工具链会优先使用`vendor`目录下面的包进行编译、测试等。如果没有找到就去`$GOPATH/src`目录下边去找

### 2.2 go module

 `go module`是Go1.11版本之后官方推出的版本管理工具，并且从Go1.13版本开始，`go module`将是Go语言默认的依赖管理工具。 

> "module" != "package"
>
> `module`和`package`，也即“模块”和“包”，并不是等价的，而是“模块”包含“包”，“包”属于”模块“

# 二、go module详解

## 2.1 go module 相关属性

* 1个开关环境变量：`GO111MODULE`
* 5个辅助环境变量：`GOPROXY`、`GONOPROXY`、`GOSUMBD`、`GONOSUMBD`和`GOPRIVATE`
* 两个辅助概念：**Go module proxy**和**Go checksum database**
* 两个主要文件：**go.mod**和**go.sum**
* 一个主要管理命令：**go mod**
* 内置在几乎所有其它子命令中：go build、go get...

**1. GO111MODULE**

要启用`go module`支持首先要设置环境变量`GO111MODULE`，通过它可以开启或关闭模块支持，它有三个可选值：`off`、`on`、`auto`，默认值是`auto`。

1. `GO111MODULE=off`禁用模块支持，编译时会从`GOPATH`和`vendor`文件夹中查找包。
2. `GO111MODULE=on`启用模块支持，编译时会忽略`GOPATH`和`vendor`文件夹，只根据 `go.mod`下载依赖。
3. `GO111MODULE=auto`，当项目在`$GOPATH/src`外且项目根目录有`go.mod`文件时，开启模块支持。

**2. go.mod**

```go
module github.com/Q1mi/studygo/blogger

go 1.12

require (
    github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
    github.com/gin-gonic/gin v1.4.0
    github.com/go-sql-driver/mysql v1.4.1
    github.com/jmoiron/sqlx v1.2.0
    github.com/satori/go.uuid v1.2.0
    google.golang.org/appengine v1.6.1 // indirect
)
```

go.mod 是启用了 Go moduels 的项目所必须的最重要的文件，它描述了当前项目（也就是当前模块）的元信息，每一行都以一个动词开头，目前有以下 5 个动词:

- module：用于定义当前项目的模块路径。 // 包名
- go：用于设置预期的 Go 版本。
- require： 用来定义依赖包及版本 。
- exclude：用于从使用中排除一个特定的模块版本。
- replace：用于将一个模块版本替换为另外一个模块版本。// 在国内访问golang.org/x的各个包都需要翻墙，此时可以在go.mod中使用replace替换成github上对应的库。                 

**3. go.sum**

 go.sum 详细罗列了当前项目直接或间接依赖的所有模块版本，并写明了那些模块版本的 SHA-256 哈希值以备 Go 在今后的操作中保证项目所依赖的那些模块版本不会被篡改。 

**4. GOPROXY**

 go get命令默认情况下，无论是在gopath mode还是module-aware mode，都是直接从vcs服务(比如github、gitlab等)下载module的。但是Go 1.11中，我们可以通过设置GOPROXY环境变量来做一些改变：让Go命令从其他地方下载module。 

 Go1.13之后`GOPROXY`默认值为`https://proxy.golang.org`，在国内是无法访问的， 推荐设置成`https://goproxy.cn`

```shell
go env -w GOPROXY=https://goproxy.cn,direct
```

**5. GOPRIVATE/GONOPROXY/GONOSUMBD**

* 这三个环境变量都是用在当前项目依赖了私有模块，   GOPRIVATE 较为特殊，它的值将作为 GONOPROXY 和 GONOSUMDB 的默认值，所以建议的最佳姿势是只是用 GOPRIVATE。 

* GOPRIVATE=GOPRIVATE，所有以“git.gbcom.com”为前缀的模块版本都将不经过Go module proxy和Go checksun database。

  ```
  go env -w GOPRIVATE=git.gbcom.com
  ```

  

**6. go mod 命令**

常用的`go mod`命令如下：

```
go mod download   下载go.mod文件依赖的module到本地cache（默认为$GOPATH/pkg/mod目录）
go mod edit        编辑go.mod文件
go mod graph       打印模块依赖图
go mod init        初始化当前文件夹, 创建go.mod文件
go mod tidy        增加缺少的module，删除无用的module
go mod vendor      将依赖复制到vendor下
go mod verify      校验依赖
go mod why         解释为什么需要依赖
```

**7. global caching**

用来存放Go module的全局缓存数据：

*  同一个模块版本的数据只缓存一份，所有其他模块共享使用。 
*  目前所有模块版本数据均缓存在 `$GOPATH/pkg/mod`和  `$GOPATH/pkg/sum` 下， 
*  以使用 `go clean-modcache` 清理所有已缓存的模块版本数据。 

# 三、如何使用go module

* 将Go版本升级到V1.11及以上（推荐V1.13及以上）

* 在项目目录下执行`go mod init` ，生成`go.mod`文件 

* 执行`go mod tidy`，拉取必需模块，删除不用模块

* 不建议使用 `go mod vendor`，因为 Go modules 正在淡化 Vendor 的概念，很有可能 Go2 就去掉了。 

* 公共库导入：

  ```shell
  # 假设有两个项目blog 和 article，blog 是应用的入口，main 所在位置，article 是一个公共库，提供工具类，供其它项目引用
  ├─article
  │      article.go
  │      go.mod
  │
  ├─blog
  │      go.mod
  │      main.go
  
  1. article 中必须有go.mod 文件（即使用module管理依赖），内容如下
     	module github.com/article
  
     	go 1.13
  2. blog中，修改go.mod，手动添加article项目依赖，内容如下
    	go 1.13
  
      require github.com/article v0.0.0-incompatible
  
  	replace github.com/article => ../article
  	# github.com/article 格式：在 go1.13 中， go module 名称规范要求路径的第一部分必须满足域名规范，可以替换成公司内部域名
  	# replace 的第二个参数指定了不从远程获取，而是本地某个路径下的模块替换 github.com/article
  3. 引用article项目
  	import "github.com/article"
  ```


* 依赖库升降级

  > - 查看版本历史：
  >
  >   go list -m -versions github.com/gin-gonic/gin
  >
  > - 使用go mod 修改依赖文件：
  >
  >   go mod edit -require="github.com/gin-gonic/gin@v1.1.4"
  >
  > - 下载更新依赖
  >
  >   go tidy

# 四、go module 使用问题记录

1. **路径问题**
   * Go Module功能开启后，下载的包将存放与$GOPATH/pkg/mod路径 
   *  环境变量GOPATH不再用于解析imports包路径，即原有的GOPATH/src/下的包，通过import是找不到的
   * 项目目录可放在任意目录下，并非一定在$GOPATH/src下。

2. **判断项目是否启用了Go Modules**
  
   * 一个项目中只要包含了`go.mod`文件，且环境中`GO111MODULE`不为off，那么go就会为这个项目启动Go Modules
   * 若当前项目中没有`go.mod`文件，且环境中`GO111MODULE`不为on，那么每一次构建代码时Go都会从头推算并拉取所需版本，但是并不会自动生成go.mod文件
3. **更新现有的模块**
   * go get -u
     * 只会更新主要模块，忽略了单元测试
   * go get -u ./..
   * go get -u all 
     * 更新所有模块，推荐使用

4. **本地导入，导入的包中不能再有本地的包**

   ```
   ├─projectA
   │      main.go
   │      go.mod
   ├─projectB
   |      projectb.go
   │      go.mod
   │      
   ├─projectc
   |      projectc.go
   │      go.mod
   ```

   在上面三个项目中，main.go引用projectb.go中的函数，同时projectb.go中引用projectc.go中的函数。
   A项目go.mod内容如下

   ```
       package main    go 1.13    require github.com/projectB v0.0.0-incompatible	replace github.com/projectB => ../projectB 
   ```

   B项目go.mod 内容如下

   ```
       package github.com/projectB    go 1.13    require github.com/projectC v0.0.0-incompatible	replace github.com/projectC => ../projectC 
   ```

   在编译时，B项目中的replace并不会生效，会报如下错误

   ```
   go: git.gbcom.com/projectB@v0.0.0 requires        git.gbcom.com/projectC@v0.0.0: unrecognized import path "github.com/projectC" (https fetch: Get https://github/projectC?go-get=1: dial tcp x.x.x.x:443: i/o timeout)
   ```

   **也就是说本地导入，导入的包不能再有本地的导入，导入本地的再有本地的导入直接失效，即使你用了replace，程序也会无视，而直接去代理服务器傻乎乎的下载你显示制定的本地包。。。**
   
   解决办法：直接push到仓库中，从使用仓库中导入。