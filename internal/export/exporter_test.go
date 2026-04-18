package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/owner/driftctl-lite/internal/drift"
	"github.com/owner/driftctl-lite/internal/export"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ResourceID: "vpc-1", ResourceType: "aws_vpc", Status: drift.StatusMissing},
		{ResourceID: "sg-2", ResourceType: "aws_sg", Status: drift.StatusModified},
		{ResourceID: "s3-3", ResourceType: "aws_s3", Status: drift.StatusInSync},
	}
}

func TestExport_JSON(t *testing.T) {
	var buf bytes.Buffer
	err := export.Export(makeResults(), export.FormatJSON, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []export.Record
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}
	if records[0].ResourceID != "vpc-1" {
		t.Errorf("expected vpc-1, got %s", records[0].ResourceID)
	}
}

func TestExport_CSV(t *testing.T) {
	var buf bytes.Buffer
	err := export.Export(makeResults(), export.FormatCSV, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 4 { // header + 3 rows
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "timestamp") {
		t.Errorf("expected CSV header, got: %s", lines[0])
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := export.Export(makeResults(), export.Format("xml"), &buf)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExport_Empty(t *testing.T) {
	var buf bytes.Buffer
	err := export.Export([]drift.Result{}, export.FormatJSON, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []export.Record
	_ = json.Unmarshal(buf.Bytes(), &records)
	if len(records) != 0 {
		t.Errorf("expected empty records")
	}
}
