package cache

import (
	"sync"

	"github.com/bagmangood/cache/pkg/queue"
)

type Cache interface {
	Read(key string) (interface{}, error)
	Write(key string, value interface{})
	Size() int
}

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

func (c *lruCache) Size() int {
	return len(c.items)
}
