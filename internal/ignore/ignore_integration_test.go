package ignore_test

import (
	"os"
	"testing"

	"github.com/driftctl-lite/internal/ignore"
)

func TestIgnore_RoundTrip(t *testing.T) {
	f, err := os.CreateTemp("", "driftignore-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	_, _ = f.WriteString("aws_lambda_function/*\naws_s3_bucket/prod-bucket\n")
	f.Close()

	l, err := ignore.LoadFromFile(f.Name())
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	results := []ignore.DriftResult{
		{ResourceType: "aws_lambda_function", ResourceID: "fn-a", Status: "missing"},
		{ResourceType: "aws_lambda_function", ResourceID: "fn-b", Status: "extra"},
		{ResourceType: "aws_s3_bucket", ResourceID: "prod-bucket", Status: "modified"},
		{ResourceType: "aws_s3_bucket", ResourceID: "dev-bucket", Status: "modified"},
	}

	out := l.FilterIgnored(results)
	if len(out) != 1 {
		t.Fatalf("expected 1 result after filtering, got %d", len(out))
	}
	if out[0].ResourceID != "dev-bucket" {
		t.Errorf("expected dev-bucket, got %s", out[0].ResourceID)
	}
}
