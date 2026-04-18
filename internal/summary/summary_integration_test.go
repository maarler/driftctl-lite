package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-lite/internal/drift"
	"github.com/owner/driftctl-lite/internal/summary"
)

func TestComputeAndPrint_RoundTrip(t *testing.T) {
	results := []drift.Result{
		{ResourceID: "i-001", ResourceType: "aws_instance", Status: drift.StatusInSync},
		{ResourceID: "i-002", ResourceType: "aws_instance", Status: drift.StatusMissing},
		{ResourceID: "sg-001", ResourceType: "aws_security_group", Status: drift.StatusExtra},
		{ResourceID: "s3-001", ResourceType: "aws_s3_bucket", Status: drift.StatusModified},
	}

	stats := summary.Compute(results)

	if stats.Total != 4 {
		t.Fatalf("expected total 4, got %d", stats.Total)
	}

	var buf bytes.Buffer
	summary.Print(&buf, stats)
	out := buf.String()

	expected := []string{
		"4 resource(s) checked",
		"In sync:  1",
		"Missing:  1",
		"Extra:    1",
		"Modified: 1",
		"Drift detected!",
	}
	for _, e := range expected {
		if !strings.Contains(out, e) {
			t.Errorf("output missing %q\nfull output:\n%s", e, out)
		}
	}
}
