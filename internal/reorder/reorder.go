// Package reorder provides sorting utilities for drift results,
// allowing callers to order results by various fields for consistent output.
package reorder

import (
	"sort"

	"driftctl-lite/internal/drift"
)

// Field represents a sortable field on a drift result.
type Field string

const (
	FieldID     Field = "id"
	FieldType   Field = "type"
	FieldStatus Field = "status"
)

// Options controls how results are sorted.
type Options struct {
	By        Field
	Ascending bool
}

// DefaultOptions returns the default sort options (by ID, ascending).
func DefaultOptions() Options {
	return Options{By: FieldID, Ascending: true}
}

// Apply sorts a slice of drift results according to the given options.
// The original slice is sorted in-place and also returned for convenience.
func Apply(results []drift.Result, opts Options) []drift.Result {
	if len(results) == 0 {
		return results
	}

	sort.SliceStable(results, func(i, j int) bool {
		var less bool
		switch opts.By {
		case FieldType:
			less = results[i].Resource.Type < results[j].Resource.Type
		case FieldStatus:
			less = statusOrder(results[i]) < statusOrder(results[j])
		default: // FieldID
			less = results[i].Resource.ID < results[j].Resource.ID
		}
		if opts.Ascending {
			return less
		}
		return !less
	})

	return results
}

// statusOrder assigns a numeric priority to each drift status for sorting.
func statusOrder(r drift.Result) int {
	switch r.Status {
	case drift.StatusMissing:
		return 0
	case drift.StatusModified:
		return 1
	case drift.StatusExtra:
		return 2
	default:
		return 3
	}
}
