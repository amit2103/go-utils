package cache

import (
	"testing"
)

type CacheValue struct {
	size int
}

func (cv *CacheValue) Size() int {
	return cv.size
}

func TestSetInsertsValue(t *testing.T) {
	cache := NewCache(100)
	data := &CacheValue{0}
	key := "key"
	cache.Set(key, data)

	v, ok := cache.Get(key)
	if !ok || v.(*CacheValue) != data {
		t.Errorf("Cache has incorrect value: %v != %v", data, v)
	}

	k := cache.Keys()
	if len(k) != 1 || k[0] != key {
		t.Errorf("Cache.Keys() returned incorrect values: %v", k)
	}
	values := cache.Items()
	if len(values) != 1 || values[0].Key != key {
		t.Errorf("Cache.Values() returned incorrect values: %v", values)
	}
}

func TestSetIfAbsent(t *testing.T) {
	cache := NewCache(100)
	data := &CacheValue{0}
	key := "key"
	cache.SetIfAbsent(key, data)

	v, ok := cache.Get(key)
	if !ok || v.(*CacheValue) != data {
		t.Errorf("Cache has incorrect value: %v != %v", data, v)
	}

	cache.SetIfAbsent(key, &CacheValue{1})

	v, ok = cache.Get(key)
	if !ok || v.(*CacheValue) != data {
		t.Errorf("Cache has incorrect value: %v != %v", data, v)
	}
}

func TestGetValueWithMultipleTypes(t *testing.T) {
	cache := NewCache(100)
	data := &CacheValue{0}
	key := "key"
	cache.Set(key, data)

	v, ok := cache.Get("key")
	if !ok || v.(*CacheValue) != data {
		t.Errorf("Cache has incorrect value for \"key\": %v != %v", data, v)
	}

	v, ok = cache.Get(string([]byte{'k', 'e', 'y'}))
	if !ok || v.(*CacheValue) != data {
		t.Errorf("Cache has incorrect value for []byte {'k','e','y'}: %v != %v", data, v)
	}
}
