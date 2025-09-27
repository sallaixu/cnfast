# 变更日志

本文档记录了 CNFast 项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 新增
- 添加详细的代码注释和文档
- 优化项目结构和代码组织
- 增强错误处理和调试信息
- 添加环境变量配置支持
- 改进代理服务选择逻辑

### 改进
- 重构 HTTP 客户端，提高稳定性
- 优化 Docker 镜像加速逻辑
- 增强 Git 命令处理能力
- 改进帮助信息显示
- 提升代码可读性和维护性

### 修复
- 修复代理服务选择问题
- 解决网络连接超时问题
- 修复命令参数解析错误
- 修复错误信息显示问题

## [1.0.0] - 2024-01-01

### 新增
- 初始版本发布
- 支持 GitHub 仓库加速
- 支持 Docker 镜像加速
- 多平台构建支持
- 命令行界面

### 功能
- **Git 加速**: 支持 clone、pull、fetch、push 操作
- **Docker 加速**: 支持 pull、push、build 操作
- **代理服务**: 自动选择最优代理
- **多平台**: 支持 Linux、macOS、Windows
- **配置**: 支持环境变量配置

### 支持的镜像源
- Docker Hub (`docker.io`)
- Google Container Registry (`gcr.io`)
- Kubernetes Registry (`k8s.gcr.io`, `registry.k8s.io`)
- GitHub Container Registry (`ghcr.io`)
- Quay.io (`quay.io`)
- NVIDIA Container Registry (`nvcr.io`)
- Cloudsmith (`docker.cloudsmith.io`)

### 安装方式
- 一键安装脚本
- 手动下载安装
- 源码编译安装

### 使用示例
```bash
# GitHub 仓库加速
cnfast git clone https://github.com/user/repo.git

# Docker 镜像加速
cnfast docker pull nginx:latest
```

## [0.9.0] - 2023-12-15

### 新增
- 项目初始开发
- 基础架构搭建
- 核心功能实现

### 功能
- 代理服务管理
- HTTP 客户端实现
- 命令行解析
- 错误处理机制

## 版本说明

### 版本号规则

- **主版本号**: 不兼容的 API 修改
- **次版本号**: 向下兼容的功能性新增
- **修订号**: 向下兼容的问题修正

### 发布周期

- **主版本**: 每年 1-2 次
- **次版本**: 每季度 1-2 次
- **修订版**: 根据需要发布

### 支持策略

- **当前版本**: 完全支持
- **前一个主版本**: 安全更新
- **更早版本**: 不再支持

## 贡献者

感谢所有为项目做出贡献的开发者：

- [项目维护者](https://github.com/sallaixu)
- [贡献者列表](https://github.com/sallaixu/cnfast/contributors)

## 致谢

感谢以下组织和项目：

- [Go 语言团队](https://golang.org/)
- [Docker 团队](https://www.docker.com/)
- [GitHub 团队](https://github.com/)
- 所有提供镜像服务的组织和企业

## 许可证

本项目采用 MIT 许可证。详情请查看 [LICENSE](LICENSE) 文件。

---

**注意**: 本文档会随着项目的发展持续更新。如果您发现任何问题或需要添加内容，请提交 Issue 或 Pull Request。
