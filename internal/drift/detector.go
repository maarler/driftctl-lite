package drift

import "github.com/example/driftctl-lite/internal/state"

// DriftType categorises the kind of drift found.
type DriftType string

const (
	DriftMissing  DriftType = "MISSING"   // declared but not live
	DriftExtra    DriftType = "EXTRA"     // live but not declared
	DriftModified DriftType = "MODIFIED"  // attribute mismatch
)

// Result describes a single drift finding.
type Result struct {
	ResourceID string
	Type       DriftType
	Details    string
}

// Detect compares declared state against live resources and returns drift results.
func Detect(declared *state.State, live []state.Resource) []Result {
	var results []Result

	declaredMap := declared.ResourceMap()
	liveMap := make(map[string]state.Resource, len(live))
	for _, r := range live {
		liveMap[r.ID] = r
	}

	// Check declared vs live
	for id, dr := range declaredMap {
		lr, exists := liveMap[id]
		if !exists {
			results = append(results, Result{ResourceID: id, Type: DriftMissing, Details: "resource not found in live infrastructure"})
			continue
		}
		for k, dv := range dr.Attributes {
			if lv, ok := lr.Attributes[k]; !ok || lv != dv {
				results = append(results, Result{ResourceID: id, Type: DriftModified, Details: "attribute '" + k + "' differs: declared='" + dv + "' live='" + lv + "'"})
			}
		}
	}

	// Check live vs declared (extra)
	for id := range liveMap {
		if _, exists := declaredMap[id]; !exists {
			results = append(results, Result{ResourceID: id, Type: DriftExtra, Details: "resource exists in live infrastructure but not declared"})
		}
	}

	return results
}
