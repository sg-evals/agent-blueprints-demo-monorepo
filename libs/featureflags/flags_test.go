package featureflags

import "testing"

func TestStoreSetAndEnabled(t *testing.T) {
	s := NewStore()
	s.Set("dark-mode", true)
	if !s.Enabled("dark-mode") {
		t.Error("expected dark-mode enabled")
	}
	if s.Enabled("nonexistent") {
		t.Error("expected nonexistent flag disabled")
	}
}

func TestStoreAll(t *testing.T) {
	s := NewStore()
	s.Set("a", true)
	s.Set("b", false)
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 flags, got %d", len(all))
	}
}
