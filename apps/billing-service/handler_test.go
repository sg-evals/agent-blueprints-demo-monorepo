package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/authz"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

func TestBillingHealth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewBillingHandler(logger, bus)
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	h.Health(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCreateInvoiceNoAuth(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewBillingHandler(logger, bus)
	req := httptest.NewRequest("POST", "/invoices", nil)
	rec := httptest.NewRecorder()
	h.CreateInvoice(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestCreateInvoiceNoAdminRole(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewBillingHandler(logger, bus)
	req := httptest.NewRequest("POST", "/invoices", nil)
	req.Header.Set("Authorization", "Bearer some-token")
	rec := httptest.NewRecorder()
	h.CreateInvoice(rec, req)
	// No user in context, so expect 403
	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

func TestCreateInvoiceWithAdmin(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewBillingHandler(logger, bus)
	req := httptest.NewRequest("POST", "/invoices", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	user := &authz.User{ID: "admin-1", Email: "admin@test.com", Roles: []authz.Role{authz.RoleAdmin}}
	ctx := authz.WithUser(req.Context(), user)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h.CreateInvoice(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
}

func TestListCharges(t *testing.T) {
	logger := observability.NewLogger("test", observability.LevelError)
	bus := eventbus.New()
	h := NewBillingHandler(logger, bus)
	req := httptest.NewRequest("GET", "/charges", nil)
	rec := httptest.NewRecorder()
	h.ListCharges(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
