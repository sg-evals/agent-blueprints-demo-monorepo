package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	logger := observability.NewLogger("orders-service", observability.LevelInfo)
	bus := eventbus.New()
	handler := NewOrdersHandler(logger, bus)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/orders", handler.CreateOrder)
	mux.HandleFunc("/orders/", handler.GetOrder)

	logger.Info("starting orders-service", map[string]interface{}{"port": port})
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
