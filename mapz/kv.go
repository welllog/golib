package mapz

type KV[K comparable, V any] map[K]V

// Get returns the value associated with the key and whether the key existed or not.
func (m KV[K, V]) Get(key K) (V, bool) {
	value, ok := m[key]
	return value, ok
}

// Set sets the value associated with the key.
func (m KV[K, V]) Set(key K, value V) {
	m[key] = value
}

// SetNx sets the value associated with the key if the key does not exist.
func (m KV[K, V]) SetNx(key K, value V) bool {
	_, ok := m[key]
	if !ok {
		m[key] = value
	}
	return !ok
}

// SetX sets the value associated with the key if the key exists.
func (m KV[K, V]) SetX(key K, value V) bool {
	_, ok := m[key]
	if ok {
		m[key] = value
	}
	return ok
}

// Delete deletes the value associated with the key.
func (m KV[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		delete(m, key)
	}
}

// Has returns whether the key exists.
func (m KV[K, V]) Has(key K) bool {
	_, ok := m[key]
	return ok
}

// Contains returns whether the key exists.
func (m KV[K, V]) Contains(key K) bool {
	_, ok := m[key]
	return ok
}

// Len returns the number of items.
func (m KV[K, V]) Len() int {
	return len(m)
}

// Keys returns all the keys.
func (m KV[K, V]) Keys() []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all the values.
func (m KV[K, V]) Values() []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Range calls fn sequentially for each key and value present in the map.
func (m KV[K, V]) Range(fn func(key K, value V) bool) {
	for k, v := range m {
		if !fn(k, v) {
			break
		}
	}
}
