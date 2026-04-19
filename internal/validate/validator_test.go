package validate_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/validate"
)

func makeResult(id, rtype, status string) drift.Result {
	return drift.Result{
		ResourceID:   id,
		ResourceType: rtype,
		Status:       status,
	}
}

func TestValidate_NoViolations(t *testing.T) {
	results := []drift.Result{
		makeResult("vpc-1", "aws_vpc", "in_sync"),
		makeResult("sg-1", "aws_security_group", "modified"),
	}
	violations, err := validate.Validate(results, validate.DefaultRules())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestValidate_EmptyID(t *testing.T) {
	results := []drift.Result{
		makeResult("", "aws_vpc", "missing"),
	}
	violations, err := validate.Validate(results, validate.DefaultRules())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].RuleName != "non-empty-id" {
		t.Errorf("expected rule non-empty-id, got %s", violations[0].RuleName)
	}
}

func TestValidate_EmptyType(t *testing.T) {
	results := []drift.Result{
		makeResult("vpc-1", "", "extra"),
	}
	violations, err := validate.Validate(results, validate.DefaultRules())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].RuleName != "non-empty-type" {
		t.Errorf("expected rule non-empty-type, got %s", violations[0].RuleName)
	}
}

func TestValidate_NoRules_ReturnsError(t *testing.T) {
	results := []drift.Result{makeResult("id-1", "aws_vpc", "in_sync")}
	_, err := validate.Validate(results, nil)
	if err == nil {
		t.Error("expected error for empty rules, got nil")
	}
}

func TestValidate_CustomRule(t *testing.T) {
	rules := []validate.Rule{
		{
			Name:    "only-in-sync",
			Message: "resource must be in sync",
			Check: func(r drift.Result) bool {
				return r.Status == "in_sync"
			},
		},
	}
	results := []drift.Result{
		makeResult("vpc-1", "aws_vpc", "in_sync"),
		makeResult("sg-1", "aws_sg", "modified"),
	}
	violations, err := validate.Validate(results, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].ResourceID != "sg-1" {
		t.Errorf("expected sg-1, got %s", violations[0].ResourceID)
	}
}
