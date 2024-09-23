# flow

## 流量请求限制

> client -> flow -> server

 配置
```
handle:
  request_threshold: 100 # 100次/5s
  block_duration: 5m # 5分钟 黑名单时间
  request_timeout: 5s # 5秒 短时间统计次数重置
server:
  enable_tcp: true # 开启tcp服务
  enable_udp: true # 开启udp服务
  tcp_port: 6001 # 监听端口TCP
  udp_port: 8888 # 监听端口UDP
  internal_tcp_addr: "127.0.0.1:5666" # 内部服务地址TCP
  internal_udp_addr: "127.0.0.1:5888" # 内部服务地址UDP
web_server:
  enable_http: true # 开启http服务
  enable_tls: false # 开启tls服务
  host_addr: "0.0.0.0:8080" # 监听端口
  cert_file: "" # 证书文件
  key_file: "" # 私钥文件
  remote_addr: "127.0.0.1:80" # 内部服务地址
log:
  level_file: "warn" # 文件日志级别
  level_console: "info" # 控制台日志级别
  path: "logs/access.log" # 日志文件路径
  max_size: 5 # 日志文件最大大小（MB）
  max_age: 10 # 日志文件最大保留天数
  max_backups: 5 # 日志文件最大备份数
  compress: true # 是否压缩日志文件
```