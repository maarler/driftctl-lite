// Package watch provides periodic drift detection by re-running detection
// on a configured interval and emitting results to a callback.
package watch

import (
	"context"
	"time"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/state"
	"driftctl-lite/internal/live"
)

// ResultHandler is called each time a drift check completes.
type ResultHandler func(results []drift.Result, err error)

// Watcher runs drift detection on a fixed interval.
type Watcher struct {
	interval time.Duration
	statePath string
	fetcher   *live.Fetcher
	handler   ResultHandler
}

// New creates a Watcher that checks for drift every interval.
func New(interval time.Duration, statePath string, fetcher *live.Fetcher, handler ResultHandler) *Watcher {
	return &Watcher{
		interval:  interval,
		statePath: statePath,
		fetcher:   fetcher,
		handler:   handler,
	}
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.check()

	for {
		select {
		case <-ticker.C:
			w.check()
		case <-ctx.Done():
			return
		}
	}
}

func (w *Watcher) check() {
	declared, err := state.LoadFromFile(w.statePath)
	if err != nil {
		w.handler(nil, err)
		return
	}

	live, err := w.fetcher.Fetch()
	if err != nil {
		w.handler(nil, err)
		return
	}

	results := drift.Detect(declared, live)
	w.handler(results, nil)
}
