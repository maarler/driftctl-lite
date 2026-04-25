// Package profile provides named drift-check profiles that bundle
// common configuration presets (filters, thresholds, output format).
package profile

import (
	"encoding/json"
	"fmt"
	"os"
)

// Profile holds a named set of reusable scan options.
type Profile struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	OutputFormat string           `json:"output_format,omitempty"` // text | json
	OnlyDrift   bool              `json:"only_drift,omitempty"`
	FilterType  string            `json:"filter_type,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	Severity    string            `json:"severity,omitempty"` // info | warning | critical
}

// Registry holds a collection of named profiles.
type Registry struct {
	Profiles map[string]Profile `json:"profiles"`
}

// LoadFromFile reads a JSON profile registry from disk.
func LoadFromFile(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Registry{Profiles: make(map[string]Profile)}, nil
		}
		return nil, fmt.Errorf("profile: read %s: %w", path, err)
	}
	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("profile: parse %s: %w", path, err)
	}
	if reg.Profiles == nil {
		reg.Profiles = make(map[string]Profile)
	}
	return &reg, nil
}

// Get returns the named profile or an error if it does not exist.
func (r *Registry) Get(name string) (Profile, error) {
	p, ok := r.Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("profile: %q not found", name)
	}
	return p, nil
}

// Default returns a sensible built-in profile.
func Default() Profile {
	return Profile{
		Name:         "default",
		Description:  "Show all resources, text output",
		OutputFormat: "text",
		OnlyDrift:    false,
	}
}
