package retry

import (
	"testing"
	"time"
)

func TestRetryBackoffPositive(t *testing.T) {
	cfg := RetryConfig{
		BaseDelay: 100 * time.Millisecond,
		MaxDelay:  10 * time.Second,
		Jitter:    false,
	}
	for attempt := 1; attempt <= 10; attempt++ {
		d := RetryBackoff(attempt, cfg)
		if d < 0 {
			t.Errorf("attempt %d: got negative duration %v", attempt, d)
		}
	}
}

func TestRetryBackoffMaxDelay(t *testing.T) {
	cfg := RetryConfig{
		BaseDelay: 1 * time.Second,
		MaxDelay:  5 * time.Second,
		Jitter:    false,
	}
	d := RetryBackoff(100, cfg)
	if d > cfg.MaxDelay {
		t.Errorf("expected delay <= %v, got %v", cfg.MaxDelay, d)
	}
}

func TestRetryBackoffClampsZeroAttempt(t *testing.T) {
	cfg := RetryConfig{
		BaseDelay: 100 * time.Millisecond,
		MaxDelay:  10 * time.Second,
		Jitter:    false,
	}
	d := RetryBackoff(0, cfg)
	if d < 0 {
		t.Fatalf("attempt 0: got negative duration %v", d)
	}
	if d != cfg.BaseDelay {
		t.Errorf("attempt 0 clamped to 1: expected %v, got %v", cfg.BaseDelay, d)
	}
}

func TestDoSuccess(t *testing.T) {
	calls := 0
	err := Do(func() error {
		calls++
		return nil
	}, DefaultConfig())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}
