package authz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHasRole(t *testing.T) {
	u := &User{ID: "1", Roles: []Role{RoleAdmin, RoleUser}}
	if !u.HasRole(RoleAdmin) {
		t.Error("expected admin role")
	}
	if u.HasRole(RoleViewer) {
		t.Error("did not expect viewer role")
	}
}

func TestUserContext(t *testing.T) {
	u := &User{ID: "42", Email: "test@example.com"}
	ctx := WithUser(context.Background(), u)
	got, ok := UserFromContext(ctx)
	if !ok || got.ID != "42" {
		t.Error("expected user from context")
	}
}

func TestRequireRoleUnauthorized(t *testing.T) {
	handler := RequireRole(RoleAdmin, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestParseAuthHeader(t *testing.T) {
	token, err := ParseAuthHeader("Bearer abc123")
	if err != nil || token != "abc123" {
		t.Errorf("expected abc123, got %q, err=%v", token, err)
	}
	_, err = ParseAuthHeader("Basic xyz")
	if err == nil {
		t.Error("expected error for non-bearer")
	}
}
