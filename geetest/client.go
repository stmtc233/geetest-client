// geetest/client.go
package geetest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	// "net/url" // 不再需要 net/url
	"time"
)

// Client 是 geetest 服务的主要 API 客户端。
type Client struct {
	BaseURL    string
	SessionID  *string
	Proxy      *string
	User_Agent *string
	httpClient *http.Client

	// 用于不同验证类型的服务
	Click *ClickService
	Slide *SlideService
}

// ClientOption 是一个用于配置 Client 的函数。
type ClientOption func(*Client)

// NewClient 创建一个新的 geetest API 客户端。
func NewClient(baseURL string, options ...ClientOption) (*Client, error) {
	c := &Client{
		BaseURL: baseURL,
		httpClient: &http.Client{ // 使用默认的 Transport，不配置代理
			Timeout: 30 * time.Second,
		},
	}

	for _, option := range options {
		option(c)
	}

	// 错误逻辑已被移除：不再在 Go 客户端层面配置代理的 Transport
	// if c.Proxy != nil { ... }

	// 将具体服务链接回主客户端
	c.Click = &ClickService{client: c}
	c.Slide = &SlideService{client: c}

	return c, nil
}

// WithSessionID 为客户端发出的所有请求设置一个默认的会话 ID。
func WithSessionID(sessionID string) ClientOption {
	return func(c *Client) {
		c.SessionID = &sessionID
	}
}

// WithProxy 为客户端发出的所有请求设置一个默认的代理。
// 这个代理值将被包含在发送到 Rust 后端的 JSON 请求体中。
// 只支持http代理！！！
// "http://user:pass@host:port"
func WithProxy(proxy string) ClientOption {
	return func(c *Client) {
		c.Proxy = &proxy
	}
}

func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.User_Agent = &userAgent
	}
}

// WithHTTPClient 允许提供一个自定义的 http.Client。
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// post 函数保持不变，它只负责发送请求到你的 Rust 后端。

// ... post 和 HealthCheck 函数保持原样 ...
func post[T any](c *Client, ctx context.Context, path string, reqBody interface{}) (*T, error) {
	// 序列化请求体
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建请求
	requestURL := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 执行请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("收到非 200 状态码: %d - %s", resp.StatusCode, string(respBody))
	}

	// 将响应反序列化到我们的通用包装器中
	var apiResp ApiResponse[T]
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("反序列化响应失败: %w. Body: %s", err, string(respBody))
	}

	// 检查 API 级别的错误
	if !apiResp.Success {
		if apiResp.Error != nil {
			return nil, fmt.Errorf("API 错误: %s", *apiResp.Error)
		}
		return nil, fmt.Errorf("API 调用失败但没有错误信息")
	}

	return apiResp.Data, nil
}

func (c *Client) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("创建健康检查请求失败: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("健康检查请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("健康检查失败，状态码: %s", resp.Status)
	}
	return nil
}
