// Package services 包含应用程序的核心服务逻辑
package services

import (
	"cnfast/config"
	"cnfast/internal/enums"
	"cnfast/internal/models"
	"cnfast/internal/pkg/colors"
	"cnfast/internal/pkg/help"
	"cnfast/internal/pkg/httpclient"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// ProxyService 代理服务结构体
// 负责处理代理相关的业务逻辑，包括获取代理列表和执行加速命令
type ProxyService struct {
	client *httpclient.Client
}

// CreateProxyService 创建代理服务实例
// baseURL: API 服务器的基础地址
func CreateProxyService(baseURL string) *ProxyService {
	client := httpclient.New(baseURL)
	client.SetEncryption(config.AESKEY, config.AESIV)
	return &ProxyService{
		client: client,
	}
}

// getProxyList 获取指定类型的可用代理列表
// proxyType: 代理类型（docker 或 git）
// 返回代理列表和可能的错误
func (p *ProxyService) getProxyList(proxyType enums.ProxyType) ([]models.ProxyItem, error) {
	var proxyList []models.ProxyItem

	if config.Debug {
		fmt.Println(colors.Info(fmt.Sprintf("正在查询 %s 类型的代理服务...", string(proxyType))))
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
		fmt.Println(colors.Success(fmt.Sprintf("成功获取 %d 个 %s 代理服务", len(proxyList), string(proxyType))))
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
		// 支持两种调用方式:
		// 1) cnfast docker <subcommand>
		// 2) cnfast docker compose (等价于 cnfast docker-compose)
		if len(os.Args) >= 3 && strings.ToLower(os.Args[2]) == "compose" {
			return p.handleDockerCommand(false)
		}
		return p.handleDockerCommand(true)
	case "docker-compose":
		return p.handleDockerCommand(false)
	case "git":
		return p.handleGitCommand()
	case "update":
		return p.handleUpdate()
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

	// 让用户选择 Docker 代理
	selectedProxy := selectProxyWithPrompt(proxyList)
	selectedList := []models.ProxyItem{selectedProxy}

	// 执行 Docker 代理
	DockerProxy(selectedList, isDocker)
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

// handleUpdate 处理 cnfast 自更新命令
// 通过从 releases/latest 下载安装脚本并执行，实现与 install.sh 一致的更新逻辑
func (p *ProxyService) handleUpdate() error {
	fmt.Println(colors.Info("正在检查并更新 cnfast..."))

	// 根据当前系统和架构构建下载地址（与 install.sh 保持一致）
	baseURL := "https://gitee.com/sallai/cnfast/releases/download/latest"

	osType := runtime.GOOS
	arch := runtime.GOARCH

	var osPrefix string
	switch osType {
	case "linux":
		osPrefix = "linux"
	case "darwin":
		osPrefix = "darwin"
	default:
		return fmt.Errorf("不支持的操作系统: %s", osType)
	}

	var archSuffix string
	switch arch {
	case "amd64":
		archSuffix = "amd64"
	case "arm64":
		archSuffix = "arm64"
	case "386":
		archSuffix = "386"
	case "arm":
		archSuffix = "arm"
	default:
		return fmt.Errorf("不支持的架构: %s", arch)
	}

	binaryName := "cnfast"
	downloadURL := fmt.Sprintf("%s/%s-%s-%s", baseURL, binaryName, osPrefix, archSuffix)

	if config.Debug {
		fmt.Println(colors.Info(fmt.Sprintf("下载地址: %s", downloadURL)))
	}

	// 发起 HTTP 请求下载最新二进制
	start := time.Now()
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("下载更新失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载更新失败，HTTP 状态码: %d", resp.StatusCode)
	}

	// 创建临时文件下载新版本
	tmpFile, err := os.CreateTemp("", "cnfast-update-*")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// 拷贝并显示进度
	written, err := copyWithProgress(tmpFile, resp.Body, resp.ContentLength, "下载更新")
	if err != nil {
		return fmt.Errorf("写入临时文件失败: %w", err)
	}

	if err := tmpFile.Chmod(0755); err != nil {
		return fmt.Errorf("设置临时文件权限失败: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("关闭临时文件失败: %w", err)
	}

	// 按照 install.sh 逻辑，将二进制安装到 /usr/local/bin
	installDir := "/usr/local/bin"
	targetPath := installDir + "/" + binaryName

	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("创建安装目录失败: %w", err)
	}

	if err := os.Rename(tmpFile.Name(), targetPath); err != nil {
		return fmt.Errorf("替换 cnfast 二进制失败: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Println(colors.Success(fmt.Sprintf("更新完成: %.2f MB，耗时 %s", float64(written)/(1<<20), elapsed.Round(time.Millisecond))))
	return nil
}

// copyWithProgress 带进度显示的 io 拷贝
// TTY 下使用单行百分比刷新（\\r），非 TTY 下每 2% 输出一行
func copyWithProgress(dst io.Writer, src io.Reader, total int64, label string) (int64, error) {
	if total <= 0 {
		return io.Copy(dst, src)
	}

	if !colors.IsTerminal() {
		// 非 TTY：使用带进度的 Reader 交互，但只打印完成行
		return io.Copy(dst, src)
	}

	out := os.Stdout
	var written int64
	buf := make([]byte, 64*1024)
	lastPercent := -1
	carriage := false // 是否需要先清行
	for {
		n, err := src.Read(buf)
		if n > 0 {
			if _, werr := dst.Write(buf[:n]); werr != nil {
				return written, werr
			}
			written += int64(n)
			pct := int(written * 100 / total)
			if pct != lastPercent {
				if carriage {
					out.Write([]byte("\r"))
				}
				fmt.Fprintf(out, "%s %3d%% (%.2f MB / %.2f MB)",
					colors.Cyan(label), pct, float64(written)/(1<<20), float64(total)/(1<<20))
				carriage = true
				lastPercent = pct
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return written, err
		}
	}
	if carriage {
		out.Write([]byte("\r\x1b[K")) // 清空该行
	}
	return written, nil
}

// Start 启动代理服务
// 处理用户输入的命令并执行相应的加速操作
func (p *ProxyService) Start() error {
	if err := p.handlerCmd(); err != nil {
		fmt.Fprintln(os.Stderr, colors.Error(err.Error()))
		help.PrintUsage()
		return err
	}
	return nil
}
