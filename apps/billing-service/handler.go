package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/authz"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

// BillingHandler handles billing-related requests.
type BillingHandler struct {
	logger *observability.Logger
	bus    *eventbus.Bus
}

// NewBillingHandler creates a new billing handler.
func NewBillingHandler(logger *observability.Logger, bus *eventbus.Bus) *BillingHandler {
	return &BillingHandler{logger: logger, bus: bus}
}

// Health returns service health status.
func (h *BillingHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "billing-service"})
}

// CreateInvoice creates a new invoice. Requires auth and admin role.
func (h *BillingHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := authz.ParseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, ok := authz.UserFromContext(r.Context())
	if !ok || !user.HasRole(authz.RoleAdmin) {
		http.Error(w, "forbidden: admin role required", http.StatusForbidden)
		return
	}

	invoice := map[string]interface{}{
		"id":     "inv-001",
		"amount": 99.99,
		"status": "created",
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"invoice_id": "inv-001",
	})

	h.bus.Publish(context.Background(), eventbus.Event{
		ID:      "evt-inv-001",
		Type:    "invoice.created",
		Payload: payload,
	})

	h.logger.Info("invoice created", map[string]interface{}{"invoice_id": "inv-001"})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)
}

// ListCharges returns a list of charges.
func (h *BillingHandler) ListCharges(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	charges := []map[string]interface{}{
		{"id": "chg-001", "amount": 49.99, "status": "completed"},
		{"id": "chg-002", "amount": 29.99, "status": "pending"},
	}

	h.logger.Info("listing charges", map[string]interface{}{"count": len(charges)})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(charges)
}
