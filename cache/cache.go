package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// Cache is a typical LRU cache implementation. Once the max size is
//reached the least recently used data is evicted
type Cache struct {
	mu sync.Mutex

	// list & table contain *entry objects.
	list  *list.List
	table map[string]*list.Element

	size      int64
	capacity  int64
	evictions int64
}

// Value is the interface values that go into Cache need to satisfy
type Value interface {
	Size() int
}

// the contract for the intem stored in cache
type Item struct {
	Key   string
	Value Value
}

type entry struct {
	key          string
	value        Value
	size         int64
	timeAccessed time.Time
}

func NewCache(capacity int64) *Cache {
	return &Cache{
		list:     list.New(),
		table:    make(map[string]*list.Element),
		capacity: capacity,
	}
}

// Get, will return the value if present and mark the value as used
func (lru *Cache) Get(key string) (v Value, ok bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	element := lru.table[key]
	if element == nil {
		return nil, false
	}
	lru.moveToFront(element)
	return element.Value.(*entry).value, true
}

// returns value but doesnt change the priority of eviction for the entry
func (lru *Cache) Peek(key string) (v Value, ok bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	element := lru.table[key]
	if element == nil {
		return nil, false
	}
	return element.Value.(*entry).value, true
}

func (lru *Cache) Set(key string, value Value) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if element := lru.table[key]; element != nil {
		lru.updateInplace(element, value)
	} else {
		lru.addNew(key, value)
	}
}

func (lru *Cache) SetIfAbsent(key string, value Value) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if element := lru.table[key]; element != nil {
		lru.moveToFront(element)
	} else {
		lru.addNew(key, value)
	}
}

func (lru *Cache) Delete(key string) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	element := lru.table[key]
	if element == nil {
		return false
	}

	lru.list.Remove(element)
	delete(lru.table, key)
	lru.size -= element.Value.(*entry).size
	return true
}

func (lru *Cache) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.list.Init()
	lru.table = make(map[string]*list.Element)
	lru.size = 0
}

func (lru *Cache) SetCapacity(capacity int64) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.capacity = capacity
	lru.checkCapacity()
}

func (lru *Cache) Stats() (length, size, capacity, evictions int64, oldest time.Time) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if lastElem := lru.list.Back(); lastElem != nil {
		oldest = lastElem.Value.(*entry).timeAccessed
	}
	return int64(lru.list.Len()), lru.size, lru.capacity, lru.evictions, oldest
}

func (lru *Cache) StatsJSON() string {
	if lru == nil {
		return "{}"
	}
	l, s, c, e, o := lru.Stats()
	return fmt.Sprintf("{\"Length\": %v, \"Size\": %v, \"Capacity\": %v, \"Evictions\": %v, \"OldestAccess\": \"%v\"}", l, s, c, e, o)
}

func (lru *Cache) Length() int64 {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return int64(lru.list.Len())
}

func (lru *Cache) Size() int64 {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.size
}

func (lru *Cache) Capacity() int64 {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.capacity
}

func (lru *Cache) Evictions() int64 {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.evictions
}

func (lru *Cache) Oldest() (oldest time.Time) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if lastElem := lru.list.Back(); lastElem != nil {
		oldest = lastElem.Value.(*entry).timeAccessed
	}
	return
}

func (lru *Cache) Keys() []string {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	keys := make([]string, 0, lru.list.Len())
	for e := lru.list.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(*entry).key)
	}
	return keys
}

func (lru *Cache) Items() []Item {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	items := make([]Item, 0, lru.list.Len())
	for e := lru.list.Front(); e != nil; e = e.Next() {
		v := e.Value.(*entry)
		items = append(items, Item{Key: v.key, Value: v.value})
	}
	return items
}

func (lru *Cache) updateInplace(element *list.Element, value Value) {
	valueSize := int64(value.Size())
	sizeDiff := valueSize - element.Value.(*entry).size
	element.Value.(*entry).value = value
	element.Value.(*entry).size = valueSize
	lru.size += sizeDiff
	lru.moveToFront(element)
	lru.checkCapacity()
}

func (lru *Cache) moveToFront(element *list.Element) {
	lru.list.MoveToFront(element)
	element.Value.(*entry).timeAccessed = time.Now()
}

func (lru *Cache) addNew(key string, value Value) {
	newEntry := &entry{key, value, int64(value.Size()), time.Now()}
	element := lru.list.PushFront(newEntry)
	lru.table[key] = element
	lru.size += newEntry.size
	lru.checkCapacity()
}

func (lru *Cache) checkCapacity() {
	//TODO refactor method
	for lru.size > lru.capacity {
		delElem := lru.list.Back()
		delValue := delElem.Value.(*entry)
		lru.list.Remove(delElem)
		delete(lru.table, delValue.key)
		lru.size -= delValue.size
		lru.evictions++
	}
}
