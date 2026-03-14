package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

func TestOrdersHealth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewOrdersHandler(logger, bus)
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	h.Health(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCreateOrderNoAuth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewOrdersHandler(logger, bus)
	req := httptest.NewRequest("POST", "/orders", nil)
	rec := httptest.NewRecorder()
	h.CreateOrder(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestCreateOrderWithAuth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewOrdersHandler(logger, bus)
	req := httptest.NewRequest("POST", "/orders", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	rec := httptest.NewRecorder()
	h.CreateOrder(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
}

func TestGetOrder(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewOrdersHandler(logger, bus)
	req := httptest.NewRequest("GET", "/orders/ord-001", nil)
	rec := httptest.NewRecorder()
	h.GetOrder(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestGetOrderNoID(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewOrdersHandler(logger, bus)
	req := httptest.NewRequest("GET", "/orders/", nil)
	rec := httptest.NewRecorder()
	h.GetOrder(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}
