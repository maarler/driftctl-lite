// Package alert provides threshold-based alerting for drift detection results.
//
// Use Evaluate to compare a slice of drift results against Thresholds,
// receiving an Alert that indicates OK, WARNING, or CRITICAL status.
//
// Example:
//
//	t := alert.DefaultThresholds()
//	a := alert.Evaluate(results, t)
//	alert.Print(a)
package alert
