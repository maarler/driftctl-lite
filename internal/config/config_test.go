package config

import (
	"encoding/json"
	"os"
	"testing"
)

func writeTempConfig(t *testing.T, cfg *Config) string {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	f, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	f.Write(data)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.StateFile != "state.json" {
		t.Errorf("expected state.json, got %s", cfg.StateFile)
	}
	if cfg.OutputFmt != "text" {
		t.Errorf("expected text, got %s", cfg.OutputFmt)
	}
}

func TestLoadFromFile_Exists(t *testing.T) {
	original := &Config{
		StateFile:  "custom_state.json",
		LiveSource: "custom_live.json",
		OutputFmt:  "json",
		OnlyDrift:  true,
	}
	path := writeTempConfig(t, original)
	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StateFile != "custom_state.json" {
		t.Errorf("expected custom_state.json, got %s", cfg.StateFile)
	}
	if !cfg.OnlyDrift {
		t.Error("expected OnlyDrift=true")
	}
}

func TestLoadFromFile_Missing(t *testing.T) {
	cfg, err := LoadFromFile("/nonexistent/config.json")
	if err != nil {
		t no error for missing file, got %v", err)
	}
	if cfg.StateFile != "state.json" {
		t.Errorf("expected default state.json, got %s", cfg.StateFile)
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := Default()
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestValidate_BadOutputFmt(t *testing.T) {
	cfg := Default()
	cfg.OutputFmt = "xml"
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for bad output_format")
	}
}

func TestValidate_EmptyStateFile(t *testing.T) {
	cfg := Default()
	cfg.StateFile = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for empty state_file")
	}
}
