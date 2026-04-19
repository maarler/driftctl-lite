package enrich_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/enrich"
)

func TestEnrich_RoundTrip(t *testing.T) {
	results := []drift.Result{
		{ID: "bucket-1", ResourceType: "aws_s3_bucket", Status: "missing"},
		{ID: "instance-1", ResourceType: "aws_ec2_instance", Status: "modified"},
		{ID: "db-1", ResourceType: "aws_rds_instance", Status: "ok"},
	}
	rules := []enrich.Rule{
		{ResourceType: "aws_s3_bucket", Metadata: enrich.Metadata{Owner: "storage-team", Environment: "prod", CostCenter: "cc-10"}},
		{ResourceType: "aws_ec2_instance", Metadata: enrich.Metadata{Owner: "compute-team", Environment: "staging"}},
	}
	out := enrich.Apply(results, rules)
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
	if out[0].Attributes["_owner"] != "storage-team" {
		t.Errorf("bucket owner mismatch: %q", out[0].Attributes["_owner"])
	}
	if out[1].Attributes["_environment"] != "staging" {
		t.Errorf("instance environment mismatch: %q", out[1].Attributes["_environment"])
	}
	if out[2].Attributes != nil && out[2].Attributes["_owner"] != "" {
		t.Error("rds instance should have no owner annotation")
	}
}
