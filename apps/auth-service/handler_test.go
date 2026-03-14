package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

func TestAuthHealth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	h := NewAuthHandler(logger)
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	h.Health(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestValidateTokenNoAuth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	h := NewAuthHandler(logger)
	req := httptest.NewRequest("POST", "/validate", nil)
	rec := httptest.NewRecorder()
	h.ValidateToken(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestValidateTokenWithAuth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	h := NewAuthHandler(logger)
	req := httptest.NewRequest("POST", "/validate", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()
	h.ValidateToken(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCreateSession(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	h := NewAuthHandler(logger)
	req := httptest.NewRequest("POST", "/sessions", nil)
	rec := httptest.NewRecorder()
	h.CreateSession(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
}

func TestCreateSessionWrongMethod(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	h := NewAuthHandler(logger)
	req := httptest.NewRequest("GET", "/sessions", nil)
	rec := httptest.NewRecorder()
	h.CreateSession(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
