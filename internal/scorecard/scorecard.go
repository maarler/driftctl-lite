// Package scorecard computes a drift health score from detected results.
package scorecard

import (
	"fmt"
	"io"
	"os"

	"github.com/driftctl-lite/internal/drift"
)

// Score holds the computed health score and breakdown.
type Score struct {
	Total   int
	InSync  int
	Drifted int
	Percent float64
}

// Compute calculates a Score from a slice of drift results.
func Compute(results []drift.Result) Score {
	total := len(results)
	if total == 0 {
		return Score{Total: 0, InSync: 0, Drifted: 0, Percent: 100.0}
	}

	drifted := 0
	for _, r := range results {
		if r.Status != drift.StatusInSync {
			drifted++
		}
	}

	inSync := total - drifted
	pct := float64(inSync) / float64(total) * 100.0

	return Score{
		Total:   total,
		InSync:  inSync,
		Drifted: drifted,
		Percent: pct,
	}
}

// Print writes a human-readable scorecard to the given writer.
func Print(s Score, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, "Drift Scorecard\n")
	fmt.Fprintf(w, "  Total resources : %d\n", s.Total)
	fmt.Fprintf(w, "  In sync         : %d\n", s.InSync)
	fmt.Fprintf(w, "  Drifted         : %d\n", s.Drifted)
	fmt.Fprintf(w, "  Health score    : %.1f%%\n", s.Percent)
}
