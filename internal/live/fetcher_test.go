package live

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/driftctl-lite/internal/state"
)

func writeTempLive(t *testing.T, resources []state.Resource) string {
	t.Helper()
	data, err := json.Marshal(resources)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	f, err := os.CreateTemp("", "live-*.json")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.Write(data); err != nil {
		t.Fatalf("write: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestFetcher_File(t *testing.T) {
	resources := []state.Resource{
		{ID: "vpc-1", Type: "aws_vpc", Attributes: map[string]string{"cidr": "10.0.0.0/16"}},
		{ID: "sg-1", Type: "aws_security_group", Attributes: map[string]string{"name": "default"}},
	}
	path := writeTempLive(t, resources)

	f := NewFetcher(SourceFile, path)
	rm, err := f.Fetch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rm) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(rm))
	}
	if rm["vpc-1"].Type != "aws_vpc" {
		t.Errorf("expected aws_vpc, got %s", rm["vpc-1"].Type)
	}
}

func TestFetcher_MissingFile(t *testing.T) {
	f := NewFetcher(SourceFile, "/nonexistent/live.json")
	_, err := f.Fetch()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFetcher_UnsupportedSource(t *testing.T) {
	f := NewFetcher("s3", "bucket/path")
	_, err := f.Fetch()
	if err == nil {
		t.Fatal("expected error for unsupported source")
	}
}
