package retry

import (
	"errors"
	"time"
)

// Options configures retry behaviour.
type Options struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64 // backoff multiplier; 1.0 = constant delay
}

// Default returns sensible retry defaults.
func Default() Options {
	return Options{
		MaxAttempts: 3,
		Delay:       200 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// ErrMaxAttemptsReached is returned when all attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("retry: max attempts reached")

// Do calls fn up to opts.MaxAttempts times, backing off between attempts.
// It returns nil on the first success, or ErrMaxAttemptsReached wrapping the
// last error once all attempts are exhausted.
func Do(opts Options, fn func() error) error {
	if opts.MaxAttempts <= 0 {
		return ErrMaxAttemptsReached
	}
	delay := opts.Delay
	var lastErr error
	for i := 0; i < opts.MaxAttempts; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if i < opts.MaxAttempts-1 {
			time.Sleep(delay)
			if opts.Multiplier > 0 {
				delay = time.Duration(float64(delay) * opts.Multiplier)
			}
		}
	}
	return errors.Join(ErrMaxAttemptsReached, lastErr)
}
