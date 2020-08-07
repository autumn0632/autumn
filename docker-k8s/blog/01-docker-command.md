

## 一、常用的docker命令

### 1. 启动容器

* docker run -it --rm ubuntu:18.04 bash
  > * `-it`：`i`表示交互式操作；`t`表示分配终端
  > * `--rm`：容器推出后将其删除。默认情况下，容器推出后不会删除
  > * `ubuntu:18.04`：用`ubuntu:18.04`镜像为基础来启动容器。  
  > * `bash`：镜像后面跟的是命令

* docker run --name webserver -d -p 80:80 nginx

  > * `--name`：给容器指定名字
  > * `-d`：后台执行容器
  > * `-p`：端口映射

* docker exec -it webserver bash

### 2. 查看镜像/容器

**镜像**

* docker images ls：显示顶层镜像

  > * `-a`：列出所有镜像
  > * `-q`：只返回镜像ID
  > * `--format` "{{.ID}}: {{.Repository}}"：按指定格式返回

* docker image ls -f dangling=true：虚悬镜像查看

* docker image prune：虚悬镜像删除

* docker image rm：镜像删除

**镜像**

* docker container ls：容器查看

### 3. 创建镜像

> docker build 工作原理：
>
> ​		Docker 在运行时分为 Docker 引擎（ 也就是服务端守护进程） 和客户端工具。 Docker 的引擎提供了一组 REST API，被称为 Docker Remote API， 而如 docker 命令这样的客户端工具， 则是通过这组 API 与Docker 引擎交互， 从而完成各种功能。 
>
> ​		而 docker build 命令构建镜像， 其实并非在本地构建， 而是在服务端， 也就是 Docker 引擎
> 中构建的。当构建的时候， 用户会指定构建镜像上下文的路径， docker build 命令得知这个路径后， 会将路径下的所有内容打包， 然后上传给 Docker 引擎。 这样 Docker 引擎收到这个上下文包后， 展开就会获得构建镜像所需的一切文件。  

* docker build -t nginx:v3 .  

  > * `-t`：tag标记、
  > * `.`：



## 二、Dockerfile 镜像定制

Dockerfile 是一个文本文件， 其内包含了一条条的 指令(Instruction)， 每一条指令
构建一层， 因此每一条指令的内容， 就是描述该层应当如何构建。  

**FROM 指定镜像醒**

---

指定基础镜像。Dockerfile 中 FROM 是必备的指令，并且必须是第一条指令 。

`scratch`这个镜像是虚拟的概念， 并不实际存在， 它表示一个空白的镜像。  意味着不以任何镜像为基础， 接下来所写的指令将作为镜像第一层开始存在。  

**RUN执行命令**

---

执行命令行命令。格式有两种：

* shell格式：`RUN <命令>`

  > RUN echo '<h1>Hello, Docker!</h1>' > /usr/share/nginx/html/index.html  

* exec 格式：`RUN ["可执行文件", "参数1", "参数2"]`

每一个`RUN`命令，就会新建一层。因此要注意使用`RUN`的方法:

* 避免冗余的`RUN`指令，将多个指令写在一个`RUN`指令中

* 镜像构建时， 一定要确保每一层只添加真正需要添加的东西， 任何无关的东西都应该清理掉。  

  > FROM debian:stretch
  > RUN buildDeps='gcc libc6-dev make wget' \
  > 		&& apt-get update \
  > 		&& apt-get install -y $buildDeps \
  > 		&& wget -O redis.tar.gz "http://download.redis.io/releases/redis-5.0.3.tar.gz" \
  > 		&& mkdir -p /usr/src/redis \
  > 		&& tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1 \
  > 		&& make -C /usr/src/redis \
  > 		&& make -C /usr/src/redis install \
  > 		&& rm -rf /var/lib/apt/lists/* \
  > 		&& rm redis.tar.gz \
  > 		&& rm -r /usr/src/redis \
  > 		&& apt-get purge -y --auto-remove $buildDeps  

  Dockerfile 支持 Shell 类的行尾添加 \ 的命令换行方式， 以及行首 # 进行注释的格式。  

**COPY复制文件**

---

`COPY` 指令将从构建上下文目录中` <源路径> `的文件/目录复制到新的一层的镜像内的` <目标路径> `位置 ， 有两种格式：

* `COPY [--chown=<user>:<group>] <源路径>... <目标路径>`

* `COPY [--chown=<user>:<group>] ["<源路径1>",... "<目标路径>"]`

使用 COPY 指令， 源文件的各种元数据都会保留。 比如读、 写、 执行权限、 文件变更时间等。有 时候还可以加上 `--chown=<user>:<group> `选项来改变文件的所属用户及所属组。  



**CMD容器启动命令**

---

命令格式：

* shell 格式： `CMD <命令>`
* exec 格式： `CMD ["可执行文件", "参数1", "参数2"...]`
* 参数列表格式： `CMD ["参数1", "参数2"...]` 。 在指定了 `ENTRYPOINT` 指令后， 用 `CMD` 指定具体的参数。      

`CMD`与`RUN`的区别：

`CMD` 指令就是用于指定默认的容器主进程的启动命令。`RUN`用于执行普通命令。

> **程序的前台与后台**
>
> 对于容器而言， 其启动程序就是容器应用进程， 容器就是为了主进程而存在的， 主进程退出， 容器就失去了存在的意义， 从而退出， 其它辅助进程不是它需要关心的东西。  所以，容器内的程序都要以前台运行，例如：
>
> ```
> CMD ["nginx", "-g", "daemon off;"]
> ```
>
>   

**ADD高级复制**

---

`ADD`指令和`COPY`的格式和性质基本一致。 但是在 COPY 基础上增加了一些功能。  

* 源路径可以是`url`，可以自动下载文件到目标路径下。不推荐使用
* 源路径为一个`tar`压缩文件（`gzip`/`bzip2`/`xz`）， 将会自动解压缩到目标路径下。

> `ADD`与`COPY`的使用原则：
>
> 所有的文件复制均使用 `COPY` 指令， 仅在需要自动解压缩的场合使用 `ADD` 。  

**ENTRYPOINT 入口点 **

****

`ENTRYPOINT` 的格式和 `RUN` 指令格式一样， 分为 exec 格式和 shell 格式。`ENTRYPOINT` 的目的和 `CMD` 一样， 都是在指定容器启动程序及参数。  

当指定了 `ENTRYPOINT` 后， `CMD` 的含义就发生了改变， 不再是直接的运行其命令， 而是将 `CMD` 的内容作为参数传给 ENTRYPOINT 指令  

**ENV 设置环境变量**

---

这个指定就是设置环境变量，如论是其它指令，还是运行时的应用，都可以直接使用这里定义的环境变量。格式如下：

* `ENV <key> <value>`
* `ENV <key1>=<value1> <key2>=<value2>...`

**VOLUME 定义匿名卷  **

---

