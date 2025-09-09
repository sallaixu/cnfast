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

// Resp 后端统一包装体
type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// BizError 业务失败
type BizError struct {
	Code    int
	Message string
}

func (e *BizError) Error() string {
	return fmt.Sprintf("biz %d: %s", e.Code, e.Message)
}

/* ---------- Client ---------- */

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	headers    map[string]string
}

// New 创建客户端
func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}
}

// SetHeader 设置公共头
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

/* ---------- 内部：统一请求 & 解析 ---------- */

func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) ([]byte, error) {
	// 1. 构造请求
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	// 2. 头
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// 3. 发送
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	// 4. 解析统一包装
	var r Resp
	if err = json.Unmarshal(bodyBytes, &r); err != nil {
		return nil, fmt.Errorf("unmarshal resp: %w", err)
	}

	switch r.Code {
	case 0: // 成功
		if r.Data == nil {
			return nil, nil
		}
		// 再序列化 data 段返回
		return json.Marshal(r.Data)
	default: // 业务错误
		return nil, &BizError{Code: r.Code, Message: r.Message}
	}
}

/* ---------- 对外 API：只返回 data 段 ---------- */

// Get JSON 返回
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

// Post JSON 返回
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
