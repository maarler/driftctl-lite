package rank_test

import (
	"testing"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/rank"
)

func makeResult(id string, status drift.Status) drift.Result {
	return drift.Result{ResourceID: id, ResourceType: "aws_instance", Status: status}
}

func TestByPriority_OrdersCorrectly(t *testing.T) {
	input := []drift.Result{
		makeResult("c", drift.StatusInSync),
		makeResult("b", drift.StatusExtra),
		makeResult("a", drift.StatusMissing),
		makeResult("d", drift.StatusModified),
	}
	out := rank.ByPriority(input)
	expected := []string{"a", "d", "b", "c"}
	for i, r := range out {
		if r.ResourceID != expected[i] {
			t.Errorf("pos %d: got %s, want %s", i, r.ResourceID, expected[i])
		}
	}
}

func TestByPriority_SameStatus_SortsByID(t *testing.T) {
	input := []drift.Result{
		makeResult("z", drift.StatusMissing),
		makeResult("a", drift.StatusMissing),
		makeResult("m", drift.StatusMissing),
	}
	out := rank.ByPriority(input)
	if out[0].ResourceID != "a" || out[1].ResourceID != "m" || out[2].ResourceID != "z" {
		t.Errorf("unexpected order: %v", out)
	}
}

func TestTopN_LimitsResults(t *testing.T) {
	input := []drift.Result{
		makeResult("a", drift.StatusMissing),
		makeResult("b", drift.StatusModified),
		makeResult("c", drift.StatusExtra),
		makeResult("d", drift.StatusInSync),
	}
	out := rank.TopN(input, 2)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
	if out[0].ResourceID != "a" {
		t.Errorf("expected a, got %s", out[0].ResourceID)
	}
}

func TestTopN_NLargerThanResults(t *testing.T) {
	input := []drift.Result{
		makeResult("a", drift.StatusMissing),
	}
	out := rank.TopN(input, 10)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

func TestByPriority_DoesNotMutateInput(t *testing.T) {
	input := []drift.Result{
		makeResult("b", drift.StatusExtra),
		makeResult("a", drift.StatusMissing),
	}
	_ = rank.ByPriority(input)
	if input[0].ResourceID != "b" {
		t.Error("input slice was mutated")
	}
}
