package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewDefaults(t *testing.T) {
	c := New()
	if c.httpClient.Timeout != 30*time.Second {
		t.Errorf("expected 30s timeout, got %v", c.httpClient.Timeout)
	}
	if c.userAgent != "agent-blueprints/1.0" {
		t.Errorf("expected default user agent, got %q", c.userAgent)
	}
}

func TestWithOptions(t *testing.T) {
	c := New(
		WithTimeout(5*time.Second),
		WithBaseURL("http://example.com"),
		WithUserAgent("test/1.0"),
	)
	if c.httpClient.Timeout != 5*time.Second {
		t.Error("timeout not set")
	}
	if c.baseURL != "http://example.com" {
		t.Error("base URL not set")
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") == "" {
			t.Error("missing user-agent")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	c := New(WithBaseURL(server.URL))
	body, status, err := c.Get(context.Background(), "/health")
	if err != nil {
		t.Fatal(err)
	}
	if status != 200 {
		t.Errorf("expected 200, got %d", status)
	}
	if string(body) != `{"status":"ok"}` {
		t.Errorf("unexpected body: %s", body)
	}
}
