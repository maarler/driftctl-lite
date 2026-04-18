package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/snyk/driftctl-lite/internal/drift"
	"github.com/snyk/driftctl-lite/internal/metrics"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ResourceID: "r1", Status: drift.StatusInSync},
		{ResourceID: "r2", Status: drift.StatusMissing},
		{ResourceID: "r3", Status: drift.StatusExtra},
		{ResourceID: "r4", Status: drift.StatusModified},
		{ResourceID: "r5", Status: drift.StatusInSync},
	}
}

func TestCollect_Counts(t *testing.T) {
	results := makeResults()
	s := metrics.Collect(results, 42*time.Millisecond)

	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
	if s.InSync != 2 {
		t.Errorf("expected InSync=2, got %d", s.InSync)
	}
	if s.Missing != 1 {
		t.Errorf("expected Missing=1, got %d", s.Missing)
	}
	if s.Extra != 1 {
		t.Errorf("expected Extra=1, got %d", s.Extra)
	}
	if s.Modified != 1 {
		t.Errorf("expected Modified=1, got %d", s.Modified)
	}
	if s.Duration != 42*time.Millisecond {
		t.Errorf("expected Duration=42ms, got %s", s.Duration)
	}
}

func TestCollect_Empty(t *testing.T) {
	s := metrics.Collect(nil, 0)
	if s.Total != 0 || s.InSync != 0 {
		t.Error("expected all zeros for empty results")
	}
}

func TestPrint_Output(t *testing.T) {
	results := makeResults()
	s := metrics.Collect(results, 10*time.Millisecond)

	var buf bytes.Buffer
	metrics.Print(s, &buf)
	out := buf.String()

	for _, want := range []string{"Total:", "In sync:", "Missing:", "Extra:", "Modified:", "Scan completed"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\ngot:\n%s", want, out)
		}
	}
}
