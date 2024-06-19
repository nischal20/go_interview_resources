package main

import (
	"container/heap"
	"container/list"
	"fmt"
	"sync"
)

// CacheItem represents a cache item
type CacheItem struct {
	key   string
	value string
	data  interface{}
}

// EvictionStrategy defines the interface for eviction policies
type EvictionStrategy interface {
	Insert(item *CacheItem)
	Update(item *CacheItem)
	Evict() *CacheItem
}

// LFU Eviction Strategy Implementation

type LFUCacheItem struct {
	key       string
	value     string
	frequency int
	index     int
}

type LFUPriorityQueue []*LFUCacheItem

func (pq LFUPriorityQueue) Len() int { return len(pq) }

func (pq LFUPriorityQueue) Less(i, j int) bool {
	if pq[i].frequency == pq[j].frequency {
		return pq[i].index < pq[j].index
	}
	return pq[i].frequency < pq[j].frequency
}

func (pq LFUPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *LFUPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*LFUCacheItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *LFUPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *LFUPriorityQueue) update(item *LFUCacheItem, value string, frequency int) {
	item.value = value
	item.frequency = frequency
	heap.Fix(pq, item.index)
}

type LFUEvictionStrategy struct {
	mu sync.Mutex
	pq LFUPriorityQueue
}

func NewLFUEvictionStrategy() *LFUEvictionStrategy {
	pq := make(LFUPriorityQueue, 0)
	heap.Init(&pq)
	return &LFUEvictionStrategy{
		pq: pq,
	}
}

func (s *LFUEvictionStrategy) Insert(item *CacheItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	lfuItem := &LFUCacheItem{
		key:       item.key,
		value:     item.value,
		frequency: 1,
	}
	item.data = lfuItem
	heap.Push(&s.pq, lfuItem)
}

func (s *LFUEvictionStrategy) Update(item *CacheItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	lfuItem := item.data.(*LFUCacheItem)
	lfuItem.frequency++
	heap.Fix(&s.pq, lfuItem.index)
}

func (s *LFUEvictionStrategy) Evict() *CacheItem {
	s.mu.Lock()
	defer s.mu.Unlock()
	evictedItem := heap.Pop(&s.pq).(*LFUCacheItem)
	return &CacheItem{
		key:   evictedItem.key,
		value: evictedItem.value,
		data:  evictedItem,
	}
}

// LRU Eviction Strategy Implementation

type LRUCacheItem struct {
	key   string
	value string
}

type LRUEvictionStrategy struct {
	mu    sync.Mutex
	ll    *list.List
	items map[string]*list.Element
}

func NewLRUEvictionStrategy() *LRUEvictionStrategy {
	return &LRUEvictionStrategy{
		ll:    list.New(),
		items: make(map[string]*list.Element),
	}
}

func (s *LRUEvictionStrategy) Insert(item *CacheItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	lruItem := &LRUCacheItem{
		key:   item.key,
		value: item.value,
	}
	element := s.ll.PushFront(lruItem)
	item.data = element
	s.items[item.key] = element
}

func (s *LRUEvictionStrategy) Update(item *CacheItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	element := item.data.(*list.Element)
	s.ll.MoveToFront(element)
}

func (s *LRUEvictionStrategy) Evict() *CacheItem {
	s.mu.Lock()
	defer s.mu.Unlock()
	element := s.ll.Back()
	if element == nil {
		return nil
	}
	s.ll.Remove(element)
	lruItem := element.Value.(*LRUCacheItem)
	delete(s.items, lruItem.key)
	return &CacheItem{
		key:   lruItem.key,
		value: lruItem.value,
		data:  element,
	}
}

// Cache Struct

type Cache struct {
	mu       sync.RWMutex
	capacity int
	items    map[string]*CacheItem
	count    int
	strategy EvictionStrategy
}

func NewCache(capacity int, strategy EvictionStrategy) *Cache {
	return &Cache{
		capacity: capacity,
		items:    make(map[string]*CacheItem),
		strategy: strategy,
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if item, found := c.items[key]; found {
		c.strategy.Update(item)
		return item.value, true
	}
	return "", false
}

func (c *Cache) Put(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.capacity == 0 {
		return
	}

	if item, found := c.items[key]; found {
		item.value = value
		c.strategy.Update(item)
		return
	}

	if c.count == c.capacity {
		evictedItem := c.strategy.Evict()
		delete(c.items, evictedItem.key)
		c.count--
	}

	newItem := &CacheItem{
		key:   key,
		value: value,
		data:  nil,
	}
	c.strategy.Insert(newItem)
	c.items[key] = newItem
	c.count++
}

// Usage Example

func main() {
	fmt.Println("Using LFU Cache")
	lfuStrategy := NewLFUEvictionStrategy()
	lfuCache := NewCache(2, lfuStrategy)

	lfuCache.Put("a", "1")
	lfuCache.Put("b", "2")
	fmt.Println(lfuCache.Get("a")) // Output: 1 true

	lfuCache.Put("c", "3")         // Evicts key "b"
	fmt.Println(lfuCache.Get("b")) // Output: "" false

	lfuCache.Put("d", "4")         // Evicts key "a"
	fmt.Println(lfuCache.Get("a")) // Output: "" false
	fmt.Println(lfuCache.Get("c")) // Output: 3 true
	fmt.Println(lfuCache.Get("d")) // Output: 4 true

	fmt.Println("Using LRU Cache")
	lruStrategy := NewLRUEvictionStrategy()
	lruCache := NewCache(2, lruStrategy)

	lruCache.Put("a", "1")
	lruCache.Put("b", "2")
	fmt.Println(lruCache.Get("a")) // Output: 1 true

	lruCache.Put("c", "3")         // Evicts key "b"
	fmt.Println(lruCache.Get("b")) // Output: "" false

	lruCache.Put("d", "4")         // Evicts key "a"
	fmt.Println(lruCache.Get("a")) // Output: "" false
	fmt.Println(lruCache.Get("c")) // Output: 3 true
	fmt.Println(lruCache.Get("d")) // Output: 4 true
}
