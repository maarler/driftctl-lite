package trend_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"driftctl-lite/internal/trend"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "trend.json")
}

func TestLoad_Missing(t *testing.T) {
	entries, err := trend.Load("/nonexistent/trend.json")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty slice")
	}
}

func TestAppendAndLoad(t *testing.T) {
	p := tempPath(t)
	e1 := trend.Entry{Timestamp: time.Now(), Total: 3, Missing: 1, Extra: 1, Modified: 1}
	e2 := trend.Entry{Timestamp: time.Now(), Total: 1, Missing: 0, Extra: 0, Modified: 1}
	if err := trend.Append(p, e1); err != nil {
		t.Fatal(err)
	}
	if err := trend.Append(p, e2); err != nil {
		t.Fatal(err)
	}
	entries, err := trend.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestAnalyze_Stable(t *testing.T) {
	tr := trend.Analyze([]trend.Entry{})
	if tr.Direction != "stable" || tr.Delta != 0 {
		t.Fatalf("unexpected trend: %+v", tr)
	}
}

func TestAnalyze_Worsening(t *testing.T) {
	entries := []trend.Entry{{Total: 2}, {Total: 5}}
	tr := trend.Analyze(entries)
	if tr.Direction != "worsening" || tr.Delta != 3 {
		t.Fatalf("unexpected trend: %+v", tr)
	}
}

func TestAnalyze_Improving(t *testing.T) {
	entries := []trend.Entry{{Total: 5}, {Total: 2}}
	tr := trend.Analyze(entries)
	if tr.Direction != "improving" || tr.Delta != -3 {
		t.Fatalf("unexpected trend: %+v", tr)
	}
}

func TestAppend_InvalidPath(t *testing.T) {
	e := trend.Entry{Timestamp: time.Now(), Total: 1}
	err := trend.Append("/no/such/dir/trend.json", e)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not-json"), 0o644)
	_, err := trend.Load(p)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestEntry_JSONRoundTrip(t *testing.T) {
	e := trend.Entry{Timestamp: time.Now().Truncate(time.Second), Total: 4, Missing: 2, Extra: 1, Modified: 1}
	data, _ := json.Marshal(e)
	var got trend.Entry
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Total != e.Total {
		t.Fatalf("total mismatch")
	}
}
