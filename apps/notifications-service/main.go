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
		port = "8084"
	}
	logger := observability.NewLogger("notifications-service", observability.LevelInfo)
	bus := eventbus.New()
	handler := NewNotificationsHandler(logger, bus)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/notifications/send", handler.SendNotification)
	mux.HandleFunc("/notifications", handler.ListNotifications)

	logger.Info("starting notifications-service", map[string]interface{}{"port": port})
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
