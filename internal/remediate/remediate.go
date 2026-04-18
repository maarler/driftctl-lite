// Package remediate provides suggested remediation actions for detected drift.
package remediate

import (
	"fmt"

	"github.com/owner/driftctl-lite/internal/drift"
)

// Action describes a suggested fix for a drifted resource.
type Action struct {
	ResourceID string
	ResourceType string
	Severity string
	Message string
}

// Suggest returns a list of remediation actions for the given drift results.
func Suggest(results []drift.Result) []Action {
	var actions []Action
	for _, r := range results {
		if r.Status == drift.StatusInSync {
			continue
		}
		actions = append(actions, Action{
			ResourceID:   r.ResourceID,
			ResourceType: r.ResourceType,
			Severity:     severity(r.Status),
			Message:      message(r),
		})
	}
	return actions
}

func severity(status drift.Status) string {
	switch status {
	case drift.StatusMissing:
		return "high"
	case drift.StatusExtra:
		return "medium"
	case drift.StatusModified:
		return "low"
	default:
		return "unknown"
	}
}

func message(r drift.Result) string {
	switch r.Status {
	case drift.StatusMissing:
		return fmt.Sprintf("Resource %q (%s) is declared but missing in live infrastructure. Re-apply your IaC configuration.", r.ResourceID, r.ResourceType)
	case drift.StatusExtra:
		return fmt.Sprintf("Resource %q (%s) exists in live infrastructure but is not declared. Import or remove it.", r.ResourceID, r.ResourceType)
	case drift.StatusModified:
		return fmt.Sprintf("Resource %q (%s) has configuration drift. Review and reconcile changes.", r.ResourceID, r.ResourceType)
	default:
		return "No action needed."
	}
}
