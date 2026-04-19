package schedule_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"driftctl-lite/internal/schedule"
)

func TestScheduler_Integration_MultiTick(t *testing.T) {
	var ticks int32
	job := schedule.Job{
		Interval: 15 * time.Millisecond,
		Handler: func(ctx context.Context) error {
			atomic.AddInt32(&ticks, 1)
			return nil
		},
	}
	s := schedule.New(job)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	s.Run(ctx) //nolint:errcheck
	got := atomic.LoadInt32(&ticks)
	if got < 2 {
		t.Fatalf("expected >=2 ticks in 60ms with 15ms interval, got %d", got)
	}
}
