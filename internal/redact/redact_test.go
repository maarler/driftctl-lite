package redact_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/redact"
)

func makeResult(id, rtype string, status drift.Status, diffs []drift.Diff) drift.Result {
	return drift.Result{ResourceID: id, ResourceType: rtype, Status: status, Diffs: diffs}
}

func TestApply_NoSensitiveFields(t *testing.T) {
	input := []drift.Result{
		makeResult("vpc-1", "aws_vpc", drift.Modified, []drift.Diff{
			{Field: "cidr", Expected: "10.0.0.0/16", Got: "10.1.0.0/16"},
		}),
	}
	out := redact.Apply(input, nil)
	if out[0].Diffs[0].Expected == "***REDACTED***" {
		t.Fatal("non-sensitive field should not be redacted")
	}
}

func TestApply_SensitiveFieldsRedacted(t *testing.T) {
	input := []drift.Result{
		makeResult("db-1", "aws_db", drift.Modified, []drift.Diff{
			{Field: "password", Expected: "hunter2", Got: "letmein"},
			{Field: "name", Expected: "prod", Got: "staging"},
		}),
	}
	out := redact.Apply(input, nil)
	if out[0].Diffs[0].Expected != "***REDACTED***" {
		t.Errorf("expected password to be redacted, got %q", out[0].Diffs[0].Expected)
	}
	if out[0].Diffs[0].Got != "***REDACTED***" {
		t.Errorf("expected password got to be redacted, got %q", out[0].Diffs[0].Got)
	}
	if out[0].Diffs[1].Expected == "***REDACTED***" {
		t.Error("non-sensitive field 'name' should not be redacted")
	}
}

func TestApply_CustomKeys(t *testing.T) {
	input := []drift.Result{
		makeResult("s3-1", "aws_s3", drift.Modified, []drift.Diff{
			{Field: "bucket_policy", Expected: "old", Got: "new"},
		}),
	}
	out := redact.Apply(input, []string{"bucket_policy"})
	if out[0].Diffs[0].Expected != "***REDACTED***" {
		t.Error("custom sensitive key should be redacted")
	}
}

func TestApply_NoDiffs(t *testing.T) {
	input := []drift.Result{
		makeResult("ec2-1", "aws_instance", drift.InSync, nil),
	}
	out := redact.Apply(input, nil)
	if out[0].Diffs != nil {
		t.Error("expected nil diffs to remain nil")
	}
}

func TestApply_PreservesOriginal(t *testing.T) {
	orig := []drift.Result{
		makeResult("sg-1", "aws_sg", drift.Modified, []drift.Diff{
			{Field: "token", Expected: "abc", Got: "xyz"},
		}),
	}
	redact.Apply(orig, nil)
	if orig[0].Diffs[0].Expected == "***REDACTED***" {
		t.Error("Apply should not mutate original results")
	}
}
