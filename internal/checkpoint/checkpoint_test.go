package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"driftctl-lite/internal/checkpoint"
	"driftctl-lite/internal/drift"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ResourceID: "vpc-1", ResourceType: "aws_vpc", Status: drift.StatusInSync},
		{ResourceID: "sg-2", ResourceType: "aws_security_group", Status: drift.StatusMissing},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	results := makeResults()

	if err := checkpoint.Save(dir, "my-check", results); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entry, err := checkpoint.Load(dir, "my-check")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entry.Name != "my-check" {
		t.Errorf("name: got %q, want %q", entry.Name, "my-check")
	}
	if len(entry.Results) != len(results) {
		t.Errorf("results len: got %d, want %d", len(entry.Results), len(results))
	}
	if entry.SavedAt.IsZero() {
		t.Error("SavedAt should not be zero")
	}
	if entry.SavedAt.After(time.Now().Add(time.Second)) {
		t.Error("SavedAt is in the future")
	}
}

func TestLoad_Missing(t *testing.T) {
	dir := t.TempDir()
	_, err := checkpoint.Load(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing checkpoint")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	err := checkpoint.Save(filepath.Join("\x00invalid"), "x", nil)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestList_Empty(t *testing.T) {
	dir := t.TempDir()
	names, err := checkpoint.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestList_ReturnsSavedCheckpoints(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"alpha", "beta", "gamma"} {
		if err := checkpoint.Save(dir, name, makeResults()); err != nil {
			t.Fatalf("Save %s: %v", name, err)
		}
	}
	names, err := checkpoint.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 checkpoints, got %d: %v", len(names), names)
	}
}

func TestList_MissingDir(t *testing.T) {
	names, err := checkpoint.List(filepath.Join(os.TempDir(), "no-such-dir-xyz"))
	if err != nil {
		t.Fatalf("List on missing dir should not error: %v", err)
	}
	if names != nil {
		t.Errorf("expected nil, got %v", names)
	}
}
