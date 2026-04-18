package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/output"
)

func makeResult(missing, extra, modified []drift.Resource) drift.Result {
	return drift.Result{
		Missing:  missing,
		Extra:    extra,
		Modified: modified,
	}
}

func TestReporter_NoDrift_Text(t *testing.T) {
	var buf bytes.Buffer
	r := &output.Reporter{Writer: &buf, Format: output.FormatText}
	result := makeResult(nil, nil, nil)
	if err := r.Report(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestReporter_Drift_Text(t *testing.T) {
	var buf bytes.Buffer
	r := &output.Reporter{Writer: &buf, Format: output.FormatText}
	result := makeResult(
		[]drift.Resource{{ID: "vpc-1", Type: "aws_vpc"}},
		[]drift.Resource{{ID: "sg-99", Type: "aws_sg"}},
		nil,
	)
	if err := r.Report(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING in output")
	}
	if !strings.Contains(out, "EXTRA") {
		t.Errorf("expected EXTRA in output")
	}
}

func TestReporter_Drift_JSON(t *testing.T) {
	var buf bytes.Buffer
	r := &output.Reporter{Writer: &buf, Format: output.FormatJSON}
	result := makeResult(
		nil,
		nil,
		[]drift.Resource{{ID: "ec2-1", Type: "aws_instance"}},
	)
	if err := r.Report(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "modified") {
		t.Errorf("expected 'modified' key in JSON output, got: %s", buf.String())
	}
}
