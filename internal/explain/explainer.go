// Package explain provides human-readable explanations for detected drift.
package explain

import (
	"fmt"
	"strings"

	"driftctl-lite/internal/drift"
)

// Explanation holds a human-readable description of a single drift result.
type Explanation struct {
	ResourceID   string
	ResourceType string
	Summary      string
	Details      []string
}

// Explain generates explanations for a slice of drift results.
func Explain(results []drift.Result) []Explanation {
	explanations := make([]Explanation, 0, len(results))
	for _, r := range results {
		exp := explainOne(r)
		explanations = append(explanations, exp)
	}
	return explanations
}

func explainOne(r drift.Result) Explanation {
	exp := Explanation{
		ResourceID:   r.ID,
		ResourceType: r.Type,
	}

	switch r.Status {
	case drift.StatusMissing:
		exp.Summary = fmt.Sprintf("Resource '%s' (%s) is declared but missing from live infrastructure.", r.ID, r.Type)
		exp.Details = []string{"Ensure the resource has been provisioned or remove it from the state file."}
	case drift.StatusExtra:
		exp.Summary = fmt.Sprintf("Resource '%s' (%s) exists in live infrastructure but is not declared.", r.ID, r.Type)
		exp.Details = []string{"Import the resource into your state or remove it from live infrastructure."}
	case drift.StatusModified:
		exp.Summary = fmt.Sprintf("Resource '%s' (%s) has configuration differences.", r.ID, r.Type)
		exp.Details = buildDiffDetails(r.Diffs)
	case drift.StatusOK:
		exp.Summary = fmt.Sprintf("Resource '%s' (%s) is in sync.", r.ID, r.Type)
	}

	return exp
}

func buildDiffDetails(diffs map[string]drift.Diff) []string {
	details := make([]string, 0, len(diffs))
	for key, d := range diffs {
		details = append(details, fmt.Sprintf("  Field '%s': expected=%s, got=%s",
			key,
			strings.TrimSpace(fmt.Sprintf("%v", d.Expected)),
			strings.TrimSpace(fmt.Sprintf("%v", d.Got))))
	}
	return details
}
