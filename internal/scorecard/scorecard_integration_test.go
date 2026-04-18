package scorecard_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/scorecard"
)

func TestScorecard_RoundTrip(t *testing.T) {
	results := []drift.Result{
		{ResourceID: "vpc-1", ResourceType: "aws_vpc", Status: drift.StatusInSync},
		{ResourceID: "sg-1", ResourceType: "aws_sg", Status: drift.StatusMissing},
		{ResourceID: "s3-1", ResourceType: "aws_s3", Status: drift.StatusModified},
		{ResourceID: "ec2-1", ResourceType: "aws_ec2", Status: drift.StatusInSync},
	}

	s := scorecard.Compute(results)

	if s.Total != 4 {
		t.Fatalf("expected total 4, got %d", s.Total)
	}
	if s.InSync != 2 {
		t.Fatalf("expected inSync 2, got %d", s.InSync)
	}
	if s.Drifted != 2 {
		t.Fatalf("expected drifted 2, got %d", s.Drifted)
	}
	if s.Percent != 50.0 {
		t.Fatalf("expected 50.0%%, got %.1f", s.Percent)
	}

	var buf bytes.Buffer
	scorecard.Print(s, &buf)
	out := buf.String()

	for _, want := range []string{"50.0%", "Total", "Drifted", "In sync"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}
