// Package group provides grouping of drift results by resource type or status.
package group

import "github.com/driftctl-lite/internal/drift"

// ByType groups drift results by their resource type.
func ByType(results []drift.Result) map[string][]drift.Result {
	groups := make(map[string][]drift.Result)
	for _, r := range results {
		groups[r.ResourceType] = append(groups[r.ResourceType], r)
	}
	return groups
}

// ByStatus groups drift results by their drift status string.
func ByStatus(results []drift.Result) map[string][]drift.Result {
	groups := make(map[string][]drift.Result)
	for _, r := range results {
		key := statusKey(r)
		groups[key] = append(groups[key], r)
	}
	return groups
}

func statusKey(r drift.Result) string {
	if r.Missing {
		return "missing"
	}
	if r.Extra {
		return "extra"
	}
	if r.Modified {
		return "modified"
	}
	return "in_sync"
}

// Summary holds counts per group key.
type Summary struct {
	Key   string
	Count int
}

// Summarize returns a slice of Summary sorted by key for a given grouping.
func Summarize(groups map[string][]drift.Result) []Summary {
	summaries := make([]Summary, 0, len(groups))
	for k, v := range groups {
		summaries = append(summaries, Summary{Key: k, Count: len(v)})
	}
	return summaries
}
