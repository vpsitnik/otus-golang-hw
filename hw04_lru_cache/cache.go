package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type itemCache struct {
	key   Key
	value interface{}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	if cache.queue.Len() >= cache.capacity {
		// clear last
		last := cache.queue.Back()
		cache.queue.Remove(last)
		delete(cache.items, last.Value.(itemCache).key)
	}

	if item, ok := cache.items[key]; ok {
		// exist
		item.Value = itemCache{key, value}
		cache.queue.MoveToFront(item)
		return true
	}

	// add new
	cache.items[key] = cache.queue.PushFront(itemCache{key, value})
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := cache.items[key]; ok {
		return item.Value.(itemCache).value, true
	}
	return nil, false
}

func (cache *lruCache) Clear() {
	*cache = *(NewCache(cache.capacity).(*lruCache))
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
