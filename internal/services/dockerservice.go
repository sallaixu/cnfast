// Package services 包含 Docker 相关的服务逻辑
package services

import (
	"cnfast/config"
	"cnfast/internal/models"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Docker 镜像加速配置
var (
	// baseAccelDomain 基础加速域名
	baseAccelDomain = "docker.521456.xyz"

	// registryToAccelDomain 镜像源到加速域名的映射
	// 将各种 Docker registry 映射到对应的加速域名
	registryToAccelDomain = map[string]string{
		"quay.io":              "quay." + baseAccelDomain,
		"gcr.io":               "gcr." + baseAccelDomain,
		"k8s.gcr.io":           "k8s-gcr." + baseAccelDomain,
		"registry.k8s.io":      "k8s." + baseAccelDomain,
		"ghcr.io":              "ghcr." + baseAccelDomain,
		"docker.cloudsmith.io": "cloudsmith." + baseAccelDomain,
		"nvcr.io":              "nvcr." + baseAccelDomain,
		"registry-1.docker.io": baseAccelDomain,
		"docker.io":            baseAccelDomain, // 默认 Docker 官方仓库
	}

	// accelDomains 需要加速的域名列表
	accelDomains = getAccelDomains()
)

// getAccelDomains 获取需要加速的域名列表
func getAccelDomains() []string {
	domains := make([]string, 0, len(registryToAccelDomain))
	for domain := range registryToAccelDomain {
		domains = append(domains, domain)
	}
	return domains
}

// SetBaseAccelDomain 设置基础加速域名并重新生成映射
// domain: 新的基础加速域名
func SetBaseAccelDomain(domain string) {
	if config.Debug {
		fmt.Printf("设置代理域名: %s\n", domain)
	}
	baseAccelDomain = domain

	// 重新生成完整的加速域名映射
	registryToAccelDomain = map[string]string{
		"quay.io":              "quay." + baseAccelDomain,
		"gcr.io":               "gcr." + baseAccelDomain,
		"k8s.gcr.io":           "k8s-gcr." + baseAccelDomain,
		"registry.k8s.io":      "k8s." + baseAccelDomain,
		"ghcr.io":              "ghcr." + baseAccelDomain,
		"docker.cloudsmith.io": "cloudsmith." + baseAccelDomain,
		"nvcr.io":              "nvcr." + baseAccelDomain,
		"registry-1.docker.io": baseAccelDomain,
		"docker.io":            baseAccelDomain,
	}

	// 更新加速域名列表
	accelDomains = getAccelDomains()
}

// DockerProxy 执行 Docker 命令并应用镜像加速
// proxyList: 代理服务列表
// dockerFlag: 是否为 docker 命令（true）还是 docker-compose 命令（false）
func DockerProxy(proxyList []models.ProxyItem, dockerFlag bool) {
	// 如果不是 docker 命令，则处理 docker-compose
	if !dockerFlag {
		DockerComposeProxy(proxyList)
		return
	}

	// 检查命令参数数量
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "错误: 参数数量不足\n")
		fmt.Fprintf(os.Stderr, "用法: cnfast docker <command> [arguments]\n")
		os.Exit(1)
	}

	// 获取最佳代理服务
	bestProxy := getBestProxy(proxyList)
	if bestProxy == nil {
		fmt.Fprintf(os.Stderr, "错误: 未找到可用的代理服务\n")
		os.Exit(1)
	}

	fmt.Printf("使用最高评分代理: %s (评分: %d)\n", bestProxy.GetDisplayName(), bestProxy.Score)
	SetBaseAccelDomain(bestProxy.ProxyUrl)

	// 支持的命令列表
	supportedCommands := []string{"pull", "push", "build"}
	command := os.Args[2]

	// 检查命令是否支持
	if !isCommandSupported(command, supportedCommands) {
		fmt.Fprintf(os.Stderr, "错误: 不支持的命令 '%s'\n", command)
		fmt.Fprintf(os.Stderr, "支持的命令: %s\n", strings.Join(supportedCommands, ", "))
		os.Exit(1)
	}

	// 构建新的参数列表
	newArgs := []string{command}
	var originalImage string    // 原始镜像名
	var acceleratedImage string // 加速后的镜像名
	var needRetagging bool      // 是否需要重新打标签

	for idx, arg := range os.Args[3:] {
		// 如果是镜像参数，进行加速替换
		if idx == 0 && (command == "pull" || command == "push") {
			originalImage = arg
			acceleratedImage = replaceImageWithSpecificDomain(arg)
			if acceleratedImage != originalImage {
				fmt.Printf("镜像加速: %s -> %s\n", originalImage, acceleratedImage)
				needRetagging = (command == "pull") // 只有 pull 命令需要重新打标签
			}
			arg = acceleratedImage
		}
		newArgs = append(newArgs, arg)
	}

	if config.Debug {
		fmt.Printf("执行命令: docker %s\n", strings.Join(newArgs, " "))
	}

	// 执行 Docker 命令
	cmd := exec.Command("docker", newArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 运行命令
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "命令执行失败: %v\n", err)
		os.Exit(1)
	}

	// 如果是 pull 命令且使用了加速，需要重新打标签
	if needRetagging {
		retagImage(acceleratedImage, originalImage)
	}
}

// getBestProxy 获取评分最高的代理服务
func getBestProxy(proxyList []models.ProxyItem) *models.ProxyItem {
	if len(proxyList) == 0 {
		return nil
	}

	best := &proxyList[0]
	for i := 1; i < len(proxyList); i++ {
		if proxyList[i].Score > best.Score {
			best = &proxyList[i]
		}
	}
	return best
}

// isCommandSupported 检查命令是否在支持列表中
func isCommandSupported(command string, supportedCommands []string) bool {
	for _, cmd := range supportedCommands {
		if cmd == command {
			return true
		}
	}
	return false
}

// replaceImageWithSpecificDomain 根据映射表替换镜像域名
// raw: 原始镜像名称
// 返回: 加速后的镜像名称
func replaceImageWithSpecificDomain(raw string) string {
	// 检查镜像名称是否包含域名
	splits := strings.SplitN(raw, "/", 2)
	if len(splits) == 1 {
		// 没有域名，使用默认加速域名
		return baseAccelDomain + "/" + raw
	}

	// 检查第一部分是否是已知的 registry 域名
	if accelDomain, exists := registryToAccelDomain[splits[0]]; exists {
		return accelDomain + "/" + splits[1]
	}

	// 如果第一部分不是已知的 registry 域名，但有斜杠分隔，
	// 说明这是 Docker Hub 上的镜像（如 homeassistant/home-assistant:stable）
	// 应该使用默认加速域名
	return baseAccelDomain + "/" + raw
}

// retagImage 将加速域名的镜像重新打标签为原始名称
// acceleratedImage: 带加速域名的镜像名
// originalImage: 原始镜像名
func retagImage(acceleratedImage, originalImage string) {
	// 1. 使用原始名称重新打标签
	tagCmd := exec.Command("docker", "tag", acceleratedImage, originalImage)
	tagCmd.Stdout = os.Stdout
	tagCmd.Stderr = os.Stderr

	if err := tagCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "警告: 重新打标签失败: %v\n", err)
		fmt.Fprintf(os.Stderr, "镜像仍然可用，但标签为: %s\n", acceleratedImage)
		return
	}

	// 2. 删除加速域名的标签（清理临时标签）
	rmiCmd := exec.Command("docker", "rmi", acceleratedImage)
	// 不显示删除输出，保持界面简洁
	if config.Debug {
		rmiCmd.Stdout = os.Stdout
		rmiCmd.Stderr = os.Stderr
	}

	if err := rmiCmd.Run(); err != nil {
		if config.Debug {
			fmt.Fprintf(os.Stderr, "警告: 删除旧标签失败: %v\n", err)
		}
		// 忽略删除失败，因为不影响镜像使用
	}
}

// DockerComposeProxy 处理 docker-compose 命令的代理
// proxyList: 代理服务列表
func DockerComposeProxy(proxyList []models.ProxyItem) {
	// TODO: 实现 docker-compose 代理功能
	fmt.Println("docker-compose 代理功能尚未实现")
}
