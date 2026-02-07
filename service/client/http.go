package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// HTTPClient 提供一個配置好的HTTP客戶端
type HTTPClient struct {
	client  *http.Client
	timeout time.Duration
}

// Option 用於配置HTTPClient的函數類型
type Option func(*HTTPClient)

// WithTimeout 設置HTTP請求超時
func WithTimeout(timeout time.Duration) Option {
	return func(c *HTTPClient) {
		c.timeout = timeout
		c.client.Timeout = timeout
	}
}

// WithMaxIdleConns 設置最大空閒連接數
func WithMaxIdleConns(max int) Option {
	return func(c *HTTPClient) {
		c.client.Transport.(*http.Transport).MaxIdleConns = max
	}
}

// WithMaxConnsPerHost 設置每個主機的最大連接數
func WithMaxConnsPerHost(max int) Option {
	return func(c *HTTPClient) {
		c.client.Transport.(*http.Transport).MaxConnsPerHost = max
	}
}

// WithIdleConnTimeout 設置空閒連接超時時間
func WithIdleConnTimeout(timeout time.Duration) Option {
	return func(c *HTTPClient) {
		c.client.Transport.(*http.Transport).IdleConnTimeout = timeout
	}
}

// NewHTTPClient 創建一個新的HTTP客戶端
func NewHTTPClient(options ...Option) *HTTPClient {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	client := &HTTPClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		timeout: 30 * time.Second,
	}

	for _, opt := range options {
		opt(client)
	}

	return client
}

// NewHTTPClientWithEnv 創建一個從環境變量讀取配置的HTTP客戶端
func NewHTTPClientWithEnv(prefix string) *HTTPClient {
	timeout := getEnvDuration(prefix+"_HTTP_TIMEOUT", 30*time.Second)
	maxIdleConns := getEnvInt(prefix+"_HTTP_MAX_IDLE_CONNS", 100)
	idleConnTimeout := getEnvDuration(prefix+"_HTTP_IDLE_CONN_TIMEOUT", 90*time.Second)

	return NewHTTPClient(
		WithTimeout(timeout),
		WithMaxIdleConns(maxIdleConns),
		WithIdleConnTimeout(idleConnTimeout),
	)
}

// Get 發送GET請求
func (c *HTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return c.doRequest(req)
}

// GetWithRetry 發送GET請求並在失敗時重試
func (c *HTTPClient) GetWithRetry(ctx context.Context, url string, maxRetries int) ([]byte, error) {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			time.Sleep(time.Duration(i) * time.Second)
		}

		data, err := c.Get(ctx, url)
		if err == nil {
			return data, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

// doRequest 執行HTTP請求
func (c *HTTPClient) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// getEnvInt 從環境變量獲取整數值
func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

// getEnvDuration 從環境變量獲取時間值
func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	return defaultVal
}

// ValidateURL 驗證URL是否有效
func ValidateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("URL must have scheme and host")
	}

	return nil
}
