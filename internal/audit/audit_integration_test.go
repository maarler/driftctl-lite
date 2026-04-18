package audit_test

import (
	"testing"

	"github.com/snyk/driftctl-lite/internal/audit"
	"github.com/snyk/driftctl-lite/internal/drift"
)

// TestAudit_MultiRun_RoundTrip simulates multiple detection runs being
// logged and then replayed, verifying cumulative counts are correct.
func TestAudit_MultiRun_RoundTrip(t *testing.T) {
	path := tempLogPath(t)
	l := audit.NewLogger(path)

	runs := [][]drift.Result{
		{
			{ResourceID: "r1", ResourceType: "aws_vpc", Status: drift.StatusInSync},
		},
		{
			{ResourceID: "r2", ResourceType: "aws_sg", Status: drift.StatusMissing},
			{ResourceID: "r3", ResourceType: "aws_sg", Status: drift.StatusModified},
		},
		{
			{ResourceID: "r4", ResourceType: "aws_s3", Status: drift.StatusExtra},
		},
	}

	for i, run := range runs {
		if err := l.Record("state.json", "file", run); err != nil {
			t.Fatalf("run %d record failed: %v", i, err)
		}
	}

	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	expectedDrift := []int{0, 2, 1}
	for i, e := range entries {
		if e.DriftCount != expectedDrift[i] {
			t.Errorf("entry %d: expected drift %d, got %d", i, expectedDrift[i], e.DriftCount)
		}
	}
}
