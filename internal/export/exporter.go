// Package export provides functionality to export drift results to various formats.
package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/owner/driftctl-lite/internal/drift"
)

// Format represents the export format.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Record is a flat representation of a drift result for export.
type Record struct {
	Timestamp  string `json:"timestamp" csv:"timestamp"`
	ResourceID string `json:"resource_id" csv:"resource_id"`
	Type       string `json:"type" csv:"type"`
	Status     string `json:"status" csv:"status"`
}

// Export writes drift results to w in the specified format.
func Export(results []drift.Result, format Format, w io.Writer) error {
	records := toRecords(results)
	switch format {
	case FormatJSON:
		return exportJSON(records, w)
	case FormatCSV:
		return exportCSV(records, w)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// ExportToFile writes drift results to a file at path in the specified format.
func ExportToFile(results []drift.Result, format Format, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("export: create file: %w", err)
	}
	defer f.Close()
	return Export(results, format, f)
}

func toRecords(results []drift.Result) []Record {
	ts := time.Now().UTC().Format(time.RFC3339)
	out := make([]Record, 0, len(results))
	for _, r := range results {
		out = append(out, Record{
			Timestamp:  ts,
			ResourceID: r.ResourceID,
			Type:       r.ResourceType,
			Status:     string(r.Status),
		})
	}
	return out
}

func exportJSON(records []Record, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func exportCSV(records []Record, w io.Writer) error {
	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"timestamp", "resource_id", "type", "status"})
	for _, r := range records {
		_ = cw.Write([]string{r.Timestamp, r.ResourceID, r.Type, r.Status})
	}
	cw.Flush()
	return cw.Error()
}
