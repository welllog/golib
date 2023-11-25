package mapz

type KV[K comparable, V any] map[K]V

func (m KV[K, V]) Get(key K) (V, bool) {
	value, ok := m[key]
	return value, ok
}

func (m KV[K, V]) Set(key K, value V) {
	m[key] = value
}

func (m KV[K, V]) SetNx(key K, value V) bool {
	_, ok := m[key]
	if !ok {
		m[key] = value
	}
	return !ok
}

func (m KV[K, V]) SetX(key K, value V) bool {
	_, ok := m[key]
	if ok {
		m[key] = value
	}
	return ok
}

func (m KV[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		delete(m, key)
	}
}

func (m KV[K, V]) Has(key K) bool {
	_, ok := m[key]
	return ok
}

func (m KV[K, V]) Len() int {
	return len(m)
}

func (m KV[K, V]) Keys() []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m KV[K, V]) Values() []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func (m KV[K, V]) Range(fn func(key K, value V) bool) {
	for k, v := range m {
		if !fn(k, v) {
			break
		}
	}
}
