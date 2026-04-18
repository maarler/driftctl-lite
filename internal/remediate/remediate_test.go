package remediate_test

import (
	"testing"

	"github.com/owner/driftctl-lite/internal/drift"
	"github.com/owner/driftctl-lite/internal/remediate"
)

func makeResult(id, rtype string, status drift.Status) drift.Result {
	return drift.Result{ResourceID: id, ResourceType: rtype, Status: status}
}

func TestSuggest_NoActions_WhenAllInSync(t *testing.T) {
	results := []drift.Result{
		makeResult("vpc-1", "aws_vpc", drift.StatusInSync),
	}
	actions := remediate.Suggest(results)
	if len(actions) != 0 {
		t.Fatalf("expected 0 actions, got %d", len(actions))
	}
}

func TestSuggest_Missing(t *testing.T) {
	results := []drift.Result{makeResult("sg-1", "aws_security_group", drift.StatusMissing)}
	actions := remediate.Suggest(results)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Severity != "high" {
		t.Errorf("expected severity high, got %s", actions[0].Severity)
	}
}

func TestSuggest_Extra(t *testing.T) {
	results := []drift.Result{makeResult("s3-1", "aws_s3_bucket", drift.StatusExtra)}
	actions := remediate.Suggest(results)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Severity != "medium" {
		t.Errorf("expected severity medium, got %s", actions[0].Severity)
	}
}

func TestSuggest_Modified(t *testing.T) {
	results := []drift.Result{makeResult("ec2-1", "aws_instance", drift.StatusModified)}
	actions := remediate.Suggest(results)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Severity != "low" {
		t.Errorf("expected severity low, got %s", actions[0].Severity)
	}
	if actions[0].ResourceID != "ec2-1" {
		t.Errorf("unexpected resource id: %s", actions[0].ResourceID)
	}
}

func TestSuggest_Mixed(t *testing.T) {
	results := []drift.Result{
		makeResult("a", "aws_vpc", drift.StatusInSync),
		makeResult("b", "aws_vpc", drift.StatusMissing),
		makeResult("c", "aws_vpc", drift.StatusExtra),
	}
	actions := remediate.Suggest(results)
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(actions))
	}
}
