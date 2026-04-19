package label_test

import (
	"testing"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/label"
	"driftctl-lite/internal/state"
)

func makeResult(id, rtype string, labels map[string]string) drift.Result {
	return drift.Result{
		Resource: state.Resource{
			ID:     id,
			Type:   rtype,
			Labels: labels,
		},
		Status: drift.StatusInSync,
	}
}

func TestParseRule_KeyOnly(t *testing.T) {
	r, err := label.ParseRule("env")
	if err != nil || r.Key != "env" || r.Value != "" {
		t.Fatalf("unexpected rule: %+v err: %v", r, err)
	}
}

func TestParseRule_KeyValue(t *testing.T) {
	r, err := label.ParseRule("env=prod")
	if err != nil || r.Key != "env" || r.Value != "prod" {
		t.Fatalf("unexpected rule: %+v err: %v", r, err)
	}
}

func TestParseRule_Empty(t *testing.T) {
	r, _ := label.ParseRule("")
	if !r.Matches(map[string]string{}) {
		t.Fatal("empty rule should match anything")
	}
}

func TestFilter_NoRules(t *testing.T) {
	results := []drift.Result{makeResult("a", "aws_s3", map[string]string{"env": "prod"})}
	out := label.Filter(results, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

func TestFilter_ByKey(t *testing.T) {
	results := []drift.Result{
		makeResult("a", "aws_s3", map[string]string{"env": "prod"}),
		makeResult("b", "aws_ec2", map[string]string{}),
	}
	r, _ := label.ParseRule("env")
	out := label.Filter(results, []label.Rule{r})
	if len(out) != 1 || out[0].Resource.ID != "a" {
		t.Fatalf("unexpected filter result: %+v", out)
	}
}

func TestFilter_ByKeyValue(t *testing.T) {
	results := []drift.Result{
		makeResult("a", "aws_s3", map[string]string{"env": "prod"}),
		makeResult("b", "aws_ec2", map[string]string{"env": "staging"}),
	}
	r, _ := label.ParseRule("env=prod")
	out := label.Filter(results, []label.Rule{r})
	if len(out) != 1 || out[0].Resource.ID != "a" {
		t.Fatalf("unexpected filter result: %+v", out)
	}
}

func TestFilter_MultipleRules(t *testing.T) {
	results := []drift.Result{
		makeResult("a", "aws_s3", map[string]string{"env": "prod", "team": "ops"}),
		makeResult("b", "aws_ec2", map[string]string{"env": "prod"}),
	}
	r1, _ := label.ParseRule("env=prod")
	r2, _ := label.ParseRule("team")
	out := label.Filter(results, []label.Rule{r1, r2})
	if len(out) != 1 || out[0].Resource.ID != "a" {
		t.Fatalf("unexpected filter result: %+v", out)
	}
}
