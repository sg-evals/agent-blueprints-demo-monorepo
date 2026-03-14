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
		port = "8081"
	}
	logger := observability.NewLogger("auth-service", observability.LevelInfo)
	handler := NewAuthHandler(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/validate", handler.ValidateToken)
	mux.HandleFunc("/sessions", handler.CreateSession)

	logger.Info("starting auth-service", map[string]interface{}{"port": port})
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
