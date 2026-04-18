// Package baseline provides save, load, and comparison of drift baselines.
//
// A baseline captures a snapshot of drift results at a specific point in time.
// Subsequent runs can be compared against the baseline to surface only new
// or changed drift — reducing noise from known, accepted drift.
//
// Usage:
//
//	baseline.Save("baseline.json", results)
//	b, _ := baseline.Load("baseline.json")
//	delta := baseline.Compare(b, newResults)
package baseline
