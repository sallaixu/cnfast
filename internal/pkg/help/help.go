// Package help 提供命令行帮助信息
package help

import (
	"cnfast/config"
	"fmt"
)

// PrintHelp 显示应用程序的帮助信息
func PrintHelp() {
	fmt.Println("CNFast - 国内开发者网络加速工具")
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Printf("版本: %s\n", config.Version)
	fmt.Println("GitHub: https://github.com/sallaixu/cnfast")
	fmt.Println("描述: 让每个想法都能连接世界")
	fmt.Println()
	fmt.Println("用法: cnfast <command> [arguments]")
	fmt.Println()
	fmt.Println("命令:")
	fmt.Println("  git <command>       执行 Git 命令并加速 GitHub 仓库访问")
	fmt.Println("    clone <repo>      克隆 GitHub 仓库")
	fmt.Println("    pull              拉取最新更改")
	fmt.Println("    fetch             获取远程更改")
	fmt.Println()
	fmt.Println("  docker <command>    执行 Docker 命令并加速镜像拉取")
	fmt.Println("    pull <image>      拉取 Docker 镜像")
	fmt.Println()
	fmt.Println("  docker-compose      执行 docker-compose 命令并加速")
	fmt.Println()
	fmt.Println("  -v, --version       显示版本信息")
	fmt.Println("  -h, --help          显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  # GitHub 仓库加速")
	fmt.Println("  cnfast git clone https://github.com/user/repo.git")
	fmt.Println("  cnfast git pull")
	fmt.Println("  cnfast git fetch")
	fmt.Println()
	fmt.Println("  # Docker 镜像加速")
	fmt.Println("  cnfast docker pull nginx:latest")
	fmt.Println("  cnfast docker pull ubuntu:20.04")
	fmt.Println()
	fmt.Println("  # 查看版本")
	fmt.Println("  cnfast --version")
	fmt.Println()
	fmt.Println("环境变量:")
	fmt.Println("  CNFAST_API_HOST     设置 API 服务器地址")
	fmt.Println("  CNFAST_DEBUG        启用调试模式 (true/false)")
	fmt.Println("  CNFAST_TIMEOUT      设置请求超时时间（秒）")
	fmt.Println()
	fmt.Println("更多信息请访问: https://github.com/sallaixu/cnfast")
}

// PrintVersion 显示版本信息
func PrintVersion() {
	fmt.Println("------------------------------------------------")
	fmt.Printf("cnfast: v%s\n", config.Version)
	fmt.Println("github: https://github.com/sallaixu/cnfast")
	fmt.Println("note  : 让每个想法都能连接世界")
	fmt.Println("------------------------------------------------")
}

// PrintUsage 显示基本用法信息
func PrintUsage() {
	fmt.Println("用法: cnfast <command> [arguments]")
	fmt.Println("使用 'cnfast --help' 查看详细帮助信息")
}
