package main

import (
	"cnfast/internal/services"
)

func main() {
	// Create a user service instance with JSONPlaceholder API base URL
	ProxyService := services.CreateProxyService("https://jsonplaceholder.typicode.com")
	ProxyService.Start()
}
