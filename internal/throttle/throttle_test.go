package throttle_test

import (
	"testing"
	"time"

	"driftctl-lite/internal/throttle"
)

func TestAllow_ZeroInterval_AlwaysAllows(t *testing.T) {
	th := throttle.New(0)
	for i := 0; i < 5; i++ {
		if err := th.Allow(); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	}
}

func TestAllow_FirstCall_Succeeds(t *testing.T) {
	th := throttle.New(10 * time.Second)
	if err := th.Allow(); err != nil {
		t.Fatalf("first call should succeed, got %v", err)
	}
}

func TestAllow_SecondCall_TooSoon_Throttled(t *testing.T) {
	th := throttle.New(10 * time.Second)
	_ = th.Allow()
	if err := th.Allow(); err != throttle.ErrThrottled {
		t.Fatalf("expected ErrThrottled, got %v", err)
	}
}

func TestAllow_AfterInterval_Succeeds(t *testing.T) {
	th := throttle.New(20 * time.Millisecond)
	_ = th.Allow()
	time.Sleep(30 * time.Millisecond)
	if err := th.Allow(); err != nil {
		t.Fatalf("expected nil after interval elapsed, got %v", err)
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	th := throttle.New(10 * time.Second)
	_ = th.Allow()
	th.Reset()
	if err := th.Allow(); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}

func TestLastRun_ZeroBeforeFirstAllow(t *testing.T) {
	th := throttle.New(5 * time.Second)
	if !th.LastRun().IsZero() {
		t.Fatal("expected zero time before any Allow call")
	}
}

func TestLastRun_UpdatedAfterAllow(t *testing.T) {
	th := throttle.New(5 * time.Second)
	before := time.Now()
	_ = th.Allow()
	if th.LastRun().Before(before) {
		t.Fatal("LastRun should be >= time before Allow")
	}
}
