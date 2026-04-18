package drift

import (
	"encoding/json"
)

// Result holds the categorised differences found during drift detection.
type Result struct {
	Missing  []Resource `json:"missing"`
	Extra    []Resource `json:"extra"`
	Modified []Resource `json:"modified"`
}

// HasDrift returns true when any differences were detected.
func (r Result) HasDrift() bool {
	return len(r.Missing)+len(r.Extra)+len(r.Modified) > 0
}

// Summary returns a short human-readable summary string.
func (r Result) Summary() string {
	if !r.HasDrift() {
		return "no drift detected"
	}
	return fmt.Sprintf("%d missing, %d extra, %d modified",
		len(r.Missing), len(r.Extra), len(r.Modified))
}

// MarshalJSON serialises the result as JSON bytes.
func (r Result) MarshalJSON() ([]byte, error) {
	type alias Result
	return json.Marshal(alias(r))
}
