package summary

import (
	"fmt"
	"io"

	"github.com/owner/driftctl-lite/internal/drift"
)

// Stats holds aggregated drift statistics.
type Stats struct {
	Total    int
	Missing  int
	Extra    int
	Modified int
	InSync   int
}

// Compute calculates stats from a drift result set.
func Compute(results []drift.Result) Stats {
	s := Stats{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case drift.StatusMissing:
			s.Missing++
		case drift.StatusExtra:
			s.Extra++
		case drift.StatusModified:
			s.Modified++
		case drift.StatusInSync:
			s.InSync++
		}
	}
	return s
}

// Print writes a human-readable summary to w.
func Print(w io.Writer, s Stats) {
	fmt.Fprintf(w, "Summary: %d resource(s) checked\n", s.Total)
	fmt.Fprintf(w, "  In sync:  %d\n", s.InSync)
	fmt.Fprintf(w, "  Missing:  %d\n", s.Missing)
	fmt.Fprintf(w, "  Extra:    %d\n", s.Extra)
	fmt.Fprintf(w, "  Modified: %d\n", s.Modified)
	if s.Missing+s.Extra+s.Modified > 0 {
		fmt.Fprintln(w, "Drift detected!")
	} else {
		fmt.Fprintln(w, "No drift detected.")
	}
}
