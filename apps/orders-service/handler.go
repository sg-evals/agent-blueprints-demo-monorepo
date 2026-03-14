package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/authz"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/retry"
)

// OrdersHandler handles order-related requests.
type OrdersHandler struct {
	logger *observability.Logger
	bus    *eventbus.Bus
}

// NewOrdersHandler creates a new orders handler.
func NewOrdersHandler(logger *observability.Logger, bus *eventbus.Bus) *OrdersHandler {
	return &OrdersHandler{logger: logger, bus: bus}
}

// Health returns service health status.
func (h *OrdersHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "orders-service"})
}

// CreateOrder creates a new order and publishes an event with retry.
func (h *OrdersHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := authz.ParseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	order := map[string]interface{}{
		"id":     "ord-001",
		"status": "created",
		"total":  149.99,
	}

	cfg := retry.RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   retry.DefaultConfig().BaseDelay,
		MaxDelay:    retry.DefaultConfig().MaxDelay,
	}

	err = retry.Do(func() error {
		payload, _ := json.Marshal(map[string]interface{}{
			"order_id": "ord-001",
		})
		return h.bus.Publish(context.Background(), eventbus.Event{
			ID:      "evt-ord-001",
			Type:    "order.created",
			Payload: payload,
		})
	}, cfg)

	if err != nil {
		h.logger.Error("failed to publish order event", map[string]interface{}{"error": err.Error()})
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("order created", map[string]interface{}{"order_id": "ord-001"})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetOrder returns a mock order by ID.
func (h *OrdersHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	orderID := parts[len(parts)-1]
	if orderID == "" || orderID == "orders" {
		http.Error(w, "order ID required", http.StatusBadRequest)
		return
	}

	order := map[string]interface{}{
		"id":     orderID,
		"status": "completed",
		"total":  149.99,
		"items":  []string{"item-1", "item-2"},
	}

	h.logger.Info("fetched order", map[string]interface{}{"order_id": orderID})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
