package custom

import (
	"sync"
)

// CacheItem represents an item in the cache.
type CacheItem struct {
    Value     interface{}
}

// Cache is a simple in-memory cache with an expiry time.
type Cache struct {
    mu    sync.Mutex
    items map[string]CacheItem
}

// NewCache creates a new Cache instance.
func NewCache() *Cache {
    return &Cache{
        items: make(map[string]CacheItem),
    }
}

// Set adds an item to the cache with the specified key, value, and expiry time.
func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.items[key] = CacheItem{
        Value:     value,
    }
}

// Get retrieves an item from the cache by its key. If the item is expired or not found, it returns nil.
func (c *Cache) Get(key string) interface{} {
    c.mu.Lock()
    defer c.mu.Unlock()

    item, found := c.items[key]
    if !found {
        return nil
    }

    return item.Value
}

// Remove removes an item from the cache by its key.
func (c *Cache) Remove(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    delete(c.items, key)
}

// Clear removes all items from the cache.
func (c *Cache) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.items = make(map[string]CacheItem)
}
