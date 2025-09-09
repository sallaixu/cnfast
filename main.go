package main

import (
	"cnfast/config"
	"cnfast/internal/services"
)

func main() {
	ProxyService := services.CreateProxyService(config.ApiHost)
	ProxyService.Start()
}
