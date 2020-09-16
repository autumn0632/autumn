# Docker网络基础

Docker 本身的技术依赖于近年来LInux内核虚拟化技术的发展。因此有必要深入了解Docker背后的网络原理和基础知识。

## 一、Linux网络基础

docker涉及到的主要的网络技术有**网络命名空间（NetWork Namespace）**、**Veth设备对**、**网桥**、**iptables**和路由。

### 1.1 网络命名空间

#### 1.1.1 命名空间

​		Docker 和虚拟机技术一样，从操作系统级上实现了资源的隔离，它本质上是宿主机上的进程（容器进程），所以资源隔离主要就是指进程资源的隔离。实现资源隔离的核心技术就是 Linux namespace。

​		Linux namespace 实现了 6 项资源隔离，基本上涵盖了一个小型操作系统的运行要素，包括主机名、用户权限、文件系统、网络、进程号、进程间通信。

| namespace | 系统调用参数  | 隔离内容                   | 内核版本 |
| --------- | ------------- | -------------------------- | -------- |
| UTS       | CLONE_NEWUTS  | 主机名和域名               | 2.6.19   |
| IPC       | CLONE_NEWIPC  | 信号量、消息队列和共享内存 | 2.6.19   |
| PID       | CLONE_NEWPID  | 进程编号                   | 2.6.24   |
| Network   | CLONE_NEWNET  | 网络设备、网络栈、端口等   | 2.6.29   |
| Mount     | CLONE_NEWNS   | 挂载点（文件系统）         | 2.4.19   |
| User      | CLONE_NEWUSER | 用户和用户组               | 3.8      |

这 6 项资源隔离分别对应 6 种系统调用，通过传入上表中的参数，调用 clone() 函数来完成。

```
clone(int (*child_func)(void *), void *child_stack, int flags, void *arg);
```

#### 1.1.1.2 网络命名空间

### 1.2 Veth设备对

### 1.3 网桥

### 1.4 iptables

### 1.5 路由



