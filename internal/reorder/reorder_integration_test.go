package reorder_test

import (
	"testing"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/reorder"
	"driftctl-lite/internal/state"
)

// TestReorder_RoundTrip verifies that applying two opposite sort orders
// produces a consistent, deterministic result.
func TestReorder_RoundTrip(t *testing.T) {
	original := []drift.Result{
		{Resource: state.Resource{ID: "z", Type: "ec2"}, Status: drift.StatusInSync},
		{Resource: state.Resource{ID: "m", Type: "s3"}, Status: drift.StatusMissing},
		{Resource: state.Resource{ID: "a", Type: "vpc"}, Status: drift.StatusExtra},
		{Resource: state.Resource{ID: "f", Type: "iam"}, Status: drift.StatusModified},
	}

	asc := reorder.Apply(append([]drift.Result{}, original...), reorder.Options{By: reorder.FieldID, Ascending: true})
	desc := reorder.Apply(append([]drift.Result{}, original...), reorder.Options{By: reorder.FieldID, Ascending: false})

	if len(asc) != len(desc) {
		t.Fatalf("length mismatch: %d vs %d", len(asc), len(desc))
	}

	// Ascending first element should equal descending last element.
	if asc[0].Resource.ID != desc[len(desc)-1].Resource.ID {
		t.Errorf("asc[0]=%q should equal desc[last]=%q",
			asc[0].Resource.ID, desc[len(desc)-1].Resource.ID)
	}

	// Re-sorting descending result ascending should restore original ascending order.
	re := reorder.Apply(append([]drift.Result{}, desc...), reorder.Options{By: reorder.FieldID, Ascending: true})
	for i := range asc {
		if asc[i].Resource.ID != re[i].Resource.ID {
			t.Errorf("pos %d: asc=%q re=%q", i, asc[i].Resource.ID, re[i].Resource.ID)
		}
	}
}
