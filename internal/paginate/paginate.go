// Package paginate provides utilities for paginating slices of drift results.
package paginate

import (
	"errors"

	"github.com/driftctl-lite/internal/drift"
)

// Options controls pagination behaviour.
type Options struct {
	Page     int // 1-based page number
	PageSize int // number of results per page
}

// DefaultOptions returns sensible pagination defaults.
func DefaultOptions() Options {
	return Options{
		Page:     1,
		PageSize: 20,
	}
}

// Page holds a single page of results together with metadata.
type Page struct {
	Results    []drift.Result
	Page       int
	PageSize   int
	TotalItems int
	TotalPages int
	HasNext    bool
	HasPrev    bool
}

// Apply slices results according to opts and returns a Page.
// An error is returned when opts contains invalid values.
func Apply(results []drift.Result, opts Options) (Page, error) {
	if opts.PageSize <= 0 {
		return Page{}, errors.New("paginate: page size must be greater than zero")
	}
	if opts.Page <= 0 {
		return Page{}, errors.New("paginate: page number must be greater than zero")
	}

	total := len(results)
	totalPages := total / opts.PageSize
	if total%opts.PageSize != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	start := (opts.Page - 1) * opts.PageSize
	if start >= total {
		start = total
	}
	end := start + opts.PageSize
	if end > total {
		end = total
	}

	return Page{
		Results:    results[start:end],
		Page:       opts.Page,
		PageSize:   opts.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
		HasNext:    opts.Page < totalPages,
		HasPrev:    opts.Page > 1,
	}, nil
}
