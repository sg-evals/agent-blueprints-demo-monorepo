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
		port = "8082"
	}
	logger := observability.NewLogger("billing-service", observability.LevelInfo)
	bus := eventbus.New()
	handler := NewBillingHandler(logger, bus)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/invoices", handler.CreateInvoice)
	mux.HandleFunc("/charges", handler.ListCharges)

	logger.Info("starting billing-service", map[string]interface{}{"port": port})
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
