package temp

var NginxConf = `
user  nginx;
worker_processes  auto;

error_log  /var/log/nginx/error.log;
#error_log  logs/error.log  notice;
#error_log  logs/error.log  info;

pid        /var/run/nginx.pid;


events {
    worker_connections  10240;
    use epoll;
}


http {
    include       mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    tcp_nopush      on;
    tcp_nodelay     on;

    server_tokens off;
    client_max_body_size 2048m;
    keepalive_timeout  65;

    gzip  on;
    gzip_min_length 1k;
    gzip_comp_level 9;
    gzip_types text/plain application/x-javascript application/javascript  text/css application/xml text/javascript application/x-httpd-php image/jpeg image/gif image/png;
    gzip_buffers  4 16k;
    gzip_vary on;

    include vhosts/*.conf;
}
`

var ArchConf = `
server {
    listen 80;
    server_name _;
    root /ddhome/project/arch/web/user;
    index index.html index.htm;

	location /admin {
		alias /ddhome/project/arch/web/admin;
	}
}


`