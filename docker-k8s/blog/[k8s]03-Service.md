

# 示例

## 1. 使用Service连接到应用

### **NodePort方式**

1. 创建应用 `run-my-nginx.yaml `

   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: my-nginx
   spec:
     selector:
       matchLabels:
         run: my-nginx
     replicas: 2
     template:
       metadata:
         labels:
           run: my-nginx
       spec:
         containers:
         - name: my-nginx
           image: nginx
           ports:
           - containerPort: 80
   ```

   创建：`kubectl apply -f run-my-nginx.yaml`

   查看：`kubectl get pods -l run=my-nginx -o wide`

2. 创建Service

   创建：`kubectl expose deployment/my-nginx`，

   等价于使用`kubectl create -f`命令创建，对应如下的yaml文件：

   ```yaml
   apiVersion: v1
   kind: Service
   metadata:
     name: my-nginx
     labels:
       run: my-nginx
   spec:
     ports:
     - port: 80
       protocol: TCP
     selector:
       run: my-nginx
   ```

   查看：`kubectl get svc my-nginx`，`kubectl describe svc my-nginx`

   ```
   [root@qsc02 app-yaml]# kubectl get svc my-nginx
   NAME       TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
   my-nginx   ClusterIP   10.98.0.97   <none>        80/TCP    33s
   ```

   

3. 修改Service type类型

   `kubectl edit svc my-nginx`， 将`type`值`ClusterIP`改成`NodePort`

   ```yaml
   # Please edit the object below. Lines beginning with a '#' will be ignored,
   # and an empty file will abort the edit. If an error occurs while saving this file will be
   # reopened with the relevant failures.
   #
   apiVersion: v1
   kind: Service
   metadata:
     creationTimestamp: 2020-08-26T12:10:35Z
     name: my-nginx
     namespace: default
     resourceVersion: "224042"
     selfLink: /api/v1/namespaces/default/services/my-nginx
     uid: 25ccb99d-e795-11ea-9b5d-005056ab02c8
   spec:
     clusterIP: 10.98.0.97
     ports:
     - port: 80
       protocol: TCP
       targetPort: 80
     selector:
       run: my-nginx
     sessionAffinity: None
     type: ClusterIP  # => NodePort
   status:
     loadBalancer: {}
   ```

   ```
   [root@qsc02 app-yaml]# kubectl get svc my-nginx 
   NAME       TYPE       CLUSTER-IP   EXTERNAL-IP   PORT(S)        AGE
   my-nginx   NodePort   10.98.0.97   <none>        80:31291/TCP   4m44s
   ```

4. 集群外访问集群任一节点的ip地址+端口号（31291）进行访问

