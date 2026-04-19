package dedupe_test

import (
	"testing"

	"github.com/example/driftctl-lite/internal/dedupe"
	"github.com/example/driftctl-lite/internal/drift"
)

func makeResult(id, typ string, status drift.Status) drift.Result {
	return drift.Result{ResourceID: id, ResourceType: typ, Status: status}
}

func TestApply_Empty(t *testing.T) {
	out := dedupe.Apply(nil, dedupe.KeepFirst)
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d", len(out))
	}
}

func TestApply_NoDuplicates(t *testing.T) {
	input := []drift.Result{
		makeResult("a", "aws_s3", drift.StatusInSync),
		makeResult("b", "aws_s3", drift.StatusMissing),
	}
	out := dedupe.Apply(input, dedupe.KeepFirst)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestApply_KeepFirst(t *testing.T) {
	input := []drift.Result{
		makeResult("a", "aws_s3", drift.StatusInSync),
		makeResult("a", "aws_s3", drift.StatusMissing),
	}
	out := dedupe.Apply(input, dedupe.KeepFirst)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if out[0].Status != drift.StatusInSync {
		t.Errorf("expected InSync, got %s", out[0].Status)
	}
}

func TestApply_KeepLast(t *testing.T) {
	input := []drift.Result{
		makeResult("a", "aws_s3", drift.StatusInSync),
		makeResult("a", "aws_s3", drift.StatusMissing),
	}
	out := dedupe.Apply(input, dedupe.KeepLast)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if out[0].Status != drift.StatusMissing {
		t.Errorf("expected Missing, got %s", out[0].Status)
	}
}

func TestApply_KeepDrift_PrefersDrift(t *testing.T) {
	input := []drift.Result{
		makeResult("a", "aws_s3", drift.StatusInSync),
		makeResult("a", "aws_s3", drift.StatusModified),
	}
	out := dedupe.Apply(input, dedupe.KeepDrift)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if out[0].Status != drift.StatusModified {
		t.Errorf("expected Modified, got %s", out[0].Status)
	}
}

func TestApply_KeepDrift_KeepsExistingDrift(t *testing.T) {
	input := []drift.Result{
		makeResult("a", "aws_s3", drift.StatusMissing),
		makeResult("a", "aws_s3", drift.StatusInSync),
	}
	out := dedupe.Apply(input, dedupe.KeepDrift)
	if out[0].Status != drift.StatusMissing {
		t.Errorf("expected Missing, got %s", out[0].Status)
	}
}
