// Package label provides utilities for labeling and filtering drift results
// based on key-value metadata attached to resources.
package label

import (
	"strings"

	"driftctl-lite/internal/drift"
)

// Rule represents a label selector in the form "key" or "key=value".
type Rule struct {
	Key   string
	Value string // empty means match any value
}

// ParseRule parses a label selector string into a Rule.
func ParseRule(s string) (Rule, error) {
	if s == "" {
		return Rule{}, nil
	}
	parts := strings.SplitN(s, "=", 2)
	r := Rule{Key: strings.TrimSpace(parts[0])}
	if len(parts) == 2 {
		r.Value = strings.TrimSpace(parts[1])
	}
	return r, nil
}

// Matches reports whether the given labels satisfy the rule.
func (r Rule) Matches(labels map[string]string) bool {
	if r.Key == "" {
		return true
	}
	v, ok := labels[r.Key]
	if !ok {
		return false
	}
	if r.Value == "" {
		return true
	}
	return v == r.Value
}

// Filter returns only those results whose resource labels match all provided rules.
func Filter(results []drift.Result, rules []Rule) []drift.Result {
	if len(rules) == 0 {
		return results
	}
	var out []drift.Result
	for _, res := range results {
		if matchesAll(res.Resource.Labels, rules) {
			out = append(out, res)
		}
	}
	return out
}

func matchesAll(labels map[string]string, rules []Rule) bool {
	for _, r := range rules {
		if !r.Matches(labels) {
			return false
		}
	}
	return true
}
