// Package diff provides utilities for computing field-level differences
// between declared and live resource attributes.
package diff

import (
	"fmt"
	"sort"
)

// FieldDiff represents a single field change.
type FieldDiff struct {
	Field    string
	Declared interface{}
	Live     interface{}
}

// String returns a human-readable representation of the field diff.
func (f FieldDiff) String() string {
	return fmt.Sprintf("  ~ %s: %v => %v", f.Field, f.Declared, f.Live)
}

// Compute returns field-level diffs between two attribute maps.
// declared is the expected state; live is the actual state.
func Compute(declared, live map[string]interface{}) []FieldDiff {
	if declared == nil {
		declared = map[string]interface{}{}
	}
	if live == nil {
		live = map[string]interface{}{}
	}

	keys := unionKeys(declared, live)
	var diffs []FieldDiff

	for _, k := range keys {
		dv, dok := declared[k]
		lv, lok := live[k]
		switch {
		case dok && !lok:
			diffs = append(diffs, FieldDiff{Field: k, Declared: dv, Live: nil})
		case !dok && lok:
			diffs = append(diffs, FieldDiff{Field: k, Declared: nil, Live: lv})
		case fmt.Sprintf("%v", dv) != fmt.Sprintf("%v", lv):
			diffs = append(diffs, FieldDiff{Field: k, Declared: dv, Live: lv})
		}
	}
	return diffs
}

func unionKeys(a, b map[string]interface{}) []string {
	seen := map[string]struct{}{}
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
