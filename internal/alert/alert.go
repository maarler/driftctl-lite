// Package alert provides threshold-based alerting for drift results.
// It evaluates drift results against configurable thresholds and
// emits alerts when limits are exceeded.
package alert

import (
	"fmt"
	"io"
	"os"

	"driftctl-lite/internal/drift"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelOK      Level = "OK"
	LevelWarning Level = "WARNING"
	LevelCritical Level = "CRITICAL"
)

// Thresholds defines the drift count limits that trigger alerts.
type Thresholds struct {
	Warning  int // number of drifted resources to trigger WARNING
	Critical int // number of drifted resources to trigger CRITICAL
}

// Alert holds the result of a threshold evaluation.
type Alert struct {
	Level   Level
	Message string
	Drifted int
}

// DefaultThresholds returns sensible default thresholds.
func DefaultThresholds() Thresholds {
	return Thresholds{
		Warning:  1,
		Critical: 5,
	}
}

// Evaluate checks results against thresholds and returns an Alert.
func Evaluate(results []drift.Result, t Thresholds) Alert {
	drifted := 0
	for _, r := range results {
		if r.Status != drift.StatusInSync {
			drifted++
		}
	}

	switch {
	case t.Critical > 0 && drifted >= t.Critical:
		return Alert{
			Level:   LevelCritical,
			Message: fmt.Sprintf("%d drifted resource(s) exceed critical threshold (%d)", drifted, t.Critical),
			Drifted: drifted,
		}
	case t.Warning > 0 && drifted >= t.Warning:
		return Alert{
			Level:   LevelWarning,
			Message: fmt.Sprintf("%d drifted resource(s) exceed warning threshold (%d)", drifted, t.Warning),
			Drifted: drifted,
		}
	default:
		return Alert{
			Level:   LevelOK,
			Message: "all resources are in sync",
			Drifted: drifted,
		}
	}
}

// Print writes a human-readable alert summary to stdout.
func Print(a Alert) {
	Fprint(os.Stdout, a)
}

// Fprint writes a human-readable alert summary to w.
func Fprint(w io.Writer, a Alert) {
	fmt.Fprintf(w, "[%s] %s\n", a.Level, a.Message)
}
