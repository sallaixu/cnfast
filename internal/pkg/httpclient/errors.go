// Package httpclient 提供 HTTP 客户端错误处理功能
package httpclient

import (
	"fmt"
	"net/http"
)

// APIError 表示结构化的 API 错误响应
// 用于处理 HTTP 请求过程中的各种错误情况
type APIError struct {
	// StatusCode HTTP 状态码
	StatusCode int

	// Message 错误消息
	Message string

	// Details 错误详情
	Details string
}

// Error 实现 error 接口
func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API 错误 %d: %s - %s", e.StatusCode, e.Message, e.Details)
	}
	return fmt.Sprintf("API 错误 %d: %s", e.StatusCode, e.Message)
}

// IsNotFound 检查是否为 404 Not Found 错误
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsServerError 检查是否为 5xx 服务器错误
func IsServerError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode >= 500 && apiErr.StatusCode < 600
	}
	return false
}

// IsClientError 检查是否为 4xx 客户端错误
func IsClientError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode >= 400 && apiErr.StatusCode < 500
	}
	return false
}

// IsTimeout 检查是否为超时错误
func IsTimeout(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusRequestTimeout
	}
	return false
}

// NewAPIError 创建新的 APIError 实例
// statusCode: HTTP 状态码
// message: 错误消息
// details: 错误详情
func NewAPIError(statusCode int, message, details string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

// GetErrorType 获取错误类型描述
func GetErrorType(err error) string {
	if apiErr, ok := err.(*APIError); ok {
		switch {
		case apiErr.StatusCode >= 500:
			return "服务器错误"
		case apiErr.StatusCode >= 400:
			return "客户端错误"
		case apiErr.StatusCode >= 300:
			return "重定向"
		case apiErr.StatusCode >= 200:
			return "成功"
		default:
			return "未知错误"
		}
	}
	return "网络错误"
}
