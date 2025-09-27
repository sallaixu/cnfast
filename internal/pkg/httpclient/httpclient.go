// Package httpclient 提供 HTTP 客户端功能
// 包含统一的请求处理和错误处理机制
package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

/* ---------- 公共类型 ---------- */

// Resp 后端统一响应包装体
// 所有 API 响应都使用此格式
type Resp struct {
	// Code 响应状态码，0 表示成功
	Code int `json:"code"`

	// Message 响应消息
	Message string `json:"message"`

	// Data 响应数据
	Data interface{} `json:"data"`
}

// BizError 业务错误类型
// 用于表示 API 返回的业务错误
type BizError struct {
	// Code 错误代码
	Code int

	// Message 错误消息
	Message string
}

// Error 实现 error 接口
func (e *BizError) Error() string {
	return fmt.Sprintf("业务错误 %d: %s", e.Code, e.Message)
}

/* ---------- Client ---------- */

// Client HTTP 客户端
// 提供统一的 HTTP 请求处理功能
type Client struct {
	// BaseURL API 服务器基础地址
	BaseURL string

	// HTTPClient 底层 HTTP 客户端
	HTTPClient *http.Client

	// headers 公共请求头
	headers map[string]string
}

// New 创建新的 HTTP 客户端实例
// baseURL: API 服务器基础地址
func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}
}

// SetHeader 设置公共请求头
// key: 头部名称
// value: 头部值
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

/* ---------- 内部：统一请求 & 解析 ---------- */

// doRequest 执行统一的 HTTP 请求处理
// ctx: 上下文
// method: HTTP 方法
// endpoint: API 端点
// body: 请求体
// 返回: 响应数据字节数组和可能的错误
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) ([]byte, error) {
	// 1. 构造请求体
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// 2. 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 3. 设置请求头
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// 4. 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 5. 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 6. 解析统一响应格式
	var response Resp
	if err = json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 7. 处理响应状态
	switch response.Code {
	case 0: // 成功
		if response.Data == nil {
			return nil, nil
		}
		// 重新序列化 data 段返回
		return json.Marshal(response.Data)
	default: // 业务错误
		return nil, &BizError{Code: response.Code, Message: response.Message}
	}
}

/* ---------- 对外 API：只返回 data 段 ---------- */

// Get 执行 GET 请求
// ctx: 上下文
// endpoint: API 端点
// result: 用于接收响应数据的结构体指针
// 返回: 可能的错误
func (c *Client) Get(ctx context.Context, endpoint string, result interface{}) error {
	data, err := c.doRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	if result != nil && len(data) > 0 {
		return json.Unmarshal(data, result)
	}

	return nil
}

// Post 执行 POST 请求
// ctx: 上下文
// endpoint: API 端点
// body: 请求体数据
// result: 用于接收响应数据的结构体指针
// 返回: 可能的错误
func (c *Client) Post(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	data, err := c.doRequest(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}

	if result != nil && len(data) > 0 {
		return json.Unmarshal(data, result)
	}

	return nil
}
