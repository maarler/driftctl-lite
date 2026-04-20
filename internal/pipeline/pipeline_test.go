package pipeline_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/pipeline"
)

func makeResults(ids ...string) []drift.Result {
	var out []drift.Result
	for _, id := range ids {
		out = append(out, drift.Result{
			ResourceID:   id,
			ResourceType: "aws_instance",
			Status:       drift.StatusInSync,
		})
	}
	return out
}

func TestPipeline_NoStages_ReturnsInput(t *testing.T) {
	p := pipeline.New()
	input := makeResults("a", "b")
	out := p.Run(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestPipeline_NilInput_ReturnsEmpty(t *testing.T) {
	p := pipeline.New()
	out := p.Run(nil)
	if out == nil {
		t.Fatal("expected non-nil slice")
	}
	if len(out) != 0 {
		t.Fatalf("expected 0 results, got %d", len(out))
	}
}

func TestPipeline_SingleStage_Transforms(t *testing.T) {
	dropAll := func(r []drift.Result) []drift.Result { return []drift.Result{} }
	p := pipeline.New().Add(dropAll)
	out := p.Run(makeResults("a", "b", "c"))
	if len(out) != 0 {
		t.Fatalf("expected 0 results after drop, got %d", len(out))
	}
}

func TestPipeline_MultipleStages_ChainCorrectly(t *testing.T) {
	// Stage 1: keep only first result
	keeFirst := func(r []drift.Result) []drift.Result {
		if len(r) == 0 {
			return r
		}
		return r[:1]
	}
	// Stage 2: mark remaining as missing
	markMissing := func(r []drift.Result) []drift.Result {
		for i := range r {
			r[i].Status = drift.StatusMissing
		}
		return r
	}
	p := pipeline.New().Add(keepFirst).Add(markMissing)
	out := p.Run(makeResults("x", "y", "z"))
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Status != drift.StatusMissing {
		t.Fatalf("expected StatusMissing, got %v", out[0].Status)
	}
}

func TestPipeline_Len(t *testing.T) {
	p := pipeline.New()
	if p.Len() != 0 {
		t.Fatalf("expected 0 stages, got %d", p.Len())
	}
	noop := func(r []drift.Result) []drift.Result { return r }
	p.Add(noop).Add(noop)
	if p.Len() != 2 {
		t.Fatalf("expected 2 stages, got %d", p.Len())
	}
}
