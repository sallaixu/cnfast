// Package main 是 CNFast 应用程序的入口点
// CNFast 是一个专为国内开发者设计的网络加速工具
// 主要功能包括 GitHub 仓库加速和 Docker 镜像加速
package main

import (
	"cnfast/config"
	"cnfast/internal/services"
	"fmt"
	"os"
)

// main 函数是应用程序的入口点
// 负责初始化代理服务并启动命令行处理
func main() {
	// 创建代理服务实例
	proxyService := services.CreateProxyService(config.ApiHost)

	// 启动服务，处理用户命令
	if err := proxyService.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting service: %v\n", err)
		os.Exit(1)
	}
}
