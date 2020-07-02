# go 语言安部署
​	go语言的安装可以根据平台下载安装包使用，也可以下载源码先编译在安装

## 一、安装

	1. 安装包下载：wget https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz
	2. 解压到安装目录：tar -zxvf go1.11.1.linux-amd64.tar.gz -C /usr/local
	3. 设置环境变量：export PATH=$PATH:$GO_INSTALL_DIR/go/bin

# 二、GOPATH与工作空间

 	1. go 命令依赖一个重要的环境变量：$GOPATH。这个目录下包括三个子目录：src、bin、package，
     ​	用来存放Go源码，Go的可运行文件，以及相应的编译之后的包文件。
 	2. GOPATH允许多个目录，当有多个目录时，用分隔符隔开，多个目录的时候Windows是分号，Linux系统是冒号，
     当有多个GOPATH时，默认会将go get的内容放在第一个目录下。
 	3. 代码目录结构规划：
     GOPATH下的src目录就是接下来开发程序的主要目录，所有的源码都是放在这个目录下面，
     一般的做法就是一个目录一个项目，例如: $GOPATH/src/mymath 表示mymath这个应用包或者可执应用**export GOPATH=/home/apple/mygo**

​		

## 三、go命令

 1. **go build**

     * 如果是普通包，执行go build后不会产生任何文件；如果需要在$GOPATH/pkg下生成
         相应的文件，需要执行go install
         * 如果是main包，执行go build后会在文件夹下生成可执行文件；执行go install 后，会在$GOPATH/bin下生成可执行文件

	2. **go run**

	3. **go vet**

	4. **go clean**

	5. **go fmt**
    格式化代码规范
    主要参数：
    ​	-w 把改写后的内容直接写入到文件中，而不是作为结果打印到标准输出。

	6. go get
    获取远程包的工具，目前go get支持多数开源社区(例如：github、googlecode、bitbucket、Launchpad)
    例如：go get github.com/astaxie/beedb
    实际上分成了两步操作：第一步是下载源码包，第二步是执行go install。

	7. go install
    这个命令在内部实际上分成了两步操作：第一步是生成结果文件(可执行文件或者.a包)，
    第二步会把编译好的结果移到$GOPATH/pkg或者$GOPATH/bin。

	8. **go doc**

    在终端打印出文档：go doc tar

    > godoc ：可启动一个服务程序，用于在浏览器查看本系统内所有的go语言源代码文档
    >
    > godoc -http=:6060 
    >
    > 如果想给包写一段文字量比较大的文档，可以在工程里包含一个叫作 doc.go 的文件，使用同样的包名，并把包的介绍使用注释加在包名声明之前。 