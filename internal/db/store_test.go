package db

import (
	"testing"
)

func TestGormStoreInit(t *testing.T) {
	s, err := NewGormStore(":memory:")
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	if s == nil {
		t.Fatalf("expected store")
	}
}
