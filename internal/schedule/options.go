package schedule

import "time"

// Options configures a Scheduler via functional options.
type Options struct {
	Interval time.Duration
}

// Option is a functional option for Options.
type Option func(*Options)

// WithInterval sets the tick interval.
func WithInterval(d time.Duration) Option {
	return func(o *Options) {
		o.Interval = d
	}
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Interval: 5 * time.Minute,
	}
}

// NewFromOptions builds a Job-based Scheduler using functional options.
func NewFromOptions(handler func() error, opts ...Option) *Scheduler {
	o := DefaultOptions()
	for _, opt := range opts {
		opt(&o)
	}
	return New(Job{
		Interval: o.Interval,
		Handler: func(_ interface{ Done() <-chan struct{} }) error {
			return handler()
		},
	})
}
