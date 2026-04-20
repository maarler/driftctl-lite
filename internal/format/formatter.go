// Package format provides utilities for normalising and pretty-printing
// drift result fields before they are handed to an output renderer.
package format

import (
	"fmt"
	"strings"

	"github.com/driftctl-lite/internal/drift"
)

// Options controls how results are formatted.
type Options struct {
	// MaxValueLen truncates long field values to this many runes (0 = no limit).
	MaxValueLen int
	// UpperCaseStatus converts the Status field to upper-case when true.
	UpperCaseStatus bool
}

// DefaultOptions returns sensible formatting defaults.
func DefaultOptions() Options {
	return Options{
		MaxValueLen:     80,
		UpperCaseStatus: true,
	}
}

// Apply formats every result in place according to opts and returns the
// (potentially modified) slice.
func Apply(results []drift.Result, opts Options) []drift.Result {
	out := make([]drift.Result, len(results))
	for i, r := range results {
		out[i] = formatOne(r, opts)
	}
	return out
}

func formatOne(r drift.Result, opts Options) drift.Result {
	if opts.UpperCaseStatus {
		r.Status = strings.ToUpper(r.Status)
	}

	if opts.MaxValueLen > 0 && len(r.Diffs) > 0 {
		formatted := make(map[string]drift.DiffDetail, len(r.Diffs))
		for k, d := range r.Diffs {
			d.Expected = truncate(fmt.Sprintf("%v", d.Expected), opts.MaxValueLen)
			d.Got = truncate(fmt.Sprintf("%v", d.Got), opts.MaxValueLen)
			formatted[k] = d
		}
		r.Diffs = formatted
	}
	return r
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "…"
}
