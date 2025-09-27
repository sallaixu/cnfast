// Package services 包含 Git 相关的服务逻辑
package services

import (
	"cnfast/config"
	"cnfast/internal/models"
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

// extractHost 从 URL 中提取主机名
// rawURL: 完整的 URL 地址
// 返回: 主机名（包含端口）
func extractHost(rawURL string) string {
	matches := reHost.FindStringSubmatch(rawURL)
	if len(matches) < 2 {
		return ""
	}
	return matches[1] // 包含端口，如 ghproxy.com:8080
}

// GitProxy 执行 Git 命令并应用 GitHub 加速
// proxyList: 代理服务列表
func GitProxy(proxyList []models.ProxyItem) {
	// 检查命令参数数量
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "错误: 参数数量不足\n")
		fmt.Fprintf(os.Stderr, "用法: cnfast git <command> [arguments]\n")
		os.Exit(1)
	}

	// 获取最佳代理服务
	bestProxy := getBestProxy(proxyList)
	if bestProxy == nil {
		fmt.Fprintf(os.Stderr, "错误: 未找到可用的代理服务\n")
		os.Exit(1)
	}

	fmt.Printf("使用最高评分代理: %s (评分: %d)\n", bestProxy.GetDisplayName(), bestProxy.Score)
	proxyPrefix = bestProxy.ProxyUrl
	proxyHost := extractHost(proxyPrefix)

	if config.Debug {
		fmt.Printf("代理主机: %s\n", proxyHost)
	}

	// 支持的命令列表
	supportedCommands := []string{"clone", "pull", "fetch", "push"}
	command := os.Args[2]

	// 检查命令是否支持
	if !isCommandSupported(command, supportedCommands) {
		fmt.Fprintf(os.Stderr, "错误: 不支持的命令 '%s'\n", command)
		fmt.Fprintf(os.Stderr, "支持的命令: %s\n", strings.Join(supportedCommands, ", "))
		os.Exit(1)
	}

	// 构建新的参数列表
	newArgs := []string{}
	for _, arg := range os.Args[2:] {
		// 如果是 GitHub URL，进行加速替换
		if isGitHubURL(arg) {
			acceleratedURL := proxyPrefix + "/" + arg
			if config.Debug {
				fmt.Printf("URL 加速: %s -> %s\n", arg, acceleratedURL)
			}
			arg = acceleratedURL
		}
		newArgs = append(newArgs, arg)
	}

	if config.Debug {
		fmt.Printf("执行命令: git %s\n", strings.Join(newArgs, " "))
	}

	// 执行 Git 命令
	cmd := exec.Command("git", newArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 运行命令
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "命令执行失败: %v\n", err)
		os.Exit(1)
	}
}

// isGitHubURL 检查 URL 是否为 GitHub URL
func isGitHubURL(url string) bool {
	return strings.HasPrefix(url, "https://github.com/") ||
		strings.HasPrefix(url, "http://github.com/")
}
