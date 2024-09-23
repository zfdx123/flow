# 使用基础镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /flow

# 将二进制文件和配置文件复制到镜像中
COPY flow /flow/flow
COPY config.yaml /flow/config.yaml

# 赋予二进制文件执行权限
RUN chmod +x /flow/flow

# 设置默认命令
CMD ["./flow", "-config", "config.yaml"]

