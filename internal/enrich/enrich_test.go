package enrich_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/enrich"
)

func makeResult(id, rtype, status string) drift.Result {
	return drift.Result{ID: id, ResourceType: rtype, Status: status}
}

func TestApply_NoRules(t *testing.T) {
	results := []drift.Result{makeResult("a", "aws_s3_bucket", "ok")}
	out := enrich.Apply(results, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Attributes != nil && out[0].Attributes["_owner"] != "" {
		t.Error("expected no owner annotation")
	}
}

func TestApply_MatchingRule(t *testing.T) {
	results := []drift.Result{makeResult("b", "aws_s3_bucket", "missing")}
	rules := []enrich.Rule{
		{ResourceType: "aws_s3_bucket", Metadata: enrich.Metadata{Owner: "team-a", Environment: "prod"}},
	}
	out := enrich.Apply(results, rules)
	if out[0].Attributes["_owner"] != "team-a" {
		t.Errorf("expected owner team-a, got %q", out[0].Attributes["_owner"])
	}
	if out[0].Attributes["_environment"] != "prod" {
		t.Errorf("expected environment prod, got %q", out[0].Attributes["_environment"])
	}
}

func TestApply_NoMatchingRule(t *testing.T) {
	results := []drift.Result{makeResult("c", "aws_ec2_instance", "modified")}
	rules := []enrich.Rule{
		{ResourceType: "aws_s3_bucket", Metadata: enrich.Metadata{Owner: "team-b"}},
	}
	out := enrich.Apply(results, rules)
	if out[0].Attributes != nil && out[0].Attributes["_owner"] != "" {
		t.Error("expected no annotation for unmatched type")
	}
}

func TestApply_CustomMeta(t *testing.T) {
	results := []drift.Result{makeResult("d", "aws_rds_instance", "ok")}
	rules := []enrich.Rule{
		{ResourceType: "aws_rds_instance", Metadata: enrich.Metadata{
			CostCenter: "cc-99",
			Custom:     map[string]string{"tier": "gold"},
		}},
	}
	out := enrich.Apply(results, rules)
	if out[0].Attributes["_cost_center"] != "cc-99" {
		t.Errorf("expected cost center cc-99, got %q", out[0].Attributes["_cost_center"])
	}
	if out[0].Attributes["_tier"] != "gold" {
		t.Errorf("expected tier gold, got %q", out[0].Attributes["_tier"])
	}
}

func TestApply_PreservesExistingAttributes(t *testing.T) {
	res := makeResult("e", "aws_s3_bucket", "ok")
	res.Attributes = map[string]string{"region": "us-east-1"}
	rules := []enrich.Rule{
		{ResourceType: "aws_s3_bucket", Metadata: enrich.Metadata{Owner: "team-c"}},
	}
	out := enrich.Apply([]drift.Result{res}, rules)
	if out[0].Attributes["region"] != "us-east-1" {
		t.Error("expected existing attribute to be preserved")
	}
	if out[0].Attributes["_owner"] != "team-c" {
		t.Error("expected owner annotation to be added")
	}
}
