package flatten_test

import (
	"testing"

	"driftctl-lite/internal/diff"
	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/flatten"
)

func makeResult(id, typ string, status drift.Status, diffs map[string]diff.FieldDiff) drift.Result {
	return drift.Result{
		ResourceID:   id,
		ResourceType: typ,
		Status:       status,
		Diffs:        diffs,
	}
}

func TestFlatten_Empty(t *testing.T) {
	records := flatten.Flatten(nil)
	if len(records) != 0 {
		t.Fatalf("expected 0 records, got %d", len(records))
	}
}

func TestFlatten_InSync_NoDiffs(t *testing.T) {
	results := []drift.Result{
		makeResult("vpc-1", "aws_vpc", drift.StatusInSync, nil),
	}
	records := flatten.Flatten(results)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Key != "-" {
		t.Errorf("expected key '-', got %s", records[0].Key)
	}
}

func TestFlatten_Modified_ExpandsDiffs(t *testing.T) {
	diffs := map[string]diff.FieldDiff{
		"cidr": {Wanted: "10.0.0.0/16", Got: "10.1.0.0/16"},
		"name": {Wanted: "prod", Got: "staging"},
	}
	results := []drift.Result{
		makeResult("vpc-1", "aws_vpc", drift.StatusModified, diffs),
	}
	records := flatten.Flatten(results)
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	// sorted keys: cidr, name
	if records[0].Key != "cidr" {
		t.Errorf("expected first key 'cidr', got %s", records[0].Key)
	}
	if records[1].Key != "name" {
		t.Errorf("expected second key 'name', got %s", records[1].Key)
	}
	if records[0].Wanted != "10.0.0.0/16" {
		t.Errorf("unexpected wanted: %s", records[0].Wanted)
	}
}

func TestFlatten_MultipleResults(t *testing.T) {
	results := []drift.Result{
		makeResult("sg-1", "aws_sg", drift.StatusMissing, nil),
		makeResult("sg-2", "aws_sg", drift.StatusExtra, nil),
	}
	records := flatten.Flatten(results)
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].ID != "sg-1" || records[1].ID != "sg-2" {
		t.Errorf("unexpected IDs: %s, %s", records[0].ID, records[1].ID)
	}
}
