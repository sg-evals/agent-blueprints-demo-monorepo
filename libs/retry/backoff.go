package retry

import (
	"math"
	"math/rand"
	"time"
)

// RetryConfig holds retry configuration.
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Jitter      bool
}

// DefaultConfig returns a default retry configuration.
func DefaultConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    10 * time.Second,
		Jitter:      true,
	}
}

// RetryBackoff calculates the backoff duration for a given attempt.
// Attempt numbering starts at 1. The returned duration is always non-negative.
func RetryBackoff(attempt int, cfg RetryConfig) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	delay := time.Duration(math.Pow(2, float64(attempt-1))) * cfg.BaseDelay
	if delay > cfg.MaxDelay {
		delay = cfg.MaxDelay
	}
	if cfg.Jitter {
		delay = time.Duration(rand.Int63n(int64(delay)))
	}
	return delay
}

// Do executes fn with retries according to cfg.
func Do(fn func() error, cfg RetryConfig) error {
	var lastErr error
	for i := 1; i <= cfg.MaxAttempts; i++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if i < cfg.MaxAttempts {
			time.Sleep(RetryBackoff(i, cfg))
		}
	}
	return lastErr
}
