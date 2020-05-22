

## **基本操作**

1. 登录：

     > mysql -u root -p

2. 指定参数登录：

     > mysql -uroot -p312312 -h127.0.0.1 -P3808

3. 查看数据库/表：

     > show databases/tables

4. 切换数据库：

     > use fstar

5. 查看数据表字段信息

     > show full columns from token_transaction;

6. 创建数据库：

     CREATE DATABASE xxx

7. 修改root密码：
    mysql -u root
    mysql> use mysql;
    mysql> UPDATE user SET Password = PASSWORD('newpass') WHERE user = 'root';
    mysql> FLUSH PRIVILEGES;

8. 关闭数据库：mysqladmin -uroot -p312312 shutdown -S /home/prod/softwares/mysql/mysqld.sock

## 数据查询

1. 查看建表语句：show create table xxx;

2. 查看表的详细信息：describe xxx

3. 查询非空数据

   > select * from url_safe where url <>'';  // <> 类始于!=
   >
   > select * from url_safe where url is(not) null;

4. 根据某字段进行分组

   > select count(*) from url_safe group by ap_mac；// 有几个不同的ap_mac

5. 过滤重复数据查询

   >  SELECT DISTINCT 字段名 FROM 表名 WHERE 查询条件 

6. 

## 数据写入

1. 更新数据

     > UPDATE release_token SET supply=1000000,create_time='2018-6-19 10:03',last_update_time='2018-6-19 10:03' where id=2;

2. 数据插入，有重复数据时进行更新

     > 如果在INSERT语句末尾指定了ON DUPLICATE KEY UPDATE，并且插入行后会导致在一个UNIQUE索引或PRIMARY KEY中出现重复值，则在出现重复值的行执行UPDATE；如果不会导致唯一值列重复的问题，则插入新行。
     >
     > insert into release_token(id, name, age) values("",'"","") ON DUPLICATE KEY UPDATE name=values(name), age=values(age)

3. 模糊查询

     > select * from tbName where name list "%刘"；

4. 查看表的索引： show index from xxx;

5. 为表添加索引：ALTER TABLE `table_name` ADD UNIQUE ( `column` ) 

6. 查询最大的记录：select * from contract_account where id =(select max(id) from contract_account); 

7. 查看创建表的详情：show create table user \G;

8. 更改数据类型： alter table 表名 modify column 字段名 类型;

9. 删除数据库：drop database xxx

10. 删除一条记录：delect from xxx where id = 1

11. 查看事务的隔离级别：show variables like 'tx_isolation'

12. 修改中文编码设置：my.conf
       [mysqld]
       character-set-server=utf8 
       [client]
       default-character-set=utf8
       [mysql]
       default-character-set=utf8

13. 查看编码格式
      数据库：show variables like 'character_set_database';
       数据表：show create table 表名;
       修改：
       alter database resttest character set utf8;

