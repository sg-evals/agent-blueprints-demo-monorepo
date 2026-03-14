package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

func TestSearchHealth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	h := NewSearchHandler(logger)
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	h.Health(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSearchWithQuery(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	h := NewSearchHandler(logger)
	req := httptest.NewRequest("GET", "/search?q=test", nil)
	rec := httptest.NewRecorder()
	h.Search(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var result map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&result)
	if result["query"] != "test" {
		t.Errorf("expected query 'test', got %v", result["query"])
	}
}

func TestSearchNoQuery(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	h := NewSearchHandler(logger)
	req := httptest.NewRequest("GET", "/search", nil)
	rec := httptest.NewRecorder()
	h.Search(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestSearchEnhanced(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	h := NewSearchHandler(logger)
	h.flags.Set("enhanced-search", true)
	req := httptest.NewRequest("GET", "/search?q=test", nil)
	rec := httptest.NewRecorder()
	h.Search(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var result map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&result)
	total := int(result["total"].(float64))
	if total != 3 {
		t.Errorf("expected 3 results with enhanced search, got %d", total)
	}
}
