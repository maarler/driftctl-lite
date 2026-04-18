// Package metrics collects and reports drift scan statistics.
package metrics

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/snyk/driftctl-lite/internal/drift"
)

// Stats holds counters for a single drift scan run.
type Stats struct {
	Total    int
	InSync   int
	Missing  int
	Extra    int
	Modified int
	ScannedAt time.Time
	Duration  time.Duration
}

// Collect builds a Stats value from a slice of drift results and elapsed duration.
func Collect(results []drift.Result, elapsed time.Duration) Stats {
	s := Stats{
		Total:     len(results),
		ScannedAt: time.Now().UTC(),
		Duration:  elapsed,
	}
	for _, r := range results {
		switch r.Status {
		case drift.StatusInSync:
			s.InSync++
		case drift.StatusMissing:
			s.Missing++
		case drift.StatusExtra:
			s.Extra++
		case drift.StatusModified:
			s.Modified++
		}
	}
	return s
}

// Print writes a human-readable metrics summary to w (defaults to os.Stdout).
func Print(s Stats, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, "Scan completed at %s (took %s)\n", s.ScannedAt.Format(time.RFC3339), s.Duration.Round(time.Millisecond))
	fmt.Fprintf(w, "  Total:    %d\n", s.Total)
	fmt.Fprintf(w, "  In sync:  %d\n", s.InSync)
	fmt.Fprintf(w, "  Missing:  %d\n", s.Missing)
	fmt.Fprintf(w, "  Extra:    %d\n", s.Extra)
	fmt.Fprintf(w, "  Modified: %d\n", s.Modified)
}
