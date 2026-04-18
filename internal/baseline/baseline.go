// Package baseline provides functionality to save and compare drift baselines.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/driftctl-lite/internal/drift"
)

// Baseline represents a saved drift result at a point in time.
type Baseline struct {
	CreatedAt time.Time           `json:"created_at"`
	Results   []drift.Result      `json:"results"`
}

// Save writes a baseline to the given file path.
func Save(path string, results []drift.Result) error {
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Results:   results,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("baseline: write: %w", err)
	}
	return nil
}

// Load reads a baseline from the given file path.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline: file not found: %s", path)
		}
		return nil, fmt.Errorf("baseline: read: %w", err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// Compare returns results that are new or changed compared to the baseline.
func Compare(base *Baseline, current []drift.Result) []drift.Result {
	index := make(map[string]drift.Result, len(base.Results))
	for _, r := range base.Results {
		index[r.ResourceID] = r
	}
	var delta []drift.Result
	for _, r := range current {
		prev, found := index[r.ResourceID]
		if !found || prev.Status != r.Status {
			delta = append(delta, r)
		}
	}
	return delta
}
