package policy_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/policy"
)

func makeResult(id, rtype string, status drift.Status) drift.Result {
	return drift.Result{ResourceID: id, ResourceType: rtype, Status: status}
}

func TestEvaluate_NoViolations_WhenAllInSync(t *testing.T) {
	results := []drift.Result{makeResult("vpc-1", "aws_vpc", drift.StatusInSync)}
	rules := []policy.Rule{{ResourceType: "aws_vpc", DisallowModified: true}}
	violations := policy.Evaluate(results, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestEvaluate_Missing(t *testing.T) {
	results := []drift.Result{makeResult("sg-1", "aws_sg", drift.StatusMissing)}
	rules := []policy.Rule{{ResourceType: "aws_sg", DisallowMissing: true}}
	violations := policy.Evaluate(results, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Resource != "sg-1" {
		t.Errorf("unexpected resource: %s", violations[0].Resource)
	}
}

func TestEvaluate_Extra(t *testing.T) {
	results := []drift.Result{makeResult("bucket-99", "aws_s3", drift.StatusExtra)}
	rules := []policy.Rule{{ResourceType: "aws_s3", DisallowExtra: true}}
	violations := policy.Evaluate(results, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestEvaluate_Modified(t *testing.T) {
	results := []drift.Result{makeResult("rds-1", "aws_rds", drift.StatusModified)}
	rules := []policy.Rule{{ResourceType: "aws_rds", DisallowModified: true}}
	violations := policy.Evaluate(results, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestEvaluate_NoMatchingRule(t *testing.T) {
	results := []drift.Result{makeResult("ec2-1", "aws_ec2", drift.StatusModified)}
	rules := []policy.Rule{{ResourceType: "aws_s3", DisallowModified: true}}
	violations := policy.Evaluate(results, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestEvaluate_MultipleViolations(t *testing.T) {
	results := []drift.Result{
		makeResult("sg-1", "aws_sg", drift.StatusMissing),
		makeResult("sg-2", "aws_sg", drift.StatusExtra),
	}
	rules := []policy.Rule{{ResourceType: "aws_sg", DisallowMissing: true, DisallowExtra: true}}
	violations := policy.Evaluate(results, rules)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}
