package window

import (
	"testing"
	"time"

	"driftctl-lite/internal/drift"
)

func makeResult(id, rtype, status string) drift.Result {
	return drift.Result{
		ResourceID:   id,
		ResourceType: rtype,
		Status:       drift.Status(status),
	}
}

func TestCollect_EmptyWindow(t *testing.T) {
	w := New(time.Minute)
	got := w.Collect()
	if len(got) != 0 {
		t.Fatalf("expected 0 results, got %d", len(got))
	}
}

func TestCollect_WithinTTL(t *testing.T) {
	w := New(time.Minute)
	w.Add([]drift.Result{
		makeResult("a", "aws_s3_bucket", "missing"),
		makeResult("b", "aws_s3_bucket", "ok"),
	})
	got := w.Collect()
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestCollect_ExpiredEntriesDropped(t *testing.T) {
	w := New(time.Minute)

	// Inject a past timestamp for the first batch.
	now := time.Now()
	w.nowFn = func() time.Time { return now.Add(-2 * time.Minute) }
	w.Add([]drift.Result{makeResult("old", "aws_s3_bucket", "missing")})

	// Restore real time for the second batch.
	w.nowFn = func() time.Time { return now }
	w.Add([]drift.Result{makeResult("new", "aws_s3_bucket", "ok")})

	got := w.Collect()
	if len(got) != 1 {
		t.Fatalf("expected 1 result after expiry, got %d", len(got))
	}
	if got[0].ResourceID != "new" {
		t.Errorf("expected result id 'new', got %q", got[0].ResourceID)
	}
}

func TestCollect_DeduplicatesOnKey(t *testing.T) {
	w := New(time.Minute)
	w.Add([]drift.Result{makeResult("x", "aws_instance", "missing")})
	w.Add([]drift.Result{makeResult("x", "aws_instance", "modified")})

	got := w.Collect()
	if len(got) != 1 {
		t.Fatalf("expected 1 deduplicated result, got %d", len(got))
	}
	// Latest entry wins.
	if string(got[0].Status) != "modified" {
		t.Errorf("expected status 'modified', got %q", got[0].Status)
	}
}

func TestLen(t *testing.T) {
	w := New(time.Minute)
	if w.Len() != 0 {
		t.Fatal("expected Len 0 on empty window")
	}
	w.Add([]drift.Result{makeResult("a", "t", "ok")})
	w.Add([]drift.Result{makeResult("b", "t", "ok")})
	if w.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", w.Len())
	}
}
