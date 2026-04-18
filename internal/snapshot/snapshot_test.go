package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/snapshot"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ResourceID: "vpc-1", ResourceType: "aws_vpc", Status: drift.StatusInSync},
		{ResourceID: "sg-2", ResourceType: "aws_security_group", Status: drift.StatusMissing},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	results := makeResults()
	if err := snapshot.Save(path, results); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(snap.Results) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(snap.Results))
	}
	if snap.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
	if time.Since(snap.CapturedAt) > 5*time.Second {
		t.Error("CapturedAt seems too old")
	}
	for i, r := range snap.Results {
		if r.ResourceID != results[i].ResourceID {
			t.Errorf("result[%d] ID: want %s got %s", i, results[i].ResourceID, r.ResourceID)
		}
		if r.Status != results[i].Status {
			t.Errorf("result[%d] Status: want %s got %s", i, results[i].Status, r.Status)
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	err := snapshot.Save("/nonexistent_dir/snap.json", makeResults())
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
