package output

import (
	"fmt"
	"io"
	"os"

	"github.com/driftctl-lite/internal/drift"
)

// Format represents the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes drift results to an output destination.
type Reporter struct {
	Writer io.Writer
	Format Format
}

// NewReporter creates a Reporter writing to stdout with the given format.
func NewReporter(format Format) *Reporter {
	return &Reporter{
		Writer: os.Stdout,
		Format: format,
	}
}

// Report writes the drift result to the reporter's writer.
func (r *Reporter) Report(result drift.Result) error {
	switch r.Format {
	case FormatJSON:
		return r.reportJSON(result)
	default:
		return r.reportText(result)
	}
}

func (r *Reporter) reportText(result drift.Result) error {
	if !result.HasDrift() {
		fmt.Fprintln(r.Writer, "✅ No drift detected.")
		return nil
	}
	for _, res := range result.Missing {
		fmt.Fprintf(r.Writer, "❌ MISSING   %s (%s)\n", res.ID, res.Type)
	}
	for _, res := range result.Extra {
		fmt.Fprintf(r.Writer, "➕ EXTRA     %s (%s)\n", res.ID, res.Type)
	}
	for _, res := range result.Modified {
		fmt.Fprintf(r.Writer, "✏️  MODIFIED  %s (%s)\n", res.ID, res.Type)
	}
	return nil
}

func (r *Reporter) reportJSON(result drift.Result) error {
	data, err := result.MarshalJSON()
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	_, err = fmt.Fprintln(r.Writer, string(data))
	return err
}
