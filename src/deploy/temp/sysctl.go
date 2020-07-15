package temp

var Sysctl = `
# net
## 开启time_wait状态快速回收
net.ipv4.tcp_tw_recycle = 1
## 开启time_wait状态重用机制
net.ipv4.tcp_tw_reuse = 1
## 减少允许time_wait状态的数量，默认180000
net.ipv4.tcp_max_tw_buckets = 6000
## 减少fin_wait_2状态的时间，默认60，防止对端长时间不响应导致占用大量的socket套接字
net.ipv4.tcp_fin_timeout = 10
## 在放弃连接之前syn重试的次数
net.ipv4.tcp_syn_retries = 1
## 定义了内核在放弃连接之前所送出的syn+ack的数据，默认5，大约会花费180秒
net.ipv4.tcp_synack_retries = 1
## 防治ddos攻击，synflood
net.ipv4.tcp_syncookies = 1
## 系统最多有多少个套接字不被关联到任一个用户句柄，所谓的孤儿连接，简单防护ddos工具，内存增大这个值也应该被增大
net.ipv4.tcp_max_orphans = 3276800
## 半连接队列，对于未获得对方确认的连接请求，可以保存在这个队列，服务器网络异常中断可以排查这个参数
net.ipv4.tcp_max_syn_backlog = 262144
## 每个端口接受的数据包的速率比内核处理这些包的速率快时，允许送到队列的数据包的最大数目
net.core.netdev_max_backlog = 262144
## 表示每个套接字所允许的最大缓冲区的大小
net.core.optmem_max = 81920
## 关闭tcp时间戳功能，tcp存在一种行为，可以缓存每个连接最新的时间戳，后续请求中如果时间戳小于缓存的时间戳，即视为无效，相应的数据包会被丢弃
net.ipv4.tcp_timestamps = 0
## 间隔多久发一次keepalive探测包，默认7200
net.ipv4.tcp_keepalive_time = 30
## 探测失败后，间隔多久后重新探测，默认75秒
net.ipv4.tcp_keepalive_intvl = 30
## 探测失败后，最多尝试几次，默认9次
net.ipv4.tcp_keepalive_probes = 3
## 默认的TCP数据接收窗口大小
net.core.rmem_default = 8388608
## 最大的TCP数据接收窗口大小
net.core.rmem_max = 16777216
## 默认发送TCP数据窗口大小
net.core.wmem_default = 8388608
## 最大的TCP发送数据窗口大小
net.core.wmem_max = 16777216
## 内存使用的下限 警戒值 上限值（内存页）
net.ipv4.tcp_mem = 94500000 915000000 927000000
## socket接收缓冲区内存使用的下限 警戒值 上限（内存页）
net.ipv4.tcp_rmem = 4096  87380   4194304
## socket发送缓冲区内存使用的下限 警戒值 上限（内存页）
net.ipv4.tcp_wmem = 4096  16384   4194304
## 启用有选择的应答,通过有选择地应答乱序接收到的报文来提高性能
net.ipv4.tcp_sack = 1
## 启用转发应答，可以进行有选择应答（SACK）从而减少拥塞情况的发生
net.ipv4.tcp_fack = 1
## 启用RFC 1323定义的window scaling，要支持超过64KB的TCP窗口，必须启用该值
net.ipv4.tcp_window_scaling = 1
## 开启反向路径过滤
net.ipv4.conf.default.rp_filter = 1
## 禁用ip源路由
net.ipv4.conf.default.accept_source_route = 0
## 开启路由转发功能
net.ipv4.ip_forward = 1

# kernel
## 以字节为单位规定单一信息队列的最大值
kernel.msgmnb = 65536
## 以字节为单位规定信息队列中任意信息的最大允许的大小
kernel.msgmax = 65536
## 以字节为单位规定一次在该系统中可以使用的共享内存总量
kernel.shmall = 4294967296
## 以字节为单位内核可允许的最大共享内存
kernel.shmmax = 68719476736
## 使用sysrq组合键是了解系统目前运行情况，为安全起见设为0关闭
kernel.sysrq = 0
## 控制core文件的文件名是否添加pid作为扩展
kernel.core_uses_pid = 1

# mem
## 优先互动性并尽量避免将进程装换出物理内存
vm.swappiness = 0
## 定义一个进程能够拥有的最多的内存区域，jvm要求高时，这个值也必须调大
vm.max_map_count = 65535
`
