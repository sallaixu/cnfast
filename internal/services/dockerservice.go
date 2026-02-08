// Package services 包含 Docker 相关的服务逻辑
package services

import (
	"cnfast/config"
	"cnfast/internal/models"

	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
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

	// 代理服务在上游已选择，这里直接使用第一个
	if len(proxyList) == 0 {
		fmt.Fprintf(os.Stderr, "错误: 未找到可用的代理服务\n")
		os.Exit(1)
	}

	bestProxy := &proxyList[0]
	fmt.Printf("使用代理: %s (评分: %d)\n", bestProxy.ProxyUrl, bestProxy.Score)
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

// runComposeConfig 尝试兼容 docker compose 与 docker-compose 两种命令
// 返回命令输出（YAML 字节）和错误
func runComposeConfig(composeFile string) ([]byte, error) {
	// 优先尝试 docker compose
	cmd := exec.Command("docker", "compose", "-f", composeFile, "config")
	cmd.Stdin = os.Stdin
	output, err := cmd.CombinedOutput()
	if err == nil {
		return output, nil
	}

	// 兼容旧版 docker-compose 二进制
	cmd = exec.Command("docker-compose", "-f", composeFile, "config")
	cmd.Stdin = os.Stdin
	output2, err2 := cmd.CombinedOutput()
	if err2 == nil {
		return output2, nil
	}

	// 同时返回 docker compose 的错误，方便调试
	return output, fmt.Errorf("docker compose 失败: %v; docker-compose 失败: %v", err, err2)
}

// DockerComposeProxy 处理 docker-compose 命令的代理
// proxyList: 代理服务列表
func DockerComposeProxy(proxyList []models.ProxyItem) {
	if len(proxyList) == 0 {
		fmt.Fprintln(os.Stderr, "错误: 未找到可用的代理服务")
		os.Exit(1)
	}

	best := getBestProxy(proxyList)
	if best == nil {
		fmt.Fprintln(os.Stderr, "错误: 未找到可用的代理服务")
		os.Exit(1)
	}

	fmt.Printf("使用代理: %s (评分: %d)\n", best.ProxyUrl, best.Score)
	SetBaseAccelDomain(best.ProxyUrl)

	// 只考虑单 compose 文件，在当前目录按常见命名查找
	composeCandidates := []string{
		"docker-compose.yml",
		"docker-compose.yaml",
		"compose.yml",
		"compose.yaml",
	}

	var composeFile string
	for _, f := range composeCandidates {
		if _, err := os.Stat(f); err == nil {
			composeFile = f
			break
		}
	}

	if composeFile == "" {
		fmt.Fprintln(os.Stderr, "错误: 当前目录未找到 docker compose 配置文件 (docker-compose.yml|docker-compose.yaml|compose.yml|compose.yaml)")
		os.Exit(1)
	}

	if config.Debug {
		fmt.Printf("使用 compose 文件: %s\n", composeFile)
	}

	// 使用 docker compose/docker-compose CLI 解析配置为 YAML
	output, err := runComposeConfig(composeFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 解析 docker compose 配置失败: %v\n", err)
		fmt.Fprintf(os.Stderr, "命令输出:\n%s\n", string(output))
		os.Exit(1)
	}

	// 解析 YAML，提取 services -> image 映射
	type composeConfig struct {
		Services map[string]struct {
			Image string `yaml:"image"`
		} `yaml:"services"`
	}

	var cfg composeConfig
	if err := yaml.Unmarshal(output, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 解析 docker compose YAML 失败: %v\n", err)
		os.Exit(1)
	}

	// 构建去重后的镜像列表，同时记录使用该镜像的 service 名称
	type imageItem struct {
		Image    string
		Services []string
	}

	imageMap := make(map[string]*imageItem)

	for svcName, svc := range cfg.Services {
		if svc.Image == "" {
			continue // 没有 image 的服务（例如仅 build）忽略
		}
		if item, ok := imageMap[svc.Image]; ok {
			item.Services = append(item.Services, svcName)
		} else {
			imageMap[svc.Image] = &imageItem{
				Image:    svc.Image,
				Services: []string{svcName},
			}
		}
	}

	if len(imageMap) == 0 {
		fmt.Println("未在 compose 配置中找到任何需要拉取的镜像")
		return
	}

	images := make([]*imageItem, 0, len(imageMap))
	for _, item := range imageMap {
		images = append(images, item)
	}

	// 按镜像名排序，输出更稳定
	sort.Slice(images, func(i, j int) bool {
		return images[i].Image < images[j].Image
	})

	fmt.Println("发现以下镜像:")
	for i, item := range images {
		svcNames := strings.Join(item.Services, ", ")
		fmt.Printf("%d) %-20s -> %s\n", i+1, svcNames, item.Image)
	}

	// 让用户选择要拉取的镜像，支持多选，默认全部
	fmt.Print("请输入要拉取的镜像序号（多个用空格分隔，直接回车默认全部）: ")
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	var indices []int
	if line == "" {
		// 默认全选
		for i := range images {
			indices = append(indices, i)
		}
	} else {
		parts := strings.Fields(line)
		for _, p := range parts {
			n, err := strconv.Atoi(p)
			if err != nil || n < 1 || n > len(images) {
				fmt.Printf("输入无效: %s，已忽略\n", p)
				continue
			}
			indices = append(indices, n-1)
		}
		if len(indices) == 0 {
			fmt.Println("没有有效的序号，已取消操作")
			return
		}
	}

	// 对选中的镜像逐个执行加速拉取
	for _, idx := range indices {
		item := images[idx]
		original := item.Image
		accelerated := replaceImageWithSpecificDomain(original)

		needRetagging := false
		if accelerated != original {
			fmt.Printf("\n镜像加速: %s -> %s\n", original, accelerated)
			needRetagging = true
		}

		if config.Debug {
			fmt.Printf("执行命令: docker pull %s\n", accelerated)
		}

		pullCmd := exec.Command("docker", "pull", accelerated)
		pullCmd.Stdin = os.Stdin
		pullCmd.Stdout = os.Stdout
		pullCmd.Stderr = os.Stderr

		if err := pullCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "拉取镜像失败 (%s): %v\n", original, err)
			// 失败时继续尝试下一个镜像
			continue
		}

		if needRetagging {
			retagImage(accelerated, original)
		}
	}
}
