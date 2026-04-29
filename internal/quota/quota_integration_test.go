package quota_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/quota"
)

// TestQuota_RoundTrip verifies that applying quota with drop enabled produces
// a stable result when applied twice (idempotency).
func TestQuota_RoundTrip(t *testing.T) {
	results := []drift.Result{
		makeResult("a", "aws_s3_bucket", "missing"),
		makeResult("b", "aws_s3_bucket", "extra"),
		makeResult("c", "aws_s3_bucket", "modified"),
		makeResult("d", "aws_iam_role", "missing"),
		makeResult("e", "aws_iam_role", "in_sync"),
	}

	opts := quota.Options{MaxPerType: 2, DropExceeding: true}

	first := quota.Apply(results, opts)
	if len(first.Results) != 4 {
		t.Fatalf("first pass: expected 4 results, got %d", len(first.Results))
	}

	// Second pass on already-trimmed results should be stable.
	second := quota.Apply(first.Results, opts)
	if len(second.Results) != len(first.Results) {
		t.Fatalf("second pass not idempotent: got %d, want %d", len(second.Results), len(first.Results))
	}
	if len(second.Violations) != 0 {
		t.Fatalf("second pass should have no violations, got %d", len(second.Violations))
	}
}
