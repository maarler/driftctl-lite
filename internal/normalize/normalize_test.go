package normalize_test

import (
	"testing"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/normalize"
)

func makeResult(id, rtype, status string, diffs map[string]drift.Diff) drift.Result {
	return drift.Result{
		ResourceID:   id,
		ResourceType: rtype,
		Status:       status,
		Diffs:        diffs,
	}
}

func TestApply_DefaultOptions_LowercasesType(t *testing.T) {
	input := []drift.Result{makeResult("id1", "AWS_S3_BUCKET", "modified", nil)}
	out := normalize.Apply(input, normalize.DefaultOptions())
	if out[0].ResourceType != "aws_s3_bucket" {
		t.Errorf("expected lowercase type, got %q", out[0].ResourceType)
	}
}

func TestApply_DefaultOptions_TrimsID(t *testing.T) {
	input := []drift.Result{makeResult("  my-bucket  ", "s3", "in_sync", nil)}
	out := normalize.Apply(input, normalize.DefaultOptions())
	if out[0].ResourceID != "my-bucket" {
		t.Errorf("expected trimmed ID, got %q", out[0].ResourceID)
	}
}

func TestApply_DefaultOptions_TrimsDiffValues(t *testing.T) {
	diffs := map[string]drift.Diff{
		"region": {Expected: " us-east-1 ", Actual: "  eu-west-1  "},
	}
	input := []drift.Result{makeResult("r1", "ec2", "modified", diffs)}
	out := normalize.Apply(input, normalize.DefaultOptions())
	d := out[0].Diffs["region"]
	if d.Expected != "us-east-1" {
		t.Errorf("expected trimmed Expected, got %q", d.Expected)
	}
	if d.Actual != "eu-west-1" {
		t.Errorf("expected trimmed Actual, got %q", d.Actual)
	}
}

func TestApply_NoOptions_LeavesFieldsUnchanged(t *testing.T) {
	opts := normalize.Options{}
	input := []drift.Result{makeResult("  ID  ", "AWS_EC2", "missing", nil)}
	out := normalize.Apply(input, opts)
	if out[0].ResourceID != "  ID  " {
		t.Errorf("expected untouched ID, got %q", out[0].ResourceID)
	}
	if out[0].ResourceType != "AWS_EC2" {
		t.Errorf("expected untouched type, got %q", out[0].ResourceType)
	}
}

func TestApply_Empty_ReturnsEmpty(t *testing.T) {
	out := normalize.Apply(nil, normalize.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d items", len(out))
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	original := []drift.Result{makeResult("  orig  ", "TYPE", "extra", nil)}
	normalize.Apply(original, normalize.DefaultOptions())
	if original[0].ResourceID != "  orig  " {
		t.Errorf("original result was mutated")
	}
}
