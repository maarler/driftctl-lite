package scorecard

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftctl-lite/internal/drift"
)

func makeResults(statuses []drift.Status) []drift.Result {
	results := make([]drift.Result, len(statuses))
	for i, s := range statuses {
		results[i] = drift.Result{ResourceID: fmt.Sprintf("res-%d", i), Status: s}
	}
	return results
}

func TestCompute_Empty(t *testing.T) {
	s := Compute(nil)
	if s.Total != 0 || s.Percent != 100.0 {
		t.Errorf("expected perfect score on empty, got %+v", s)
	}
}

func TestCompute_AllInSync(t *testing.T) {
	results := []drift.Result{
		{ResourceID: "a", Status: drift.StatusInSync},
		{ResourceID: "b", Status: drift.StatusInSync},
	}
	s := Compute(results)
	if s.Drifted != 0 || s.Percent != 100.0 {
		t.Errorf("unexpected score: %+v", s)
	}
}

func TestCompute_Mixed(t *testing.T) {
	results := []drift.Result{
		{ResourceID: "a", Status: drift.StatusInSync},
		{ResourceID: "b", Status: drift.StatusMissing},
		{ResourceID: "c", Status: drift.StatusModified},
		{ResourceID: "d", Status: drift.StatusExtra},
	}
	s := Compute(results)
	if s.Total != 4 || s.Drifted != 3 || s.InSync != 1 {
		t.Errorf("unexpected counts: %+v", s)
	}
	if s.Percent != 25.0 {
		t.Errorf("expected 25.0%%, got %.1f", s.Percent)
	}
}

func TestPrint_Output(t *testing.T) {
	s := Score{Total: 4, InSync: 3, Drifted: 1, Percent: 75.0}
	var buf bytes.Buffer
	Print(s, &buf)
	out := buf.String()
	if !strings.Contains(out, "75.0%") {
		t.Errorf("expected percent in output, got: %s", out)
	}
	if !strings.Contains(out, "Drift Scorecard") {
		t.Errorf("expected header in output, got: %s", out)
	}
}
