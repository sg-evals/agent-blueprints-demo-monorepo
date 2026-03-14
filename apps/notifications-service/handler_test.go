package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

func TestNotificationsHealth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewNotificationsHandler(logger, bus)
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	h.Health(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSendNotification(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewNotificationsHandler(logger, bus)
	body := `{"to":"user@example.com","message":"hello","channel":"email"}`
	req := httptest.NewRequest("POST", "/notifications/send", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.SendNotification(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
}

func TestSendNotificationMissingFields(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewNotificationsHandler(logger, bus)
	body := `{"to":"","message":""}`
	req := httptest.NewRequest("POST", "/notifications/send", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.SendNotification(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestListNotifications(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewNotificationsHandler(logger, bus)
	req := httptest.NewRequest("GET", "/notifications", nil)
	rec := httptest.NewRecorder()
	h.ListNotifications(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
