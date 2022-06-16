package lru

import "container/list"

type Cache struct {
	// max bytes
	maxBytes int64
	// used bytes
	nbytes int64
	ll     *list.List
	cache  map[string]*list.Element
	// executed when an entry is purged
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value use len to count how many bytes it takes
type Value interface {
	Len() int
}

func New(maxBytes int64, nbytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		nbytes:    nbytes,
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	//获取到队列的的队首节点
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}
