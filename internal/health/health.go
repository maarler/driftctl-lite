// Package health provides a simple health check for driftctl-lite,
// summarising the overall drift status as a machine-readable signal.
package health

import (
	"fmt"
	"io"
	"os"

	"driftctl-lite/internal/drift"
)

// Status represents the health of the infrastructure.
type Status string

const (
	StatusHealthy   Status = "HEALTHY"
	StatusDegraded  Status = "DEGRADED"
)

// Report holds the outcome of a health evaluation.
type Report struct {
	Status  Status
	Total   int
	Drifted int
}

// Evaluate inspects drift results and returns a Report.
func Evaluate(results []drift.Result) Report {
	total := len(results)
	drifted := 0
	for _, r := range results {
		if r.Status != drift.StatusInSync {
			drifted++
		}
	}
	status := StatusHealthy
	if drifted > 0 {
		status = StatusDegraded
	}
	return Report{
		Status:  status,
		Total:   total,
		Drifted: drifted,
	}
}

// Print writes the health report to w (defaults to os.Stdout).
func Print(r Report, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, "Health: %s | total=%d drifted=%d\n", r.Status, r.Total, r.Drifted)
}
