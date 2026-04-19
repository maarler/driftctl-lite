// Package rollup aggregates drift results into a compact summary grouped by resource type.
package rollup

import (
	"fmt"
	"io"
	"os"
	"sort"

	"driftctl-lite/internal/drift"
)

// TypeRollup holds aggregated counts for a single resource type.
type TypeRollup struct {
	Type    string
	Total   int
	Missing int
	Extra   int
	Modified int
	InSync  int
}

// Report is a slice of TypeRollup entries.
type Report []TypeRollup

// Compute aggregates drift results by resource type.
func Compute(results []drift.Result) Report {
	m := map[string]*TypeRollup{}
	for _, r := range results {
		entry, ok := m[r.ResourceType]
		if !ok {
			entry = &TypeRollup{Type: r.ResourceType}
			m[r.ResourceType] = entry
		}
		entry.Total++
		switch r.Status {
		case drift.StatusMissing:
			entry.Missing++
		case drift.StatusExtra:
			entry.Extra++
		case drift.StatusModified:
			entry.Modified++
		default:
			entry.InSync++
		}
	}
	report := make(Report, 0, len(m))
	for _, v := range m {
		report = append(report, *v)
	}
	sort.Slice(report, func(i, j int) bool {
		return report[i].Type < report[j].Type
	})
	return report
}

// Print writes the rollup report to stdout.
func Print(r Report) {
	Fprint(os.Stdout, r)
}

// Fprint writes the rollup report to w.
func Fprint(w io.Writer, r Report) {
	fmt.Fprintf(w, "%-20s %6s %7s %5s %8s %6s\n", "TYPE", "TOTAL", "MISSING", "EXTRA", "MODIFIED", "IN_SYNC")
	for _, e := range r {
		fmt.Fprintf(w, "%-20s %6d %7d %5d %8d %6d\n", e.Type, e.Total, e.Missing, e.Extra, e.Modified, e.InSync)
	}
}
