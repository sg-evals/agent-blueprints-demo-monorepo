package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	logger := observability.NewLogger("search-service", observability.LevelInfo)
	handler := NewSearchHandler(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/search", handler.Search)

	logger.Info("starting search-service", map[string]interface{}{"port": port})
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
