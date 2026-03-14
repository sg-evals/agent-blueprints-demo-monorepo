package featureflags

import "sync"

// Store holds feature flag state.
type Store struct {
	mu    sync.RWMutex
	flags map[string]bool
}

// NewStore creates a new feature flag store.
func NewStore() *Store {
	return &Store{
		flags: make(map[string]bool),
	}
}

// Set sets a feature flag value.
func (s *Store) Set(name string, enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flags[name] = enabled
}

// Enabled checks if a feature flag is enabled.
func (s *Store) Enabled(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.flags[name]
}

// All returns all feature flags.
func (s *Store) All() map[string]bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]bool, len(s.flags))
	for k, v := range s.flags {
		result[k] = v
	}
	return result
}
