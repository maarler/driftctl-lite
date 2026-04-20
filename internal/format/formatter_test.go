package format_test

import (
	"strings"
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/format"
)

func makeResult(id, status string, diffs map[string]drift.DiffDetail) drift.Result {
	return drift.Result{
		ResourceID:   id,
		ResourceType: "aws_s3_bucket",
		Status:       status,
		Diffs:        diffs,
	}
}

func TestApply_UpperCaseStatus(t *testing.T) {
	input := []drift.Result{makeResult("r1", "missing", nil)}
	out := format.Apply(input, format.DefaultOptions())
	if out[0].Status != "MISSING" {
		t.Errorf("expected MISSING, got %s", out[0].Status)
	}
}

func TestApply_NoUpperCase(t *testing.T) {
	opts := format.Options{UpperCaseStatus: false, MaxValueLen: 0}
	input := []drift.Result{makeResult("r1", "extra", nil)}
	out := format.Apply(input, opts)
	if out[0].Status != "extra" {
		t.Errorf("expected extra, got %s", out[0].Status)
	}
}

func TestApply_TruncatesLongValues(t *testing.T) {
	long := strings.Repeat("a", 100)
	diffs := map[string]drift.DiffDetail{
		"name": {Expected: long, Got: long},
	}
	input := []drift.Result{makeResult("r1", "modified", diffs)}
	opts := format.Options{MaxValueLen: 10, UpperCaseStatus: false}
	out := format.Apply(input, opts)

	d, ok := out[0].Diffs["name"]
	if !ok {
		t.Fatal("expected diff key 'name'")
	}
	if len([]rune(d.Expected.(string))) > 11 { // 10 chars + ellipsis
		t.Errorf("expected value not truncated: %s", d.Expected)
	}
}

func TestApply_NoLimit_KeepsFullValue(t *testing.T) {
	long := strings.Repeat("b", 200)
	diffs := map[string]drift.DiffDetail{
		"tag": {Expected: long, Got: ""},
	}
	input := []drift.Result{makeResult("r2", "modified", diffs)}
	opts := format.Options{MaxValueLen: 0, UpperCaseStatus: false}
	out := format.Apply(input, opts)

	d := out[0].Diffs["tag"]
	if d.Expected != long {
		t.Errorf("expected full value to be preserved")
	}
}

func TestApply_Empty(t *testing.T) {
	out := format.Apply(nil, format.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d items", len(out))
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	input := []drift.Result{makeResult("r1", "in_sync", nil)}
	original := input[0].Status
	format.Apply(input, format.DefaultOptions())
	if input[0].Status != original {
		t.Error("Apply must not mutate the input slice")
	}
}
