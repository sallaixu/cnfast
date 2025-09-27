# CNFast 用户指南

## 快速开始

### 安装 CNFast

#### 方法一：一键安装脚本

```bash
# 使用 curl 安装
curl -fsSL https://raw.githubusercontent.com/sallai/release/main/install.sh | bash

# 或使用 wget 安装
wget -qO- https://raw.githubusercontent.com/sallai/release/main/install.sh | bash
```

#### 方法二：手动安装

1. 从 [Release 页面](https://github.com/sallai/release/releases) 下载对应平台的二进制文件
2. 解压并移动到系统 PATH 目录：

```bash
# Linux/macOS
tar -zxvf cnfast_linux_amd64.tar.gz
sudo mv cnfast /usr/local/bin/

# Windows
# 解压 cnfast_windows_amd64.zip
# 将 cnfast.exe 移动到 PATH 目录
```

#### 方法三：从源码构建

```bash
git clone https://github.com/sallaixu/cnfast.git
cd cnfast
make local
```

### 验证安装

```bash
cnfast --version
```

## 基本使用

### GitHub 仓库加速

#### 克隆仓库

```bash
# 使用 CNFast 加速克隆
cnfast git clone https://github.com/microsoft/vscode.git

# 等同于
git clone https://github.com/microsoft/vscode.git
```

#### 拉取更新

```bash
# 进入已克隆的仓库目录
cd vscode

# 使用 CNFast 加速拉取
cnfast git pull

# 等同于
git pull
```

#### 获取远程更改

```bash
# 使用 CNFast 加速获取
cnfast git fetch

# 等同于
git fetch
```

### Docker 镜像加速

#### 拉取镜像

```bash
# 拉取官方镜像
cnfast docker pull nginx:latest

# 拉取 Kubernetes 镜像
cnfast docker pull k8s.gcr.io/pause:3.2

# 拉取 GitHub 镜像
cnfast docker pull ghcr.io/octocat/hello-world:latest
```

#### 推送镜像

```bash
# 推送镜像到仓库
cnfast docker push your-registry/your-image:tag
```

#### 构建镜像

```bash
# 构建镜像
cnfast docker build -t your-image:tag .
```

## 高级功能

### 环境变量配置

#### 设置 API 服务器

```bash
export CNFAST_API_HOST="https://your-api-server.com"
```

#### 启用调试模式

```bash
export CNFAST_DEBUG=true
```

#### 设置超时时间

```bash
export CNFAST_TIMEOUT=60
```

### 代理服务选择

CNFast 会自动选择评分最高的代理服务，但您也可以手动指定：

```bash
# 查看可用代理
cnfast status

# 使用特定代理（如果支持）
export CNFAST_PROXY_URL="https://specific-proxy.com"
```

### 网络诊断

#### 检查网络状态

```bash
cnfast status
```

#### 测试连接速度

```bash
cnfast test
```

## 使用场景

### 场景一：开发环境搭建

```bash
# 1. 克隆项目仓库
cnfast git clone https://github.com/your-org/your-project.git
cd your-project

# 2. 拉取 Docker 镜像
cnfast docker pull node:18-alpine
cnfast docker pull postgres:13

# 3. 启动开发环境
docker-compose up -d
```

### 场景二：CI/CD 流水线

```yaml
# .github/workflows/deploy.yml
name: Deploy
on: [push]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Setup CNFast
        run: |
          curl -fsSL https://raw.githubusercontent.com/sallai/release/main/install.sh | bash
      - name: Pull Docker images
        run: |
          cnfast docker pull nginx:latest
          cnfast docker pull node:18-alpine
```

### 场景三：团队协作

```bash
# 1. 团队成员安装 CNFast
curl -fsSL https://raw.githubusercontent.com/sallai/release/main/install.sh | bash

# 2. 配置环境变量
echo 'export CNFAST_DEBUG=true' >> ~/.bashrc
source ~/.bashrc

# 3. 使用加速功能
cnfast git clone https://github.com/team/project.git
cnfast docker pull team/application:latest
```

## 故障排除

### 常见问题

#### 1. 命令未找到

**问题**: `cnfast: command not found`

**解决方案**:
```bash
# 检查是否在 PATH 中
which cnfast

# 手动添加到 PATH
export PATH=$PATH:/usr/local/bin
```

#### 2. 网络连接失败

**问题**: `获取代理服务失败`

**解决方案**:
```bash
# 检查网络连接
ping cnfast-api.521456.xyz

# 检查防火墙设置
sudo ufw status

# 使用代理服务器
export CNFAST_API_HOST="https://alternative-api.com"
```

#### 3. 权限问题

**问题**: `Permission denied`

**解决方案**:
```bash
# Linux/macOS
sudo chmod +x /usr/local/bin/cnfast

# Windows
# 以管理员身份运行命令提示符
```

#### 4. Docker 命令失败

**问题**: `docker: command not found`

**解决方案**:
```bash
# 安装 Docker
# Ubuntu/Debian
sudo apt-get install docker.io

# macOS
brew install docker

# Windows
# 下载 Docker Desktop
```

### 调试技巧

#### 启用详细日志

```bash
export CNFAST_DEBUG=true
cnfast git clone https://github.com/user/repo.git
```

#### 检查代理状态

```bash
cnfast status
```

#### 测试网络连接

```bash
# 测试 GitHub 连接
curl -I https://github.com

# 测试 Docker Hub 连接
curl -I https://registry-1.docker.io
```

## 性能优化

### 网络优化

#### 1. 使用稳定的网络环境

- 避免使用不稳定的网络连接
- 选择网络延迟较低的时间段使用
- 使用有线网络而非无线网络

#### 2. 配置网络参数

```bash
# 增加超时时间
export CNFAST_TIMEOUT=120

# 使用本地代理
export CNFAST_API_HOST="http://localhost:8080"
```

### 系统优化

#### 1. 清理缓存

```bash
# 清理 Docker 缓存
docker system prune -a

# 清理 Git 缓存
git gc --prune=now
```

#### 2. 优化系统资源

```bash
# 增加文件描述符限制
ulimit -n 65536

# 优化网络参数
echo 'net.core.rmem_max = 16777216' >> /etc/sysctl.conf
echo 'net.core.wmem_max = 16777216' >> /etc/sysctl.conf
```

## 最佳实践

### 1. 使用建议

- **定期更新**: 保持 CNFast 版本最新
- **网络环境**: 在稳定的网络环境下使用
- **资源管理**: 避免同时运行多个大型任务
- **安全考虑**: 不要在公共网络上使用敏感操作

### 2. 性能建议

- **批量操作**: 将多个操作合并执行
- **缓存利用**: 充分利用本地缓存
- **网络优化**: 使用 CDN 和镜像站点
- **资源监控**: 监控系统资源使用情况

### 3. 安全建议

- **权限控制**: 使用最小权限原则
- **网络安全**: 在安全的网络环境下使用
- **数据保护**: 保护敏感数据和凭据
- **定期检查**: 定期检查系统安全性

## 技术支持

### 获取帮助

1. **查看帮助信息**
```bash
cnfast --help
```

2. **查看版本信息**
```bash
cnfast --version
```

3. **检查状态**
```bash
cnfast status
```

### 联系支持

- **GitHub Issues**: [提交问题](https://github.com/sallaixu/cnfast/issues)
- **文档**: [查看文档](https://github.com/sallaixu/cnfast/wiki)
- **社区**: [参与讨论](https://github.com/sallaixu/cnfast/discussions)

### 贡献代码

如果您想为项目做出贡献，请：

1. Fork 项目仓库
2. 创建功能分支
3. 提交更改
4. 创建 Pull Request

更多信息请查看 [贡献指南](CONTRIBUTING.md)。
