package cache

import (
	"sync"

	"github.com/bagmangood/cache/pkg/queue"
)

// Cache is a generic interface for caches, with simple accessors for read/write and the
// number of elements currently in the cache. Keys are required to be strings, but the
// values can be any tye.
type Cache interface {
	Read(key string) (interface{}, error)
	Write(key string, value interface{})
	Size() int
}

// NewLRU returns a new Cache with the capacity specified. It is thread-safe and obeys
// Least Recently Used when over capacity. LRU in this context means read or write
// operations on a particular key.
func NewLRU(capacity int) Cache {
	return &lruCache{
		items: make(map[string]*item),
		queue: queue.New(capacity),
	}
}

type item struct {
	value interface{}
	node  *queue.Node
}

type lruCache struct {
	items map[string]*item
	queue *queue.Queue
	mutex sync.RWMutex
}

// Read returns the specified value or NotFound if it is not present in the cache.
func (c *lruCache) Read(key string) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	itm, ok := c.items[key]

	if !ok {
		return nil, NotFound
	}

	// bump the value to the end of the queue and return
	c.queue.Touch(itm.node)
	return itm.value, nil
}

// Write inserts the provided value into the cache.
func (c *lruCache) Write(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	newNode, toPurge := c.queue.Add(key)

	if toPurge != "" {
		delete(c.items, toPurge)
	}

	c.items[key] = &item{
		value: value,
		node:  newNode,
	}
}

// Size is the current number of elements stored in the cache.
func (c *lruCache) Size() int {
	return len(c.items)
}
