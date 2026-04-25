package profile_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/driftctl-lite/internal/profile"
)

func writeTempRegistry(t *testing.T, reg profile.Registry) string {
	t.Helper()
	data, err := json.Marshal(reg)
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "profiles.json")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadFromFile_Missing(t *testing.T) {
	reg, err := profile.LoadFromFile("/no/such/file.json")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if len(reg.Profiles) != 0 {
		t.Errorf("expected empty registry")
	}
}

func TestLoadFromFile_Valid(t *testing.T) {
	src := profile.Registry{
		Profiles: map[string]profile.Profile{
			"ci": {Name: "ci", OutputFormat: "json", OnlyDrift: true},
		},
	}
	path := writeTempRegistry(t, src)
	reg, err := profile.LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, err := reg.Get("ci")
	if err != nil {
		t.Fatalf("expected ci profile: %v", err)
	}
	if p.OutputFormat != "json" {
		t.Errorf("expected json format, got %q", p.OutputFormat)
	}
	if !p.OnlyDrift {
		t.Error("expected only_drift=true")
	}
}

func TestGet_Missing(t *testing.T) {
	reg := &profile.Registry{Profiles: make(map[string]profile.Profile)}
	_, err := reg.Get("nonexistent")
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestDefault(t *testing.T) {
	p := profile.Default()
	if p.Name != "default" {
		t.Errorf("expected name=default, got %q", p.Name)
	}
	if p.OutputFormat != "text" {
		t.Errorf("expected text format, got %q", p.OutputFormat)
	}
}
