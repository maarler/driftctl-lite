// Package dedupe provides deduplication of drift results by resource ID.
package dedupe

import "github.com/example/driftctl-lite/internal/drift"

// Strategy controls how duplicates are resolved.
type Strategy string

const (
	// KeepFirst retains the first occurrence of a duplicate.
	KeepFirst Strategy = "first"
	// KeepLast retains the last occurrence of a duplicate.
	KeepLast Strategy = "last"
	// KeepDrift prefers the entry that has drift over one that is in sync.
	KeepDrift Strategy = "drift"
)

// Apply removes duplicate drift.Result entries sharing the same resource ID.
// The strategy determines which duplicate is kept.
func Apply(results []drift.Result, strategy Strategy) []drift.Result {
	if len(results) == 0 {
		return results
	}

	seen := make(map[string]int) // id -> index in out
	out := make([]drift.Result, 0, len(results))

	for _, r := range results {
		idx, exists := seen[r.ResourceID]
		if !exists {
			seen[r.ResourceID] = len(out)
			out = append(out, r)
			continue
		}

		switch strategy {
		case KeepLast:
			out[idx] = r
		case KeepDrift:
			if r.Status != drift.StatusInSync {
				out[idx] = r
			}
		// KeepFirst: do nothing
		}
	}

	return out
}
