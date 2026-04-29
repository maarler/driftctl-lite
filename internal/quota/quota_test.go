package quota_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/quota"
)

func makeResult(id, typ, status string) drift.Result {
	return drift.Result{
		Resource: drift.Resource{ID: id, Type: typ},
		Status:   status,
	}
}

func TestApply_NoLimit_ReturnsAll(t *testing.T) {
	results := []drift.Result{
		makeResult("a", "aws_s3_bucket", "missing"),
		makeResult("b", "aws_s3_bucket", "missing"),
		makeResult("c", "aws_s3_bucket", "missing"),
	}
	opts := quota.DefaultOptions()
	report := quota.Apply(results, opts)
	if len(report.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(report.Results))
	}
	if len(report.Violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(report.Violations))
	}
}

func TestApply_WithinQuota_NoViolations(t *testing.T) {
	results := []drift.Result{
		makeResult("a", "aws_s3_bucket", "missing"),
		makeResult("b", "aws_s3_bucket", "missing"),
	}
	opts := quota.Options{MaxPerType: 3}
	report := quota.Apply(results, opts)
	if len(report.Violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(report.Violations))
	}
	if len(report.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(report.Results))
	}
}

func TestApply_ExceedsQuota_FlagsViolation(t *testing.T) {
	results := []drift.Result{
		makeResult("a", "aws_s3_bucket", "missing"),
		makeResult("b", "aws_s3_bucket", "missing"),
		makeResult("c", "aws_s3_bucket", "extra"),
	}
	opts := quota.Options{MaxPerType: 2, DropExceeding: false}
	report := quota.Apply(results, opts)
	if len(report.Results) != 3 {
		t.Fatalf("expected 3 results kept (flag only), got %d", len(report.Results))
	}
	if len(report.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(report.Violations))
	}
	v := report.Violations[0]
	if v.Type != "aws_s3_bucket" {
		t.Errorf("expected type aws_s3_bucket, got %s", v.Type)
	}
	if v.Allowed != 2 {
		t.Errorf("expected allowed=2, got %d", v.Allowed)
	}
	if v.Actual != 3 {
		t.Errorf("expected actual=3, got %d", v.Actual)
	}
}

func TestApply_DropExceeding_RemovesExtra(t *testing.T) {
	results := []drift.Result{
		makeResult("a", "aws_s3_bucket", "missing"),
		makeResult("b", "aws_s3_bucket", "missing"),
		makeResult("c", "aws_s3_bucket", "extra"),
		makeResult("d", "aws_iam_role", "modified"),
	}
	opts := quota.Options{MaxPerType: 2, DropExceeding: true}
	report := quota.Apply(results, opts)
	if len(report.Results) != 3 {
		t.Fatalf("expected 3 results (2 s3 + 1 iam), got %d", len(report.Results))
	}
	if len(report.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(report.Violations))
	}
}

func TestApply_InSyncResultsAlwaysKept(t *testing.T) {
	results := []drift.Result{
		makeResult("x", "aws_s3_bucket", "in_sync"),
		makeResult("y", "aws_s3_bucket", "in_sync"),
		makeResult("z", "aws_s3_bucket", "in_sync"),
	}
	opts := quota.Options{MaxPerType: 1, DropExceeding: true}
	report := quota.Apply(results, opts)
	if len(report.Results) != 3 {
		t.Fatalf("expected all 3 in-sync results kept, got %d", len(report.Results))
	}
	if len(report.Violations) != 0 {
		t.Fatalf("expected no violations for in-sync results, got %d", len(report.Violations))
	}
}
