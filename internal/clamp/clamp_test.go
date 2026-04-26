package clamp_test

import (
	"testing"

	"github.com/driftctl-lite/internal/clamp"
	"github.com/driftctl-lite/internal/drift"
)

func makeResult(id, typ, status string, priority int, diffs map[string]drift.DiffValue) drift.Result {
	return drift.Result{
		ResourceID:   id,
		ResourceType: typ,
		Status:       status,
		Meta:         map[string]interface{}{"priority": priority},
		Diffs:        diffs,
	}
}

func TestApply_Empty(t *testing.T) {
	results := clamp.Apply(nil, clamp.DefaultOptions())
	if len(results) != 0 {
		t.Fatalf("expected empty, got %d", len(results))
	}
}

func TestApply_PriorityWithinBounds(t *testing.T) {
	r := makeResult("r1", "aws_s3", "modified", 50, nil)
	out := clamp.Apply([]drift.Result{r}, clamp.DefaultOptions())
	if out[0].Meta["priority"].(int) != 50 {
		t.Fatalf("expected 50, got %v", out[0].Meta["priority"])
	}
}

func TestApply_PriorityBelowMin(t *testing.T) {
	r := makeResult("r1", "aws_s3", "modified", -5, nil)
	opts := clamp.DefaultOptions()
	out := clamp.Apply([]drift.Result{r}, opts)
	if out[0].Meta["priority"].(int) != opts.MinPriority {
		t.Fatalf("expected %d, got %v", opts.MinPriority, out[0].Meta["priority"])
	}
}

func TestApply_PriorityAboveMax(t *testing.T) {
	r := makeResult("r1", "aws_s3", "modified", 200, nil)
	opts := clamp.DefaultOptions()
	out := clamp.Apply([]drift.Result{r}, opts)
	if out[0].Meta["priority"].(int) != opts.MaxPriority {
		t.Fatalf("expected %d, got %v", opts.MaxPriority, out[0].Meta["priority"])
	}
}

func TestApply_DiffKeysCapped(t *testing.T) {
	diffs := map[string]drift.DiffValue{
		"a": {Declared: "1", Live: "2"},
		"b": {Declared: "3", Live: "4"},
		"c": {Declared: "5", Live: "6"},
	}
	r := makeResult("r2", "aws_ec2", "modified", 10, diffs)
	opts := clamp.Options{MinPriority: 0, MaxPriority: 100, MaxDiffKeys: 2}
	out := clamp.Apply([]drift.Result{r}, opts)
	if len(out[0].Diffs) != 2 {
		t.Fatalf("expected 2 diff keys, got %d", len(out[0].Diffs))
	}
}

func TestApply_DiffKeysUnderLimit_NotTrimmed(t *testing.T) {
	diffs := map[string]drift.DiffValue{
		"x": {Declared: "a", Live: "b"},
	}
	r := makeResult("r3", "aws_rds", "modified", 20, diffs)
	opts := clamp.DefaultOptions()
	out := clamp.Apply([]drift.Result{r}, opts)
	if len(out[0].Diffs) != 1 {
		t.Fatalf("expected 1 diff key, got %d", len(out[0].Diffs))
	}
}

func TestApply_ZeroMaxDiffKeys_NoTrimming(t *testing.T) {
	diffs := map[string]drift.DiffValue{
		"a": {}, "b": {}, "c": {}, "d": {}, "e": {},
	}
	r := makeResult("r4", "aws_lambda", "modified", 30, diffs)
	opts := clamp.Options{MinPriority: 0, MaxPriority: 100, MaxDiffKeys: 0}
	out := clamp.Apply([]drift.Result{r}, opts)
	if len(out[0].Diffs) != 5 {
		t.Fatalf("expected 5 diff keys, got %d", len(out[0].Diffs))
	}
}
