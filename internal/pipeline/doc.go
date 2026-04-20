// Package pipeline provides a lightweight composable stage-based processor
// for slices of drift.Result. Each Stage is a pure function that receives
// and returns []drift.Result, making it easy to chain filtering, enrichment,
// deduplication, classification, or any other transformation in a defined order.
//
// Example usage:
//
//	p := pipeline.New().
//		Add(myFilterStage).
//		Add(myEnrichStage)
//	output := p.Run(results)
package pipeline
