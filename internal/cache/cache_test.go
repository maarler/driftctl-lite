package cache_test

import (
	"os"
	"testing"

	"driftctl-lite/internal/cache"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestCache_SetAndGet(t *testing.T) {
	c := cache.New(tempDir(t))
	resources := map[string]string{"aws_s3_bucket.my": "present", "aws_instance.web": "running"}

	if err := c.Set("live", resources); err != nil {
		t.Fatalf("Set: %v", err)
	}

	entry, err := c.Get("live")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if len(entry.Resources) != 2 {
		t.Errorf("expected 2 resources, got %d", len(entry.Resources))
	}
	if entry.FetchedAt.IsZero() {
		t.Error("FetchedAt should not be zero")
	}
}

func TestCache_GetMissing(t *testing.T) {
	c := cache.New(tempDir(t))
	entry, err := c.Get("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Error("expected nil for missing key")
	}
}

func TestCache_Invalidate(t *testing.T) {
	c := cache.New(tempDir(t))
	_ = c.Set("live", map[string]string{"res": "val"})

	if err := c.Invalidate("live"); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}
	entry, _ := c.Get("live")
	if entry != nil {
		t.Error("expected nil after invalidation")
	}
}

func TestCache_InvalidateMissing(t *testing.T) {
	c := cache.New(tempDir(t))
	if err := c.Invalidate("ghost"); err != nil {
		t.Errorf("Invalidate on missing key should not error: %v", err)
	}
}

func TestCache_OverwriteExisting(t *testing.T) {
	c := cache.New(tempDir(t))

	_ = c.Set("live", map[string]string{"aws_s3_bucket.old": "present"})
	_ = c.Set("live", map[string]string{"aws_s3_bucket.new": "present", "aws_instance.web": "running"})

	entry, err := c.Get("live")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if len(entry.Resources) != 2 {
		t.Errorf("expected 2 resources after overwrite, got %d", len(entry.Resources))
	}
	if _, ok := entry.Resources["aws_s3_bucket.old"]; ok {
		t.Error("old resource should not be present after overwrite")
	}
}
