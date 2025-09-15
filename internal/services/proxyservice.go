package services

import (
	"cnfast/internal/enums"
	"cnfast/internal/models"
	"cnfast/internal/pkg/help"
	"cnfast/internal/pkg/httpclient"
	"context"
	"fmt"
	"os"
	"strings"
)

var (
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
func (p ProxyService) getProxyList(proxyType enums.ProxyType) ([]models.ProxyItem, error) {
	var proxyList []models.ProxyItem
	fmt.Println("query proxy service " + string(proxyType))
	err := p.client.Get(context.Background(), "/api/proxy/list?type="+string(proxyType), &proxyList)
	return proxyList, err
}

func (p ProxyService) handlerCmd() bool {
	// 没有参数打印help
	if len(os.Args) == 1 {
		help.PrintHelp()
		return true
	}
	// 获取加速地址
	var list []models.ProxyItem
	var err error

	firstArg := strings.ToLower(os.Args[1])
	flag := false
	docker_flag := false
	switch firstArg {
	case "docker":
		docker_flag = true
		fallthrough
	case "docker-compose":
		fallthrough
	case "docker compose":
		flag = true
		list, err = p.getProxyList(enums.ServiceDocker)
		if err != nil {
			fmt.Printf("get proxy service error! \n %s", err)
			return true
		}
		DockerProxy(list, docker_flag)
	case "git":
		flag = true
		list, err = p.getProxyList(enums.ServiceGit)
		if err != nil {
			fmt.Printf("get proxy service error! \n %s", err)
			return true
		}
		GitProxy(list)
	case "-v":
		fallthrough
	case "-version":
		fallthrough
	case "v":
		flag = true
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

	if p.handlerCmd() {
		return
	} else {
		fmt.Printf("the commond is not supported !\n")
	}
}
