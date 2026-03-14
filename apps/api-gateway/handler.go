package main

import (
	"encoding/json"
	"net/http"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/authz"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/featureflags"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/httpclient"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

// GatewayHandler handles API gateway requests.
type GatewayHandler struct {
	logger *observability.Logger
	client *httpclient.Client
	flags  *featureflags.Store
}

// NewGatewayHandler creates a new gateway handler.
func NewGatewayHandler(logger *observability.Logger) *GatewayHandler {
	flags := featureflags.NewStore()
	flags.Set("rate-limiting", true)
	return &GatewayHandler{
		logger: logger,
		client: httpclient.New(httpclient.WithUserAgent("api-gateway/1.0")),
		flags:  flags,
	}
}

// Health returns service health status.
func (h *GatewayHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "api-gateway"})
}

// ProxyRequest routes requests to downstream services.
func (h *GatewayHandler) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	token, err := authz.ParseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	h.logger.Info("proxying request", map[string]interface{}{"path": r.URL.Path, "token_len": len(token)})

	if h.flags.Enabled("rate-limiting") {
		h.logger.Debug("rate limiting enabled", nil)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "routed", "path": r.URL.Path})
}
