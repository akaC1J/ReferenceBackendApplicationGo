package cache

import "sync"

type Node[K comparable, T any] struct {
	key   K
	value T
	prev  *Node[K, T]
	next  *Node[K, T]
}

type LruCache[K comparable, T any] struct {
	capacity int
	cache    map[K]*Node[K, T]
	head     *Node[K, T]
	tail     *Node[K, T]
	mx       sync.Mutex
}

func NewLruCache[K comparable, T any](capacity int) *LruCache[K, T] {
	return &LruCache[K, T]{
		capacity: capacity,
		cache:    make(map[K]*Node[K, T]),
		head:     nil,
		tail:     nil,
		mx:       sync.Mutex{},
	}
}

func (lc *LruCache[K, T]) Get(key K) (T, bool) {
	lc.mx.Lock()
	defer lc.mx.Unlock()
	node, ok := lc.cache[key]
	if !ok {
		var zero T
		return zero, false
	}

	lc.moveToHead(node)
	return node.value, true
}

func (lc *LruCache[K, T]) Put(key K, value T) {
	lc.mx.Lock()
	defer lc.mx.Unlock()
	if lc.capacity == 0 {
		return
	}
	node, ok := lc.cache[key]
	if ok {
		node.value = value
		lc.moveToHead(node)
		return
	}

	if len(lc.cache) >= lc.capacity {
		delete(lc.cache, lc.tail.key)
		lc.removeTail()
	}

	node = &Node[K, T]{key: key, value: value}

	lc.cache[key] = node
	lc.addToHead(node)
}

func (lc *LruCache[K, T]) moveToHead(node *Node[K, T]) {
	lc.removeNode(node)
	lc.addToHead(node)
}

func (lc *LruCache[K, T]) removeNode(node *Node[K, T]) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		lc.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		lc.tail = node.prev
	}
}

func (lc *LruCache[K, T]) addToHead(node *Node[K, T]) {
	node.next = lc.head
	node.prev = nil
	if lc.head != nil {
		lc.head.prev = node
	}
	lc.head = node
	if lc.tail == nil {
		lc.tail = node
	}
}

func (lc *LruCache[K, T]) removeTail() {
	if lc.tail != nil {
		lc.tail = lc.tail.prev
		if lc.tail == nil {
			lc.head = nil
		} else {
			lc.tail.next = nil
		}
	}
}
