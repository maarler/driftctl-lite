package reorder_test

import (
	"testing"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/reorder"
	"driftctl-lite/internal/state"
)

func makeResult(id, rtype string, status drift.DriftStatus) drift.Result {
	return drift.Result{
		Resource: state.Resource{ID: id, Type: rtype},
		Status:   status,
	}
}

func TestApply_Empty(t *testing.T) {
	out := reorder.Apply(nil, reorder.DefaultOptions())
	if out != nil && len(out) != 0 {
		t.Fatalf("expected empty, got %v", out)
	}
}

func TestApply_ByID_Ascending(t *testing.T) {
	results := []drift.Result{
		makeResult("c", "aws_s3", drift.StatusInSync),
		makeResult("a", "aws_s3", drift.StatusMissing),
		makeResult("b", "aws_s3", drift.StatusExtra),
	}
	out := reorder.Apply(results, reorder.Options{By: reorder.FieldID, Ascending: true})
	ids := []string{out[0].Resource.ID, out[1].Resource.ID, out[2].Resource.ID}
	expect := []string{"a", "b", "c"}
	for i, v := range expect {
		if ids[i] != v {
			t.Errorf("pos %d: want %q got %q", i, v, ids[i])
		}
	}
}

func TestApply_ByID_Descending(t *testing.T) {
	results := []drift.Result{
		makeResult("a", "aws_s3", drift.StatusInSync),
		makeResult("c", "aws_s3", drift.StatusMissing),
		makeResult("b", "aws_s3", drift.StatusExtra),
	}
	out := reorder.Apply(results, reorder.Options{By: reorder.FieldID, Ascending: false})
	if out[0].Resource.ID != "c" || out[1].Resource.ID != "b" || out[2].Resource.ID != "a" {
		t.Errorf("unexpected order: %v", out)
	}
}

func TestApply_ByType(t *testing.T) {
	results := []drift.Result{
		makeResult("1", "vpc", drift.StatusInSync),
		makeResult("2", "ec2", drift.StatusMissing),
		makeResult("3", "s3", drift.StatusExtra),
	}
	out := reorder.Apply(results, reorder.Options{By: reorder.FieldType, Ascending: true})
	if out[0].Resource.Type != "ec2" || out[1].Resource.Type != "s3" || out[2].Resource.Type != "vpc" {
		t.Errorf("unexpected type order: %v", out)
	}
}

func TestApply_ByStatus(t *testing.T) {
	results := []drift.Result{
		makeResult("1", "t", drift.StatusInSync),
		makeResult("2", "t", drift.StatusExtra),
		makeResult("3", "t", drift.StatusMissing),
		makeResult("4", "t", drift.StatusModified),
	}
	out := reorder.Apply(results, reorder.Options{By: reorder.FieldStatus, Ascending: true})
	if out[0].Status != drift.StatusMissing {
		t.Errorf("first should be missing, got %v", out[0].Status)
	}
	if out[1].Status != drift.StatusModified {
		t.Errorf("second should be modified, got %v", out[1].Status)
	}
	if out[2].Status != drift.StatusExtra {
		t.Errorf("third should be extra, got %v", out[2].Status)
	}
	if out[3].Status != drift.StatusInSync {
		t.Errorf("fourth should be in-sync, got %v", out[3].Status)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := reorder.DefaultOptions()
	if opts.By != reorder.FieldID {
		t.Errorf("default field should be ID, got %v", opts.By)
	}
	if !opts.Ascending {
		t.Error("default should be ascending")
	}
}
