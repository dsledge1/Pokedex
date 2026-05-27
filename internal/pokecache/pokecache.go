package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	storage map[string]cacheEntry
	hold    sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		storage: make(map[string]cacheEntry),
		hold:    sync.Mutex{},
	}

	go c.reapLoop(interval)

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.hold.Lock()
	defer c.hold.Unlock()
	c.storage[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {

	c.hold.Lock()

	defer c.hold.Unlock()
	entry, ok := c.storage[key]

	if !ok {

		return []byte{}, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		c.hold.Lock()
		for key, value := range c.storage {

			if time.Since(value.createdAt) > interval {

				delete(c.storage, key)
			}
		}
		c.hold.Unlock()
	}

}
