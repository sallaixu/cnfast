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
	fmt.Println("  git <command>          执行 Git 命令并加速 GitHub 仓库访问")
	fmt.Println("    clone <repo>         克隆 GitHub 仓库")
	fmt.Println("    pull                 拉取最新更改")
	fmt.Println("    down <url> [file]    使用代理加速下载 GitHub Release 文件")
	fmt.Println()
	fmt.Println("  docker <command>       执行 Docker 命令并加速镜像拉取")
	fmt.Println("    pull <image>         拉取 Docker 镜像（支持加速域名与自动 retag）")
	fmt.Println("    push <image>         推送 Docker 镜像（使用加速域名）")
	fmt.Println("    build ...            构建镜像，保留原始行为")
	fmt.Println()
	fmt.Println("  docker-compose         解析 docker-compose.yml 中的镜像并加速拉取")
	fmt.Println("  docker compose         等价于 docker-compose，用于兼容 Docker 新版命令")
	fmt.Println()
	fmt.Println("  update                 检查并更新到最新版本")
	fmt.Println()
	fmt.Println("  -v, --version          显示版本信息")
	fmt.Println("  -h, --help             显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  # GitHub 仓库加速")
	fmt.Println("  cnfast git clone https://github.com/user/repo.git")
	fmt.Println()
	fmt.Println("  # Docker 镜像加速")
	fmt.Println("  cnfast docker pull nginx:latest")
	fmt.Println("  cnfast docker pull ubuntu:20.04")
	fmt.Println()
	fmt.Println("  # docker-compose 镜像加速")
	fmt.Println("  cnfast docker-compose")
	fmt.Println("  cnfast docker compose")
	fmt.Println()
	fmt.Println("  # 更新 cnfast 自身")
	fmt.Println("  cnfast update")
	fmt.Println()
	fmt.Println("  # 查看版本")
	fmt.Println("  cnfast --version")
	fmt.Println()
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
