package services

import (
	"bufio"
	"cnfast/internal/models"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// 基础加速域名（可配置）
var baseAccelDomain = "docker.521456.xyz"

// 镜像源到加速域名的完整映射（在定义时就拼接完整）
var registryToAccelDomain = map[string]string{
	"quay.io":              "quay." + baseAccelDomain,
	"gcr.io":               "gcr." + baseAccelDomain,
	"k8s.gcr.io":           "k8s-gcr." + baseAccelDomain,
	"registry.k8s.io":      "k8s." + baseAccelDomain,
	"ghcr.io":              "ghcr." + baseAccelDomain,
	"docker.cloudsmith.io": "cloudsmith." + baseAccelDomain,
	"nvcr.io":              "nvcr." + baseAccelDomain,
	"registry-1.docker.io": baseAccelDomain,
	"docker.io":            baseAccelDomain, // 默认docker官方仓库
}

// 需要加速的域名列表（从映射中提取）
var accelDomains = getAccelDomains()

func getAccelDomains() []string {
	domains := make([]string, 0, len(registryToAccelDomain))
	for domain := range registryToAccelDomain {
		domains = append(domains, domain)
	}
	return domains
}

// SetBaseAccelDomain 设置基础加速域名并重新生成映射
func SetBaseAccelDomain(domain string) {
	fmt.Printf("set proxy domain:%s \n", domain)
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

func DockerProxy(proxyList []models.ProxyItem, dockerFlag bool) {

	if !dockerFlag {
		DockerComposeProxy(proxyList)
		return
	}

	// 检查git命令合法性
	if len(os.Args) < 3 {
		fmt.Printf("args lenght less than 2 \n")
		os.Exit(1)
	}

	// 获取代理地址
	fmt.Printf("use max score proxy: %s\n", proxyList[0].ProxyUrl)
	proxyHost := proxyList[0].ProxyUrl
	SetBaseAccelDomain(proxyHost)
	// 支持的命令列表
	supportCmd := []string{"pull"}

	// 检查是否支持该命令
	command := os.Args[2]
	found := false
	for _, s := range supportCmd {
		if s == command {
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("not support command: %s\n", command)
		fmt.Printf("supported commands: %s\n", strings.Join(supportCmd, ", "))
		os.Exit(1)
	}

	// 构建新的参数列表
	newArgs := []string{command}
	for idx, arg := range os.Args[3:] {
		// 如果是镜像参数，进行替换加速
		if idx == 0 {
			acceleratedImage := replaceImageWithSpecificDomain(arg)
			fmt.Printf("accelerate image: %s -> %s\n", arg, acceleratedImage)
			arg = acceleratedImage
		}
		newArgs = append(newArgs, arg)
	}

	fmt.Printf("executing: docker %s\n", strings.Join(newArgs, " "))

	// 执行docker命令
	cmd := exec.Command("docker", newArgs...)
	pr, pw := io.Pipe()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = pw

	// 启动读goroutine过滤代理host信息
	done := make(chan struct{})
	go func() {
		defer close(done)
		sc := bufio.NewScanner(pr)
		for sc.Scan() {
			line := sc.Text()
			// 过滤掉含代理host的行
			if strings.Contains(line, proxyHost) {
				continue
			}
			fmt.Fprintln(os.Stderr, line)
		}
	}()

	// 运行命令
	if err := cmd.Run(); err != nil {
		pw.Close()
		<-done
		fmt.Printf("command failed: %v\n", err)
		os.Exit(1)
	}
	pw.Close()
	<-done

}

// replaceImageWithSpecificDomain 根据映射表替换镜像域名
func replaceImageWithSpecificDomain(raw string) string {
	// 检查每个需要加速的域名
	splits := strings.SplitN(raw, "/", 2)
	if len(splits) == 1 {
		raw = baseAccelDomain + "/" + raw
	} else {
		val, ok := registryToAccelDomain[splits[0]]
		if ok {
			raw = val + "/" + splits[1]
		} else {
			raw = val + "/" + raw
		}
	}
	// 不需要加速的镜像原样返回
	return raw
}

func DockerComposeProxy(proxyList []models.ProxyItem) {

}
