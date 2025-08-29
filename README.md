# CNFast - 国内开发者网络加速工具

<p align="center">
  <img src="https://img.shields.io/badge/Version-1.0.0-brightgreen.svg" alt="Version">
  <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg" alt="Platform">
</p>

<p align="center">
  <b>CNFast</b> 是一个专为国内开发者设计的网络加速工具，解决访问GitHub、Docker Hub等国外资源缓慢或无法访问的问题。
</p>

## ✨ 特性

- 🚀 **极速克隆**：加速GitHub仓库的克隆、拉取和推送操作
- 🐳 **镜像加速**：优化Docker镜像拉取速度，支持多 registry
- 🔒 **稳定可靠**：基于稳定的代理技术，保证连接成功率
- 🛠️ **简单易用**：命令行工具，一键加速，无需复杂配置
- 🌐 **多平台支持**：支持Linux、macOS和Windows系统

## 📦 安装方式

### 一键安装脚本

```bash
# 使用curl安装
curl -fsSL https://raw.githubusercontent.com/sallai/release/main/install.sh | bash

# 或使用wget安装
wget -qO- https://raw.githubusercontent.com/sallai/release/main/install.sh | bash
```

### 手动安装

1. 从 [Release页面](https://github.com/sallai/release/releases) 下载对应平台的二进制文件
2. 解压并移动到系统PATH目录：
```bash
tar -zxvf cnfast_linux_amd64.tar.gz
sudo mv cnfast /usr/local/bin/
```

## 🚀 使用方法

### GitHub仓库加速

```bash
# 克隆仓库
cnfast clone https://github.com/sallaixu/cnfast.git

# 或者使用git命令前缀模式
cnfast git clone https://github.com/sallaixu/cnfast.git

# 拉取更新
cnfast git pull origin main

# 推送更改
cnfast git push origin main
```

### Docker镜像加速

```bash
# 拉取Docker镜像
cnfast pull nginx:latest

# 或者使用docker命令前缀模式
cnfast docker pull nginx:latest

# 推送镜像
cnfast docker push your-registry/your-image:tag
```

### 其他功能

```bash
# 查看当前版本
cnfast --version

# 查看帮助信息
cnfast --help

# 检查网络状态
cnfast status
```

## 🔧 配置说明

CNFast支持通过环境变量进行配置：

```bash
# 设置代理模式 (默认: auto)
export CNFAST_PROXY_MODE=auto  # 可选: direct, proxy, auto

# 设置自定义加速端点
export CNFAST_GITHUB_ENDPOINT=https://your-github-mirror.com
export CNFAST_DOCKER_ENDPOINT=https://your-docker-mirror.com

# 启用详细日志
export CNFAST_DEBUG=true
```

## 🏗️ 工作原理

CNFast通过智能路由技术，自动选择最优的国内镜像节点：

1. **GitHub加速**：使用国内镜像站加速仓库访问，解决clone/pull/push缓慢问题
2. **Docker加速**：自动配置国内registry mirror，加速镜像拉取和推送
3. **智能切换**：根据网络状况自动选择最佳线路，保证连接稳定性

## 📊 性能对比

| 操作 | 直接访问 | 使用CNFast | 提升 |
|------|----------|------------|------|
| GitHub克隆 | 15-50 KB/s | 5-10 MB/s | 100x+ |
| Docker拉取 | 20-100 KB/s | 10-50 MB/s | 100x+ |
| 连接成功率 | 60-80% | 99% | 显著提升 |

## ❓ 常见问题

### Q: CNFast是否免费？
A: 是的，CNFast是完全免费的开源工具。

### Q: 支持哪些GitHub操作？
A: 支持clone、pull、push、fetch等所有git操作。

### Q: 是否支持私有仓库？
A: 支持，CNFast会保持原有的认证信息不变。

### Q: 如何更新CNFast？
A: 运行 `cnfast update` 可以自动检查并更新到最新版本。

## 🤝 参与贡献

我们欢迎任何形式的贡献！包括但不限于：
- 提交bug报告或功能请求
- 提交代码改进
- 完善文档
- 分享使用经验

请阅读 [贡献指南](CONTRIBUTING.md) 了解如何开始。

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

感谢所有为项目做出贡献的开发者，以及提供镜像服务的组织和企业。

---

**CNFast** - 让开发更流畅，让学习更高效！ 🚀