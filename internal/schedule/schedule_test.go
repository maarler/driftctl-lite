package schedule_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"driftctl-lite/internal/schedule"
)

func TestRun_ZeroInterval_CallsOnce(t *testing.T) {
	var called int32
	job := schedule.Job{
		Interval: 0,
		Handler: func(ctx context.Context) error {
			atomic.AddInt32(&called, 1)
			return nil
		},
	}
	s := schedule.New(job)
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&called) != 1 {
		t.Fatalf("expected 1 call, got %d", called)
	}
}

func TestRun_CancelStopsLoop(t *testing.T) {
	var count int32
	job := schedule.Job{
		Interval: 20 * time.Millisecond,
		Handler: func(ctx context.Context) error {
			atomic.AddInt32(&count, 1)
			return nil
		},
	}
	s := schedule.New(job)
	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	err := s.Run(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	if atomic.LoadInt32(&count) < 2 {
		t.Fatalf("expected at least 2 ticks, got %d", count)
	}
}

func TestRun_HandlerError_Stops(t *testing.T) {
	handlerErr := errors.New("boom")
	job := schedule.Job{
		Interval: 10 * time.Millisecond,
		Handler: func(ctx context.Context) error {
			return handlerErr
		},
	}
	s := schedule.New(job)
	ctx := context.Background()
	err := s.Run(ctx)
	if !errors.Is(err, handlerErr) {
		t.Fatalf("expected handlerErr, got %v", err)
	}
}
