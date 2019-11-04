package cache

import (
	"container/list"
	"sync"
)

// Store is the underlying hash map used for fast access
// and deletion for the LRU Cache
type Store interface {
	Get(k string) (*list.Element, bool)
	Set(k string, e *list.Element)
	Delete(k string)
}

// InMemoryStore is an in-memory implementation of Store.
// Which is save for concurrent use
type InMemoryStore struct {
	sync.RWMutex
	store map[string]*list.Element
}

// NewInMemoryStore creates a new in-memory store
func NewInMemoryStore() Store {
	store := make(map[string]*list.Element)
	return &InMemoryStore{store: store}
}

// Get returns the value for the given key from the store
func (i *InMemoryStore) Get(key string) (*list.Element, bool) {
	i.RLock()
	defer i.RUnlock()

	e, exist := i.store[key]
	return e, exist
}

// Set saves the given value for the key in the store
func (i *InMemoryStore) Set(key string, value *list.Element) {
	i.Lock()
	i.store[key] = value
	i.Unlock()
}

// Delete removes the given value for the key from the store
func (i *InMemoryStore) Delete(key string) {
	i.Lock()
	delete(i.store, key)
	i.Unlock()
}
