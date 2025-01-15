package hw04lrucache

import (
	"sync"
)

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
	mutex    sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mutex:    sync.Mutex{},
	}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	newCacheItem := &cacheItem{key, value}

	if item, ok := cache.items[key]; ok {
		item.Value = newCacheItem
		cache.queue.MoveToFront(item)
		return true
	}

	if cache.capacity == cache.queue.Len() {
		lastItem := cache.queue.Back()
		cache.queue.Remove(lastItem)
		delete(cache.items, lastItem.Value.(*cacheItem).key)
	}

	newItem := cache.queue.PushFront(newCacheItem)
	cache.items[key] = newItem

	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	item, ok := cache.items[key]
	if !ok {
		return nil, false
	}

	cache.queue.MoveToFront(item)

	return item.Value.(*cacheItem).value, true
}

func (cache *lruCache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}
