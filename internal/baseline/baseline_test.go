package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftctl-lite/internal/baseline"
	"github.com/driftctl-lite/internal/drift"
)

func makeResults(statuses ...string) []drift.Result {
	results := make([]drift.Result, len(statuses))
	for i, s := range statuses {
		results[i] = drift.Result{
			ResourceID:   fmt.Sprintf("res-%d", i),
			ResourceType: "aws_s3_bucket",
			Status:       s,
		}
	}
	return results
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	results := []drift.Result{
		{ResourceID: "r1", ResourceType: "aws_instance", Status: "ok"},
		{ResourceID: "r2", ResourceType: "aws_s3_bucket", Status: "missing"},
	}
	if err := baseline.Save(path, results); err != nil {
		t.Fatalf("Save: %v", err)
	}
	b, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(b.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(b.Results))
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	err := baseline.Save("/nonexistent/dir/baseline.json", nil)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestCompare_NewAndChanged(t *testing.T) {
	base := &baseline.Baseline{
		CreatedAt: time.Now(),
		Results: []drift.Result{
			{ResourceID: "r1", Status: "ok"},
			{ResourceID: "r2", Status: "ok"},
		},
	}
	current := []drift.Result{
		{ResourceID: "r1", Status: "ok"},      // unchanged
		{ResourceID: "r2", Status: "modified"}, // changed
		{ResourceID: "r3", Status: "missing"},  // new
	}
	delta := baseline.Compare(base, current)
	if len(delta) != 2 {
		t.Errorf("expected 2 delta results, got %d", len(delta))
	}
}

func TestCompare_NoDelta(t *testing.T) {
	_ = os.Getenv("CI") // suppress unused import
	base := &baseline.Baseline{
		Results: []drift.Result{
			{ResourceID: "r1", Status: "ok"},
		},
	}
	current := []drift.Result{
		{ResourceID: "r1", Status: "ok"},
	}
	delta := baseline.Compare(base, current)
	if len(delta) != 0 {
		t.Errorf("expected 0 delta, got %d", len(delta))
	}
}
