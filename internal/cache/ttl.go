package cache

import "time"

// IsFresh reports whether the entry was fetched within the given TTL.
func IsFresh(entry *Entry, ttl time.Duration) bool {
	if entry == nil {
		return false
	}
	if ttl <= 0 {
		return false
	}
	return time.Since(entry.FetchedAt) < ttl
}
