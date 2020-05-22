1. 预期状态

   用户通过创建对象来告诉Kubernetes系统自己希望集群里的工作负载是什么样子；这就是集群的预期状态。

2. Kubernetes对象：

   任意时刻系统里的实体状态是由Kubernetes对象来表示的。用户可以直接和Kubernetes对象交互，而不用和容器直接交互。基本的Kubernetes对象如下：

   * `pod`：是节点上的最小的部署单元。它是一组必须一起运行的容器。一般来说，但不是必须的，Pod通常包含一个容器。多个容器时，一定是超亲密关系的，共享某些资源。
   * `service`：定义一组逻辑上存在关系的Pod以及访问它们的相关策略。
   * `Volume`：定义Pod里所有容器都能访问的目录。
   * `NameSpace`：由物理集群支撑的虚拟集群。

3. Kubernetes控制器：

   控制器是基于Kubernetes的基础对象构建并提供额外的特性，包括：

   * `Replicaset`：确保在给定时间运行着特定数量的Pod副本
   * `Deployment`：确保在给定时间运行着特定数量的Pod副本
   * `Statefulset`：用来控制部署顺序以及卷的访问等等。
   * `Daemonset`：用来在集群的所有节点或者特定节点运行Pod的拷贝。
   * `Job`：用来执行一些任务并且在成功完成工作之后或者在给定时间之后退出。

4. Kubernetes控制平面：

   控制平面的工作就是确保集群的当前状态和用户的预期状态一致。要实现这一目标，Kubernetes会自动执行一系列任务——比如，启动或者重启容器，扩展某个给定应用的副本数量，等等。

   * **Kubernetes Master**：

     作为Kubernetes控制平面的一部分，Kubernetes master的工作是持续维护集群的预期状态。

     * `kube-apiserver`：整个集群的单点管理点。API server实现了RESTful的接口，用于和工具以及库函数的通信。`kubectl`命令直接和API server交互。
     * `kube-controller-manager`：通过管理不同类型的控制器来规范集群的状态。
     * `kube-scheduler`：在集群里的可用节点上调度工作负载。

   * **Kubernetes node**：

     节点是集群里运行工作负载的工作机器（物理的VM或者物理服务器等）。这些节点由Kubernetes master控制，并且通过持续监控来维护应用程序的预期状态。集群里的每个Kubernetes节点运行两个进程：

     * `kubelet`：是节点和Kubernetes Master之间的通信接口。
     * `kube-proxy`：是网络路由，它将每个节点上通过Kubernetes API定义的服务暴露出去。它还能够执行简单的TCP和UDP的流转发。

