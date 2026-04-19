// Package schedule provides periodic drift-check scheduling.
package schedule

import (
	"context"
	"time"
)

// Job holds configuration for a scheduled drift check.
type Job struct {
	Interval time.Duration
	Handler  func(ctx context.Context) error
}

// Scheduler runs a Job on a fixed interval until the context is cancelled.
type Scheduler struct {
	job Job
}

// New creates a new Scheduler.
func New(job Job) *Scheduler {
	return &Scheduler{job: job}
}

// Run starts the scheduler loop. It blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) error {
	if s.job.Interval <= 0 {
		return s.job.Handler(ctx)
	}
	ticker := time.NewTicker(s.job.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.job.Handler(ctx); err != nil {
				return err
			}
		}
	}
}
