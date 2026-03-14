package main

import (
	"context"
	"fmt"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/retry"
)

// IngestWorker processes incoming data events.
type IngestWorker struct {
	bus    *eventbus.Bus
	logger *observability.Logger
}

// NewIngestWorker creates a new ingest worker.
func NewIngestWorker(bus *eventbus.Bus, logger *observability.Logger) *IngestWorker {
	return &IngestWorker{bus: bus, logger: logger}
}

// ProcessEvent processes a single ingest event with retries.
func (w *IngestWorker) ProcessEvent(ctx context.Context, event eventbus.Event) error {
	cfg := retry.RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   retry.DefaultConfig().BaseDelay,
		MaxDelay:    retry.DefaultConfig().MaxDelay,
	}
	return retry.Do(func() error {
		w.logger.Info("processing event", map[string]interface{}{"event_id": event.ID, "type": event.Type})
		if event.Type == "" {
			return fmt.Errorf("empty event type")
		}
		return nil
	}, cfg)
}
