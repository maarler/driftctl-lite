// Package profile manages named drift-check profiles.
//
// A profile is a reusable bundle of scan options (output format, filters,
// severity threshold, tag selectors) stored in a JSON registry file.
// Profiles can be loaded from disk and merged into the active Config,
// allowing teams to share consistent scan presets across CI and local runs.
//
// Usage:
//
//	reg, err := profile.LoadFromFile(".driftprofiles.json")
//	p, err := reg.Get("ci")
//	profile.MergeIntoConfig(p, cfg)
package profile
