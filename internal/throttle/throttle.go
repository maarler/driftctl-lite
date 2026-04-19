// Package throttle provides rate-limiting for drift scan runs.
package throttle

import (
	"errors"
	"sync"
	"time"
)

// ErrThrottled is returned when a scan is attempted too soon after the last one.
var ErrThrottled = errors.New("throttle: scan attempted too soon, please wait before retrying")

// Throttler enforces a minimum interval between successive scan runs.
type Throttler struct {
	mu       sync.Mutex
	interval time.Duration
	lastRun  time.Time
}

// New creates a Throttler with the given minimum interval between runs.
func New(interval time.Duration) *Throttler {
	return &Throttler{interval: interval}
}

// Allow returns nil if enough time has elapsed since the last run, or
// ErrThrottled if the caller should wait longer. On success it records
// the current time as the new last-run timestamp.
func (t *Throttler) Allow() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.interval <= 0 {
		return nil
	}

	now := time.Now()
	if !t.lastRun.IsZero() && now.Sub(t.lastRun) < t.interval {
		return ErrThrottled
	}
	t.lastRun = now
	return nil
}

// Reset clears the last-run timestamp, allowing the next call to Allow to
// succeed regardless of when it is called.
func (t *Throttler) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastRun = time.Time{}
}

// LastRun returns the timestamp of the most recent successful Allow call.
func (t *Throttler) LastRun() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastRun
}
