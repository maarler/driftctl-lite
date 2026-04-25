// Package normalize standardizes drift result fields to ensure consistent
// comparisons and output regardless of how upstream data was formatted.
//
// Typical usage:
//
//	results = normalize.Apply(results, normalize.DefaultOptions())
package normalize
