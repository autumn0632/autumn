1. 什么是redis：

   * Redis是一个基于内存的高性能key-value数据库。
   * 整个数据库统统加载在内存当中进行操作，定期通过异步操作把数据库数据flush到硬盘上进行保存。
   * redis全称：**Remote Dictionary Server**
   * redis 也经常用来做分布式锁
   * redis 支持事务 、持久化、LUA脚本、LRU驱动事件、多种集群方案。

2. 为什么使用redis

   * 高性能：用户第一次访问数据库中的某些数据。这个过程会比较慢，因为是从硬盘上读取的。将该用户访问的数据存在数缓存中，这样下一次再访问这些数据的时候就可以直接从缓存中获取了。操作缓存就是直接操作内存，所以速度相当快。
   * 高并发：直接操作缓存能够承受的请求是远远大于直接访问数据库的。

3. redis和memcached的区别：

   * **redis支持更丰富的数据类型（支持更复杂的应用场景）**：Redis不仅仅支持简单的k/v类型的数据，同时还提供list，set，zset，hash等数据结构的存储。memcache支持简单的数据类型，String。
   * **Redis支持数据的持久化，可以将内存中的数据保持在磁盘中，重启的时候可以再次加载进行使用,而Memecache把数据全部存在内存之中。**
   * **Memcached是多线程，非阻塞IO复用的网络模型；Redis使用单线程的多路 IO 复用模型。**

4. redis是单进程单线程的：

   redis利用队列技术将并发访问变为串行访问，消除了传统数据库串行控制的开销

5. redis的数据类型：

   * 字符串String

     > **常用命令:** set,get,decr,incr,mget 等。
     >
     > String数据结构是简单的key-value类型，value其实不仅可以是String，也可以是数字。

   * 字典Hash

     > **常用命令：** hget,hset,hgetall 等。
     >
     > Hash 是一个 string 类型的 field 和 value 的映射表，hash 特别适合用于存储对象，直接仅仅修改这个对象中的某个字段的值。
     >
     > 通过 **"数组 + 链表"** 的链地址法来解决部分 **哈希冲突**，同时这样的结构也吸收了两种不同数据结构的优点。

   * 列表List

     > **常用命令:** lpush,rpush,lpop,rpop,lrange等.
     >
     > `LPUSH` 和 `RPUSH` 分别可以向 list 的左边（头部）和右边（尾部）添加一个新元素；
     >
     > `LINDEX` 命令可以从 list 中取出指定下表的元素；
     >
     > 可以通过 lrange 命令，就是从某个元素开始读取多少个元素，可以基于 list 实现分页查询，这个很棒的一个功能，基于 redis 实现简单的高性能分页，可以做类似微博那种下拉不断分页的东西（一页一页的往下走），性能高。

   * 集合Set

     > **常用命令：**sadd,spop,smembers,sunion 等
     >
     > 可以基于 set 轻易实现交集、并集、差集的操作。
     >
     > sinterstore key1 key2 key3   -  将交集存在key1内

   * 有序集合SortedSet

     > **常用命令：** zadd,zrange,zrem,zcard等
     >
     > 和set相比，sorted set增加了一个权重参数score，使得集合中的元素能够按score进行有序排列。

6. redis的数据淘汰策略：

   redis 支持对存储在数据库中的值设置一个过期时间。当到达过期时间之后，redis通过**定期删除+惰性删除。**两种方式进行删除：

   * **定期删除：**

     redis默认每隔 100ms 就**随机抽取**一些设置了过期时间的key，检查其是否过期，如果过期就删除。注意这里是随机抽取的。因为如果 redis 存了几十万个 key ，每隔100ms就遍历所有的设置过期时间的 key 的话，就会给 CPU 带来很大的负载！

   * **惰性删除：**

     定期删除可能会导致很多过期 key 到了时间并没有被删除掉。这些key除非系统去查一下那个 key，才会被redis给删除掉。这就是所谓的惰性删除。

   如果定期删除漏掉了很多过期 key，然后你也没及时去查，也就没走惰性删除，就会有大量过期key堆积在内存里，导致redis内存块耗尽了。这时候就需要**内存淘汰机制**

   **redis 提供 6种数据淘汰策略：**

   *  **volatile-lru**：从已设置过期时间的数据集（server.db[i].expires）中挑选最近最少使用的数据淘汰
   * **volatile-ttl**：从已设置过期时间的数据集（server.db[i].expires）中挑选将要过期的数据淘汰
   * **volatile-random**：从已设置过期时间的数据集（server.db[i].expires）中任意选择数据淘汰

   * **allkeys-lru**：当内存不足以容纳新写入数据时，在键空间中，移除最近最少使用的key（这个是最常用的）.

   * **allkeys-random**：从数据集（server.db[i].dict）中任意选择数据淘汰

   * **no-enviction**：禁止驱逐数据，也就是说当内存不足以容纳新写入数据时，新写入操作会报错。这个基本不会用到

7. redis 数据持久化

   redis支持两种不同的持久化操作：**一种持久化方式叫快照（snapshotting，RDB）,另一种方式是只追加文件（append-only file,AOF）**

   * **快照（snapshotting）持久化（RDB）**

     Redis可以通过创建快照来获得存储在内存里面的数据在某个时间点上的副本。

     快照持久化是Redis默认采用的持久化方式，在redis.conf配置文件中默认有此下配置：

     ```shell
     save 900 1          #在900秒(15分钟)之后，如果至少有1个key发生变化，Redis就会自动触发BGSAVE命令创建快照。
     save 300 10         #在300秒(5分钟)之后，如果至少有10个key发生变化，Redis就会自动触发BGSAVE命令创建快照。
     save 60 10000       #在60秒(1分钟)之后，如果至少有10000个key发生变化，Redis就会自动触发BGSAVE命令创建快照。
     ```

   * **AOF（append-only file）持久化**

     与快照持久化相比，AOF持久化 的实时性更好，因此已成为主流的持久化方案。默认情况下Redis没有开启AOF（append only file）方式的持久化，可以通过appendonly参数开启：

     ```shell
     appendonly yes
     ```

     **开启AOF持久化后每执行一条会更改Redis中的数据的命令，Redis就会将该命令写入硬盘中的AOF文件。**

     在Redis的配置文件中存在三种不同的 AOF 持久化方式，它们分别是：

     ```shell
     appendfsync always     #每次有数据修改发生时都会写入AOF文件,这样会严重降低Redis的速度
     appendfsync everysec   #每秒钟同步一次，显示地将多个写命令同步到硬盘
     appendfsync no         #让操作系统决定何时进行同步
     ```

8. redis适合的场景

   * 会话缓存（Session Cache）

   * 全页缓存（FPC）

   * 队列

     提供 list 和 set 操作，这使得Redis能作为一个很好的消息队列平台来使用。

   * 发布/订阅

9. 大批量数据插入方法

   * 普通set 集合方法
   * **pipline 管道技术**
   * **把需要插入的数据分块批量插入**

10. 管道技术

   Redis是一种基于客户端-服务端模型以及请求/响应协议的TCP服务。这意味着通常情况下一个请求会遵循以下步骤：

   - 客户端向服务端发送一个查询请求，并监听Socket返回，通常是以阻塞模式，等待服务端响应。
   - 服务端处理命令，并将结果返回给客户端

   Redis 管道技术可以在服务端未响应时，客户端可以继续向服务端发送请求，并最终一次性读取所有服务端的响应。