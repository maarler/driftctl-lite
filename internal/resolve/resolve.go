// Package resolve provides utilities for resolving drift results by
// marking previously drifted resources as acknowledged or resolved,
// and tracking which resources still require attention.
package resolve

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/driftctl-lite/internal/drift"
)

// Status represents the resolution state of a drift result.
type Status string

const (
	StatusOpen       Status = "open"
	StatusAcknowledged Status = "acknowledged"
	StatusResolved   Status = "resolved"
)

// Entry records a resolution decision for a single drift result.
type Entry struct {
	ResourceID   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	Status       Status    `json:"status"`
	Reason       string    `json:"reason,omitempty"`
	ResolvedAt   time.Time `json:"resolved_at"`
}

// Registry holds all resolution entries indexed by "type/id".
type Registry struct {
	Entries map[string]Entry `json:"entries"`
}

// key returns the lookup key for a resource.
func key(resourceType, resourceID string) string {
	return resourceType + "/" + resourceID
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{Entries: make(map[string]Entry)}
}

// Acknowledge marks a drift result as acknowledged with an optional reason.
func (r *Registry) Acknowledge(result drift.Result, reason string) {
	k := key(result.ResourceType, result.ResourceID)
	r.Entries[k] = Entry{
		ResourceID:   result.ResourceID,
		ResourceType: result.ResourceType,
		Status:       StatusAcknowledged,
		Reason:       reason,
		ResolvedAt:   time.Now().UTC(),
	}
}

// Resolve marks a drift result as fully resolved.
func (r *Registry) Resolve(result drift.Result, reason string) {
	k := key(result.ResourceType, result.ResourceID)
	r.Entries[k] = Entry{
		ResourceID:   result.ResourceID,
		ResourceType: result.ResourceType,
		Status:       StatusResolved,
		Reason:       reason,
		ResolvedAt:   time.Now().UTC(),
	}
}

// Filter returns only results whose status is "open" (not acknowledged or resolved).
func (r *Registry) Filter(results []drift.Result) []drift.Result {
	var open []drift.Result
	for _, res := range results {
		k := key(res.ResourceType, res.ResourceID)
		entry, found := r.Entries[k]
		if !found || entry.Status == StatusOpen {
			open = append(open, res)
		}
	}
	return open
}

// Save persists the registry to a JSON file at the given path.
func (r *Registry) Save(path string) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("resolve: marshal registry: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("resolve: write registry: %w", err)
	}
	return nil
}

// Load reads a registry from a JSON file. If the file does not exist,
// an empty registry is returned without error.
func Load(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NewRegistry(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("resolve: read registry: %w", err)
	}
	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("resolve: unmarshal registry: %w", err)
	}
	if reg.Entries == nil {
		reg.Entries = make(map[string]Entry)
	}
	return &reg, nil
}
