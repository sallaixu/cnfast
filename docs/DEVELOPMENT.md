# CNFast 开发文档

## 项目结构

```
cnfast/
├── cmd/                    # 命令行工具入口
├── config/                 # 配置管理
│   └── appConfig.go       # 应用程序配置
├── internal/              # 内部包
│   ├── enums/             # 枚举类型
│   │   └── servicetype.go # 服务类型枚举
│   ├── models/            # 数据模型
│   │   └── proxy.go       # 代理模型
│   ├── pkg/               # 公共包
│   │   ├── help/          # 帮助信息
│   │   │   └── help.go    # 帮助模块
│   │   └── httpclient/    # HTTP 客户端
│   │       ├── httpclient.go # HTTP 客户端实现
│   │       └── errors.go   # 错误处理
│   └── services/          # 服务层
│       ├── proxyservice.go # 代理服务
│       ├── dockerservice.go # Docker 服务
│       └── gitservice.go  # Git 服务
├── build/                 # 构建输出
├── docs/                  # 项目文档
├── main.go               # 主入口
├── go.mod               # Go 模块文件
├── Makefile            # 构建脚本
└── README.md           # 项目说明
```

## 核心组件

### 1. 配置管理 (config)

负责管理应用程序的配置信息，包括：

- API 服务器地址
- 调试模式设置
- 超时配置
- 版本信息

**主要功能：**
- 环境变量读取
- 默认值设置
- 类型转换

### 2. 数据模型 (models)

定义应用程序中使用的数据结构：

- `ProxyItem` - 代理服务项
- `ProxyList` - 代理服务列表

**主要功能：**
- 数据验证
- 业务逻辑封装
- 类型安全

### 3. 枚举类型 (enums)

定义应用程序中的枚举值：

- `ProxyType` - 代理类型枚举
- 服务类型常量

**主要功能：**
- 类型安全
- 常量管理
- 验证功能

### 4. HTTP 客户端 (httpclient)

提供统一的 HTTP 请求处理：

- 请求构建
- 响应解析
- 错误处理
- 超时管理

**主要功能：**
- GET/POST 请求
- JSON 序列化/反序列化
- 统一错误处理
- 请求头管理

### 5. 服务层 (services)

实现核心业务逻辑：

- `ProxyService` - 代理服务管理
- `DockerService` - Docker 加速
- `GitService` - Git 加速

**主要功能：**
- 命令解析
- 代理选择
- 命令执行
- 输出处理

## 开发指南

### 环境要求

- Go 1.23.4 或更高版本
- Git
- Docker (可选)

### 开发环境设置

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

### 代码规范

#### 1. 命名规范

- **包名**: 小写，简短，有意义
- **函数名**: 驼峰命名，公开函数首字母大写
- **变量名**: 驼峰命名，私有变量首字母小写
- **常量**: 全大写，下划线分隔

#### 2. 注释规范

- **包注释**: 每个包都应该有包级别的注释
- **函数注释**: 公开函数必须有注释
- **类型注释**: 公开类型必须有注释
- **变量注释**: 重要变量应该有注释

#### 3. 错误处理

- 使用 `fmt.Errorf` 包装错误
- 提供有意义的错误信息
- 使用 `errors.Is` 和 `errors.As` 检查错误类型

### 测试指南

#### 1. 单元测试

为每个包创建对应的测试文件：

```go
// internal/services/proxyservice_test.go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateProxyService(t *testing.T) {
    service := CreateProxyService("https://api.example.com")
    assert.NotNil(t, service)
    assert.NotNil(t, service.client)
}
```

#### 2. 集成测试

测试各个组件之间的交互：

```go
func TestProxyServiceIntegration(t *testing.T) {
    // 测试代理服务的完整流程
}
```

#### 3. 性能测试

使用 Go 的基准测试功能：

```go
func BenchmarkProxyService(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // 性能测试代码
    }
}
```

### 构建和部署

#### 1. 本地构建

```bash
# 构建当前平台
make local

# 构建所有平台
make all

# 构建特定平台
make linux
make windows
make darwin
```

#### 2. 交叉编译

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o cnfast-linux-amd64

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o cnfast-windows-amd64.exe
```

#### 3. 发布流程

1. 更新版本号
2. 运行测试
3. 构建二进制文件
4. 创建发布包
5. 上传到 GitHub Releases

### 贡献指南

#### 1. 提交代码

- 使用有意义的提交信息
- 遵循 Conventional Commits 规范
- 每个提交应该是一个完整的功能

#### 2. 代码审查

- 所有代码都需要经过审查
- 审查者应该检查代码质量和功能正确性
- 使用 GitHub Pull Request 进行代码审查

#### 3. 问题报告

- 使用 GitHub Issues 报告问题
- 提供详细的复现步骤
- 包含系统信息和错误日志

### 性能优化

#### 1. 内存优化

- 避免内存泄漏
- 及时释放资源
- 使用对象池复用对象

#### 2. 网络优化

- 连接复用
- 请求超时设置
- 重试机制

#### 3. 并发优化

- 使用 goroutine 处理并发请求
- 避免竞态条件
- 合理使用锁机制

### 调试技巧

#### 1. 日志记录

```go
import "log"

// 使用标准日志
log.Printf("Debug: %s", message)

// 使用调试模式
if config.Debug {
    fmt.Printf("Debug: %s\n", message)
}
```

#### 2. 错误追踪

```go
import "runtime/debug"

func handleError(err error) {
    if err != nil {
        log.Printf("Error: %v\nStack: %s", err, debug.Stack())
    }
}
```

#### 3. 性能分析

```go
import _ "net/http/pprof"

// 启动性能分析服务器
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

## 常见问题

### 1. 编译错误

- 检查 Go 版本
- 验证模块依赖
- 清理构建缓存

### 2. 运行时错误

- 检查网络连接
- 验证配置文件
- 查看错误日志

### 3. 性能问题

- 使用性能分析工具
- 检查内存使用
- 优化网络请求

## 参考资料

- [Go 官方文档](https://golang.org/doc/)
- [Go 最佳实践](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go 测试指南](https://golang.org/doc/tutorial/add-a-test)
- [Go 模块系统](https://golang.org/doc/modules/)
