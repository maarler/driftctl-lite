package health_test

import (
	"bytes"
	"strings"
	"testing"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/health"
)

func makeResults(statuses ...drift.ResourceStatus) []drift.Result {
	var out []drift.Result
	for i, s := range statuses {
		out = append(out, drift.Result{
			ResourceID:   fmt.Sprintf("res-%d", i),
			ResourceType: "aws_instance",
			Status:       s,
		})
	}
	return out
}

import "fmt"

func TestEvaluate_AllInSync(t *testing.T) {
	results := makeResults(drift.StatusInSync, drift.StatusInSync)
	r := health.Evaluate(results)
	if r.Status != health.StatusHealthy {
		t.Fatalf("expected HEALTHY, got %s", r.Status)
	}
	if r.Drifted != 0 {
		t.Fatalf("expected 0 drifted, got %d", r.Drifted)
	}
}

func TestEvaluate_WithDrift(t *testing.T) {
	results := makeResults(drift.StatusInSync, drift.StatusMissing, drift.StatusModified)
	r := health.Evaluate(results)
	if r.Status != health.StatusDegraded {
		t.Fatalf("expected DEGRADED, got %s", r.Status)
	}
	if r.Drifted != 2 {
		t.Fatalf("expected 2 drifted, got %d", r.Drifted)
	}
	if r.Total != 3 {
		t.Fatalf("expected total 3, got %d", r.Total)
	}
}

func TestEvaluate_Empty(t *testing.T) {
	r := health.Evaluate(nil)
	if r.Status != health.StatusHealthy {
		t.Fatalf("expected HEALTHY for empty input, got %s", r.Status)
	}
}

func TestPrint_Output(t *testing.T) {
	r := health.Report{Status: health.StatusDegraded, Total: 5, Drifted: 2}
	var buf bytes.Buffer
	health.Print(r, &buf)
	out := buf.String()
	if !strings.Contains(out, "DEGRADED") {
		t.Errorf("expected DEGRADED in output, got: %s", out)
	}
	if !strings.Contains(out, "drifted=2") {
		t.Errorf("expected drifted=2 in output, got: %s", out)
	}
}
