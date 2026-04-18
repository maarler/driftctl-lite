// Package audit provides drift audit logging — recording drift detection
// runs and their results to a persistent log file for later review.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/snyk/driftctl-lite/internal/drift"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp  time.Time        `json:"timestamp"`
	StateFile  string           `json:"state_file"`
	Source     string           `json:"source"`
	TotalItems int              `json:"total_items"`
	DriftCount int              `json:"drift_count"`
	Results    []drift.Result   `json:"results"`
}

// Logger writes audit entries to a file.
type Logger struct {
	path string
}

// NewLogger creates a Logger that appends to the given file path.
func NewLogger(path string) *Logger {
	return &Logger{path: path}
}

// Record appends an audit entry derived from the provided results.
func (l *Logger) Record(stateFile, source string, results []drift.Result) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	driftCount := 0
	for _, r := range results {
		if r.Status != drift.StatusInSync {
			driftCount++
		}
	}

	entry := Entry{
		Timestamp:  time.Now().UTC(),
		StateFile:  stateFile,
		Source:     source,
		TotalItems: len(results),
		DriftCount: driftCount,
		Results:    results,
	}

	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}
