package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus"
	"github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability"
)

func TestProcessEvent(t *testing.T) {
	bus := eventbus.New()
	logger := observability.NewLogger("test", observability.LevelError)
	w := NewIngestWorker(bus, logger)

	payload, _ := json.Marshal(map[string]interface{}{
		"source": "test",
	})

	event := eventbus.Event{
		ID:      "evt-001",
		Type:    "data.incoming",
		Payload: payload,
	}

	err := w.ProcessEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestProcessEventEmptyType(t *testing.T) {
	bus := eventbus.New()
	logger := observability.NewLogger("test", observability.LevelError)
	w := NewIngestWorker(bus, logger)

	event := eventbus.Event{
		ID:   "evt-002",
		Type: "",
	}

	err := w.ProcessEvent(context.Background(), event)
	if err == nil {
		t.Fatal("expected error for empty event type")
	}
}

func TestProcessEventWithPayload(t *testing.T) {
	bus := eventbus.New()
	logger := observability.NewLogger("test", observability.LevelError)
	w := NewIngestWorker(bus, logger)

	payload, _ := json.Marshal(map[string]interface{}{
		"batch_size": 100,
		"source":     "upstream",
	})

	event := eventbus.Event{
		ID:      "evt-003",
		Type:    "data.batch",
		Payload: payload,
	}

	err := w.ProcessEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
