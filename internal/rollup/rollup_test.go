package rollup

import (
	"bytes"
	"strings"
	"testing"

	"driftctl-lite/internal/drift"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ResourceID: "sg-1", ResourceType: "aws_security_group", Status: drift.StatusMissing},
		{ResourceID: "sg-2", ResourceType: "aws_security_group", Status: drift.StatusInSync},
		{ResourceID: "sg-3", ResourceType: "aws_security_group", Status: drift.StatusModified},
		{ResourceID: "s3-1", ResourceType: "aws_s3_bucket", Status: drift.StatusExtra},
		{ResourceID: "s3-2", ResourceType: "aws_s3_bucket", Status: drift.StatusInSync},
	}
}

func TestCompute_Empty(t *testing.T) {
	r := Compute(nil)
	if len(r) != 0 {
		t.Fatalf("expected empty report, got %d entries", len(r))
	}
}

func TestCompute_Counts(t *testing.T) {
	r := Compute(makeResults())
	if len(r) != 2 {
		t.Fatalf("expected 2 type entries, got %d", len(r))
	}
	// sorted: aws_s3_bucket first
	s3 := r[0]
	if s3.Type != "aws_s3_bucket" {
		t.Fatalf("unexpected type: %s", s3.Type)
	}
	if s3.Total != 2 || s3.Extra != 1 || s3.InSync != 1 {
		t.Errorf("s3 counts wrong: %+v", s3)
	}
	sg := r[1]
	if sg.Total != 3 || sg.Missing != 1 || sg.Modified != 1 || sg.InSync != 1 {
		t.Errorf("sg counts wrong: %+v", sg)
	}
}

func TestCompute_AllInSync(t *testing.T) {
	results := []drift.Result{
		{ResourceID: "a", ResourceType: "aws_vpc", Status: drift.StatusInSync},
	}
	r := Compute(results)
	if r[0].InSync != 1 || r[0].Missing != 0 {
		t.Errorf("unexpected counts: %+v", r[0])
	}
}

func TestPrint_Output(t *testing.T) {
	r := Compute(makeResults())
	var buf bytes.Buffer
	Fprint(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "aws_s3_bucket") {
		t.Error("expected aws_s3_bucket in output")
	}
	if !strings.Contains(out, "aws_security_group") {
		t.Error("expected aws_security_group in output")
	}
	if !strings.Contains(out, "TYPE") {
		t.Error("expected header in output")
	}
}
