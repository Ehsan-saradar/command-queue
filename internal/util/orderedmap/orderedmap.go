package orderedmap

import (
	"sort"
	"sync"
)

type node struct {
	key       string
	value     interface{}
	timeStamp int64
	prev      *node
	next      *node
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

func (om *OrderedMap) Set(key string, value interface{}, timestamp int64) {
	newNode := &node{
		key:       key,
		value:     value,
		timeStamp: timestamp,
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
	nodes := om.getAllNodes()
	// sort nodes by timestamp
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].timeStamp < nodes[j].timeStamp
	})
	keys := make([]string, 0, len(nodes))
	values := make([]interface{}, 0, len(nodes))
	for _, n := range nodes {
		keys = append(keys, n.key)
		values = append(values, n.value)
	}
	return keys, values
}

func (om *OrderedMap) getAllNodes() []node {
	om.mutex.RLock()
	defer om.mutex.RUnlock()
	nodes := make([]node, 0, len(om.values))
	for n := om.head; n != nil; n = n.next {
		nodes = append(nodes, node{
			key:       n.key,
			value:     n.value,
			timeStamp: n.timeStamp,
		})
	}
	return nodes
}
