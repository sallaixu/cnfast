// Package services 包含 Git 相关的服务逻辑
package services

import (
	"cnfast/config"
	"cnfast/internal/models"
	"cnfast/internal/pkg/util"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Git 代理配置
var (
	// reHost 用于提取 URL 中主机名的正则表达式
	reHost = regexp.MustCompile(`^https?://([^/]+)`)

	// proxyPrefix 代理服务前缀
	proxyPrefix = "https://proxy.pipers.cn/"
)

// GitProxy 执行 Git 命令并应用 GitHub 加速
// proxyList: 代理服务列表
func GitProxy(proxyList []models.ProxyItem) {
	// 检查命令参数数量
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "错误: 参数数量不足\n")
		fmt.Fprintf(os.Stderr, "用法: cnfast git <command> [arguments]\n")
		os.Exit(1)
	}

	// 支持的命令列表
	supportedCommands := []string{"clone", "pull", "down"}
	command := os.Args[2]

	// 检查命令是否支持
	if !isCommandSupported(command, supportedCommands) {
		fmt.Fprintf(os.Stderr, "错误: 不支持的命令 '%s'\n", command)
		fmt.Fprintf(os.Stderr, "支持的命令: %s\n", strings.Join(supportedCommands, ", "))
		os.Exit(1)
	}

	// 处理 down 命令特殊逻辑
	if command == "down" {
		executeDownloadWithProxyRetry(proxyList)
		return
	}

	// 尝试执行 Git 命令，支持代理重试
	executeGitWithProxyRetry(proxyList, command)
}

// executeGitWithProxyRetry 执行 Git 命令，支持代理重试
func executeGitWithProxyRetry(proxyList []models.ProxyItem, command string) {
	// 使用通用的代理重试框架
	ExecuteWithProxyRetry(proxyList, func(proxy models.ProxyItem) (*exec.Cmd, string, error) {
		// 构建加速后的参数
		newArgs := buildGitArgs(proxy.ProxyUrl, command)

		if config.Debug {
			fmt.Printf("执行命令: git %s\n", strings.Join(newArgs, " "))
		}

		// 提取主机名（用于隐藏敏感信息）
		host := util.ExtractHostFromURL(proxy.ProxyUrl)

		// 执行 Git 命令
		cmd := exec.Command("git", append(newArgs, "--progress")...)

		return cmd, host, nil
	}, "执行")
}

// buildGitArgs 构建 Git 命令参数
func buildGitArgs(proxyUrl, command string) []string {
	newArgs := []string{}
	for _, arg := range os.Args[2:] {
		// 如果是 GitHub URL，进行加速替换
		if isGitHubURL(arg) {
			acceleratedURL := proxyUrl + "/" + arg
			if config.Debug {
				fmt.Printf("URL 加速: %s -> %s\n", arg, acceleratedURL)
			}
			arg = acceleratedURL
		}
		newArgs = append(newArgs, arg)
	}
	return newArgs
}

// isGitHubURL 检查 URL 是否为 GitHub URL
func isGitHubURL(url string) bool {
	return strings.HasPrefix(url, "https://github.com/") ||
		strings.HasPrefix(url, "http://github.com/")
}

// executeDownloadWithProxyRetry 使用代理下载文件，支持重试
func executeDownloadWithProxyRetry(proxyList []models.ProxyItem) {
	// 检查下载 URL 参数
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "错误: 缺少下载链接地址\n")
		fmt.Fprintf(os.Stderr, "用法: cnfast git down <下载链接地址> [输出文件名]\n")
		os.Exit(1)
	}

	downloadURL := os.Args[3]

	// 检查是否为 GitHub URL
	if !isGitHubURL(downloadURL) {
		fmt.Fprintf(os.Stderr, "错误: 仅支持 GitHub 链接下载\n")
		fmt.Fprintf(os.Stderr, "链接格式: https://github.com/...\n")
		os.Exit(1)
	}

	// 使用通用的代理重试框架
	ExecuteWithProxyRetry(proxyList, func(proxy models.ProxyItem) (*exec.Cmd, string, error) {
		// 构建代理后的下载地址
		proxiedURL := proxy.ProxyUrl + "/" + downloadURL

		if config.Debug {
			fmt.Printf("下载地址: %s\n", proxiedURL)
		}

		// 提取主机名用于隐藏敏感信息
		host := util.ExtractHostFromURL(proxy.ProxyUrl)

		// 构建 curl 命令参数
		curlArgs := []string{
			"-L",             // 跟随重定向
			"--progress-bar", // 显示进度条
			"-O",             // 使用远程文件名
			proxiedURL,
		}

		// 如果用户指定了输出文件名
		if len(os.Args) >= 5 {
			outputFile := os.Args[4]
			curlArgs = []string{
				"-L",
				"--progress-bar",
				"-o", outputFile,
				proxiedURL,
			}
		}

		if config.Debug {
			// 隐藏敏感信息的命令显示
			safeArgs := make([]string, len(curlArgs))
			copy(safeArgs, curlArgs)
			for j, arg := range safeArgs {
				safeArgs[j] = strings.ReplaceAll(arg, host, "***")
			}
			fmt.Printf("执行命令: curl %s\n", strings.Join(safeArgs, " "))
		}

		// 执行 curl 命令
		cmd := exec.Command("curl", curlArgs...)

		return cmd, host, nil
	}, "下载")
}
