# CNFast - 国内开发者网络加速工具

<p align="center">
  <img src="https://img.shields.io/badge/Version-1.0.0-brightgreen.svg" alt="Version">
  <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg" alt="Platform">
  <img src="https://img.shields.io/badge/Go-1.23.4-blue.svg" alt="Go Version">
</p>

<p align="center">
  <b>CNFast</b> 是一个专为国内开发者设计的网络加速工具，通过智能代理技术解决访问GitHub、Docker Hub等国外资源缓慢或无法访问的问题。
</p>

<p align="center">
  <a href="#快速开始">快速开始</a> •
  <a href="#功能特性">功能特性</a> •
  <a href="#安装使用">安装使用</a> •
  <a href="#文档">文档</a> •
  <a href="#贡献">贡献</a>
</p>

## ✨ 功能特性

- 🚀 **极速克隆**：加速GitHub仓库的克隆与拉取，并支持代理下载 GitHub Release 文件
- 🐳 **镜像加速**：优化Docker镜像拉取速度，支持多 registry
- 🔒 **稳定可靠**：基于稳定的代理技术，保证连接成功率
- 🛠️ **简单易用**：命令行工具，一键加速，无需复杂配置
- 🌐 **多平台支持**：支持Linux、macOS和Windows系统
- ⚡ **择优推荐**：交互式选择评分最高的代理服务
- 🔧 **高度可配置**：支持环境变量配置（API 地址、调试模式、超时时间）
- 📊 **调试支持**：提供 `CNFAST_DEBUG` 调试模式，输出详细执行信息

## 🚀 快速开始

### 一键安装

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
# Linux/macOS
tar -zxvf cnfast_linux_amd64.tar.gz
sudo mv cnfast /usr/local/bin/

# Windows
# 解压 cnfast_windows_amd64.zip
# 将 cnfast.exe 移动到 PATH 目录
```

### 验证安装

```bash
cnfast --version
```

## 📖 使用方法

### GitHub 仓库加速

```bash
# 克隆仓库
cnfast git clone https://github.com/microsoft/vscode.git

# 拉取更新（需先进入已克隆的仓库目录）
cnfast git pull

# 代理下载 GitHub Release 文件
cnfast git down https://github.com/user/repo/releases/download/v1.0.0/file.tar.gz
```

> 注：当前版本仅支持 `clone`、`pull`、`down` 三个 git 子命令。

### Docker 镜像加速

```bash
# 拉取官方镜像
cnfast docker pull nginx:latest

# 拉取 Kubernetes 镜像
cnfast docker pull k8s.gcr.io/pause:3.2

# 拉取 GitHub 镜像
cnfast docker pull ghcr.io/octocat/hello-world:latest

# 推送镜像
cnfast docker push your-registry/your-image:tag
```

### 其他功能

```bash
# 查看版本信息
cnfast --version

# 查看帮助信息
cnfast --help

# 更新 cnfast 自身
cnfast update
```

## ⚙️ 配置说明

CNFast 支持通过环境变量进行配置：

```bash
# 设置 API 服务器地址
export CNFAST_API_HOST="https://cnfast-api.521456.xyz"

# 启用调试模式
export CNFAST_DEBUG=true

# 设置请求超时时间（秒）
export CNFAST_TIMEOUT=30
```

### 支持的镜像源

- **Docker Hub** (`docker.io`)
- **Google Container Registry** (`gcr.io`)
- **Kubernetes Registry** (`k8s.gcr.io`, `registry.k8s.io`)
- **GitHub Container Registry** (`ghcr.io`)
- **Quay.io** (`quay.io`)
- **NVIDIA Container Registry** (`nvcr.io`)
- **Cloudsmith** (`docker.cloudsmith.io`)

## 🏗️ 工作原理

CNFast 通过智能代理技术，为 GitHub 与 Docker 访问加速：

1. **代理发现**：启动时从 API 服务获取可用的代理列表，按评分排序
2. **交互选择**：Git 命令执行前提示用户从代理列表中选择一个代理服务
3. **URL 改写**：Git 克隆/拉取时在 GitHub URL 前拼接代理前缀；Docker 拉取时将镜像名中的 registry 域名替换为加速域名
4. **失败重试**：Git 操作失败时可选择尝试下一个可用代理

## 📊 性能对比

| 操作类型 | 直接访问 | 使用 CNFast | 性能提升 |
|----------|----------|-------------|----------|
| GitHub 克隆 | 15-50 KB/s | 5-10 MB/s | **100x+** |
| Docker 拉取 | 20-100 KB/s | 10-50 MB/s | **100x+** |
| 连接成功率 | 60-80% | 99% | **显著提升** |

## ❓ 常见问题

### Q: CNFast 是否免费？
A: 是的，CNFast 是完全免费的开源工具。

### Q: 支持哪些 GitHub 操作？
A: 支持 `clone`（克隆）、`pull`（拉取更新）和 `down`（代理下载 Release 文件）。暂不支持 `push`、`fetch` 等其它 git 子命令。

### Q: 是否支持私有仓库？
A: CNFast 不修改原有认证信息，认证行为与直接使用 git 一致，但加速效果取决于所选的代理服务对私有仓库的支持情况。

### Q: 如何更新 CNFast？
A: 执行 `cnfast update`，或重新下载最新版本的二进制文件/使用安装脚本更新。

### Q: 支持哪些操作系统？
A: 源码支持 Linux、macOS，Windows 支持见 [docs/USER_GUIDE.md](docs/USER_GUIDE.md)。

### Q: 如何获取帮助？
A: 使用 `cnfast --help` 查看帮助信息，或访问 [GitHub Issues](https://github.com/sallaixu/cnfast/issues)。

## 🤝 参与贡献

我们欢迎任何形式的贡献！包括但不限于：

- 🐛 提交 bug 报告或功能请求
- 💻 提交代码改进
- 📚 完善文档
- 💡 分享使用经验
- 🌟 给项目点星

请阅读 [贡献指南](CONTRIBUTING.md) 了解如何开始。

## 📚 文档

- [用户指南](docs/USER_GUIDE.md) - 详细的使用说明
- [API 文档](docs/API.md) - API 接口和使用方法
- [开发文档](docs/DEVELOPMENT.md) - 开发者指南
- [变更日志](CHANGELOG.md) - 版本更新记录

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

感谢所有为项目做出贡献的开发者，以及提供镜像服务的组织和企业。

## 📞 联系我们

- **GitHub Issues**: [提交问题](https://github.com/sallaixu/cnfast/issues)
- **GitHub Discussions**: [参与讨论](https://github.com/sallaixu/cnfast/discussions)
- **项目主页**: [https://github.com/sallaixu/cnfast](https://github.com/sallaixu/cnfast)

---

**CNFast** - 让开发更流畅，让学习更高效！ 🚀