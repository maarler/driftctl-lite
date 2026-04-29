// Package quota enforces per-resource-type drift count limits.
// Results exceeding the quota threshold are flagged or dropped.
package quota

import "github.com/driftctl-lite/internal/drift"

// Options configures quota enforcement behaviour.
type Options struct {
	// MaxPerType is the maximum number of drifted resources allowed per type.
	// A value of 0 means no limit.
	MaxPerType int
	// DropExceeding removes results that exceed the quota instead of flagging them.
	DropExceeding bool
}

// DefaultOptions returns sensible quota defaults (no limit).
func DefaultOptions() Options {
	return Options{MaxPerType: 0, DropExceeding: false}
}

// Report holds the outcome of a quota evaluation.
type Report struct {
	Results   []drift.Result
	Violations []Violation
}

// Violation describes a single quota breach.
type Violation struct {
	Type    string
	Allowed int
	Actual  int
}

// Apply enforces quota rules against the provided results.
// Results that are in-sync are always kept unchanged.
func Apply(results []drift.Result, opts Options) Report {
	if opts.MaxPerType <= 0 {
		return Report{Results: results}
	}

	counts := make(map[string]int)
	var kept []drift.Result
	var violations []Violation
	violated := make(map[string]bool)

	for _, r := range results {
		if !r.HasDrift() {
			kept = append(kept, r)
			continue
		}
		counts[r.Resource.Type]++
		if counts[r.Resource.Type] > opts.MaxPerType {
			if !violated[r.Resource.Type] {
				violated[r.Resource.Type] = true
			}
			if opts.DropExceeding {
				continue
			}
		}
		kept = append(kept, r)
	}

	for typ := range violated {
		violations = append(violations, Violation{
			Type:    typ,
			Allowed: opts.MaxPerType,
			Actual:  counts[typ],
		})
	}

	return Report{Results: kept, Violations: violations}
}
