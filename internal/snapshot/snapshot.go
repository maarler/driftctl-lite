// Package snapshot provides functionality to capture and persist
// the current drift detection results to disk for later comparison.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/driftctl-lite/internal/drift"
)

// Snapshot holds a timestamped set of drift results.
type Snapshot struct {
	CapturedAt time.Time     `json:"captured_at"`
	Results    []drift.Result `json:"results"`
}

// Save writes the snapshot to the given file path as JSON.
func Save(path string, results []drift.Result) error {
	snap := Snapshot{
		CapturedAt: time.Now().UTC(),
		Results:    results,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot: file not found: %s", path)
		}
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}
