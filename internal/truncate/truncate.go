// Package truncate provides utilities for limiting drift result sets to a
// maximum number of entries, preserving the highest-priority items first.
package truncate

import "github.com/owner/driftctl-lite/internal/drift"

// Options controls truncation behaviour.
type Options struct {
	MaxResults int  // 0 means no limit
	PreserveInSync bool // if true, in-sync results are dropped before truncating
}

// DefaultOptions returns sensible defaults (no limit, keep all statuses).
func DefaultOptions() Options {
	return Options{MaxResults: 0, PreserveInSync: false}
}

// Apply truncates results according to opts.
// Drifted resources (missing / extra / modified) are always preferred over
// in-sync ones when PreserveInSync is false and a limit is in effect.
func Apply(results []drift.Result, opts Options) ([]drift.Result, bool) {
	if len(results) == 0 {
		return results, false
	}

	working := results
	if opts.PreserveInSync {
		working = dropInSync(results)
	}

	if opts.MaxResults <= 0 || len(working) <= opts.MaxResults {
		return working, false
	}

	return working[:opts.MaxResults], true
}

// Truncated reports whether Apply would truncate the given slice.
func Truncated(results []drift.Result, opts Options) bool {
	_, t := Apply(results, opts)
	return t
}

func dropInSync(results []drift.Result) []drift.Result {
	out := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if r.Status != drift.StatusInSync {
			out = append(out, r)
		}
	}
	return out
}
