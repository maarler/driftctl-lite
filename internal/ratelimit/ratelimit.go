// Package ratelimit provides a simple token-bucket rate limiter for controlling
// how frequently drift checks can be triggered.
package ratelimit

import (
	"errors"
	"sync"
	"time"
)

// ErrRateLimited is returned when the rate limit is exceeded.
var ErrRateLimited = errors.New("rate limit exceeded: too many requests")

// Limiter controls the rate of operations using a token bucket.
type Limiter struct {
	mu       sync.Mutex
	tokens   int
	cap      int
	refillAt time.Duration
	last     time.Time
	now      func() time.Time
}

// New creates a Limiter with the given capacity and refill interval.
// capacity is the max burst size; refillInterval is how often one token is added.
func New(capacity int, refillInterval time.Duration) *Limiter {
	return &Limiter{
		tokens:   capacity,
		cap:      capacity,
		refillAt: refillInterval,
		last:     time.Now(),
		now:      time.Now,
	}
}

// newWithClock is used in tests to inject a custom clock.
func newWithClock(capacity int, refillInterval time.Duration, clock func() time.Time) *Limiter {
	l := New(capacity, refillInterval)
	l.now = clock
	l.last = clock()
	return l
}

// Allow returns nil if the operation is permitted, or ErrRateLimited if not.
func (l *Limiter) Allow() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	elapsed := now.Sub(l.last)
	if l.refillAt > 0 {
		added := int(elapsed / l.refillAt)
		if added > 0 {
			l.tokens += added
			if l.tokens > l.cap {
				l.tokens = l.cap
			}
			l.last = l.last.Add(time.Duration(added) * l.refillAt)
		}
	}

	if l.tokens <= 0 {
		return ErrRateLimited
	}
	l.tokens--
	return nil
}

// Remaining returns the current number of available tokens.
func (l *Limiter) Remaining() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.tokens
}
