package checkpoint_test

import (
	"testing"

	"driftctl-lite/internal/checkpoint"
	"driftctl-lite/internal/drift"
)

// TestCheckpoint_MultiSave_RoundTrip verifies that multiple checkpoints
// can be saved independently and loaded back with correct isolation.
func TestCheckpoint_MultiSave_RoundTrip(t *testing.T) {
	dir := t.TempDir()

	before := []drift.Result{
		{ResourceID: "i-001", ResourceType: "aws_instance", Status: drift.StatusInSync},
	}
	after := []drift.Result{
		{ResourceID: "i-001", ResourceType: "aws_instance", Status: drift.StatusModified},
		{ResourceID: "i-002", ResourceType: "aws_instance", Status: drift.StatusExtra},
	}

	if err := checkpoint.Save(dir, "before", before); err != nil {
		t.Fatalf("Save before: %v", err)
	}
	if err := checkpoint.Save(dir, "after", after); err != nil {
		t.Fatalf("Save after: %v", err)
	}

	names, err := checkpoint.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 checkpoints, got %d", len(names))
	}

	eb, err := checkpoint.Load(dir, "before")
	if err != nil {
		t.Fatalf("Load before: %v", err)
	}
	if len(eb.Results) != 1 || eb.Results[0].Status != drift.StatusInSync {
		t.Errorf("before checkpoint mismatch: %+v", eb.Results)
	}

	ea, err := checkpoint.Load(dir, "after")
	if err != nil {
		t.Fatalf("Load after: %v", err)
	}
	if len(ea.Results) != 2 {
		t.Errorf("after checkpoint should have 2 results, got %d", len(ea.Results))
	}
}
