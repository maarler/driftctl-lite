package mask_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/mask"
)

func makeResult(id, typ, status string, diffs map[string][2]string) drift.Result {
	return drift.Result{ID: id, Type: typ, Status: status, Diffs: diffs}
}

func TestApply_NoSensitiveFields(t *testing.T) {
	results := []drift.Result{
		makeResult("r1", "vm", "modified", map[string][2]string{
			"name": {"old", "new"},
		}),
	}
	opts := mask.Options{Fields: []string{"password"}, Placeholder: "***"}
	out := mask.Apply(results, opts)
	if out[0].Diffs["name"] != ([2]string{"old", "new"}) {
		t.Errorf("expected name diff to be unchanged")
	}
}

func TestApply_SensitiveFieldMasked(t *testing.T) {
	results := []drift.Result{
		makeResult("r1", "vm", "modified", map[string][2]string{
			"password": {"hunter2", "s3cr3t"},
			"region":   {"us-east-1", "eu-west-1"},
		}),
	}
	opts := mask.DefaultOptions()
	out := mask.Apply(results, opts)
	if out[0].Diffs["password"] != ([2]string{"***", "***"}) {
		t.Errorf("expected password to be masked, got %v", out[0].Diffs["password"])
	}
	if out[0].Diffs["region"] != ([2]string{"us-east-1", "eu-west-1"}) {
		t.Errorf("expected region to be unchanged")
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	results := []drift.Result{
		makeResult("r2", "db", "modified", map[string][2]string{
			"token": {"abc", "xyz"},
		}),
	}
	opts := mask.Options{Fields: []string{"token"}, Placeholder: "<REDACTED>"}
	out := mask.Apply(results, opts)
	if out[0].Diffs["token"] != ([2]string{"<REDACTED>", "<REDACTED>"}) {
		t.Errorf("unexpected placeholder value: %v", out[0].Diffs["token"])
	}
}

func TestApply_EmptyFields_NoChange(t *testing.T) {
	results := []drift.Result{
		makeResult("r3", "bucket", "modified", map[string][2]string{
			"secret": {"a", "b"},
		}),
	}
	opts := mask.Options{Fields: []string{}}
	out := mask.Apply(results, opts)
	if out[0].Diffs["secret"] != ([2]string{"a", "b"}) {
		t.Errorf("expected no masking when Fields is empty")
	}
}

func TestApply_NoDiffs_PassesThrough(t *testing.T) {
	results := []drift.Result{
		makeResult("r4", "sg", "in_sync", nil),
	}
	out := mask.Apply(results, mask.DefaultOptions())
	if out[0].Diffs != nil {
		t.Errorf("expected nil diffs to remain nil")
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	origDiffs := map[string][2]string{"api_key": {"real", "new"}}
	results := []drift.Result{
		makeResult("r5", "fn", "modified", origDiffs),
	}
	mask.Apply(results, mask.DefaultOptions())
	if results[0].Diffs["api_key"] != ([2]string{"real", "new"}) {
		t.Errorf("Apply must not mutate the original result")
	}
}
