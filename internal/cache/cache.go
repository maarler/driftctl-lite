package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry holds a cached snapshot of live resources.
type Entry struct {
	FetchedAt time.Time          `json:"fetched_at"`
	Resources map[string]string  `json:"resources"`
}

// Cache persists and retrieves live resource snapshots.
type Cache struct {
	dir string
}

// New creates a Cache backed by dir.
func New(dir string) *Cache {
	return &Cache{dir: dir}
}

func (c *Cache) path(key string) string {
	return filepath.Join(c.dir, key+".json")
}

// Set writes an entry to disk.
func (c *Cache) Set(key string, resources map[string]string) error {
	if err := os.MkdirAll(c.dir, 0o755); err != nil {
		return err
	}
	entry := Entry{FetchedAt: time.Now().UTC(), Resources: resources}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path(key), data, 0o644)
}

// Get reads an entry from disk. Returns (nil, nil) when not found.
func (c *Cache) Get(key string) (*Entry, error) {
	data, err := os.ReadFile(c.path(key))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// Invalidate removes a cached entry.
func (c *Cache) Invalidate(key string) error {
	err := os.Remove(c.path(key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
