package cache

import (
	"testing"
	"time"
)

func TestCacheSetAddNewValue(t *testing.T) {
	store := NewInMemoryStore()
	c := New(2, store)

	expiresAt := time.Now().Add(10 * time.Second).Unix()
	c.Set("name", "Dika", int(expiresAt))

	v, exist := c.Get("name")
	if !exist || v != "Dika" {
		t.Errorf("expected: %s, got: %s", "Dika", v)
	}
}

func TestCacheSetOverridesExistingKey(t *testing.T) {
	store := NewInMemoryStore()
	c := New(1, store)

	expiresAt := time.Now().Add(10 * time.Second).Unix()
	c.Set("name", "Dika", int(expiresAt))
	c.Set("name", "UpdatedDika", int(expiresAt))

	v, exist := c.Get("name")
	if !exist || v != "UpdatedDika" {
		t.Errorf("expected: %s, got: %s", "UpdatedDika", v)
	}
}

func TestCacheSetRemovesLRUValue(t *testing.T) {
	store := NewInMemoryStore()
	c := New(1, store)

	expiresAt := time.Now().Add(10 * time.Second).Unix()
	c.Set("name", "Dika", int(expiresAt))
	c.Set("email", "dika@mail.com", int(expiresAt))

	if _, exist := c.Get("name"); exist {
		t.Errorf("name should have been removed")
	}

	email, exist := c.Get("email")
	if !exist || email != "dika@mail.com" {
		t.Errorf("expected: %s, got: %s", "dika@mail.com", email)
	}
}

func TestCacheSetSafeForConcurrentUse(t *testing.T) {
	store := NewInMemoryStore()
	c := New(1, store)

	expiresAt := time.Now().Add(10 * time.Second).Unix()
	go c.Set("name", "Dika", int(expiresAt))
	go c.Set("email", "dika@mail.com", int(expiresAt))

	if c.Size() != 1 {
		t.Errorf("cache size should be %d, got %d. Only one 'set' operation should proceed",
			1, c.Size())
	}
}

func TestCacheGet(t *testing.T) {
	store := NewInMemoryStore()
	c := New(2, store)

	expiresAt := time.Now().Add(10 * time.Second).Unix()
	c.Set("name", "Dika", int(expiresAt))

	v, exist := c.Get("name")
	if !exist || v != "Dika" {
		t.Errorf("expected: %s, got: %s", "Dika", v)
	}
}

func TestCacheGetReturnsFalseWhenKeyNotFound(t *testing.T) {
	store := NewInMemoryStore()
	c := New(2, store)

	_, exist := c.Get("name")
	if exist {
		t.Errorf("expected: %T, got: %T", false, exist)
	}
}

func TestCacheGetReturnsFalseWhenKeyHasExpired(t *testing.T) {
	store := NewInMemoryStore()
	c := New(2, store)

	c.Set("name", "Dika", 0)

	_, exist := c.Get("name")
	if exist {
		t.Errorf("Key should have expired. expected: %T, got: %T",
			false, exist)
	}
}

func TestDelete(t *testing.T) {
	store := NewInMemoryStore()
	c := New(2, store)

	expiresAt := time.Now().Add(10 * time.Second).Unix()
	c.Set("name", "Dika", int(expiresAt))

	c.Delete("name")

	if _, exist := c.Get("name"); exist {
		t.Errorf("Key should have been deleted. expected: %T, got: %T",
			false, exist)
	}
}

func TestDeleteNoop(t *testing.T) {
	store := NewInMemoryStore()
	c := New(1, store)

	c.Delete("name")

	if _, exist := c.Get("name"); exist {
		t.Errorf("Key should have been deleted. expected: %T, got: %T",
			false, exist)
	}
}

func TestHas(t *testing.T) {
	store := NewInMemoryStore()
	c := New(1, store)

	expiresAt := time.Now().Add(10 * time.Second).Unix()
	c.Set("name", "Dika", int(expiresAt))

	exist := c.Has("name")
	if !exist {
		t.Errorf("expected %T, got %T", true, exist)
	}
}

func TestHasReturnsFalseNotFound(t *testing.T) {
	store := NewInMemoryStore()
	c := New(1, store)

	exist := c.Has("name")
	if exist {
		t.Errorf("expected %T, got %T", false, exist)
	}
}
