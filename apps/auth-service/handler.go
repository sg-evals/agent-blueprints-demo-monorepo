package main

import (
	"encoding/json"
	"net/http"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/authz"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/httpclient"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

// AuthHandler handles authentication requests.
type AuthHandler struct {
	logger *observability.Logger
	client *httpclient.Client
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(logger *observability.Logger) *AuthHandler {
	return &AuthHandler{
		logger: logger,
		client: httpclient.New(httpclient.WithUserAgent("auth-service/1.0")),
	}
}

// Health returns service health status.
func (h *AuthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "auth-service"})
}

// ValidateToken validates a bearer token.
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, err := authz.ParseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	h.logger.Info("token validated", map[string]interface{}{"token_len": len(token)})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":   true,
		"subject": "user-123",
		"expires": "2026-12-31T23:59:59Z",
	})
}

// CreateSession creates a new session.
func (h *AuthHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session := map[string]interface{}{
		"session_id": "sess-abc-123",
		"user_id":    "user-123",
		"created_at": "2026-01-01T00:00:00Z",
		"ttl":        3600,
	}

	h.logger.Info("session created", map[string]interface{}{"session_id": "sess-abc-123"})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}
