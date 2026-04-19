package truncate_test

import (
	"testing"

	"github.com/owner/driftctl-lite/internal/drift"
	"github.com/owner/driftctl-lite/internal/truncate"
)

func makeResult(id string, status drift.Status) drift.Result {
	return drift.Result{ResourceID: id, ResourceType: "aws_s3_bucket", Status: status}
}

func TestApply_NoLimit(t *testing.T) {
	results := []drift.Result{
		makeResult("a", drift.StatusMissing),
		makeResult("b", drift.StatusInSync),
	}
	out, truncated := truncate.Apply(results, truncate.DefaultOptions())
	if truncated {
		t.Fatal("expected no truncation")
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestApply_LimitEnforced(t *testing.T) {
	results := []drift.Result{
		makeResult("a", drift.StatusMissing),
		makeResult("b", drift.StatusExtra),
		makeResult("c", drift.StatusModified),
	}
	out, truncated := truncate.Apply(results, truncate.Options{MaxResults: 2})
	if !truncated {
		t.Fatal("expected truncation")
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestApply_PreserveInSync_DropsInSync(t *testing.T) {
	results := []drift.Result{
		makeResult("a", drift.StatusInSync),
		makeResult("b", drift.StatusMissing),
		makeResult("c", drift.StatusInSync),
	}
	out, _ := truncate.Apply(results, truncate.Options{PreserveInSync: true})
	for _, r := range out {
		if r.Status == drift.StatusInSync {
			t.Errorf("in-sync result %q should have been dropped", r.ResourceID)
		}
	}
}

func TestApply_Empty(t *testing.T) {
	out, truncated := truncate.Apply(nil, truncate.Options{MaxResults: 5})
	if truncated || len(out) != 0 {
		t.Fatal("expected empty result without truncation")
	}
}

func TestTruncated_Helper(t *testing.T) {
	results := []drift.Result{
		makeResult("x", drift.StatusModified),
		makeResult("y", drift.StatusModified),
	}
	if !truncate.Truncated(results, truncate.Options{MaxResults: 1}) {
		t.Fatal("expected Truncated to return true")
	}
	if truncate.Truncated(results, truncate.Options{MaxResults: 10}) {
		t.Fatal("expected Truncated to return false")
	}
}
