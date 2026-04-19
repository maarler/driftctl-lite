// Package redact provides utilities to mask sensitive field values
// in drift results before output or storage.
package redact

import "github.com/driftctl-lite/internal/drift"

// DefaultSensitiveKeys holds common field names considered sensitive.
var DefaultSensitiveKeys = []string{
	"password", "secret", "token", "api_key", "private_key", "credentials",
}

const masked = "***REDACTED***"

// Apply returns a copy of results with sensitive field values masked.
func Apply(results []drift.Result, sensitiveKeys []string) []drift.Result {
	if len(sensitiveKeys) == 0 {
		sensitiveKeys = DefaultSensitiveKeys
	}
	keySet := make(map[string]struct{}, len(sensitiveKeys))
	for _, k := range sensitiveKeys {
		keySet[k] = struct{}{}
	}

	out := make([]drift.Result, len(results))
	for i, r := range results {
		out[i] = drift.Result{
			ResourceID:   r.ResourceID,
			ResourceType: r.ResourceType,
			Status:       r.Status,
			Diffs:        redactDiffs(r.Diffs, keySet),
		}
	}
	return out
}

func redactDiffs(diffs []drift.Diff, keySet map[string]struct{}) []drift.Diff {
	if diffs == nil {
		return nil
	}
	out := make([]drift.Diff, len(diffs))
	for i, d := range diffs {
		if _, sensitive := keySet[d.Field]; sensitive {
			out[i] = drift.Diff{Field: d.Field, Expected: masked, Got: masked}
		} else {
			out[i] = d
		}
	}
	return out
}
