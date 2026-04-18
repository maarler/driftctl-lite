package ignore

import (
	"os"
	"testing"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "driftignore-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadFromFile_Missing(t *testing.T) {
	l, err := LoadFromFile("/nonexistent/path/.driftignore")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(l.rules) != 0 {
		t.Errorf("expected empty rules, got %d", len(l.rules))
	}
}

func TestLoadFromFile_ParsesRules(t *testing.T) {
	path := writeTempIgnore(t, "# comment\naws_s3_bucket/my-bucket\naws_instance/*\n")
	l, err := LoadFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(l.rules))
	}
}

func TestMatches_Exact(t *testing.T) {
	path := writeTempIgnore(t, "aws_s3_bucket/my-bucket\n")
	l, _ := LoadFromFile(path)
	if !l.Matches("aws_s3_bucket", "my-bucket") {
		t.Error("expected match")
	}
	if l.Matches("aws_s3_bucket", "other-bucket") {
		t.Error("unexpected match")
	}
}

func TestMatches_Wildcard(t *testing.T) {
	path := writeTempIgnore(t, "aws_instance/*\n")
	l, _ := LoadFromFile(path)
	if !l.Matches("aws_instance", "i-12345") {
		t.Error("expected wildcard match")
	}
}

func TestFilterIgnored(t *testing.T) {
	path := writeTempIgnore(t, "aws_s3_bucket/ignored-bucket\n")
	l, _ := LoadFromFile(path)
	results := []DriftResult{
		{ResourceType: "aws_s3_bucket", ResourceID: "ignored-bucket", Status: "modified"},
		{ResourceType: "aws_s3_bucket", ResourceID: "kept-bucket", Status: "modified"},
	}
	out := l.FilterIgnored(results)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].ResourceID != "kept-bucket" {
		t.Errorf("unexpected resource: %s", out[0].ResourceID)
	}
}
