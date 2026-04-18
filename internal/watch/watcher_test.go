package watch_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/live"
	"driftctl-lite/internal/watch"
)

func writeTempState(t *testing.T, resources map[string]map[string]string) string {
	t.Helper()
	data, _ := json.Marshal(resources)
	p := filepath.Join(t.TempDir(), "state.json")
	os.WriteFile(p, data, 0644)
	return p
}

func writeTempLive(t *testing.T, resources map[string]map[string]string) string {
	t.Helper()
	data, _ := json.Marshal(resources)
	p := filepath.Join(t.TempDir(), "live.json")
	os.WriteFile(p, data, 0644)
	return p
}

func TestWatcher_CallsHandlerOnTick(t *testing.T) {
	res := map[string]map[string]string{
		"aws_s3_bucket.logs": {"region": "us-east-1"},
	}
	statePath := writeTempState(t, res)
	livePath := writeTempLive(t, res)

	fetcher := live.NewFetcher("file://" + livePath)

	var mu sync.Mutex
	calls := 0
	var lastResults []drift.Result

	w := watch.New(50*time.Millisecond, statePath, fetcher, func(results []drift.Result, err error) {
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		mu.Lock()
		calls++
		lastResults = results
		mu.Unlock()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	mu.Lock()
	defer mu.Unlock()
	if calls < 2 {
		t.Errorf("expected at least 2 handler calls, got %d", calls)
	}
	for _, r := range lastResults {
		if r.Status != "IN_SYNC" {
			t.Errorf("expected IN_SYNC, got %s", r.Status)
		}
	}
}

func TestWatcher_HandlesInvalidState(t *testing.T) {
	fetcher := live.NewFetcher("file:///dev/null")
	var gotErr error
	w := watch.New(50*time.Millisecond, "/nonexistent/state.json", fetcher, func(_ []drift.Result, err error) {
		gotErr = err
	})
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	w.Run(ctx)
	if gotErr == nil {
		t.Error("expected error for missing state file")
	}
}
