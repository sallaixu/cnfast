package services

import (
	"cnfast/internal/enums"
	"cnfast/internal/pkg/httpclient"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	proxyPrefix = "https://proxy.pipers.cn/"
)

type ProxyService struct {
	client *httpclient.Client
}

// 创建代理服务示例
func CreateProxyService(baseURL string) *ProxyService {
	return &ProxyService{
		client: httpclient.New(baseURL),
	}
}

// 获取可用代理列表
func (s *ProxyService) getProxyList(proxyType enums.ProxyType) ([]ProxyService, error) {
	fmt.Println("query proxy lsit " + string(proxyType))
	result := []ProxyService{}
	return result, nil
}

func checkCmd() bool {

	firstArg := os.Args[1]
	flag := false
	switch firstArg {
	case "-v":
		fmt.Println("")
		fmt.Println("------------------------------------------------")
		fmt.Println("cnfast: v1.0.0")
		fmt.Println("github: https://github.com/sallaixu/cnfast")
		fmt.Println("note  : 让每个想法都能连接世界")
		fmt.Println("------------------------------------------------")
		fmt.Println("")
		flag = true
	}
	return flag
}

// 启动服务
func (p ProxyService) Start() {

	if checkCmd() {
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("args lenght less than 2")
		os.Exit(1)
	}
	// 保留 git 子命令
	newArgs := []string{os.Args[1]}
	supportCmd := []string{"clone", "pull", "fetch"}

	found := false
	for _, s := range supportCmd {
		if s == newArgs[0] {
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("not support command %s", newArgs[0])
		fmt.Printf("supported command %s", strings.Join(supportCmd, ", "))
		os.Exit(1)
	}
	for _, arg := range os.Args[2:] {
		if strings.HasPrefix(arg, "https://github.com/") ||
			strings.HasPrefix(arg, "http://github.com/") {
			arg = proxyPrefix + arg
		}
		newArgs = append(newArgs, arg)
	}

	cmd := exec.Command("git", newArgs...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	_ = cmd.Run()
}
