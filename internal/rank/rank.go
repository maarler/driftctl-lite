// Package rank orders drift results by severity and impact.
package rank

import (
	"sort"

	"driftctl-lite/internal/drift"
)

// Priority assigns a numeric priority to a drift status (lower = higher priority).
func Priority(status drift.Status) int {
	switch status {
	case drift.StatusMissing:
		return 0
	case drift.StatusModified:
		return 1
	case drift.StatusExtra:
		return 2
	case drift.StatusInSync:
		return 3
	default:
		return 4
	}
}

// ByPriority sorts results so highest-priority drift appears first.
// Within the same priority, results are ordered by resource ID.
func ByPriority(results []drift.Result) []drift.Result {
	out := make([]drift.Result, len(results))
	copy(out, results)
	sort.SliceStable(out, func(i, j int) bool {
		pi := Priority(out[i].Status)
		pj := Priority(out[j].Status)
		if pi != pj {
			return pi < pj
		}
		return out[i].ResourceID < out[j].ResourceID
	})
	return out
}

// TopN returns at most n highest-priority results.
func TopN(results []drift.Result, n int) []drift.Result {
	sorted := ByPriority(results)
	if n > len(sorted) {
		n = len(sorted)
	}
	return sorted[:n]
}
