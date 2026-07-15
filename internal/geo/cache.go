package geo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Cache persists geocoding results to a JSON file so repeated runs don't
// re-query the geocoding service for locations already resolved. It is safe
// for concurrent use.
type Cache struct {
	mu   sync.RWMutex
	path string
	data map[string]Coordinates
}

// NewCache loads the cache from path if it exists, or returns an empty cache
// ready to be populated if the file is not found. Any other read or decode
// error is returned to the caller.
func NewCache(path string) (*Cache, error) {
	c := &Cache{
		path: path,
		data: make(map[string]Coordinates),
	}

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return c, nil
		}
		return nil, fmt.Errorf("open cache file: %w", err)
	}
	defer file.Close() //nolint:errcheck // deferred close, error unrecoverable

	if err := json.NewDecoder(file).Decode(&c.data); err != nil {
		return nil, fmt.Errorf("decode cache file: %w", err)
	}

	return c, nil
}

// Get returns the cached Coordinates for location and true if present.
func (c *Cache) Get(location string) (Coordinates, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	coords, ok := c.data[location]
	return coords, ok
}

// Set stores coords for location, overwriting any existing entry.
func (c *Cache) Set(location string, coords Coordinates) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[location] = coords
}

// Save writes the current cache contents to disk as indented JSON,
// creating the parent directory if it does not already exist.
func (c *Cache) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if dir := filepath.Dir(c.path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create cache directory: %w", err)
		}
	}

	file, err := os.Create(c.path)
	if err != nil {
		return fmt.Errorf("create cache file: %w", err)
	}
	defer file.Close() //nolint:errcheck // deferred close, error unrecoverable

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(c.data); err != nil {
		return fmt.Errorf("encode cache: %w", err)
	}

	return nil
}
