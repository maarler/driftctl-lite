// Package checkpoint saves and restores named snapshots of drift results.
//
// Each checkpoint is stored as a JSON file in a user-specified directory,
// identified by a short name (e.g. "pre-deploy", "post-deploy").
// Checkpoints can be listed and loaded for later comparison or audit.
package checkpoint
