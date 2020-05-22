# autons

## 1. 项目背景：

使用docker技术，通过docker-compose -f xxx.yaml up -d，就可以启动一个django web项目

优势：1. 开发环境既是生产环境，避免环境不一致导致运行运行失败

​			2. 部署方便，启动方便，监控方便（目前还未涉及监控）

## 2. 项目实现

1. 通过编写Dockerfile实现一个运行时django镜像，并发布。

   ```shell
   # Dockerfile
   FROM python:3
   ENV PYTHONUNBUFFERED 1
   RUN mkdir /code
   COPY requirements.txt /code
   WORKDIR /code
   RUN pip install -r requirements.txt
   ```

   > 在Dockerfile目录下执行命令：
   >
   > 生成镜像文件：docker build -t autons:v1 .
   >
   > 打tag：docker tag autons:v1 autumn66/autons:laster
   >
   > docker login
   >
   > docker push autumn66/autons:laster

2. 编写docker-compose.yaml，使用docker-compose 启动

   ```yaml
   # docker-compose.yaml
   version: "3"
     
   services:
       autons-web:
           container_name: autons-web
           image: autumn66/autons-base:laster
           volumes:
               - .:/code
           ports:
               - 8000:8000
           working_dir: /code
           command: python manage.py runserver 0.0.0.0:8000
           
   # 数据库使用宿主机的mysql，注意在django settings中正确配置主机地址（docker0地址）
   ```

   > 执行命令：docker-compose up -d