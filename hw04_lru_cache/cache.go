package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, v interface{}) bool {
	/*
		если элемент присутствует в словаре, то обновить его значение и переместить элемент в начало очереди;
		если элемента нет в словаре, то добавить в словарь и в начало очереди (при этом, если размер очереди
		больше ёмкости кэша, то необходимо удалить последний элемент из очереди и его значение из словаря);
		возвращаемое значение - флаг, присутствовал ли элемент в кэше.
	*/
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.items[key]; ok {
		l.queue.MoveToFront(l.items[key])
		l.queue.Front().Value = &cacheItem{key: key, value: v}
		l.items[key].Value = l.queue.Front().Value
		return true
	}

	if l.queue.Len() == l.capacity {
		delete(l.items, l.queue.Back().Value.(*cacheItem).key)
		l.queue.Remove(l.queue.Back())
	}

	l.items[key] = l.queue.PushFront(&cacheItem{key: key, value: v})

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	/*
		если элемент присутствует в словаре, то переместить элемент в начало очереди и вернуть его значение и true;
		если элемента нет в словаре, то вернуть nil и false (работа с кешом похожа на работу с map)
	*/
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.items[key]; ok {
		l.queue.MoveToFront(l.items[key])
		l.items[key] = l.queue.Front()
		res := l.items[key].Value.(*cacheItem).value
		return res, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	// l = &lruCache{}
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}
