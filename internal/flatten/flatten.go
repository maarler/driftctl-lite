// Package flatten provides utilities to flatten nested drift results
// into a single-level key-value representation for reporting and export.
package flatten

import (
	"fmt"
	"sort"

	"driftctl-lite/internal/drift"
)

// Record is a flat representation of a single drift result.
type Record struct {
	ID     string
	Type   string
	Status string
	Key    string
	Wanted string
	Got    string
}

// Flatten converts a slice of drift results into flat records.
// Each diff field becomes its own record row.
func Flatten(results []drift.Result) []Record {
	var records []Record

	for _, r := range results {
		if len(r.Diffs) == 0 {
			records = append(records, Record{
				ID:     r.ResourceID,
				Type:   r.ResourceType,
				Status: string(r.Status),
				Key:    "-",
				Wanted: "-",
				Got:    "-",
			})
			continue
		}

		keys := make([]string, 0, len(r.Diffs))
		for k := range r.Diffs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			d := r.Diffs[k]
			records = append(records, Record{
				ID:     r.ResourceID,
				Type:   r.ResourceType,
				Status: string(r.Status),
				Key:    k,
				Wanted: fmt.Sprintf("%v", d.Wanted),
				Got:    fmt.Sprintf("%v", d.Got),
			})
		}
	}

	return records
}
