package main

type OrderedMap struct {
	Map map[interface{}]interface{}
	// Map     map[string]map[string]*task
	order []interface{}
	key   interface{}
	index int
	size  int
}

func NewOrderedMap() *OrderedMap {
	m := &OrderedMap{
		Map:   map[interface{}]interface{}{},
		order: []interface{}{},
		index: 0,
		size:  0,
	}

	return m
}

func (m *OrderedMap) NewIter() (interface{}, interface{}, func() (interface{}, interface{})) {
	key, val := m.key, m.Value(m.key)

	next := m.Next

	return key, val, next
}

func (m *OrderedMap) Add(key interface{}, value interface{}) {
	m.order = append(m.order, key)
	m.key = m.order[0]
	m.Map[key] = value
}

func (m *OrderedMap) Value(key interface{}) interface{} {
	if m.Contains(key) {
		return m.Map[key]
	}
	return nil
}

func (m *OrderedMap) Contains(key interface{}) bool {
	_, ok := m.Map[key]
	if ok {
		return true
	}

	return false
}

func (m *OrderedMap) Next() (interface{}, interface{}) {
	m.index++
	if m.index > m.size {
		m.key = m.order[0]

		return nil, nil
	}

	m.key = m.order[m.index]

	return m.key, m.Map[m.key]
}
