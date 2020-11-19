FROM redis:latest
WORKDIR /data/
RUN sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list

RUN apt-get clean

# 不能直接 apt-get install curl ,因为容器里面默认apt的包是空的，所以需要更新到本地
RUN apt-get update 

# docker 是基于Ubuntu的，所以里面基本默认都带有apt-get这个工具
RUN apt-get install -y curl \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p  /usr/local/etc/redis/ 

CMD [ "redis-server","/usr/local/etc/redis/redis.conf"]
