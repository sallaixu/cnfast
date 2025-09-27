// Package services 包含应用程序的核心服务逻辑
package services

import (
	"cnfast/config"
	"cnfast/internal/enums"
	"cnfast/internal/models"
	"cnfast/internal/pkg/help"
	"cnfast/internal/pkg/httpclient"
	"context"
	"fmt"
	"os"
	"strings"
)

// ProxyService 代理服务结构体
// 负责处理代理相关的业务逻辑，包括获取代理列表和执行加速命令
type ProxyService struct {
	client *httpclient.Client
}

// CreateProxyService 创建代理服务实例
// baseURL: API 服务器的基础地址
func CreateProxyService(baseURL string) *ProxyService {
	return &ProxyService{
		client: httpclient.New(baseURL),
	}
}

// getProxyList 获取指定类型的可用代理列表
// proxyType: 代理类型（docker 或 git）
// 返回代理列表和可能的错误
func (p *ProxyService) getProxyList(proxyType enums.ProxyType) ([]models.ProxyItem, error) {
	var proxyList []models.ProxyItem

	if config.Debug {
		fmt.Printf("正在查询 %s 类型的代理服务...\n", string(proxyType))
	}

	// 构建 API 请求路径
	endpoint := fmt.Sprintf("/api/proxy/list?type=%s", string(proxyType))

	// 发送 HTTP 请求获取代理列表
	err := p.client.Get(context.Background(), endpoint, &proxyList)
	if err != nil {
		return nil, fmt.Errorf("获取代理列表失败: %w", err)
	}

	// 验证代理列表
	if len(proxyList) == 0 {
		return nil, fmt.Errorf("未找到可用的 %s 代理服务", string(proxyType))
	}

	if config.Debug {
		fmt.Printf("成功获取 %d 个 %s 代理服务\n", len(proxyList), string(proxyType))
	}

	return proxyList, nil
}

// handlerCmd 处理命令行参数并执行相应的操作
// 返回 true 表示命令已处理，false 表示命令不支持
func (p *ProxyService) handlerCmd() error {
	// 没有参数时显示帮助信息
	if len(os.Args) == 1 {
		help.PrintHelp()
		return nil
	}

	firstArg := strings.ToLower(os.Args[1])

	switch firstArg {
	case "docker":
		return p.handleDockerCommand(true)
	case "docker-compose", "docker compose":
		return p.handleDockerCommand(false)
	case "git":
		return p.handleGitCommand()
	case "-v", "--version", "v", "version":
		help.PrintVersion()
		return nil
	case "-h", "--help", "h", "help":
		help.PrintHelp()
		return nil
	default:
		return fmt.Errorf("不支持的命令: %s\n使用 'cnfast --help' 查看可用命令", firstArg)
	}
}

// handleDockerCommand 处理 Docker 相关命令
func (p *ProxyService) handleDockerCommand(isDocker bool) error {
	// 获取 Docker 代理列表
	proxyList, err := p.getProxyList(enums.ServiceDocker)
	if err != nil {
		return fmt.Errorf("获取 Docker 代理服务失败: %w", err)
	}

	// 执行 Docker 代理
	DockerProxy(proxyList, isDocker)
	return nil
}

// handleGitCommand 处理 Git 相关命令
func (p *ProxyService) handleGitCommand() error {
	// 获取 Git 代理列表
	proxyList, err := p.getProxyList(enums.ServiceGit)
	if err != nil {
		return fmt.Errorf("获取 Git 代理服务失败: %w", err)
	}

	// 执行 Git 代理
	GitProxy(proxyList)
	return nil
}

// Start 启动代理服务
// 处理用户输入的命令并执行相应的加速操作
func (p *ProxyService) Start() error {
	if err := p.handlerCmd(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		help.PrintUsage()
		return err
	}
	return nil
}
