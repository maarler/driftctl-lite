package filter

import "github.com/driftctl-lite/internal/drift"

// Options holds filtering criteria for drift results.
type Options struct {
	Types      []string // resource types to include; empty means all
	OnlyDrift  bool     // if true, exclude resources with no drift
}

// Apply returns a new Result containing only the entries that match opts.
func Apply(result drift.Result, opts Options) drift.Result {
	typeSet := make(map[string]bool, len(opts.Types))
	for _, t := range opts.Types {
		typeSet[t] = true
	}

	filtered := drift.Result{
		Missing:  []drift.ResourceDiff{},
		Extra:    []drift.ResourceDiff{},
		Modified: []drift.ResourceDiff{},
	}

	appendIf := func(diffs []drift.ResourceDiff, hasDrift bool) []drift.ResourceDiff {
		out := []drift.ResourceDiff{}
		for _, d := range diffs {
			if len(typeSet) > 0 && !typeSet[d.Type] {
				continue
			}
			if opts.OnlyDrift && !hasDrift {
				continue
			}
			out = append(out, d)
		}
		return out
	}

	filtered.Missing = appendIf(result.Missing, true)
	filtered.Extra = appendIf(result.Extra, true)
	filtered.Modified = appendIf(result.Modified, true)

	return filtered
}
