package common

import "sync"

// --- ram cache interface ---
type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Del(key string)
}

// --- LRU Cache ---
type LinkNode struct {
	key       string
	value     interface{}
	pre, next *LinkNode
}

type LRUCache struct {
	cap        uint32
	data       map[string]*LinkNode
	head, tail *LinkNode
	lock       sync.Mutex
}

func (l *LRUCache) removeNode(node *LinkNode) {
	node.pre.next = node.next
	node.next.pre = node.pre
}

func (l *LRUCache) addNode(node *LinkNode) {
	node.next = l.head.next
	node.next.pre = node
	node.pre = l.head
	l.head.next = node
}

func (l *LRUCache) moveToHead(node *LinkNode) {
	l.removeNode(node)
	l.addNode(node)
}

func (l *LRUCache) Get(key string) interface{} {
	l.lock.Lock()
	defer l.lock.Unlock()

	if n, exists := l.data[key]; exists {
		l.moveToHead(n)
		return n.value
	} else {
		return nil
	} // else>
}

func (l *LRUCache) Set(key string, value interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if n, exist := l.data[key]; exist {
		n.value = value
		l.moveToHead(n)
	} else {
		newNode := &LinkNode{
			key:   key,
			value: value,
			pre:   nil,
			next:  nil,
		} // node
		if uint32(len(l.data)) == l.cap {
			delete(l.data, l.tail.pre.key)
			l.removeNode(l.tail.pre)
		} // if>>
		l.addNode(newNode)
		l.data[key] = newNode
	} // else>
}

func (l *LRUCache) Del(key string) {

}

func NewLRUCache(cap uint32) *LRUCache {
	head := &LinkNode{"", nil, nil, nil}
	tail := &LinkNode{"", nil, nil, nil}
	head.next = tail
	tail.pre = head
	return &LRUCache{
		cap:  cap,
		data: make(map[string]*LinkNode),
		head: head,
		tail: tail,
		lock: sync.Mutex{},
	}
}

func NewDefaultLRUCache() Cache {
	const defaultRAMCacheCapacity = 256
	return NewLRUCache(defaultRAMCacheCapacity)
}
