package structure

type valueNode[K comparable, V any] struct {
	value      V
	prev, next *K
}

type OrderedMap[K comparable, V any] struct {
	data        map[K]valueNode[K, V]
	first, last *K
}

func NewOrderedMap[K comparable, V any](cap int) *OrderedMap[K, V] {
	data := make(map[K]valueNode[K, V], cap)
	return &OrderedMap[K, V]{data, nil, nil}
}

func (m *OrderedMap[K, V]) Size() int {
	return len(m.data)
}

func (m *OrderedMap[K, V]) Get(key K) V {
	return m.data[key].value
}

func (m *OrderedMap[K, V]) Has(key K) bool {
	_, ok := m.data[key]

	return ok
}

func (m *OrderedMap[K, V]) Set(key K, value V) {
	if node, ok := m.data[key]; ok {
		// Key already exists, only update its value
		m.data[key] = valueNode[K, V]{value, node.prev, node.next}
		return
	}

	node := valueNode[K, V]{value, m.last, nil}
	m.data[key] = node

	// Update last node
	if m.last != nil {
		lastNode := m.data[*m.last]
		m.data[*m.last] = valueNode[K, V]{lastNode.value, lastNode.prev, &key}
	}
	m.last = &key

	// If this is the first element inserted, update first as well
	if m.first == nil {
		m.first = &key
	}
}

func (m *OrderedMap[K, V]) Delete(key K) {
	node, ok := m.data[key]

	if !ok {
		return
	}

	delete(m.data, key)
	prev, next := node.prev, node.next
	if prev != nil {
		// Set the previous node's `next` to the deleted node's `next`
		prevNode := m.data[*prev]
		m.data[*prev] = valueNode[K, V]{prevNode.value, prevNode.prev, next}
	} else {
		// the deleted node was the first one
		m.first = node.next
	}
	if next != nil {
		// Set the next node's `prev` to the deleted node's `prev`
		nextNode := m.data[*next]
		m.data[*next] = valueNode[K, V]{nextNode.value, prev, nextNode.next}
	} else {
		// the deleted node was the last one
		m.last = node.prev
	}
}

func (m *OrderedMap[K, V]) IterateEntries(fn func(K, V)) {

	for k := m.first; k != nil; {
		node := m.data[*k]
		fn(*k, node.value)
		k = node.next
	}
}

func (m *OrderedMap[K, V]) IterateKeys(fn func(K)) {
	m.IterateEntries(func(k K, v V) { fn(k) })
}

func (m *OrderedMap[K, V]) IterateValues(fn func(V)) {
	m.IterateEntries(func(k K, v V) { fn(v) })
}

type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

func (m *OrderedMap[K, V]) Entries() []MapEntry[K, V] {
	entries := make([]MapEntry[K, V], m.Size())

	for k, i := m.first, 0; k != nil; i++ {
		node := m.data[*k]

		entries[i] = MapEntry[K, V]{*k, node.value}

		k = node.next
	}

	return entries
}

func (m *OrderedMap[K, V]) Keys() []K {
	keys := make([]K, m.Size())

	for k, i := m.first, 0; k != nil; i++ {
		node := m.data[*k]

		keys[i] = *k

		k = node.next
	}

	return keys
}

func (m *OrderedMap[K, V]) Values() []V {
	values := make([]V, m.Size())

	for k, i := m.first, 0; k != nil; i++ {
		node := m.data[*k]

		values[i] = node.value

		k = node.next
	}

	return values
}
