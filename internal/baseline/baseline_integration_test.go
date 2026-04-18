package baseline_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/driftctl-lite/internal/baseline"
	"github.com/driftctl-lite/internal/drift"
)

func TestBaseline_SaveCompare_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "base.json")

	initial := []drift.Result{
		{ResourceID: "vpc-1", ResourceType: "aws_vpc", Status: "ok"},
		{ResourceID: "sg-1", ResourceType: "aws_security_group", Status: "ok"},
	}
	if err := baseline.Save(path, initial); err != nil {
		t.Fatalf("initial Save: %v", err)
	}

	b, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	updated := []drift.Result{
		{ResourceID: "vpc-1", ResourceType: "aws_vpc", Status: "ok"},
		{ResourceID: "sg-1", ResourceType: "aws_security_group", Status: "modified"},
		{ResourceID: "igw-1", ResourceType: "aws_internet_gateway", Status: "extra"},
	}

	delta := baseline.Compare(b, updated)
	if len(delta) != 2 {
		t.Errorf("expected 2 delta items, got %d: %v", len(delta), delta)
	}

	ids := map[string]bool{}
	for _, r := range delta {
		ids[r.ResourceID] = true
	}
	for _, expected := range []string{"sg-1", "igw-1"} {
		if !ids[expected] {
			t.Errorf("expected %s in delta", expected)
		}
	}
	_ = fmt.Sprintf("delta count: %d", len(delta))
}
