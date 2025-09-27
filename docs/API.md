# CNFast API 文档

## 概述

CNFast 是一个专为国内开发者设计的网络加速工具，通过智能代理技术加速 GitHub 仓库访问和 Docker 镜像拉取。

## 核心功能

### 1. GitHub 仓库加速

CNFast 支持以下 Git 操作的加速：

- `clone` - 克隆仓库
- `pull` - 拉取最新更改
- `fetch` - 获取远程更改
- `push` - 推送更改

#### 使用示例

```bash
# 克隆 GitHub 仓库
cnfast git clone https://github.com/user/repo.git

# 拉取最新更改
cnfast git pull

# 获取远程更改
cnfast git fetch
```

### 2. Docker 镜像加速

CNFast 支持以下 Docker 操作的加速：

- `pull` - 拉取镜像
- `push` - 推送镜像
- `build` - 构建镜像

#### 支持的镜像源

- Docker Hub (`docker.io`)
- Google Container Registry (`gcr.io`)
- Kubernetes Registry (`k8s.gcr.io`, `registry.k8s.io`)
- GitHub Container Registry (`ghcr.io`)
- Quay.io (`quay.io`)
- NVIDIA Container Registry (`nvcr.io`)
- Cloudsmith (`docker.cloudsmith.io`)

#### 使用示例

```bash
# 拉取 Docker 镜像
cnfast docker pull nginx:latest
cnfast docker pull ubuntu:20.04

# 拉取 Kubernetes 镜像
cnfast docker pull k8s.gcr.io/pause:3.2
```

## 配置选项

### 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `CNFAST_API_HOST` | API 服务器地址 | `https://cnfast-api.521456.xyz` |
| `CNFAST_DEBUG` | 启用调试模式 | `false` |
| `CNFAST_TIMEOUT` | 请求超时时间（秒） | `30` |

### 配置示例

```bash
# 设置自定义 API 服务器
export CNFAST_API_HOST="https://your-api-server.com"

# 启用调试模式
export CNFAST_DEBUG=true

# 设置超时时间
export CNFAST_TIMEOUT=60
```

## 代理服务

### 代理类型

CNFast 支持两种类型的代理服务：

1. **Git 代理** - 用于加速 GitHub 仓库访问
2. **Docker 代理** - 用于加速 Docker 镜像拉取

### 代理选择

系统会自动选择评分最高的代理服务，确保最佳的网络性能。

### 代理评分

代理服务根据以下因素进行评分：

- 网络延迟
- 连接稳定性
- 传输速度
- 可用性

## 错误处理

### 常见错误

1. **网络连接错误**
   - 检查网络连接
   - 验证 API 服务器地址

2. **代理服务不可用**
   - 系统会自动尝试其他代理
   - 检查代理服务状态

3. **命令不支持**
   - 查看支持的命令列表
   - 使用 `cnfast --help` 获取帮助

### 调试模式

启用调试模式可以获取详细的执行信息：

```bash
export CNFAST_DEBUG=true
cnfast git clone https://github.com/user/repo.git
```

## 性能优化

### 网络加速原理

1. **智能路由** - 自动选择最优的国内镜像节点
2. **连接复用** - 减少连接建立时间
3. **缓存机制** - 提高重复请求的响应速度

### 性能对比

| 操作类型 | 直接访问 | 使用 CNFast | 性能提升 |
|----------|----------|-------------|----------|
| GitHub 克隆 | 15-50 KB/s | 5-10 MB/s | 100x+ |
| Docker 拉取 | 20-100 KB/s | 10-50 MB/s | 100x+ |
| 连接成功率 | 60-80% | 99% | 显著提升 |

## 最佳实践

### 1. 使用建议

- 定期更新 CNFast 到最新版本
- 在网络状况良好时使用
- 避免同时运行多个加速任务

### 2. 故障排除

- 检查网络连接状态
- 验证代理服务可用性
- 查看错误日志信息

### 3. 性能优化

- 使用稳定的网络环境
- 避免网络高峰期使用
- 定期清理本地缓存

## 技术支持

如果您在使用过程中遇到问题，请：

1. 查看错误日志
2. 检查网络连接
3. 联系技术支持

更多信息请访问：[GitHub 仓库](https://github.com/sallaixu/cnfast)
