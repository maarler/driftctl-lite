package explain_test

import (
	"strings"
	"testing"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/explain"
)

func makeResult(id, typ string, status drift.Status, diffs map[string]drift.Diff) drift.Result {
	return drift.Result{ID: id, Type: typ, Status: status, Diffs: diffs}
}

func TestExplain_OK(t *testing.T) {
	results := []drift.Result{makeResult("vpc-1", "aws_vpc", drift.StatusOK, nil)}
	exps := explain.Explain(results)
	if len(exps) != 1 {
		t.Fatalf("expected 1 explanation, got %d", len(exps))
	}
	if !strings.Contains(exps[0].Summary, "in sync") {
		t.Errorf("expected 'in sync' in summary, got: %s", exps[0].Summary)
	}
}

func TestExplain_Missing(t *testing.T) {
	results := []drift.Result{makeResult("sg-1", "aws_sg", drift.StatusMissing, nil)}
	exps := explain.Explain(results)
	if !strings.Contains(exps[0].Summary, "missing from live") {
		t.Errorf("unexpected summary: %s", exps[0].Summary)
	}
	if len(exps[0].Details) == 0 {
		t.Error("expected at least one detail for missing resource")
	}
}

func TestExplain_Extra(t *testing.T) {
	results := []drift.Result{makeResult("i-123", "aws_instance", drift.StatusExtra, nil)}
	exps := explain.Explain(results)
	if !strings.Contains(exps[0].Summary, "not declared") {
		t.Errorf("unexpected summary: %s", exps[0].Summary)
	}
}

func TestExplain_Modified(t *testing.T) {
	diffs := map[string]drift.Diff{
		"instance_type": {Expected: "t2.micro", Got: "t3.small"},
	}
	results := []drift.Result{makeResult("i-456", "aws_instance", drift.StatusModified, diffs)}
	exps := explain.Explain(results)
	if !strings.Contains(exps[0].Summary, "configuration differences") {
		t.Errorf("unexpected summary: %s", exps[0].Summary)
	}
	if len(exps[0].Details) != 1 {
		t.Errorf("expected 1 diff detail, got %d", len(exps[0].Details))
	}
	if !strings.Contains(exps[0].Details[0], "instance_type") {
		t.Errorf("expected field name in detail: %s", exps[0].Details[0])
	}
}

func TestExplain_Empty(t *testing.T) {
	exps := explain.Explain([]drift.Result{})
	if len(exps) != 0 {
		t.Errorf("expected 0 explanations, got %d", len(exps))
	}
}
