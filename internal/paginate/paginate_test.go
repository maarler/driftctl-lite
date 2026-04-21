package paginate_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/paginate"
)

func makeResults(n int) []drift.Result {
	out := make([]drift.Result, n)
	for i := range out {
		out[i] = drift.Result{ResourceID: fmt.Sprintf("res-%d", i)}
	}
	return out
}

func TestApply_FirstPage(t *testing.T) {
	results := makeResults(25)
	page, err := paginate.Apply(results, paginate.Options{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Results) != 10 {
		t.Errorf("expected 10 results, got %d", len(page.Results))
	}
	if page.TotalItems != 25 {
		t.Errorf("expected TotalItems=25, got %d", page.TotalItems)
	}
	if page.TotalPages != 3 {
		t.Errorf("expected TotalPages=3, got %d", page.TotalPages)
	}
	if !page.HasNext {
		t.Error("expected HasNext=true")
	}
	if page.HasPrev {
		t.Error("expected HasPrev=false")
	}
}

func TestApply_LastPage(t *testing.T) {
	results := makeResults(25)
	page, err := paginate.Apply(results, paginate.Options{Page: 3, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Results) != 5 {
		t.Errorf("expected 5 results on last page, got %d", len(page.Results))
	}
	if page.HasNext {
		t.Error("expected HasNext=false on last page")
	}
	if !page.HasPrev {
		t.Error("expected HasPrev=true on last page")
	}
}

func TestApply_PageBeyondTotal(t *testing.T) {
	results := makeResults(5)
	page, err := paginate.Apply(results, paginate.Options{Page: 10, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Results) != 0 {
		t.Errorf("expected 0 results for out-of-range page, got %d", len(page.Results))
	}
}

func TestApply_EmptyResults(t *testing.T) {
	page, err := paginate.Apply([]drift.Result{}, paginate.Options{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.TotalPages != 1 {
		t.Errorf("expected TotalPages=1 for empty input, got %d", page.TotalPages)
	}
	if len(page.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(page.Results))
	}
}

func TestApply_InvalidPageSize(t *testing.T) {
	_, err := paginate.Apply(makeResults(5), paginate.Options{Page: 1, PageSize: 0})
	if err == nil {
		t.Error("expected error for zero page size")
	}
}

func TestApply_InvalidPageNumber(t *testing.T) {
	_, err := paginate.Apply(makeResults(5), paginate.Options{Page: 0, PageSize: 10})
	if err == nil {
		t.Error("expected error for zero page number")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := paginate.DefaultOptions()
	if opts.Page != 1 || opts.PageSize != 20 {
		t.Errorf("unexpected defaults: %+v", opts)
	}
}
