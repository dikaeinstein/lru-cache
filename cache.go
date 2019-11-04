package cache

import (
	"container/list"
	"time"
)

// Cache is a LRU replacement cache
type Cache struct {
	store Store
	ll    *list.List
	size  int
}

// New returns an initialized LRU Cache
func New(size int, store Store) *Cache {
	ll := list.New()

	return &Cache{store: store, ll: ll, size: size}
}

// Node maps key to value
type node struct {
	key       string
	value     string
	expiresAt int
}

// Set inserts the value for the given key in the cache. It overrides old value
// if this keys already exists in the cache
func (c *Cache) Set(key, value string, expiresAt int) {
	current, exists := c.store.Get(key)
	if exists {
		// override existing values
		current.Value.(*node).value = value
		current.Value.(*node).expiresAt = expiresAt

		c.ll.MoveToFront(current)
		return
	}

	// If cache is full remove LRU node(key, value)
	if c.size != 0 && (c.ll.Len() == c.size) {
		lastElement := c.ll.Back()
		c.ll.Remove(lastElement)
		c.store.Delete(lastElement.Value.(*node).key)
	}

	n := &node{key, value, expiresAt}
	e := c.ll.PushFront(n)
	c.store.Set(key, e)
}

// Get returns the value for the given key. It returns the value and true if found
// but false if otherwise
func (c *Cache) Get(key string) (string, bool) {
	element, exists := c.store.Get(key)
	if !exists {
		return "", false
	}

	expire := element.Value.(*node).expiresAt
	if expire == 0 || expire < int(time.Now().Unix()) {
		return "", false
	}

	c.ll.MoveToFront(element)
	return element.Value.(*node).value, true
}

// Delete removes the key-value from the cache.
// It is a noop if no match is found.
func (c *Cache) Delete(key string) {
	element, exists := c.store.Get(key)
	if !exists {
		return
	}

	c.ll.Remove(element)
	c.store.Delete(key)
}

// Size returns the size of the cache
func (c *Cache) Size() int {
	return c.size
}

// Has checks if key exists in cache
func (c *Cache) Has(key string) bool {
	_, exist := c.store.Get(key)

	return exist
}
