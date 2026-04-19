package ratelimit

import (
	"testing"
	"time"
)

func TestAllow_WithinCapacity(t *testing.T) {
	l := New(3, time.Second)
	for i := 0; i < 3; i++ {
		if err := l.Allow(); err != nil {
			t.Fatalf("expected nil on call %d, got %v", i+1, err)
		}
	}
}

func TestAllow_ExceedsCapacity(t *testing.T) {
	l := New(2, time.Second)
	_ = l.Allow()
	_ = l.Allow()
	if err := l.Allow(); err != ErrRateLimited {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	l := newWithClock(1, 500*time.Millisecond, clock)

	// Consume the only token.
	if err := l.Allow(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should be limited now.
	if err := l.Allow(); err != ErrRateLimited {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}

	// Advance clock by 600ms — should refill one token.
	now = now.Add(600 * time.Millisecond)
	if err := l.Allow(); err != nil {
		t.Fatalf("expected token after refill, got %v", err)
	}
}

func TestRemaining_DecreasesOnAllow(t *testing.T) {
	l := New(5, time.Second)
	if r := l.Remaining(); r != 5 {
		t.Fatalf("expected 5 remaining, got %d", r)
	}
	_ = l.Allow()
	if r := l.Remaining(); r != 4 {
		t.Fatalf("expected 4 remaining, got %d", r)
	}
}

func TestAllow_ZeroCapacity_AlwaysLimited(t *testing.T) {
	l := New(0, time.Second)
	if err := l.Allow(); err != ErrRateLimited {
		t.Fatalf("expected ErrRateLimited for zero capacity, got %v", err)
	}
}
