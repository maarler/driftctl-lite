package retry

import (
	"errors"
	"testing"
	"time"
)

var errTemp = errors.New("temporary error")

func fastOpts(attempts int) Options {
	return Options{MaxAttempts: attempts, Delay: 0, Multiplier: 1.0}
}

func TestDo_SucceedsFirstAttempt(t *testing.T) {
	calls := 0
	err := Do(fastOpts(3), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesAndSucceeds(t *testing.T) {
	calls := 0
	err := Do(fastOpts(3), func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	calls := 0
	err := Do(fastOpts(3), func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
	if !errors.Is(err, errTemp) {
		t.Fatalf("expected wrapped errTemp, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ZeroAttempts(t *testing.T) {
	err := Do(fastOpts(0), func() error { return nil })
	if !errors.Is(err, ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
}

func TestDefault_Values(t *testing.T) {
	o := Default()
	if o.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts 3, got %d", o.MaxAttempts)
	}
	if o.Delay != 200*time.Millisecond {
		t.Errorf("unexpected delay %v", o.Delay)
	}
	if o.Multiplier != 2.0 {
		t.Errorf("unexpected multiplier %v", o.Multiplier)
	}
}
