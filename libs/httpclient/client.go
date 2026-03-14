package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is an instrumented HTTP client.
type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
}

// Option configures the client.
type Option func(*Client)

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

// WithBaseURL sets the base URL for requests.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithUserAgent sets the User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// New creates a new instrumented HTTP client.
func New(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "agent-blueprints/1.0",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string) ([]byte, int, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading response: %w", err)
	}
	return body, resp.StatusCode, nil
}

// BaseURL returns the configured base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}
