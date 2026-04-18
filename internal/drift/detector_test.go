package drift_test

import (
	"testing"

	"github.com/example/driftctl-lite/internal/drift"
	"github.com/example/driftctl-lite/internal/state"
)

func makeState(resources ...state.Resource) *state.State {
	return &state.State{Resources: resources}
}

func TestDetect_NoDrift(t *testing.T) {
	r := state.Resource{ID: "vpc-1", Type: "aws_vpc", Attributes: map[string]string{"cidr": "10.0.0.0/16"}}
	results := drift.Detect(makeState(r), []state.Resource{r})
	if len(results) != 0 {
		t.Errorf("expected no drift, got %d results", len(results))
	}
}

func TestDetect_Missing(t *testing.T) {
	r := state.Resource{ID: "vpc-1", Type: "aws_vpc", Attributes: map[string]string{}}
	results := drift.Detect(makeState(r), []state.Resource{})
	if len(results) != 1 || results[0].Type != drift.DriftMissing {
		t.Errorf("expected MISSING drift, got %+v", results)
	}
}

func TestDetect_Extra(t *testing.T) {
	live := state.Resource{ID: "sg-99", Type: "aws_security_group", Attributes: map[string]string{}}
	results := drift.Detect(makeState(), []state.Resource{live})
	if len(results) != 1 || results[0].Type != drift.DriftExtra {
		t.Errorf("expected EXTRA drift, got %+v", results)
	}
}

func TestDetect_Modified(t *testing.T) {
	declared := state.Resource{ID: "vpc-1", Type: "aws_vpc", Attributes: map[string]string{"cidr": "10.0.0.0/16"}}
	live := state.Resource{ID: "vpc-1", Type: "aws_vpc", Attributes: map[string]string{"cidr": "192.168.0.0/16"}}
	results := drift.Detect(makeState(declared), []state.Resource{live})
	if len(results) != 1 || results[0].Type != drift.DriftModified {
		t.Errorf("expected MODIFIED drift, got %+v", results)
	}
}
