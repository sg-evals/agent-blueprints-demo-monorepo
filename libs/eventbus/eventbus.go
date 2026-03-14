package eventbus

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// Event represents a domain event.
type Event struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Source    string          `json:"source"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// Handler processes an event.
type Handler func(ctx context.Context, event Event) error

// Bus is an in-memory event bus.
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

// New creates a new event bus.
func New() *Bus {
	return &Bus{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe registers a handler for an event type.
func (b *Bus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish sends an event to all registered handlers.
func (b *Bus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers := b.handlers[event.Type]
	b.mu.RUnlock()

	for _, h := range handlers {
		if err := h(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// SubscriberCount returns the number of subscribers for an event type.
func (b *Bus) SubscriberCount(eventType string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers[eventType])
}
