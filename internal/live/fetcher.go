package live

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/driftctl-lite/internal/state"
)

// Source represents where live state is fetched from.
type Source string

const (
	SourceFile Source = "file"
)

// Fetcher retrieves live infrastructure state.
type Fetcher struct {
	Source Source
	Path   string
}

// NewFetcher creates a Fetcher for the given source and path.
func NewFetcher(source Source, path string) *Fetcher {
	return &Fetcher{Source: source, Path: path}
}

// Fetch returns the live resource map keyed by resource ID.
func (f *Fetcher) Fetch() (map[string]state.Resource, error) {
	switch f.Source {
	case SourceFile:
		return fetchFromFile(f.Path)
	default:
		return nil, fmt.Errorf("unsupported live source: %s", f.Source)
	}
}

func fetchFromFile(path string) (map[string]state.Resource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading live state file: %w", err)
	}

	var resources []state.Resource
	if err := json.Unmarshal(data, &resources); err != nil {
		return nil, fmt.Errorf("parsing live state file: %w", err)
	}

	rm := make(map[string]state.Resource, len(resources))
	for _, r := range resources {
		rm[r.ID] = r
	}
	return rm, nil
}
