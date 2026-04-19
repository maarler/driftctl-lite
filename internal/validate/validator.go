// Package validate provides schema validation for drift results.
package validate

import (
	"errors"
	"fmt"

	"github.com/driftctl-lite/internal/drift"
)

// Rule defines a validation rule applied to a drift result.
type Rule struct {
	Name    string
	Message string
	Check   func(r drift.Result) bool
}

// Violation represents a failed validation rule for a specific resource.
type Violation struct {
	ResourceID string
	RuleName   string
	Message    string
}

// Validate applies all rules to each result and returns any violations.
func Validate(results []drift.Result, rules []Rule) ([]Violation, error) {
	if len(rules) == 0 {
		return nil, errors.New("no validation rules provided")
	}

	var violations []Violation
	for _, r := range results {
		for _, rule := range rules {
			if !rule.Check(r) {
				violations = append(violations, Violation{
					ResourceID: r.ResourceID,
					RuleName:   rule.Name,
					Message:    fmt.Sprintf("%s: %s", r.ResourceID, rule.Message),
				})
			}
		}
	}
	return violations, nil
}

// DefaultRules returns a standard set of validation rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name:    "non-empty-id",
			Message: "resource ID must not be empty",
			Check: func(r drift.Result) bool {
				return r.ResourceID != ""
			},
		},
		{
			Name:    "non-empty-type",
			Message: "resource type must not be empty",
			Check: func(r drift.Result) bool {
				return r.ResourceType != ""
			},
		},
	}
}
