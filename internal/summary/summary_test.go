package summary

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-lite/internal/drift"
)

func makeResults(statuses []drift.ResourceStatus) []drift.Result {
	var results []drift.Result
	for i, s := range statuses {
		results = append(results, drift.Result{
			ResourceID:   fmt.Sprintf("res-%d", i),
			ResourceType: "aws_instance",
			Status:       s,
		})
	}
	return results
}

func TestCompute_AllInSync(t *testing.T) {
	results := makeResults([]drift.ResourceStatus{
		drift.StatusInSync, drift.StatusInSync,
	})
	s := Compute(results)
	if s.Total != 2 || s.InSync != 2 || s.Missing != 0 {
		t.Errorf("unexpected stats: %+v", s)
	}
}

func TestCompute_Mixed(t *testing.T) {
	results := makeResults([]drift.ResourceStatus{
		drift.StatusInSync, drift.StatusMissing, drift.StatusExtra, drift.StatusModified,
	})
	s := Compute(results)
	if s.Total != 4 || s.InSync != 1 || s.Missing != 1 || s.Extra != 1 || s.Modified != 1 {
		t.Errorf("unexpected stats: %+v", s)
	}
}

func TestPrint_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	Print(&buf, Stats{Total: 2, InSync: 2})
	if !strings.Contains(buf.String(), "No drift detected.") {
		t.Errorf("expected no drift message, got: %s", buf.String())
	}
}

func TestPrint_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	Print(&buf, Stats{Total: 3, InSync: 1, Missing: 1, Modified: 1})
	if !strings.Contains(buf.String(), "Drift detected!") {
		t.Errorf("expected drift message, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "Missing:  1") {
		t.Errorf("expected missing count, got: %s", buf.String())
	}
}
