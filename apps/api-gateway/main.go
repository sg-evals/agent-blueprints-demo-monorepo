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
		port = "8080"
	}
	logger := observability.NewLogger("api-gateway", observability.LevelInfo)
	handler := NewGatewayHandler(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/api/", handler.ProxyRequest)

	logger.Info("starting api-gateway", map[string]interface{}{"port": port})
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
