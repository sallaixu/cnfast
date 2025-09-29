// Package services 包含 Git 相关的服务逻辑
package services

import (
	"bufio"
	"bytes"
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
	supportedCommands := []string{"clone", "pull", "fetch", "push"}
	command := os.Args[2]

	// 检查命令是否支持
	if !isCommandSupported(command, supportedCommands) {
		fmt.Fprintf(os.Stderr, "错误: 不支持的命令 '%s'\n", command)
		fmt.Fprintf(os.Stderr, "支持的命令: %s\n", strings.Join(supportedCommands, ", "))
		os.Exit(1)
	}

	// 尝试执行 Git 命令，支持代理重试
	executeGitWithProxyRetry(proxyList, command)
}

// executeGitWithProxyRetry 执行 Git 命令，支持代理重试
func executeGitWithProxyRetry(proxyList []models.ProxyItem, command string) {
	// 按评分排序代理列表
	sortedProxies := sortProxiesByScore(proxyList)

	// 尝试每个代理
	for i, proxy := range sortedProxies {
		fmt.Printf("使用代理: %s (评分: %d)\n", proxy.GetDisplayName(), proxy.Score)

		// 构建加速后的参数
		newArgs := buildGitArgs(proxy.ProxyUrl, command)

		if config.Debug {
			fmt.Printf("执行命令: git %s\n", strings.Join(newArgs, " "))
		}
		// 提取主机名
		host := util.ExtractHostFromURL(proxy.ProxyUrl)
		// 执行 Git 命令
		cmd := exec.Command("git", newArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		// 运行命令
		err := cmd.Run()
		if err == nil {
			// 命令执行成功，直接返回
			return
		} else {
			fmt.Printf("命令执行失败: %v\n", strings.ReplaceAll(stderr.String(), host, "***"))
		}

		// 命令执行失败，检查是否还有更多代理可以尝试
		if i < len(sortedProxies)-1 {
			fmt.Fprintf(os.Stderr, "\n❌ 代理 %s 执行失败: %v\n", proxy.GetDisplayName(), err)

			// 询问用户是否尝试下一个代理
			if askUserToRetry() {
				fmt.Printf("\n🔄 尝试下一个代理...\n\n")
				continue
			} else {
				fmt.Println("用户取消操作")
				os.Exit(1)
			}
		} else {
			// 所有代理都失败了
			fmt.Fprintf(os.Stderr, "\n❌ 所有代理都执行失败，最后一个错误: %v\n", err)
			os.Exit(1)
		}
	}
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

// sortProxiesByScore 按评分排序代理列表
func sortProxiesByScore(proxyList []models.ProxyItem) []models.ProxyItem {
	// 创建副本避免修改原列表
	sorted := make([]models.ProxyItem, len(proxyList))
	copy(sorted, proxyList)

	// 简单的冒泡排序，按评分降序排列
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].Score < sorted[j+1].Score {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

// askUserToRetry 询问用户是否重试
func askUserToRetry() bool {
	fmt.Println("请根据错误信息判断是否是代理失效导致。代理问题可尝试使用其他代理")
	fmt.Print("是否尝试使用其他代理？(y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// isGitHubURL 检查 URL 是否为 GitHub URL
func isGitHubURL(url string) bool {
	return strings.HasPrefix(url, "https://github.com/") ||
		strings.HasPrefix(url, "http://github.com/")
}
