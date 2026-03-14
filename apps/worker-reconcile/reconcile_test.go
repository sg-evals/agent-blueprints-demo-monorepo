package main

import (
	"context"
	"testing"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

func TestReconcileSuccess(t *testing.T) {
	bus := eventbus.New()
	logger := observability.NewLogger("test", observability.LevelError)
	w := NewReconcileWorker(bus, logger)
	err := w.Reconcile(context.Background(), "source-123")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestReconcileEmptySource(t *testing.T) {
	bus := eventbus.New()
	logger := observability.NewLogger("test", observability.LevelError)
	w := NewReconcileWorker(bus, logger)
	err := w.Reconcile(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestReconcileWithBackoff(t *testing.T) {
	bus := eventbus.New()
	logger := observability.NewLogger("test", observability.LevelError)
	w := NewReconcileWorker(bus, logger)
	for attempt := 1; attempt <= 5; attempt++ {
		if err := w.ReconcileWithBackoff(attempt); err != nil {
			t.Errorf("attempt %d: unexpected error: %v", attempt, err)
		}
	}
}

func TestRetryBackoffZero(t *testing.T) {
	bus := eventbus.New()
	logger := observability.NewLogger("test", observability.LevelError)
	w := NewReconcileWorker(bus, logger)
	// Attempt 0 must not produce a negative backoff delay.
	if err := w.ReconcileWithBackoff(0); err != nil {
		t.Fatalf("attempt 0 produced invalid backoff: %v", err)
	}
}
