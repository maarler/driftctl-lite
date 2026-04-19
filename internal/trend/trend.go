// Package trend tracks drift counts over time and reports directional changes.
package trend

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single drift measurement at a point in time.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Total     int       `json:"total"`
	Missing   int       `json:"missing"`
	Extra     int       `json:"extra"`
	Modified  int       `json:"modified"`
}

// Trend summarises direction of change between two entries.
type Trend struct {
	Delta     int
	Direction string // "improving", "worsening", "stable"
}

// Append loads existing entries from path, appends e, and saves back.
func Append(path string, e Entry) error {
	entries, _ := Load(path)
	entries = append(entries, e)
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("trend marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads all entries from path. Returns empty slice if file missing.
func Load(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("trend read: %w", err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("trend parse: %w", err)
	}
	return entries, nil
}

// Analyze compares the last two entries and returns a Trend.
// Returns stable trend if fewer than two entries exist.
func Analyze(entries []Entry) Trend {
	if len(entries) < 2 {
		return Trend{Delta: 0, Direction: "stable"}
	}
	prev := entries[len(entries)-2]
	curr := entries[len(entries)-1]
	delta := curr.Total - prev.Total
	dir := "stable"
	if delta < 0 {
		dir = "improving"
	} else if delta > 0 {
		dir = "worsening"
	}
	return Trend{Delta: delta, Direction: dir}
}

// Print writes a human-readable trend summary to stdout.
func Print(t Trend) {
	fmt.Printf("Trend: %s (delta: %+d)\n", t.Direction, t.Delta)
}
