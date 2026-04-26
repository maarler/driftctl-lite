// Package clamp enforces min/max bounds on numeric drift-result fields
// such as severity scores, diff counts, and priority values.
package clamp

import "github.com/driftctl-lite/internal/drift"

// Options controls the clamping behaviour.
type Options struct {
	// MinPriority is the lowest allowed priority value (inclusive).
	MinPriority int
	// MaxPriority is the highest allowed priority value (inclusive).
	MaxPriority int
	// MaxDiffKeys caps the number of diff keys retained per result.
	// 0 means unlimited.
	MaxDiffKeys int
}

// DefaultOptions returns sensible production defaults.
func DefaultOptions() Options {
	return Options{
		MinPriority: 0,
		MaxPriority: 100,
		MaxDiffKeys: 50,
	}
}

// Apply clamps each result in the slice according to opts and returns the
// updated slice. The input slice is mutated in-place for efficiency.
func Apply(results []drift.Result, opts Options) []drift.Result {
	for i := range results {
		results[i] = clampOne(results[i], opts)
	}
	return results
}

func clampOne(r drift.Result, opts Options) drift.Result {
	// Clamp Priority if the field exists on the result metadata.
	if p, ok := r.Meta["priority"]; ok {
		if v, ok := p.(int); ok {
			if v < opts.MinPriority {
				v = opts.MinPriority
			}
			if v > opts.MaxPriority {
				v = opts.MaxPriority
			}
			if r.Meta == nil {
				r.Meta = map[string]interface{}{}
			}
			r.Meta["priority"] = v
		}
	}

	// Cap diff keys when MaxDiffKeys > 0.
	if opts.MaxDiffKeys > 0 && len(r.Diffs) > opts.MaxDiffKeys {
		truncated := make(map[string]drift.DiffValue, opts.MaxDiffKeys)
		count := 0
		for k, v := range r.Diffs {
			if count >= opts.MaxDiffKeys {
				break
			}
			truncated[k] = v
			count++
		}
		r.Diffs = truncated
	}

	return r
}
