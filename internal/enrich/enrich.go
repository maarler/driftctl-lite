// Package enrich attaches metadata to drift results for richer reporting.
package enrich

import "github.com/driftctl-lite/internal/drift"

// Metadata holds extra annotations to attach to a result.
type Metadata struct {
	Owner       string
	Environment string
	CostCenter  string
	Custom      map[string]string
}

// Rule maps a resource type to metadata.
type Rule struct {
	ResourceType string
	Metadata     Metadata
}

// Apply enriches each drift result with metadata from matching rules.
// Rules are matched by ResourceType; the first matching rule wins.
func Apply(results []drift.Result, rules []Rule) []drift.Result {
	if len(rules) == 0 {
		return results
	}
	index := make(map[string]Metadata, len(rules))
	for _, r := range rules {
		index[r.ResourceType] = r.Metadata
	}
	enriched := make([]drift.Result, len(results))
	for i, res := range results {
		meta, ok := index[res.ResourceType]
		if ok {
			res = attachMeta(res, meta)
		}
		enriched[i] = res
	}
	return enriched
}

func attachMeta(res drift.Result, meta Metadata) drift.Result {
	if res.Attributes == nil {
		res.Attributes = make(map[string]string)
	}
	if meta.Owner != "" {
		res.Attributes["_owner"] = meta.Owner
	}
	if meta.Environment != "" {
		res.Attributes["_environment"] = meta.Environment
	}
	if meta.CostCenter != "" {
		res.Attributes["_cost_center"] = meta.CostCenter
	}
	for k, v := range meta.Custom {
		res.Attributes["_"+k] = v
	}
	return res
}
