package config

import (
	"testing"
)

func TestParseFlags_Defaults(t *testing.T) {
	f, err := ParseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.ConfigFile != ".driftctl.json" {
		t.Errorf("expected .driftctl.json, got %s", f.ConfigFile)
	}
	if f.OnlyDrift {
		t.Error("expected OnlyDrift=false by default")
	}
}

func TestParseFlags_Overrides(t *testing.T) {
	f, err := ParseFlags([]string{
		"-state", "my_state.json",
		"-output", "json",
		"-only-drift",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.StateFile != "my_state.json" {
		t.Errorf("expected my_state.json, got %s", f.StateFile)
	}
	if f.OutputFmt != "json" {
		t.Errorf("expected json, got %s", f.OutputFmt)
	}
	if !f.OnlyDrift {
		t.Error("expected OnlyDrift=true")
	}
}

func TestMerge_FlagsOverrideConfig(t *testing.T) {
	cfg := Default()
	flags := &Flags{
		StateFile: "override.json",
		OutputFmt: "json",
		OnlyDrift: true,
	}
	merged := Merge(cfg, flags)
	if merged.StateFile != "override.json" {
		t.Errorf("expected override.json, got %s", merged.StateFile)
	}
	if merged.OutputFmt != "json" {
		t.Errorf("expected json, got %s", merged.OutputFmt)
	}
	if !merged.OnlyDrift {
		t.Error("expected OnlyDrift=true after merge")
	}
	if merged.LiveSource != cfg.LiveSource {
		t.Errorf("live source should remain default, got %s", merged.LiveSource)
	}
}

func TestMerge_EmptyFlagsKeepConfig(t *testing.T) {
	cfg := &Config{
		StateFile:  "kept.json",
		LiveSource: "kept_live.json",
		OutputFmt:  "text",
	}
	merged := Merge(cfg, &Flags{})
	if merged.StateFile != "kept.json" {
		t.Errorf("expected kept.json, got %s", merged.StateFile)
	}
}
