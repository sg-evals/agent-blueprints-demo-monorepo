package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/httpclient"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

// NotificationsHandler handles notification requests.
type NotificationsHandler struct {
	logger *observability.Logger
	bus    *eventbus.Bus
	client *httpclient.Client
}

// NewNotificationsHandler creates a new notifications handler.
func NewNotificationsHandler(logger *observability.Logger, bus *eventbus.Bus) *NotificationsHandler {
	return &NotificationsHandler{
		logger: logger,
		bus:    bus,
		client: httpclient.New(httpclient.WithUserAgent("notifications-service/1.0")),
	}
}

// Health returns service health status.
func (h *NotificationsHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "notifications-service"})
}

// SendNotification sends a mock notification.
func (h *NotificationsHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		To      string `json:"to"`
		Message string `json:"message"`
		Channel string `json:"channel"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if payload.To == "" || payload.Message == "" {
		http.Error(w, "to and message are required", http.StatusBadRequest)
		return
	}

	eventPayload, _ := json.Marshal(map[string]interface{}{
		"to":      payload.To,
		"channel": payload.Channel,
	})

	h.bus.Publish(context.Background(), eventbus.Event{
		ID:      "evt-notif-001",
		Type:    "notification.sent",
		Payload: eventPayload,
	})

	h.logger.Info("notification sent", map[string]interface{}{"to": payload.To, "channel": payload.Channel})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "sent", "to": payload.To})
}

// ListNotifications returns a list of mock notifications.
func (h *NotificationsHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	notifications := []map[string]interface{}{
		{"id": "notif-001", "to": "user@example.com", "channel": "email", "status": "delivered"},
		{"id": "notif-002", "to": "user-456", "channel": "push", "status": "pending"},
		{"id": "notif-003", "to": "user@example.com", "channel": "sms", "status": "delivered"},
	}

	h.logger.Info("listing notifications", map[string]interface{}{"count": len(notifications)})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}
