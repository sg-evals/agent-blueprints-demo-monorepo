package observability

import "sync/atomic"

// Counter is a simple atomic counter metric.
type Counter struct {
	name  string
	value int64
}

// NewCounter creates a new counter.
func NewCounter(name string) *Counter {
	return &Counter{name: name}
}

// Inc increments the counter by 1.
func (c *Counter) Inc() {
	atomic.AddInt64(&c.value, 1)
}

// Add adds delta to the counter.
func (c *Counter) Add(delta int64) {
	atomic.AddInt64(&c.value, delta)
}

// Value returns the current counter value.
func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// Name returns the counter name.
func (c *Counter) Name() string {
	return c.name
}
