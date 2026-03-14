package main

import (
	"context"
	"fmt"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/retry"
)

// ReconcileWorker handles data reconciliation.
type ReconcileWorker struct {
	bus    *eventbus.Bus
	logger *observability.Logger
}

// NewReconcileWorker creates a new reconcile worker.
func NewReconcileWorker(bus *eventbus.Bus, logger *observability.Logger) *ReconcileWorker {
	return &ReconcileWorker{bus: bus, logger: logger}
}

// Reconcile performs reconciliation with retry logic.
func (w *ReconcileWorker) Reconcile(ctx context.Context, sourceID string) error {
	cfg := retry.DefaultConfig()
	return retry.Do(func() error {
		w.logger.Info("reconciling", map[string]interface{}{"source_id": sourceID})
		if sourceID == "" {
			return fmt.Errorf("empty source ID")
		}
		return nil
	}, cfg)
}

// ReconcileWithBackoff demonstrates direct backoff usage.
func (w *ReconcileWorker) ReconcileWithBackoff(attempt int) error {
	cfg := retry.DefaultConfig()
	cfg.Jitter = false
	delay := retry.RetryBackoff(attempt, cfg)
	w.logger.Info("backoff calculated", map[string]interface{}{
		"attempt": attempt,
		"delay":   delay.String(),
	})
	if delay < 0 {
		return fmt.Errorf("invalid negative backoff delay: %v", delay)
	}
	return nil
}
