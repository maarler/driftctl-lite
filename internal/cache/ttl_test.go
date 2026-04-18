package cache_test

import (
	"testing"
	"time"

	"driftctl-lite/internal/cache"
)

func TestIsFresh_NilEntry(t *testing.T) {
	if cache.IsFresh(nil, time.Minute) {
		t.Error("nil entry should not be fresh")
	}
}

func TestIsFresh_ZeroTTL(t *testing.T) {
	entry := &cache.Entry{FetchedAt: time.Now()}
	if cache.IsFresh(entry, 0) {
		t.Error("zero TTL should never be fresh")
	}
}

func TestIsFresh_RecentEntry(t *testing.T) {
	entry := &cache.Entry{FetchedAt: time.Now().Add(-10 * time.Second)}
	if !cache.IsFresh(entry, time.Minute) {
		t.Error("recent entry should be fresh within 1m TTL")
	}
}

func TestIsFresh_StaleEntry(t *testing.T) {
	entry := &cache.Entry{FetchedAt: time.Now().Add(-2 * time.Minute)}
	if cache.IsFresh(entry, time.Minute) {
		t.Error("old entry should not be fresh")
	}
}
