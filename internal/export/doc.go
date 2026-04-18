// Package export provides utilities for exporting drift detection results
// to external formats such as JSON and CSV.
//
// Typical usage:
//
//	err := export.Export(results, export.FormatCSV, os.Stdout)
//
// or write directly to a file:
//
//	err := export.ExportToFile(results, export.FormatJSON, "drift-report.json")
package export
