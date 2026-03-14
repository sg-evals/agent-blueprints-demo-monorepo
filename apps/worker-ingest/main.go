package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

func main() {
	logger := observability.NewLogger("worker-ingest", observability.LevelInfo)
	bus := eventbus.New()
	worker := NewIngestWorker(bus, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus.Subscribe("data.incoming", func(_ context.Context, event eventbus.Event) error {
		return worker.ProcessEvent(ctx, event)
	})

	logger.Info("worker-ingest started, waiting for events", nil)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("received signal %v, shutting down", sig)
}
