// Package checkpoint provides functionality to save and restore
// named checkpoints of drift results for incremental analysis.
package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"driftctl-lite/internal/drift"
)

// Entry represents a saved checkpoint with metadata.
type Entry struct {
	Name      string             `json:"name"`
	SavedAt   time.Time          `json:"saved_at"`
	Results   []drift.Result     `json:"results"`
}

// Save writes a named checkpoint to the given directory.
func Save(dir, name string, results []drift.Result) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("checkpoint: mkdir %s: %w", dir, err)
	}
	entry := Entry{
		Name:    name,
		SavedAt: time.Now().UTC(),
		Results: results,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("checkpoint: marshal: %w", err)
	}
	path := filepath.Join(dir, name+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("checkpoint: write %s: %w", path, err)
	}
	return nil
}

// Load reads a named checkpoint from the given directory.
func Load(dir, name string) (*Entry, error) {
	path := filepath.Join(dir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("checkpoint: %q not found", name)
		}
		return nil, fmt.Errorf("checkpoint: read %s: %w", path, err)
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("checkpoint: unmarshal: %w", err)
	}
	return &entry, nil
}

// List returns the names of all checkpoints stored in dir.
func List(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("checkpoint: list %s: %w", dir, err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}
