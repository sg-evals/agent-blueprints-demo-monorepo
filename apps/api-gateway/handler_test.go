package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

func TestHealth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	h := NewGatewayHandler(logger)
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	h.Health(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestProxyRequestNoAuth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	h := NewGatewayHandler(logger)
	req := httptest.NewRequest("GET", "/api/orders", nil)
	rec := httptest.NewRecorder()
	h.ProxyRequest(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestProxyRequestWithAuth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	h := NewGatewayHandler(logger)
	req := httptest.NewRequest("GET", "/api/orders", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	rec := httptest.NewRecorder()
	h.ProxyRequest(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
