package main

import (
	"encoding/json"
	"net/http"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/featureflags"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/httpclient"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

// SearchHandler handles search requests.
type SearchHandler struct {
	logger *observability.Logger
	client *httpclient.Client
	flags  *featureflags.Store
}

// NewSearchHandler creates a new search handler.
func NewSearchHandler(logger *observability.Logger) *SearchHandler {
	flags := featureflags.NewStore()
	flags.Set("enhanced-search", false)
	return &SearchHandler{
		logger: logger,
		client: httpclient.New(httpclient.WithUserAgent("search-service/1.0")),
		flags:  flags,
	}
}

// Health returns service health status.
func (h *SearchHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "search-service"})
}

// Search performs a mock search with feature-flagged behavior.
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	h.logger.Info("searching", map[string]interface{}{"query": query, "enhanced": h.flags.Enabled("enhanced-search")})

	results := []map[string]interface{}{
		{"id": "doc-001", "title": "Getting Started", "score": 0.95},
		{"id": "doc-002", "title": "API Reference", "score": 0.87},
	}

	if h.flags.Enabled("enhanced-search") {
		results = append(results, map[string]interface{}{
			"id": "doc-003", "title": "Advanced Topics", "score": 0.72,
		})
		h.logger.Debug("enhanced search returned additional results", nil)
	}

	response := map[string]interface{}{
		"query":    query,
		"results":  results,
		"total":    len(results),
		"enhanced": h.flags.Enabled("enhanced-search"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
