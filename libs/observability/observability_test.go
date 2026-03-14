package observability

import "testing"

func TestLoggerDoesNotPanic(t *testing.T) {
	l := NewLogger("test-service", LevelInfo)
	l.Info("starting", map[string]interface{}{"port": 8080})
	l.Error("failed", map[string]interface{}{"err": "timeout"})
	l.Debug("should be suppressed", nil)
}

func TestCounter(t *testing.T) {
	c := NewCounter("requests_total")
	c.Inc()
	c.Inc()
	c.Add(3)
	if c.Value() != 5 {
		t.Errorf("expected 5, got %d", c.Value())
	}
	if c.Name() != "requests_total" {
		t.Errorf("expected requests_total, got %s", c.Name())
	}
}
