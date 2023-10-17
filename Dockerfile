# 使用一个基础镜像，例如Alpine Linux，因为它很小而且安全
FROM alpine:latest

# 定义变量
ARG GITHUB_USERNAME
ARG REPO_NAME
ARG VERSION
ARG BINARY_NAME

# 构造二进制文件的下载URL
ARG BINARY_URL=https://github.com/${GITHUB_USERNAME}/${REPO_NAME}/releases/download/${VERSION}/${BINARY_NAME}

# 安装所需的工具
RUN apk --no-cache add wget tar unzip

# 下载二进制文件
RUN wget -O /tmp/${BINARY_NAME} ${BINARY_URL}

# 检查文件扩展名并解压缩
RUN case "${BINARY_NAME}" in \
        *.tar.gz) tar -xzf /tmp/${BINARY_NAME} -C /usr/local/bin ;; \
        *.zip) unzip /tmp/${BINARY_NAME} -d /usr/local/bin ;; \
    esac \
    && chmod +x /usr/local/bin/* \
    && rm /tmp/${BINARY_NAME}

# 假设二进制文件名为ddns，设置二进制文件为容器的入口点
ENTRYPOINT ["/usr/local/bin/ddns"]
