// Package mask provides field-level masking for drift results,
// replacing selected attribute values with a placeholder to prevent
// accidental exposure of sensitive infrastructure details in output.
package mask

import "github.com/driftctl-lite/internal/drift"

const defaultPlaceholder = "***"

// Options controls masking behaviour.
type Options struct {
	// Fields lists attribute key names whose values will be replaced.
	Fields []string
	// Placeholder is the string written in place of masked values.
	// Defaults to "***" when empty.
	Placeholder string
}

// DefaultOptions returns a sensible default set of masked field names.
func DefaultOptions() Options {
	return Options{
		Fields:      []string{"password", "secret", "token", "api_key", "private_key"},
		Placeholder: defaultPlaceholder,
	}
}

// Apply returns a copy of results with matching diff keys masked.
// Original results are not modified.
func Apply(results []drift.Result, opts Options) []drift.Result {
	if len(opts.Fields) == 0 {
		return results
	}
	ph := opts.Placeholder
	if ph == "" {
		ph = defaultPlaceholder
	}
	set := make(map[string]struct{}, len(opts.Fields))
	for _, f := range opts.Fields {
		set[f] = struct{}{}
	}

	out := make([]drift.Result, len(results))
	for i, r := range results {
		out[i] = maskOne(r, set, ph)
	}
	return out
}

func maskOne(r drift.Result, fields map[string]struct{}, ph string) drift.Result {
	if len(r.Diffs) == 0 {
		return r
	}
	masked := make(map[string][2]string, len(r.Diffs))
	for k, v := range r.Diffs {
		if _, sensitive := fields[k]; sensitive {
			masked[k] = [2]string{ph, ph}
		} else {
			masked[k] = v
		}
	}
	r.Diffs = masked
	return r
}
