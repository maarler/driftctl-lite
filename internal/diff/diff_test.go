package diff

import (
	"testing"
)

func TestCompute_NoDiff(t *testing.T) {
	declared := map[string]interface{}{"region": "us-east-1", "size": "t2.micro"}
	live := map[string]interface{}{"region": "us-east-1", "size": "t2.micro"}
	diffs := Compute(declared, live)
	if len(diffs) != 0 {
		t.Fatalf("expected no diffs, got %d", len(diffs))
	}
}

func TestCompute_Modified(t *testing.T) {
	declared := map[string]interface{}{"size": "t2.micro"}
	live := map[string]interface{}{"size": "t2.large"}
	diffs := Compute(declared, live)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Field != "size" {
		t.Errorf("unexpected field: %s", diffs[0].Field)
	}
	if diffs[0].Declared != "t2.micro" || diffs[0].Live != "t2.large" {
		t.Errorf("unexpected values: %v => %v", diffs[0].Declared, diffs[0].Live)
	}
}

func TestCompute_MissingInLive(t *testing.T) {
	declared := map[string]interface{}{"tag": "prod"}
	live := map[string]interface{}{}
	diffs := Compute(declared, live)
	if len(diffs) != 1 || diffs[0].Live != nil {
		t.Fatalf("expected one diff with nil live, got %+v", diffs)
	}
}

func TestCompute_ExtraInLive(t *testing.T) {
	declared := map[string]interface{}{}
	live := map[string]interface{}{"extra": "value"}
	diffs := Compute(declared, live)
	if len(diffs) != 1 || diffs[0].Declared != nil {
		t.Fatalf("expected one diff with nil declared, got %+v", diffs)
	}
}

func TestCompute_NilMaps(t *testing.T) {
	diffs := Compute(nil, nil)
	if len(diffs) != 0 {
		t.Fatalf("expected no diffs for nil maps, got %d", len(diffs))
	}
}

func TestFieldDiff_String(t *testing.T) {
	f := FieldDiff{Field: "region", Declared: "us-east-1", Live: "eu-west-1"}
	s := f.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
