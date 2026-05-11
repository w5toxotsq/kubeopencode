package analyzer

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// CacheEntry holds a cached analysis result with an expiration time.
type CacheEntry struct {
	Result    string
	CreatedAt time.Time
	TTL       time.Duration
}

// IsExpired returns true if the cache entry has exceeded its TTL.
func (e *CacheEntry) IsExpired() bool {
	return time.Since(e.CreatedAt) > e.TTL
}

// AnalysisCache provides a simple in-memory cache for analysis results.
type AnalysisCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// NewAnalysisCache creates a new AnalysisCache with the given TTL.
func NewAnalysisCache(ttl time.Duration) *AnalysisCache {
	return &AnalysisCache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
}

// cacheKey generates a deterministic key from resource kind, namespace, and name.
func cacheKey(kind, namespace, name string) string {
	raw := fmt.Sprintf("%s/%s/%s", kind, namespace, name)
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", sum)
}

// Get retrieves a cached result. Returns the result and true if found and not expired.
func (c *AnalysisCache) Get(kind, namespace, name string) (string, bool) {
	key := cacheKey(kind, namespace, name)
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()

	if !ok || entry.IsExpired() {
		return "", false
	}
	return entry.Result, true
}

// Set stores a result in the cache.
func (c *AnalysisCache) Set(kind, namespace, name, result string) {
	key := cacheKey(kind, namespace, name)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = &CacheEntry{
		Result:    result,
		CreatedAt: time.Now(),
		TTL:       c.ttl,
	}
}

// Invalidate removes a specific entry from the cache.
func (c *AnalysisCache) Invalidate(kind, namespace, name string) {
	key := cacheKey(kind, namespace, name)
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Flush removes all entries from the cache.
func (c *AnalysisCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}

// Size returns the number of entries currently in the cache.
func (c *AnalysisCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
