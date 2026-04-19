package group_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/group"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ResourceID: "a", ResourceType: "aws_s3_bucket", Missing: true},
		{ResourceID: "b", ResourceType: "aws_s3_bucket", Modified: true},
		{ResourceID: "c", ResourceType: "aws_instance", Extra: true},
		{ResourceID: "d", ResourceType: "aws_instance"},
	}
}

func TestByType(t *testing.T) {
	results := makeResults()
	groups := group.ByType(results)

	if len(groups["aws_s3_bucket"]) != 2 {
		t.Errorf("expected 2 s3 results, got %d", len(groups["aws_s3_bucket"]))
	}
	if len(groups["aws_instance"]) != 2 {
		t.Errorf("expected 2 instance results, got %d", len(groups["aws_instance"]))
	}
}

func TestByStatus(t *testing.T) {
	results := makeResults()
	groups := group.ByStatus(results)

	if len(groups["missing"]) != 1 {
		t.Errorf("expected 1 missing, got %d", len(groups["missing"]))
	}
	if len(groups["modified"]) != 1 {
		t.Errorf("expected 1 modified, got %d", len(groups["modified"]))
	}
	if len(groups["extra"]) != 1 {
		t.Errorf("expected 1 extra, got %d", len(groups["extra"]))
	}
	if len(groups["in_sync"]) != 1 {
		t.Errorf("expected 1 in_sync, got %d", len(groups["in_sync"]))
	}
}

func TestSummarize(t *testing.T) {
	results := makeResults()
	groups := group.ByStatus(results)
	summaries := group.Summarize(groups)

	if len(summaries) != 4 {
		t.Errorf("expected 4 summaries, got %d", len(summaries))
	}
	for _, s := range summaries {
		if s.Count != 1 {
			t.Errorf("expected count 1 for key %s, got %d", s.Key, s.Count)
		}
	}
}
