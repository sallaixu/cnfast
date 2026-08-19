// Package services 包含通用的命令执行逻辑
package services

import (
	"bufio"
	"cnfast/internal/models"
	"cnfast/internal/pkg/colors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CommandBuilder 命令构建函数类型
// 返回: cmd 命令对象, error 错误
type CommandBuilder func(proxy models.ProxyItem) (*exec.Cmd, string, error)

// ExecuteWithProxyRetry 使用代理列表重试执行命令的通用框架
// proxyList: 代理服务列表
// cmdBuilder: 命令构建函数，根据代理构建具体的命令
// actionName: 操作名称（如 "执行"、"下载" 等）
func ExecuteWithProxyRetry(proxyList []models.ProxyItem, cmdBuilder CommandBuilder, actionName string) {
	// 按评分排序代理列表
	sortedProxies := sortProxiesByScore(proxyList)

	// 尝试每个代理
	for i, proxy := range sortedProxies {
		fmt.Println(colors.Info(fmt.Sprintf("使用代理: %s (评分: %d)", colors.Cyan(proxy.GetDisplayName()), proxy.Score)))

		// 构建命令
		cmd, _, err := cmdBuilder(proxy)
		if err != nil {
			fmt.Println(colors.Error(fmt.Sprintf("构建命令失败: %v", err)))
			return
		}

		// 执行命令并输出（不再隐藏敏感信息）
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err == nil {
			fmt.Println(colors.Success(fmt.Sprintf("代理 %s %s成功", proxy.ID, actionName)))
			return
		}

		// 命令执行失败，检查是否还有更多代理可以尝试
		if i < len(sortedProxies)-1 {
			// 询问用户是否尝试下一个代理
			if askUserToRetry() {
				fmt.Println()
				fmt.Println(colors.Info("尝试下一个代理..."))
				fmt.Println()
				continue
			} else {
				fmt.Println(colors.Warn("用户取消操作"))
				os.Exit(1)
			}
		} else {
			// 所有代理都失败了
			fmt.Fprintf(os.Stderr, "\n%s\n", colors.Error(fmt.Sprintf("所有代理都%s失败，最后一个错误: %v", actionName, err)))
			os.Exit(1)
		}
	}
}

// ExecuteCommandWithOutput 已不再使用，保留占位以兼容旧代码（无实际逻辑）
func ExecuteCommandWithOutput(cmd *exec.Cmd, sensitiveInfo string) error {
	// 直接运行命令，输出完全由调用方配置
	return cmd.Run()
}

// askUserToRetry 询问用户是否重试
func askUserToRetry() bool {
	fmt.Print("\n❌是否尝试使用其他代理？(仅代理问题需要)(y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
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
