package pokecache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache(1 * time.Second)

	cache.Add("key1", []byte("value1"))

	val, ok := cache.Get("key1")
	if !ok {
		t.Errorf("Expected to find key1 in cache")
	}
	if string(val) != "value1" {
		t.Errorf("Expected value 'value1', got '%s'", string(val))
	}
}

func TestCacheExpiration(t *testing.T) {
	cache := NewCache(1 * time.Second)

	cache.Add("key1", []byte("value1"))

	time.Sleep(2 * time.Second)

	_, ok := cache.Get("key1")
	if ok {
		t.Errorf("Expected key1 to be expired from cache")
	}
}
