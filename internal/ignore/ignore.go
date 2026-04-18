package ignore

import (
	"bufio"
	"os"
	"strings"
)

// Rule represents a single ignore rule.
type Rule struct {
	ResourceType string
	ResourceID   string
}

// List holds a set of ignore rules.
type List struct {
	rules []Rule
}

// LoadFromFile reads ignore rules from a file.
// Each line has the format: resource_type/resource_id or resource_type/*
func LoadFromFile(path string) (*List, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &List{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var rules []Rule
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "/", 2)
		if len(parts) != 2 {
			continue
		}
		rules = append(rules, Rule{ResourceType: parts[0], ResourceID: parts[1]})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &List{rules: rules}, nil
}

// Matches returns true if the given type+id is covered by any rule.
func (l *List) Matches(resourceType, resourceID string) bool {
	for _, r := range l.rules {
		if r.ResourceType != resourceType {
			continue
		}
		if r.ResourceID == "*" || r.ResourceID == resourceID {
			return true
		}
	}
	return false
}

// FilterIgnored removes drift results whose resource matches the ignore list.
func (l *List) FilterIgnored(results []DriftResult) []DriftResult {
	var out []DriftResult
	for _, r := range results {
		if !l.Matches(r.ResourceType, r.ResourceID) {
			out = append(out, r)
		}
	}
	return out
}

// DriftResult is a minimal interface-compatible struct mirroring drift.Result.
type DriftResult struct {
	ResourceType string
	ResourceID   string
	Status       string
}
