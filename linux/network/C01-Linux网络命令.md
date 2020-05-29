## netstat

```shell
prod@QQ-BJ-test-wj01:~$ netstat -h
usage: netstat [-vWeenNcCF] [<Af>] -r         netstat {-V|--version|-h|--help}
       netstat [-vWnNcaeol] [<Socket> ...]
       netstat { [-vWeenNac] -i | [-cWnNe] -M | -s }

        -r, --route              显示路由表
        -i, --interfaces         display interface table
        -g, --groups             display multicast group memberships
        -s, --statistics         display networking statistics (like SNMP)
        -M, --masquerade         display masqueraded connections

        -v, --verbose            显示详细信息
        -W, --wide               don't truncate IP addresses
        -n, --numeric            不解析名称
        --numeric-hosts          不解析主机名
        --numeric-ports          忽略端口名称
        --numeric-users          忽略用户名
        -N, --symbolic           resolve hardware names
        -e, --extend             显示更多信息
        -p, --programs           display PID/Program name for sockets
        -c, --continuous         continuous listing

        -l, --listening          display listening server sockets
        -a, --all, --listening   display all sockets (default: connected)
        -o, --timers             display timers
        -F, --fib                display Forwarding Information Base (default)
        -C, --cache              display routing cache instead of FIB

  <Socket>={-t|--tcp} {-u|--udp} {-w|--raw} {-x|--unix} --ax25 --ipx --netrom
  <AF>=Use '-6|-4' or '-A <af>' or '--<af>'；默认： inet
  列出所有支持的协议：
    inet (DARPA Internet) inet6 (IPv6) ax25 (AMPR AX.25) 
    netrom (AMPR NET/ROM) ipx (Novell IPX) ddp (Appletalk DDP) 
    x25 (CCITT X.25) 
```

* **netstat -a\[u][t]**：a：列出所有端口，包括监听和未监听的；不加-a，默认显示connected，u：udp端口；t：tcp端口

* **netstat -l\[u][t]**：列出所有监听端口，u：udp；t：tcp

* **netstat -s\[u][t]**：显示所有端口的统计信息

* **netstat -pn**：p：显示pid和进程名称；n：不显示主机，端口和用户名

* **netstat -r**： 显示核心路由信息；与 `route`显示相同

* 查看连接某服务器最多的ip地址：

  ```shell
  netstat -nat | grep "192.168.1.15:22" |awk '{print $5}'|awk -F: '{print $1}'|sort|uniq -c|sort -nr|head -20
  ```

## router

