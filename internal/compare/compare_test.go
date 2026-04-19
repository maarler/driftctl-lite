package compare_test

import (
	"testing"

	"github.com/driftctl-lite/internal/compare"
	"github.com/driftctl-lite/internal/drift"
)

func makeResult(id, rtype string, status drift.DriftStatus) drift.Result {
	return drift.Result{
		ResourceID:   id,
		ResourceType: rtype,
		Status:       status,
	}
}

func TestCompare_AllResolved(t *testing.T) {
	prev := []drift.Result{makeResult("r1", "aws_s3_bucket", drift.StatusMissing)}
	curr := []drift.Result{makeResult("r1", "aws_s3_bucket", drift.StatusInSync)}
	d := compare.Compare(prev, curr)
	if len(d.Resolved) != 1 || len(d.New) != 0 || len(d.Persisted) != 0 {
		t.Fatalf("expected 1 resolved, got %+v", d)
	}
}

func TestCompare_NewDrift(t *testing.T) {
	prev := []drift.Result{makeResult("r1", "aws_s3_bucket", drift.StatusInSync)}
	curr := []drift.Result{
		makeResult("r1", "aws_s3_bucket", drift.StatusInSync),
		makeResult("r2", "aws_instance", drift.StatusModified),
	}
	d := compare.Compare(prev, curr)
	if len(d.New) != 1 || d.New[0].ResourceID != "r2" {
		t.Fatalf("expected 1 new drift, got %+v", d)
	}
}

func TestCompare_Persisted(t *testing.T) {
	prev := []drift.Result{makeResult("r1", "aws_s3_bucket", drift.StatusMissing)}
	curr := []drift.Result{makeResult("r1", "aws_s3_bucket", drift.StatusMissing)}
	d := compare.Compare(prev, curr)
	if len(d.Persisted) != 1 {
		t.Fatalf("expected 1 persisted, got %+v", d)
	}
}

func TestCompare_Mixed(t *testing.T) {
	prev := []drift.Result{
		makeResult("r1", "aws_s3_bucket", drift.StatusMissing),
		makeResult("r2", "aws_instance", drift.StatusModified),
	}
	curr := []drift.Result{
		makeResult("r1", "aws_s3_bucket", drift.StatusInSync),
		makeResult("r2", "aws_instance", drift.StatusModified),
		makeResult("r3", "aws_vpc", drift.StatusExtra),
	}
	d := compare.Compare(prev, curr)
	if len(d.Resolved) != 1 || len(d.Persisted) != 1 || len(d.New) != 1 {
		t.Fatalf("unexpected delta: %+v", d)
	}
}

func TestSummary(t *testing.T) {
	d := compare.Delta{
		Resolved:  make([]drift.Result, 2),
		New:       make([]drift.Result, 1),
		Persisted: make([]drift.Result, 3),
	}
	s := compare.Summary(d)
	if s != "Resolved: 2, New: 1, Persisted: 3" {
		t.Fatalf("unexpected summary: %s", s)
	}
}
