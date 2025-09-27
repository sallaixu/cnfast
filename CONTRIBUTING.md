# 贡献指南

感谢您对 CNFast 项目的关注！我们欢迎任何形式的贡献，包括但不限于：

- 提交 bug 报告或功能请求
- 提交代码改进
- 完善文档
- 分享使用经验

## 贡献方式

### 1. 报告问题

如果您发现了 bug 或有功能建议，请通过以下方式报告：

- 使用 [GitHub Issues](https://github.com/sallaixu/cnfast/issues) 提交问题
- 提供详细的复现步骤
- 包含系统信息和错误日志

#### 问题报告模板

```markdown
**问题描述**
简要描述问题

**复现步骤**
1. 执行命令
2. 观察结果
3. 出现错误

**期望行为**
描述期望的正确行为

**系统信息**
- 操作系统: 
- CNFast 版本: 
- Go 版本: 

**错误日志**
粘贴相关的错误日志
```

### 2. 提交代码

#### 开发流程

1. **Fork 项目**
   ```bash
   # 在 GitHub 上 Fork 项目
   # 然后克隆到本地
   git clone https://github.com/your-username/cnfast.git
   cd cnfast
   ```

2. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **进行开发**
   - 编写代码
   - 添加测试
   - 更新文档

4. **提交更改**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

5. **推送分支**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **创建 Pull Request**
   - 在 GitHub 上创建 Pull Request
   - 填写详细的描述
   - 等待代码审查

#### 代码规范

##### 1. 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**类型说明：**
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

**示例：**
```
feat: add support for custom proxy configuration

- Add new configuration option for custom proxy
- Update documentation with usage examples
- Add tests for new functionality

Closes #123
```

##### 2. 代码风格

- 使用 `gofmt` 格式化代码
- 遵循 Go 官方代码规范
- 使用有意义的变量和函数名
- 添加必要的注释

##### 3. 测试要求

- 为新功能添加单元测试
- 确保所有测试通过
- 测试覆盖率不低于 80%

#### 代码审查

所有代码都需要经过审查才能合并：

1. **审查者检查**
   - 代码质量和功能正确性
   - 测试覆盖度
   - 文档完整性

2. **作者响应**
   - 及时回复审查意见
   - 修改代码问题
   - 更新相关文档

3. **合并标准**
   - 至少一个审查者批准
   - 所有检查通过
   - 无冲突

### 3. 完善文档

#### 文档类型

- **API 文档**: 描述 API 接口和使用方法
- **用户指南**: 提供详细的使用说明
- **开发文档**: 帮助开发者理解代码结构
- **示例代码**: 提供实际使用案例

#### 文档规范

- 使用 Markdown 格式
- 保持内容准确和最新
- 提供清晰的示例
- 使用适当的标题层级

### 4. 分享经验

#### 使用案例

分享您的使用经验：

- 在 [Discussions](https://github.com/sallaixu/cnfast/discussions) 中分享
- 编写博客文章
- 制作视频教程

#### 社区参与

- 回答其他用户的问题
- 参与技术讨论
- 推广项目

## 开发环境

### 环境要求

- Go 1.23.4 或更高版本
- Git
- Docker (可选)

### 设置开发环境

1. **克隆仓库**
   ```bash
   git clone https://github.com/sallaixu/cnfast.git
   cd cnfast
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **运行测试**
   ```bash
   go test ./...
   ```

4. **构建项目**
   ```bash
   make local
   ```

### 开发工具

推荐使用以下工具：

- **IDE**: VS Code, GoLand, Vim
- **调试器**: Delve
- **性能分析**: pprof
- **代码检查**: golint, go vet

## 项目结构

```
cnfast/
├── cmd/                    # 命令行工具入口
├── config/                 # 配置管理
├── internal/              # 内部包
│   ├── enums/             # 枚举类型
│   ├── models/            # 数据模型
│   ├── pkg/               # 公共包
│   └── services/          # 服务层
├── build/                 # 构建输出
├── docs/                  # 项目文档
├── test/                  # 测试文件
├── main.go               # 主入口
├── go.mod               # Go 模块文件
├── Makefile            # 构建脚本
└── README.md           # 项目说明
```

## 发布流程

### 版本管理

使用语义化版本控制：

- **主版本号**: 不兼容的 API 修改
- **次版本号**: 向下兼容的功能性新增
- **修订号**: 向下兼容的问题修正

### 发布步骤

1. **更新版本号**
   ```bash
   # 更新 config/appConfig.go 中的版本号
   Version = "1.1.0"
   ```

2. **更新 CHANGELOG**
   ```bash
   # 记录新功能和修复
   echo "## [1.1.0] - 2024-01-01" >> CHANGELOG.md
   ```

3. **创建标签**
   ```bash
   git tag -a v1.1.0 -m "Release version 1.1.0"
   git push origin v1.1.0
   ```

4. **构建发布包**
   ```bash
   make package
   ```

5. **创建 GitHub Release**
   - 上传构建的二进制文件
   - 添加发布说明
   - 标记为最新版本

## 行为准则

### 社区准则

我们致力于为每个人提供友好、安全的环境：

- **尊重他人**: 保持礼貌和尊重
- **包容性**: 欢迎不同背景的贡献者
- **建设性**: 提供建设性的反馈
- **专业性**: 保持专业和客观

### 禁止行为

以下行为是不被允许的：

- 骚扰或歧视
- 恶意攻击
- 垃圾信息
- 违反法律法规

## 许可证

本项目采用 MIT 许可证。通过贡献代码，您同意将您的贡献也采用 MIT 许可证。

## 联系方式

如果您有任何问题或建议，请通过以下方式联系：

- **GitHub Issues**: [提交问题](https://github.com/sallaixu/cnfast/issues)
- **GitHub Discussions**: [参与讨论](https://github.com/sallaixu/cnfast/discussions)
- **邮箱**: [项目维护者邮箱]

感谢您的贡献！🎉
