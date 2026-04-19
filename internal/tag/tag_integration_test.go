package tag_test

import (
	"testing"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/tag"
)

// TestTagFilter_RoundTrip exercises ParseRule -> Filter together.
func TestTagFilter_RoundTrip(t *testing.T) {
	results := []drift.Result{
		{
			ResourceID:   "i-001",
			ResourceType: "aws_instance",
			Status:       drift.StatusModified,
			Declared:     map[string]interface{}{"env": "prod", "owner": "alice"},
		},
		{
			ResourceID:   "i-002",
			ResourceType: "aws_instance",
			Status:       drift.StatusInSync,
			Declared:     map[string]interface{}{"env": "staging", "owner": "bob"},
		},
		{
			ResourceID:   "i-003",
			ResourceType: "aws_instance",
			Status:       drift.StatusMissing,
			Declared:     map[string]interface{}{"env": "prod", "owner": "bob"},
		},
	}

	ruleStr := "env=prod"
	rule, err := tag.ParseRule(ruleStr)
	if err != nil {
		t.Fatalf("ParseRule: %v", err)
	}

	out := tag.Filter(results, []tag.Rule{rule})
	if len(out) != 2 {
		t.Fatalf("expected 2 prod results, got %d", len(out))
	}
	for _, r := range out {
		if r.Declared["env"] != "prod" {
			t.Errorf("unexpected env value: %v", r.Declared["env"])
		}
	}
}
