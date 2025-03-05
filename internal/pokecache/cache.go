package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

type Cache struct {
	mx            sync.Mutex
	CacheEntryMap map[string]CacheEntry
	interval      time.Duration
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		CacheEntryMap: make(map[string]CacheEntry),
		interval:      interval,
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.CacheEntryMap[key] = CacheEntry{
		CreatedAt: time.Now(),
		Val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	entry, ok := c.CacheEntryMap[key]
	return entry.Val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mx.Lock()
		now := time.Now()
		for key, entry := range c.CacheEntryMap {
			if now.Sub(entry.CreatedAt) > c.interval {
				delete(c.CacheEntryMap, key)
			}
		}
		c.mx.Unlock()
	}
}
