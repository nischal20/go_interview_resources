package main

import (
	"container/list"
	"fmt"
)

// Pair holds the key-value pair
type Pair struct {
	key   int
	value int
}

// LRUCache defines the structure of the cache
type LRUCache struct {
	capacity int
	cache    map[int]*list.Element
	list     *list.List
}

// NewLRUCache creates a new LRUCache with the given capacity
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		list:     list.New(),
	}
}

// Get retrieves a value from the cache by key
func (c *LRUCache) Get(key int) (int, bool) {
	if node, ok := c.cache[key]; ok {
		c.list.MoveToFront(node)
		return node.Value.(*Pair).value, true
	}
	return 0, false
}

// Put inserts a key-value pair into the cache
func (c *LRUCache) Put(key, value int) {
	if node, ok := c.cache[key]; ok {
		c.list.MoveToFront(node)
		node.Value.(*Pair).value = value
		return
	}

	if c.list.Len() == c.capacity {
		back := c.list.Back()
		if back != nil {
			pair := back.Value.(*Pair)
			delete(c.cache, pair.key)
			c.list.Remove(back)
		}
	}

	newNode := &Pair{key: key, value: value}
	node := c.list.PushFront(newNode)
	c.cache[key] = node
}

func main() {
	cache := NewLRUCache(3)

	cache.Put(1, 1)
	cache.Put(2, 2)
	fmt.Println(cache.Get(1)) // returns 1
	cache.Put(3, 3)           // evicts key 2
	fmt.Println(cache.Get(2)) // returns 0 (not found)
	fmt.Println(cache.Get(1))
	cache.Put(4, 4)           // evicts key 1
	fmt.Println(cache.Get(1)) // returns 0 (not found)
	fmt.Println(cache.Get(3)) // returns 3
	fmt.Println(cache.Get(4)) // returns 4
}
