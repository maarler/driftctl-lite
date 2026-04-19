package classify_test

import (
	"testing"

	"github.com/driftctl-lite/internal/classify"
	"github.com/driftctl-lite/internal/drift"
)

func makeResult(id, rtype string, status drift.Status) drift.Result {
	return drift.Result{ResourceID: id, ResourceType: rtype, Status: status}
}

func TestApply_InSync_IsInfo(t *testing.T) {
	results := []drift.Result{makeResult("r1", "aws_s3_bucket", drift.StatusInSync)}
	got := classify.Apply(results, classify.Options{})
	if got[0].Level != classify.LevelInfo {
		t.Fatalf("expected info, got %s", got[0].Level)
	}
}

func TestApply_Missing_Warning(t *testing.T) {
	results := []drift.Result{makeResult("r1", "aws_s3_bucket", drift.StatusMissing)}
	got := classify.Apply(results, classify.Options{})
	if got[0].Level != classify.LevelWarning {
		t.Fatalf("expected warning, got %s", got[0].Level)
	}
}

func TestApply_Missing_CriticalType(t *testing.T) {
	results := []drift.Result{makeResult("r1", "aws_iam_role", drift.StatusMissing)}
	opts := classify.Options{CriticalTypes: []string{"aws_iam_role"}}
	got := classify.Apply(results, opts)
	if got[0].Level != classify.LevelCritical {
		t.Fatalf("expected critical, got %s", got[0].Level)
	}
}

func TestApply_Extra_Warning(t *testing.T) {
	results := []drift.Result{makeResult("r2", "aws_sg", drift.StatusExtra)}
	got := classify.Apply(results, classify.Options{})
	if got[0].Level != classify.LevelWarning {
		t.Fatalf("expected warning, got %s", got[0].Level)
	}
}

func TestApply_Modified_CriticalType(t *testing.T) {
	results := []drift.Result{makeResult("r3", "aws_iam_role", drift.StatusModified)}
	opts := classify.Options{CriticalTypes: []string{"aws_iam_role"}}
	got := classify.Apply(results, opts)
	if got[0].Level != classify.LevelCritical {
		t.Fatalf("expected critical, got %s", got[0].Level)
	}
}

func TestApply_Empty(t *testing.T) {
	got := classify.Apply(nil, classify.Options{})
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %d items", len(got))
	}
}
