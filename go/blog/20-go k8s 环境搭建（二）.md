> Kubernetes 软件安装：

# Minikube

## 1. minikube简介

Minikube 是一个轻量级的Kubernetes实现，会在本机创建一台虚拟机，并部署一个只包含一个节点的简单集群。Minikube CLI提供了集群的基本引导操作，包括启动、停止、状态和删除。

Minikube支持以下Kubernetes功能：

- DNS
- NodePorts（可使用“minikube service”命令来管理）
- ConfigMaps和Secrets
- 仪表板（Dashboards，minikube dashboard）
- 容器运行时：Docker，*rkt*，*CRI-O*和*containerd*
- Enabling CNI（容器网络接口）
- Ingress
- LoadBalancer（负载均衡，可以使用“minikube tunnel”命令来启用）
- Multi-cluster（多集群，可以使用“minikube start -p <name>”命令来启用）
- Persistent Volumes
- RBAC
- 通过命令配置apiserver和kubelet

## 2. minikube 安装