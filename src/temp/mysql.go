package temp

var MySQLCNF = `
[client]
default-character-set=utf8mb4
port=3306
socket=/ddhome/local/mysql/data/mysql.sock

[mysqld]
basedir=/usr/local/mysql
datadir=/ddhome/local/mysql/data
socket=/ddhome/local/mysql/data/mysql.sock
skip-name-resolve
max_connections=4000
table_open_cache=200
log_bin = mysql-bin
binlog_format = mixed
expire_logs_days = 30
innodb_buffer_pool_size=3G
innodb_file_per_table=ON

[mysqld_safe]
log-error=/ddhome/local/mysql/data/mysqld.log
pid-file=/ddhome/local/mysql/data/mysqld.pid
`
