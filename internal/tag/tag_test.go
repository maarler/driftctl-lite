package tag_test

import (
	"testing"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/tag"
)

func makeResult(declared map[string]interface{}) drift.Result {
	return drift.Result{
		ResourceID:   "res-1",
		ResourceType: "aws_instance",
		Status:       drift.StatusInSync,
		Declared:     declared,
	}
}

func TestParseRule_KeyOnly(t *testing.T) {
	r, err := tag.ParseRule("env")
	if err != nil || r.Key != "env" || r.Value != "" {
		t.Fatalf("unexpected: %v %v", r, err)
	}
}

func TestParseRule_KeyValue(t *testing.T) {
	r, err := tag.ParseRule("env=prod")
	if err != nil || r.Key != "env" || r.Value != "prod" {
		t.Fatalf("unexpected: %v %v", r, err)
	}
}

func TestParseRule_Empty(t *testing.T) {
	_, err := tag.ParseRule("")
	if err == nil {
		t.Fatal("expected error for empty rule")
	}
}

func TestFilter_NoRules(t *testing.T) {
	results := []drift.Result{makeResult(map[string]interface{}{"env": "prod"})}
	out := tag.Filter(results, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

func TestFilter_ByKeyOnly(t *testing.T) {
	results := []drift.Result{
		makeResult(map[string]interface{}{"env": "prod"}),
		makeResult(map[string]interface{}{"team": "platform"}),
	}
	rules := []tag.Rule{{Key: "env"}}
	out := tag.Filter(results, rules)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

func TestFilter_ByKeyValue(t *testing.T) {
	results := []drift.Result{
		makeResult(map[string]interface{}{"env": "prod"}),
		makeResult(map[string]interface{}{"env": "staging"}),
	}
	rules := []tag.Rule{{Key: "env", Value: "prod"}}
	out := tag.Filter(results, rules)
	if len(out) != 1 || out[0].Declared["env"] != "prod" {
		t.Fatalf("unexpected results: %v", out)
	}
}

func TestFilter_MultipleRules(t *testing.T) {
	results := []drift.Result{
		makeResult(map[string]interface{}{"env": "prod", "team": "platform"}),
		makeResult(map[string]interface{}{"env": "prod"}),
	}
	rules := []tag.Rule{{Key: "env", Value: "prod"}, {Key: "team"}}
	out := tag.Filter(results, rules)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

func TestFilter_NoMatch(t *testing.T) {
	results := []drift.Result{
		makeResult(map[string]interface{}{"env": "dev"}),
	}
	rules := []tag.Rule{{Key: "env", Value: "prod"}}
	out := tag.Filter(results, rules)
	if len(out) != 0 {
		t.Fatalf("expected 0, got %d", len(out))
	}
}
