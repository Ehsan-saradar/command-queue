package orderedmap

import (
	"sync"
)

type node struct {
	key   string
	value interface{}
	prev  *node
	next  *node
}

type OrderedMap struct {
	head   *node
	tail   *node
	values map[string]*node
	mutex  sync.RWMutex // Mutex for concurrent access
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		values: make(map[string]*node),
		mutex:  sync.RWMutex{},
	}
}

func (om *OrderedMap) Set(key string, value interface{}) {
	newNode := &node{
		key:   key,
		value: value,
	}
	om.mutex.Lock()
	defer om.mutex.Unlock()
	om.values[key] = newNode

	if om.head == nil {
		om.head = newNode
		om.tail = newNode
	} else {
		om.tail.next = newNode
		newNode.prev = om.tail
		om.tail = newNode
	}
}

func (om *OrderedMap) Get(key string) (interface{}, bool) {
	om.mutex.RLock()
	defer om.mutex.RUnlock()

	if n, ok := om.values[key]; ok {
		return n.value, true
	}
	return nil, false
}

func (om *OrderedMap) DeleteItem(key string) {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	if n, ok := om.values[key]; ok {
		if n.prev != nil {
			n.prev.next = n.next
		} else {
			om.head = n.next
		}
		if n.next != nil {
			n.next.prev = n.prev
		} else {
			om.tail = n.prev
		}
		delete(om.values, key)
	}
}

func (om *OrderedMap) keys() []string {
	om.mutex.RLock()
	defer om.mutex.RUnlock()

	keys := make([]string, 0, len(om.values))
	for n := om.head; n != nil; n = n.next {
		keys = append(keys, n.key)
	}
	return keys
}

func (om *OrderedMap) GetAll() ([]string, []interface{}) {
	om.mutex.RLock()
	defer om.mutex.RUnlock()
	keys := make([]string, 0, len(om.values))
	values := make([]interface{}, 0, len(om.values))
	for n := om.head; n != nil; n = n.next {
		keys = append(keys, n.key)
		values = append(values, n.value)
	}
	return keys, values
}
