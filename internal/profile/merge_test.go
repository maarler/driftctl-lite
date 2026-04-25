package profile_test

import (
	"testing"

	"github.com/driftctl-lite/internal/config"
	"github.com/driftctl-lite/internal/profile"
)

func baseConfig() *config.Config {
	c := config.Default()
	return &c
}

func TestMergeIntoConfig_OutputFormat(t *testing.T) {
	p := profile.Profile{OutputFormat: "json"}
	cfg := baseConfig()
	profile.MergeIntoConfig(p, cfg)
	if cfg.Output != "json" {
		t.Errorf("expected output=json, got %q", cfg.Output)
	}
}

func TestMergeIntoConfig_OnlyDrift(t *testing.T) {
	p := profile.Profile{OnlyDrift: true}
	cfg := baseConfig()
	cfg.OnlyDrift = false
	profile.MergeIntoConfig(p, cfg)
	if !cfg.OnlyDrift {
		t.Error("expected OnlyDrift=true after merge")
	}
}

func TestMergeIntoConfig_EmptyProfile_NoChange(t *testing.T) {
	p := profile.Profile{}
	cfg := baseConfig()
	origOutput := cfg.Output
	profile.MergeIntoConfig(p, cfg)
	if cfg.Output != origOutput {
		t.Errorf("expected output unchanged, got %q", cfg.Output)
	}
}

func TestMergeIntoConfig_FilterType(t *testing.T) {
	p := profile.Profile{FilterType: "aws_s3_bucket"}
	cfg := baseConfig()
	profile.MergeIntoConfig(p, cfg)
	if cfg.FilterType != "aws_s3_bucket" {
		t.Errorf("expected filter_type=aws_s3_bucket, got %q", cfg.FilterType)
	}
}
