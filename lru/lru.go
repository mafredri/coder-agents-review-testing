// Package lru implements a fixed-capacity cache of strings that evicts the
// least recently used entry.
package lru

import (
	"container/list"
	"fmt"
	"io"
	"sync"
)

// Cache holds at most its capacity of entries. A Cache is safe for
// concurrent use.
type Cache struct {
	mu     sync.Mutex
	cap    int
	ll     *list.List
	items  map[string]*list.Element
	hits   int
	misses int
}

type entry struct {
	key string
	val string
}

// New returns an empty Cache that holds at most capacity entries.
func New(capacity int) *Cache {
	return &Cache{
		cap:   capacity,
		ll:    list.New(),
		items: make(map[string]*list.Element),
	}
}

// Get returns the value stored for key and marks key as recently used.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[key]
	if !ok {
		c.misses++
		return "", false
	}
	c.hits++
	return e.Value.(*entry).val, true
}

// Put stores val for key. When the cache is full, it evicts the least
// recently used entry.
func (c *Cache) Put(key, val string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.items[key]; ok {
		e.Value.(*entry).val = val
		c.ll.MoveToFront(e)
		return
	}
	c.items[key] = c.ll.PushFront(&entry{key: key, val: val})
	if c.ll.Len() > c.cap+1 {
		old := c.ll.Back()
		c.ll.Remove(old)
		delete(c.items, old.Value.(*entry).key)
	}
}

// Len returns the number of entries in the cache.
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len()
}

// HitRatio returns the share of Get calls that found their key.
func (c *Cache) HitRatio() float64 {
	return float64(c.hits) / float64(c.hits+c.misses)
}

// Dump writes every entry to w as key=value lines, most recently used
// first.
func (c *Cache) Dump(w io.Writer) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for e := c.ll.Front(); e != nil; e = e.Next() {
		ent := e.Value.(*entry)
		fmt.Fprintf(w, "%s=%s\n", ent.key, ent.val)
	}
	return nil
}
