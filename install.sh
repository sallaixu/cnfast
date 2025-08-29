#!/bin/bash

# cnfast 一键安装脚本
# 适用于 Linux 系统

set -e  # 遇到错误立即退出

# 定义变量
DOWNLOAD_URL="https://gitee.com/sallai/cnfast/raw/main/build/linux-amd64/cnfast"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="cnfast"

# 检查是否以 root 权限运行
if [ "$(id -u)" -ne 0 ]; then
    echo "请使用 sudo 或以 root 用户身份运行此脚本"
    exit 1
fi

# 检查系统架构
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        echo "检测到 x86_64 架构"
        ;;
    aarch64|arm64)
        echo "检测到 ARM64 架构"
        ;;
    *)
        echo "不支持的架构: $ARCH"
        exit 1
        ;;
esac

# 创建临时目录
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# 下载 cnfast 二进制文件
echo "正在下载 cnfast..."
if command -v wget &> /dev/null; then
    wget -q "$DOWNLOAD_URL" -O "$BINARY_NAME"
elif command -v curl &> /dev/null; then
    curl -s -L "$DOWNLOAD_URL" -o "$BINARY_NAME"
else
    echo "错误: 需要 wget 或 curl 来下载文件"
    exit 1
fi

# 检查下载是否成功
if [ ! -f "$BINARY_NAME" ]; then
    echo "下载失败，请检查网络连接和下载链接"
    exit 1
fi

# 使二进制文件可执行
chmod +x "$BINARY_NAME"

# 安装到系统目录
echo "正在安装到 $INSTALL_DIR..."
mv "$BINARY_NAME" "$INSTALL_DIR/"

# 清理临时文件
rm -rf "$TMP_DIR"

# 验证安装
if command -v "$BINARY_NAME" &> /dev/null; then
    echo "安装成功！"
    echo "你可以运行 'cnfast' 来使用程序"
else
    echo "安装可能存在问题，请手动检查"
fi