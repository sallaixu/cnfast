// Package help 提供命令行帮助信息
package help

import (
	"cnfast/config"
	"cnfast/internal/pkg/colors"
	"cnfast/internal/pkg/util"
	"fmt"
)

// visualLen 计算去除 ANSI 转义码后的显示宽度（CJK 字符按 2 计）
func visualLen(s string) int { return util.VisualLen(s) }

// pad 将字符串右侧补齐到指定宽度（空格填充），保持颜色不影响对齐
func pad(s string, width int) string {
	n := visualLen(s)
	if n >= width {
		return s
	}
	return s + repeat(" ", width-n)
}

// repeat 返回重复 n 次的字符串
func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}

// PrintHelp 显示应用程序的帮助信息
func PrintHelp() {
	fmt.Println(colors.Bold(colors.Cyan("CNFast")) + " - " + colors.Bold("国内开发者网络加速工具"))
	fmt.Println(colors.Faint("====================================="))
	fmt.Println()
	fmt.Printf("%s %s\n", pad(colors.Bold("版本:"), 8), colors.Cyan(config.Version))
	fmt.Printf("%s %s\n", pad(colors.Bold("GitHub:"), 8), colors.Cyan("https://github.com/sallaixu/cnfast"))
	fmt.Printf("%s %s\n", pad(colors.Bold("描述:"), 8), "让每个想法都能连接世界")
	fmt.Println()
	fmt.Println(colors.Bold("用法:") + " cnfast <command> [arguments]")
	fmt.Println()
	fmt.Println(colors.Bold(colors.Cyan("命令:")))

	printRow := func(prefix, cmd, desc string) {
		line := "  " + prefix + colors.Green(cmd)
		if desc != "" {
			line += pad("", 2) + colors.Faint(desc)
		}
		fmt.Println(line)
	}

	fmt.Println("  " + colors.Cyan("git") + " " + colors.Faint("<command>") + "  " + colors.Faint("执行 Git 命令并加速 GitHub 仓库访问"))
	printRow("  ", "clone <repo>       ", "克隆 GitHub 仓库")
	printRow("  ", "pull               ", "拉取最新更改")
	printRow("  ", "down <url> [file]  ", "使用代理加速下载 GitHub Release 文件")
	fmt.Println()

	fmt.Println("  " + colors.Cyan("docker") + " " + colors.Faint("<command>") + "  " + colors.Faint("执行 Docker 命令并加速镜像拉取"))
	printRow("  ", "pull <image>       ", "拉取 Docker 镜像（支持加速域名与自动 retag）")
	printRow("  ", "push <image>       ", "推送 Docker 镜像（使用加速域名）")
	printRow("  ", "build ...          ", "构建镜像，保留原始行为")
	fmt.Println()

	fmt.Println("  " + colors.Yellow("docker-compose") + pad("", 1) + "  " + colors.Faint("解析 docker-compose.yml 中的镜像并加速拉取"))
	fmt.Println("  " + colors.Yellow("docker compose") + "  " + colors.Faint("等价于 docker-compose，用于兼容 Docker 新版命令"))
	fmt.Println()

	fmt.Println("  " + colors.Yellow("update") + pad("", 18) + "  " + colors.Faint("检查并更新到最新版本"))
	fmt.Println()

	fmt.Println("  " + colors.Green("-v, --version") + pad("", 10) + "  " + colors.Faint("显示版本信息"))
	fmt.Println("  " + colors.Green("-h, --help") + pad("", 13) + "  " + colors.Faint("显示此帮助信息"))
	fmt.Println()

	fmt.Println(colors.Bold(colors.Cyan("示例:")))
	fmt.Println("  " + colors.Faint("# GitHub 仓库加速"))
	fmt.Println("  " + colors.Label("cnfast") + " " + colors.Cyan("git") + " " + colors.Green("clone") + " https://github.com/user/repo.git")
	fmt.Println()
	fmt.Println("  " + colors.Faint("# Docker 镜像加速"))
	fmt.Println("  " + colors.Label("cnfast") + " " + colors.Cyan("docker") + " " + colors.Green("pull") + " nginx:latest")
	fmt.Println("  " + colors.Label("cnfast") + " " + colors.Cyan("docker") + " " + colors.Green("pull") + " ubuntu:20.04")
	fmt.Println()
	fmt.Println("  " + colors.Faint("# docker-compose 镜像加速"))
	fmt.Println("  " + colors.Label("cnfast") + " " + colors.Cyan("docker-compose"))
	fmt.Println("  " + colors.Label("cnfast") + " " + colors.Cyan("docker") + " " + colors.Cyan("compose"))
	fmt.Println()
	fmt.Println("  " + colors.Faint("# 更新 cnfast 自身"))
	fmt.Println("  " + colors.Label("cnfast") + " " + colors.Cyan("update"))
	fmt.Println()
	fmt.Println("  " + colors.Faint("# 查看版本"))
	fmt.Println("  " + colors.Label("cnfast") + " " + colors.Cyan("--version"))
	fmt.Println()
	fmt.Println()
	fmt.Println("更多信息请访问: " + colors.Cyan("https://github.com/sallaixu/cnfast"))
}

// PrintVersion 显示版本信息
func PrintVersion() {
	fmt.Println(colors.Faint("------------------------------------------------"))
	fmt.Printf("%s %s\n", colors.Bold("cnfast:"), colors.Cyan("v"+config.Version))
	fmt.Println("github: " + colors.Cyan("https://github.com/sallaixu/cnfast"))
	fmt.Println("note  : 让每个想法都能连接世界")
	fmt.Println(colors.Faint("------------------------------------------------"))
}

// PrintUsage 显示基本用法信息
func PrintUsage() {
	fmt.Println("用法: cnfast <command> [arguments]")
	fmt.Println("使用 'cnfast --help' 查看详细帮助信息")
}
