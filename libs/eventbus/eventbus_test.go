package eventbus

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestPublishSubscribe(t *testing.T) {
	bus := New()
	received := false
	bus.Subscribe("order.created", func(ctx context.Context, e Event) error {
		received = true
		return nil
	})

	err := bus.Publish(context.Background(), Event{
		ID:        "1",
		Type:      "order.created",
		Source:    "orders-service",
		Timestamp: time.Now(),
		Payload:   json.RawMessage(`{"order_id": "123"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !received {
		t.Error("handler not called")
	}
}

func TestSubscriberCount(t *testing.T) {
	bus := New()
	bus.Subscribe("test", func(ctx context.Context, e Event) error { return nil })
	bus.Subscribe("test", func(ctx context.Context, e Event) error { return nil })
	if bus.SubscriberCount("test") != 2 {
		t.Errorf("expected 2 subscribers, got %d", bus.SubscriberCount("test"))
	}
}

func TestPublishNoSubscribers(t *testing.T) {
	bus := New()
	err := bus.Publish(context.Background(), Event{Type: "unknown"})
	if err != nil {
		t.Errorf("expected no error for unsubscribed event, got %v", err)
	}
}
