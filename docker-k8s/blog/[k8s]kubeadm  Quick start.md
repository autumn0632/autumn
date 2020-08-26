# kubeadm 快速部署kubernetes集群

​		Kubernetes集群环境包括Master节点/node节点，Master 节点主要包含了三个Kubernetes项目中最最最重要的组件：`apiserver`，`scheduler`，`controller-manager`,

* `apiserver`：提供了管理集群的API接口

* `scheduler`：负责分配调度Pod到集群内的node节点

* `controller-manager`：由一系列的控制器组成，通过apiserver监控整个集群的状态

  ​	kuberadm是kubernetes原生的部署工具，简单快捷方便，适于新手搭建学习。但这种方式不建议直接上生产环境。



## Master节点

### 1.0 关闭交换空间

执行`swapoff -a`， 并在`/etc/fstab`中删除对`swap`的加载，然后重启服务器

### 1.1关闭防火墙

```bash
systemctl stop firewalld && systemctl disable firewalld
```

### 1.2 检查selinux是否关闭

```csharp
[root@kubernetes01 ~]# setenforce 0
setenforce: SELinux is disabled
```

### 1.3 提前处理路由问题

```bash
cat > /etc/sysctl.d/k8s.conf << EOF
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1    
vm.swappiness=0
EOF
sysctl --system
```

### 1.4 安装docker-ce

```bash
# yum安装docekr-ce，版本是v18.06.1

# docker版本卸载
yum remove -y docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine


yum -y install yum-utils device-mapper-persistent-data lvm2
yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
yum -y install docker-ce-18.06.1.ce
systemctl enable docker.service ;systemctl start docker.service 
docker --version
```

### 1.5 安装kubelet、kubectl、kubeadm

```bash
# 1. 配置为阿里云yum源
cat > /etc/yum.repos.d/kubernetes.repo << EOF
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
EOF
# 2. 安装key文件
wget https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
rpm -import rpm-package-key.gpg
# 3. 卸载旧版本
yum remove -y kubectl kubelet kubeadm
# 4. 安装
yum install -y kubelet-1.12.1 kubectl-1.12.1 kubeadm-1.12.1 kubernetes-cni-0.6.0
# 安装完成后，慎用 yum update -y 进行升级
```

### 1.6 下载kubernetes相关组件的docker镜像

```shell
# 由于网络原因, 使用如下脚本下载方式
[root@kubernetes01 ~]# cat pull_k8s_images.sh 
#!/bin/bash
images=(kube-proxy:v1.12.1 kube-scheduler:v1.12.1 kube-controller-manager:v1.12.1
kube-apiserver:v1.12.1
etcd:3.2.24 coredns:1.2.2 pause:3.1 )
for imageName in ${images[@]} ; do
docker pull anjia0532/google-containers.${imageName}
docker tag anjia0532/google-containers.$imageName k8s.gcr.io/$imageName
docker rmi anjia0532/google-containers.$imageName
done

$ sh pull_k8s_images.sh 

$ docker images 
REPOSITORY                           TAG                 IMAGE ID            CREATED             SIZE
k8s.gcr.io/kube-proxy                v1.12.1             61afff57f010        5 months ago        96.6MB
k8s.gcr.io/kube-apiserver            v1.12.1             dcb029b5e3ad        5 months ago        194MB
k8s.gcr.io/kube-scheduler            v1.12.1             d773ad20fd80        5 months ago        58.3MB
k8s.gcr.io/kube-controller-manager   v1.12.1             aa2dd57c7329        5 months ago        164MB
k8s.gcr.io/etcd                      3.2.24              3cab8e1b9802        6 months ago        220MB
k8s.gcr.io/coredns                   1.2.2               367cdc8433a4        7 months ago        39.2MB
k8s.gcr.io/pause                     3.1                 da86e6ba6ca1        15 months ago       742kB

```

### 1.7 使用kubeadm部署kubernetes集群master节点

```bash
[root@kubernetes01 ~]# kubeadm init --kubernetes-version=v1.12.1 
preflight核验没有问题后过一段时间，看到这样的提示算是完成了对Kubernetes Master节点的部署。
Your Kubernetes master has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

You can now join any number of machines by running the following on each node
as root:

  kubeadm join 10.5.0.206:6443 --token bh3pih.cuir6xpjl7zn7pf2 --discovery-token-ca-cert-hash sha256:ae00fc1ad4a680c01be4deaae6f6e4cf554867664bc5c16e0b3f98d4f2adcf2c

在开始使用之前，需要以常规用户身份运行以下命令: 上面那段英文中有说明注意查看！因为Kubernetes集群默认是需要加密访问的！
so执行这段命令👇
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

### 1.8 健康状态检查

```bash
# 1.查看主要组件的健康状态
[root@kubernetes01 ~]# kubectl get cs
NAME                 STATUS    MESSAGE              ERROR
scheduler            Healthy   ok                   
controller-manager   Healthy   ok                   
etcd-0               Healthy   {"health": "true"}   
# 2.查看master节点状态
[root@kubernetes01 ~]# kubectl get nodes
NAME           STATUS     ROLES    AGE     VERSION
kubernetes01   NotReady   master   4m15s   v1.12.1
```

### 1.9 部署网路插件weave

```bash
[root@kubernetes01 ~]# kubectl apply -f https://git.io/weave-kube-1.6
serviceaccount/weave-net created
serviceaccount/weave-net created
clusterrole.rbac.authorization.k8s.io/weave-net created
clusterrolebinding.rbac.authorization.k8s.io/weave-net created
role.rbac.authorization.k8s.io/weave-net created
rolebinding.rbac.authorization.k8s.io/weave-net created
daemonset.extensions/weave-net created
等一会儿，查看Master节点状态，STATUS已经变了，这是因为部署的网络组件生效了
[root@kubernetes01 ~]# kubectl get nodes
NAME                STATUS   ROLES    AGE   VERSION
kubernetes-master   Ready    master   21m   v1.12.1
```

### 1.10 部署可视化插件

需要注意的是，由于 Dashboard 是一个 Web Server，很多人经常会在自己的公有云上无意地暴露 Dashboard 的端口，从而造成安全隐患。所以，1.7 版本之后的 Dashboard 项目部署完成后，默认只能通过 Proxy 的方式在本地访问。访问方式待研究

```bash
# 1.获取可视化插件docker镜像，修改tag
docker pull anjia0532/google-containers.kubernetes-dashboard-amd64:v1.10.0
docker tag  anjia0532/google-containers.kubernetes-dashboard-amd64:v1.10.0 k8s.gcr.io/kubernetes-dashboard-amd64:v1.10.1
docker rmi  anjia0532/google-containers.kubernetes-dashboard-amd64:v1.10.0 
# 2.获取并修改可视化插件YAML文件的最后部分，便于后期通过token登陆可视化页面，这里需要特别注意的是暴露了30001端口，这如果在生产环境是极不安全的！
[root@kubernetes01 ~]# wget https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml
[root@kubernetes01 ~]# tail -n 20 kubernetes-dashboard.yaml
        effect: NoSchedule

---
# ------------------- Dashboard Service ------------------- #

kind: Service
apiVersion: v1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kube-system
spec:
  type: NodePort
  ports:
    - port: 443
      targetPort: 8443
      nodePort: 30001
  selector:
    k8s-app: kubernetes-dashboard
# 3.部署可视化插件
[root@kubernetes01 ~]# kubectl apply -f kubernetes-dashboard.yaml
secret/kubernetes-dashboard-certs created
serviceaccount/kubernetes-dashboard created
role.rbac.authorization.k8s.io/kubernetes-dashboard-minimal created
rolebinding.rbac.authorization.k8s.io/kubernetes-dashboard-minimal created
deployment.apps/kubernetes-dashboard created
service/kubernetes-dashboard configured
# 4.查看可视化插件对应的Pod状态
[root@kubernetes01 ~]# kubectl get pods -n kube-system |  grep dash
kubernetes-dashboard-65c76f6c97-f29nm   1/1     Running   0          3m8s
# 5.获取token值
[root@kubernetes01 ~]# kubectl -n kube-system describe $(kubectl -n kube-system get secret -n kube-system -o name | grep namespace) | grep token
Name:         namespace-controller-token-mt4sh
Type:  kubernetes.io/service-account-token
token:      eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJuYW1lc3BhY2UtY29udHJvbGxlci10b2tlbi1tdDRzaCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJuYW1lc3BhY2UtY29udHJvbGxlciIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImY5YzE3YWQzLTUxYzItMTFlOS05NWZiLTAwMTYzZTBlNDRiYyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTpuYW1lc3BhY2UtY29udHJvbGxlciJ9.W2flckBO8CrzGyJzw2aJH5obQSjy4PNSll7uHOiIXPk4dnOTEzI-BfM4C9QrNDjbNTu8gIdLHntLj1181Sf_sRMidB_vhUPg6CFA1zy3XmYH21eVqjSxEBNXMSfrJHBgXnBzaHieaXqF55_etABB0j4xLM7V-bRsQ9AB0G3cv1IYU_gYG3BozksvAObmDEY4GgCI7f0-nu2YRqOMPJPhXWzKOGUvBBPyj171Xo06QvF6p9zpTMSoLa3aV-gU4XA2nMf2_aDdgFrGVI4p95ziewyu0o-W-DiEnXW1hRtwgg-PRe3QPU9ps3TALlr3U8rwh3xVmlqnRuNGVDqzmclVdQ
访问https://10.5.0.206:30001通过token登陆控制面板,注意是https协议！
```

### 1.11部署容器存储插件

Rook项目是基于Ceph的Kubernetes存储插件，一个可用于生产级别的做持久化存储的插件，

```bash
cd /usr/local/src
yum -y install git
git clone https://github.com/rook/rook.git
cd /usr/local/src/rook/cluster/examples/kubernetes/ceph
kubectl apply -f common.yaml
kubectl apply -f operator.yaml
kubectl apply -f cluster.yaml
```

在master节点上处理，要注意污点处理



## Node节点

安装`docer-ce`、





## 命令速查

> 1. kubectl describe pod kubernetes-dashboard-65c76f6c97-hmbd7 --namespace=kube-system 
> 2. kubectl delete -f kubernetes-dashboard.yaml 
> 3. kubectl get pods --all-namespaces
> 4. kubectl delete deployments,svc my-nginx

