package rollup_test

import (
	"testing"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/rollup"
)

func TestRollup_RoundTrip(t *testing.T) {
	results := []drift.Result{
		{ResourceID: "i-1", ResourceType: "aws_instance", Status: drift.StatusMissing},
		{ResourceID: "i-2", ResourceType: "aws_instance", Status: drift.StatusModified},
		{ResourceID: "i-3", ResourceType: "aws_instance", Status: drift.StatusInSync},
		{ResourceID: "rds-1", ResourceType: "aws_db_instance", Status: drift.StatusExtra},
	}

	report := rollup.Compute(results)

	if len(report) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(report))
	}

	totals := map[string]int{}
	for _, e := range report {
		totals[e.Type] = e.Total
	}
	if totals["aws_instance"] != 3 {
		t.Errorf("expected 3 aws_instance, got %d", totals["aws_instance"])
	}
	if totals["aws_db_instance"] != 1 {
		t.Errorf("expected 1 aws_db_instance, got %d", totals["aws_db_instance"])
	}
}
