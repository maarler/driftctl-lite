// Package policy provides rule-based enforcement on drift results.
package policy

import (
	"fmt"
	"strings"

	"github.com/driftctl-lite/internal/drift"
)

// Rule defines a single enforcement rule.
type Rule struct {
	ResourceType string `json:"resource_type"`
	DisallowExtra bool   `json:"disallow_extra"`
	DisallowMissing bool `json:"disallow_missing"`
	DisallowModified bool `json:"disallow_modified"`
}

// Violation represents a policy breach.
type Violation struct {
	Rule    string
	Resource string
	Reason  string
}

func (v Violation) String() string {
	return fmt.Sprintf("[%s] %s: %s", v.Rule, v.Resource, v.Reason)
}

// Evaluate checks drift results against the given rules and returns violations.
func Evaluate(results []drift.Result, rules []Rule) []Violation {
	var violations []Violation

	ruleMap := make(map[string]Rule)
	for _, r := range rules {
		ruleMap[strings.ToLower(r.ResourceType)] = r
	}

	for _, res := range results {
		ruleKey := strings.ToLower(res.ResourceType)
		rule, ok := ruleMap[ruleKey]
		if !ok {
			continue
		}

		switch res.Status {
		case drift.StatusMissing:
			if rule.DisallowMissing {
				violations = append(violations, Violation{
					Rule:     res.ResourceType + ":disallow_missing",
					Resource: res.ResourceID,
					Reason:   "resource is missing from live infrastructure",
				})
			}
		case drift.StatusExtra:
			if rule.DisallowExtra {
				violations = append(violations, Violation{
					Rule:     res.ResourceType + ":disallow_extra",
					Resource: res.ResourceID,
					Reason:   "resource exists in live but not in state",
				})
			}
		case drift.StatusModified:
			if rule.DisallowModified {
				violations = append(violations, Violation{
					Rule:     res.ResourceType + ":disallow_modified",
					Resource: res.ResourceID,
					Reason:   "resource attributes differ from declared state",
				})
			}
		}
	}

	return violations
}
