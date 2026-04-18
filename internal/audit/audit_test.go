package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/snyk/driftctl-lite/internal/audit"
	"github.com/snyk/driftctl-lite/internal/drift"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ResourceID: "vpc-1", ResourceType: "aws_vpc", Status: drift.StatusInSync},
		{ResourceID: "sg-1", ResourceType: "aws_sg", Status: drift.StatusMissing},
		{ResourceID: "s3-1", ResourceType: "aws_s3", Status: drift.StatusModified},
	}
}

func tempLogPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.log")
}

func TestRecord_CreatesFile(t *testing.T) {
	path := tempLogPath(t)
	l := audit.NewLogger(path)
	if err := l.Record("state.json", "file", makeResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("log file not created: %v", err)
	}
}

func TestRecord_And_ReadAll(t *testing.T) {
	path := tempLogPath(t)
	l := audit.NewLogger(path)

	results := makeResults()
	if err := l.Record("state.json", "file", results); err != nil {
		t.Fatalf("record: %v", err)
	}
	if err := l.Record("state2.json", "file", results[:1]); err != nil {
		t.Fatalf("record2: %v", err)
	}

	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	e := entries[0]
	if e.StateFile != "state.json" {
		t.Errorf("state file mismatch: %s", e.StateFile)
	}
	if e.TotalItems != 3 {
		t.Errorf("expected 3 total, got %d", e.TotalItems)
	}
	if e.DriftCount != 2 {
		t.Errorf("expected 2 drifted, got %d", e.DriftCount)
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
	if e.Timestamp.After(time.Now().Add(time.Second)) {
		t.Error("timestamp is in the future")
	}
}

func TestReadAll_MissingFile(t *testing.T) {
	entries, err := audit.ReadAll("/nonexistent/path/audit.log")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}
