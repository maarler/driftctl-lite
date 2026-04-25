// Package reorder provides deterministic ordering of drift detection results.
//
// Results can be sorted by resource ID, type, or drift status. Sorting is
// stable, preserving the relative order of equal elements. This is useful
// for producing consistent CLI output and for snapshot comparisons.
//
// Example:
//
//	sorted := reorder.Apply(results, reorder.Options{
//		By:        reorder.FieldStatus,
//		Ascending: true,
//	})
package reorder
