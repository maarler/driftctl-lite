package cache_test

import (
	"testing"
	"time"

	"driftctl-lite/internal/cache"
)

// TestCacheTTL_RoundTrip verifies that Set + Get + IsFresh works end-to-end.
func TestCacheTTL_RoundTrip(t *testing.T) {
	c := cache.New(tempDir(t))

	resources := map[string]string{
		"aws_s3_bucket.logs": "present",
	}

	if err := c.Set("integration", resources); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	entry, err := c.Get("integration")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if !cache.IsFresh(entry, 5*time.Minute) {
		t.Error("freshly written entry should be fresh")
	}

	if cache.IsFresh(entry, -1) {
		t.Error("negative TTL should never be fresh")
	}

	if err := c.Invalidate("integration"); err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}

	after, err := c.Get("integration")
	if err != nil {
		t.Fatalf("Get after invalidate failed: %v", err)
	}
	if after != nil {
		t.Error("entry should be nil after invalidation")
	}
}
