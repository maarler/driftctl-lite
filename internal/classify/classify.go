// Package classify categorises drift results by severity level.
package classify

import "github.com/driftctl-lite/internal/drift"

// Level represents a severity classification.
type Level string

const (
	LevelCritical Level = "critical"
	LevelWarning   Level = "warning"
	LevelInfo      Level = "info"
)

// Classification holds a result alongside its assigned level.
type Classification struct {
	Result drift.Result
	Level  Level
}

// Options controls classification thresholds.
type Options struct {
	// Types that are always critical when missing.
	CriticalTypes []string
}

// Apply classifies each result and returns the list of classifications.
func Apply(results []drift.Result, opts Options) []Classification {
	critical := make(map[string]bool, len(opts.CriticalTypes))
	for _, t := range opts.CriticalTypes {
		critical[t] = true
	}

	out := make([]Classification, 0, len(results))
	for _, r := range results {
		out = append(out, Classification{
			Result: r,
			Level:  classify(r, critical),
		})
	}
	return out
}

func classify(r drift.Result, critical map[string]bool) Level {
	switch r.Status {
	case drift.StatusMissing:
		if critical[r.ResourceType] {
			return LevelCritical
		}
		return LevelWarning
	case drift.StatusExtra:
		return LevelWarning
	case drift.StatusModified:
		if critical[r.ResourceType] {
			return LevelCritical
		}
		return LevelWarning
	default:
		return LevelInfo
	}
}
