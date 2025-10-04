// Package services 包含通用的命令执行逻辑
package services

import (
	"bufio"
	"bytes"
	"cnfast/internal/models"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// CommandBuilder 命令构建函数类型
// 返回: cmd 命令对象, sensitiveInfo 需要隐藏的敏感信息, error 错误
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
		fmt.Printf("使用代理: %s (评分: %d)\n", proxy.GetDisplayName(), proxy.Score)

		// 构建命令
		cmd, sensitiveInfo, err := cmdBuilder(proxy)
		if err != nil {
			fmt.Printf("构建命令失败: %v\n", err)
			return
		}

		// 执行命令并处理输出
		err = ExecuteCommandWithOutput(cmd, sensitiveInfo)

		if err == nil {
			fmt.Printf("✅ 代理 %s %s成功\n", proxy.ID, actionName)
			return
		}

		// 命令执行失败，检查是否还有更多代理可以尝试
		if i < len(sortedProxies)-1 {
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
			fmt.Fprintf(os.Stderr, "\n❌ 所有代理都%s失败，最后一个错误: %v\n", actionName, err)
			os.Exit(1)
		}
	}
}

// ExecuteCommandWithOutput 执行命令并实时处理输出，隐藏敏感信息
// cmd: 要执行的命令
// sensitiveInfo: 需要在输出中隐藏的敏感信息（如代理地址）
// 返回: error 执行错误
func ExecuteCommandWithOutput(cmd *exec.Cmd, sensitiveInfo string) error {
	// 设置标准输入
	cmd.Stdin = os.Stdin

	// 创建管道
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("创建stdout管道失败: %v\n", err)
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("创建stderr管道失败: %v\n", err)
		return err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		fmt.Printf("启动命令失败: %v\n", err)
		return err
	}

	// 实时读取并处理 stdout
	go streamPipeToOutput(stdoutPipe, os.Stdout, sensitiveInfo, "stdout")

	// 实时读取并处理 stderr
	go streamPipeToOutput(stderrPipe, os.Stderr, sensitiveInfo, "stderr")

	// 等待命令完成
	return cmd.Wait()
}

// streamPipeToOutput 实时读取管道内容并输出，同时隐藏敏感信息
// pipe: 输入管道
// output: 输出目标（如 os.Stdout 或 os.Stderr）
// sensitiveInfo: 需要隐藏的敏感信息
// pipeName: 管道名称（用于错误提示）
func streamPipeToOutput(pipe io.ReadCloser, output *os.File, sensitiveInfo string, pipeName string) {
	buf := make([]byte, 1024)
	for {
		n, err := pipe.Read(buf)
		if n > 0 {
			// 直接输出原始字节，保留控制字符
			content := buf[:n]
			// 替换敏感信息
			if sensitiveInfo != "" {
				processed := bytes.ReplaceAll(content, []byte(sensitiveInfo), []byte("***"))
				output.Write(processed)
			} else {
				output.Write(content)
			}
		}
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "读取%s错误: %v\n", pipeName, err)
			}
			break
		}
	}
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
