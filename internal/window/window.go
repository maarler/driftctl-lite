// Package window provides a sliding time-window aggregator for drift results.
// It collects snapshots over a configurable duration and returns the union
// of all unique drift results seen within that window.
package window

import (
	"sync"
	"time"

	"driftctl-lite/internal/drift"
)

// entry holds a batch of results alongside the time it was recorded.
type entry struct {
	recordedAt time.Time
	results    []drift.Result
}

// Window is a thread-safe sliding-window aggregator.
type Window struct {
	mu       sync.Mutex
	ttl      time.Duration
	entries  []entry
	nowFn    func() time.Time
}

// New creates a Window that retains results for the given TTL duration.
func New(ttl time.Duration) *Window {
	return &Window{
		ttl:   ttl,
		nowFn: time.Now,
	}
}

// Add appends a new batch of results to the window.
func (w *Window) Add(results []drift.Result) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries = append(w.entries, entry{
		recordedAt: w.nowFn(),
		results:    results,
	})
}

// Collect returns the deduplicated union of all results still within the TTL.
// Results are keyed by "type/id" to eliminate duplicates; the most recently
// added entry wins on collision.
func (w *Window) Collect() []drift.Result {
	w.mu.Lock()
	defer w.mu.Unlock()

	cutoff := w.nowFn().Add(-w.ttl)
	seen := make(map[string]drift.Result)

	var fresh []entry
	for _, e := range w.entries {
		if e.recordedAt.After(cutoff) {
			fresh = append(fresh, e)
			for _, r := range e.results {
				k := r.ResourceType + "/" + r.ResourceID
				seen[k] = r
			}
		}
	}
	w.entries = fresh

	out := make([]drift.Result, 0, len(seen))
	for _, r := range seen {
		out = append(out, r)
	}
	return out
}

// Len returns the number of batches currently held in the window.
func (w *Window) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.entries)
}
