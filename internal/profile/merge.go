package profile

import "github.com/driftctl-lite/internal/config"

// MergeIntoConfig overlays a Profile's settings onto cfg.
// Only non-zero profile fields override the config values.
func MergeIntoConfig(p Profile, cfg *config.Config) {
	if p.OutputFormat != "" {
		cfg.Output = p.OutputFormat
	}
	if p.OnlyDrift {
		cfg.OnlyDrift = true
	}
	if p.FilterType != "" {
		cfg.FilterType = p.FilterType
	}
}
