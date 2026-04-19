// Package compare provides utilities for comparing state snapshots over time.
package compare

import (
	"fmt"
	"strings"

	"github.com/driftctl-lite/internal/drift"
)

// Delta represents the change between two drift result sets.
type Delta struct {
	Resolved  []drift.Result // were drifted, now in sync
	New       []drift.Result // newly drifted
	Persisted []drift.Result // still drifted
}

// Compare computes the delta between a previous and current set of drift results.
func Compare(previous, current []drift.Result) Delta {
	prevMap := indexByID(previous)
	currMap := indexByID(current)

	delta := Delta{}

	for id, prev := range prevMap {
		if prev.Status == drift.StatusInSync {
			continue
		}
		if curr, ok := currMap[id]; ok {
			if curr.Status == drift.StatusInSync {
				delta.Resolved = append(delta.Resolved, curr)
			} else {
				delta.Persisted = append(delta.Persisted, curr)
			}
		} else {
			delta.Resolved = append(delta.Resolved, prev)
		}
	}

	for id, curr := range currMap {
		if curr.Status == drift.StatusInSync {
			continue
		}
		if _, ok := prevMap[id]; !ok {
			delta.New = append(delta.New, curr)
		}
	}

	return delta
}

// Summary returns a human-readable summary of the delta.
func Summary(d Delta) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Resolved: %d, New: %d, Persisted: %d", len(d.Resolved), len(d.New), len(d.Persisted))
	return sb.String()
}

func indexByID(results []drift.Result) map[string]drift.Result {
	m := make(map[string]drift.Result, len(results))
	for _, r := range results {
		m[r.ResourceID] = r
	}
	return m
}
