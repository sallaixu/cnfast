// Package config 包含应用程序的配置信息
package config

import (
	"os"
	"strconv"
)

// 应用程序配置
var (
	// ApiHost API 服务器地址
	ApiHost = getEnvOrDefault("CNFAST_API_HOST", "https://cnfast-api.521456.xyz")

	// Debug 是否启用调试模式
	Debug = getBoolEnvOrDefault("CNFAST_DEBUG", false)

	// Timeout HTTP 请求超时时间（秒）
	Timeout = getIntEnvOrDefault("CNFAST_TIMEOUT", 30)

	// Version 应用程序版本
	Version = "1.0.0"

	// 加密参数
	AESKEY string
    AESIV  string
)

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getBoolEnvOrDefault 获取布尔类型环境变量，如果不存在则返回默认值
func getBoolEnvOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getIntEnvOrDefault 获取整数类型环境变量，如果不存在则返回默认值
func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
