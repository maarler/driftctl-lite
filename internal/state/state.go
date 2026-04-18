package state

import (
	"encoding/json"
	"fmt"
	"os"
)

// Resource represents a single declared infrastructure resource.
type Resource struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes map[string]string `json:"attributes"`
}

// State holds the full declared state loaded from a state file.
type State struct {
	Resources []Resource `json:"resources"`
}

// LoadFromFile reads and parses a JSON state file into a State struct.
func LoadFromFile(path string) (*State, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening state file: %w", err)
	}
	defer f.Close()

	var s State
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("decoding state file: %w", err)
	}
	return &s, nil
}

// ResourceMap returns a map of resource ID to Resource for quick lookup.
func (s *State) ResourceMap() map[string]Resource {
	m := make(map[string]Resource, len(s.Resources))
	for _, r := range s.Resources {
		m[r.ID] = r
	}
	return m
}
