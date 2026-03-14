package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

func main() {
	logger := observability.NewLogger("worker-reconcile", observability.LevelInfo)
	bus := eventbus.New()
	worker := NewReconcileWorker(bus, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus.Subscribe("reconcile.request", func(_ context.Context, event eventbus.Event) error {
		var payload map[string]interface{}
		json.Unmarshal(event.Payload, &payload)
		sourceID, _ := payload["source_id"].(string)
		return worker.Reconcile(ctx, sourceID)
	})

	logger.Info("worker-reconcile started, waiting for events", nil)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down worker-reconcile")
}
