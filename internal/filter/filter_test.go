package filter_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/filter"
)

func makeResult() drift.Result {
	return drift.Result{
		Missing: []drift.ResourceDiff{
			{Type: "aws_s3_bucket", ID: "bucket-1"},
		},
		Extra: []drift.ResourceDiff{
			{Type: "aws_instance", ID: "i-123"},
		},
		Modified: []drift.ResourceDiff{
			{Type: "aws_s3_bucket", ID: "bucket-2"},
			{Type: "aws_instance", ID: "i-456"},
		},
	}
}

func TestFilter_NoOptions(t *testing.T) {
	r := filter.Apply(makeResult(), filter.Options{})
	if len(r.Missing) != 1 || len(r.Extra) != 1 || len(r.Modified) != 2 {
		t.Errorf("expected all resources, got missing=%d extra=%d modified=%d", len(r.Missing), len(r.Extra), len(r.Modified))
	}
}

func TestFilter_ByType(t *testing.T) {
	r := filter.Apply(makeResult(), filter.Options{Types: []string{"aws_s3_bucket"}})
	if len(r.Missing) != 1 {
		t.Errorf("expected 1 missing, got %d", len(r.Missing))
	}
	if len(r.Extra) != 0 {
		t.Errorf("expected 0 extra, got %d", len(r.Extra))
	}
	if len(r.Modified) != 1 {
		t.Errorf("expected 1 modified, got %d", len(r.Modified))
	}
}

func TestFilter_OnlyDrift(t *testing.T) {
	// OnlyDrift=true should still return all drift entries when hasDrift is true
	r := filter.Apply(makeResult(), filter.Options{OnlyDrift: true})
	if len(r.Missing) != 1 || len(r.Extra) != 1 || len(r.Modified) != 2 {
		t.Errorf("unexpected counts: missing=%d extra=%d modified=%d", len(r.Missing), len(r.Extra), len(r.Modified))
	}
}

func TestFilter_ByTypeAndOnlyDrift(t *testing.T) {
	r := filter.Apply(makeResult(), filter.Options{Types: []string{"aws_instance"}, OnlyDrift: true})
	if len(r.Missing) != 0 {
		t.Errorf("expected 0 missing, got %d", len(r.Missing))
	}
	if len(r.Extra) != 1 {
		t.Errorf("expected 1 extra, got %d", len(r.Extra))
	}
	if len(r.Modified) != 1 {
		t.Errorf("expected 1 modified, got %d", len(r.Modified))
	}
}
