// Package tag provides filtering and validation of drift results by resource tags.
package tag

import (
	"fmt"
	"strings"

	"driftctl-lite/internal/drift"
)

// Rule represents a tag key/value filter.
type Rule struct {
	Key   string
	Value string // empty means match any value
}

// ParseRule parses a tag rule string in the form "key=value" or "key".
func ParseRule(s string) (Rule, error) {
	if s == "" {
		return Rule{}, fmt.Errorf("tag rule must not be empty")
	}
	parts := strings.SplitN(s, "=", 2)
	r := Rule{Key: parts[0]}
	if len(parts) == 2 {
		r.Value = parts[1]
	}
	return r, nil
}

// Filter returns only those results whose declared attributes contain a tag
// matching every provided rule.
func Filter(results []drift.Result, rules []Rule) []drift.Result {
	if len(rules) == 0 {
		return results
	}
	var out []drift.Result
	for _, r := range results {
		if matchesAll(r, rules) {
			out = append(out, r)
		}
	}
	return out
}

func matchesAll(r drift.Result, rules []Rule) bool {
	for _, rule := range rules {
		if !matchesRule(r, rule) {
			return false
		}
	}
	return true
}

func matchesRule(r drift.Result, rule Rule) bool {
	val, ok := r.Declared[rule.Key]
	if !ok {
		return false
	}
	if rule.Value == "" {
		return true
	}
	return fmt.Sprintf("%v", val) == rule.Value
}
