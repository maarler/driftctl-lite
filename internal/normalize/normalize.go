// Package normalize provides utilities for standardizing drift result fields
// before comparison or display — trimming whitespace, lowercasing types, etc.
package normalize

import (
	"strings"

	"driftctl-lite/internal/drift"
)

// Options controls which normalization steps are applied.
type Options struct {
	// LowercaseType converts resource type to lowercase.
	LowercaseType bool
	// TrimIDs strips leading/trailing whitespace from resource IDs.
	TrimIDs bool
	// TrimDiffValues strips whitespace from diff Expected/Actual values.
	TrimDiffValues bool
}

// DefaultOptions returns a sensible default normalization configuration.
func DefaultOptions() Options {
	return Options{
		LowercaseType:  true,
		TrimIDs:        true,
		TrimDiffValues: true,
	}
}

// Apply normalizes a slice of drift results according to the given options.
// It returns a new slice; original results are not mutated.
func Apply(results []drift.Result, opts Options) []drift.Result {
	out := make([]drift.Result, 0, len(results))
	for _, r := range results {
		out = append(out, normalizeOne(r, opts))
	}
	return out
}

func normalizeOne(r drift.Result, opts Options) drift.Result {
	if opts.TrimIDs {
		r.ResourceID = strings.TrimSpace(r.ResourceID)
	}
	if opts.LowercaseType {
		r.ResourceType = strings.ToLower(strings.TrimSpace(r.ResourceType))
	}
	if opts.TrimDiffValues && r.Diffs != nil {
		normalized := make(map[string]drift.Diff, len(r.Diffs))
		for k, d := range r.Diffs {
			d.Expected = strings.TrimSpace(d.Expected)
			d.Actual = strings.TrimSpace(d.Actual)
			normalized[k] = d
		}
		r.Diffs = normalized
	}
	return r
}
