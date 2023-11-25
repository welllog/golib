package mapz

import "sync"

// SafeKV is a thread-safe map
type SafeKV[K comparable, V any] struct {
	entries KV[K, V]
	mu      sync.RWMutex
}

// NewSafeKV creates a new SafeKV
func NewSafeKV[K comparable, V any](cap int) *SafeKV[K, V] {
	return &SafeKV[K, V]{
		entries: make(KV[K, V], cap),
	}
}

// Get returns the value associated with the key and whether the key existed
func (s *SafeKV[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	value, ok := s.entries[key]
	s.mu.RUnlock()

	return value, ok
}

// GetWithMap returns the value associated with the key and whether the key existed
func (s *SafeKV[K, V]) GetWithMap(m map[K]V) {
	s.mu.RLock()
	for k := range m {
		v, ok := s.entries[k]
		if ok {
			m[k] = v
		}
	}
	s.mu.RUnlock()
}

// GetWithLock calls fn with the value associated with the key if the key existed
func (s *SafeKV[K, V]) GetWithLock(key K, fn func(V)) {
	s.mu.RLock()
	value, ok := s.entries[key]
	if ok {
		fn(value)
	}
	s.mu.RUnlock()
}

// Set sets the value associated with the key
func (s *SafeKV[K, V]) Set(key K, value V) {
	s.mu.Lock()
	s.entries[key] = value
	s.mu.Unlock()
}

// SetNx sets the value associated with the key if the key does not exist
func (s *SafeKV[K, V]) SetNx(key K, value V) bool {
	var ok bool
	s.mu.Lock()
	if _, ok = s.entries[key]; !ok {
		s.entries[key] = value
	}
	s.mu.Unlock()
	return !ok
}

// SetX sets the value associated with the key if the key exists
func (s *SafeKV[K, V]) SetX(key K, value V) bool {
	var ok bool
	s.mu.Lock()
	if _, ok = s.entries[key]; ok {
		s.entries[key] = value
	}
	s.mu.Unlock()
	return ok
}

// Delete deletes the value associated with the key
func (s *SafeKV[K, V]) Delete(keys ...K) {
	s.mu.Lock()
	for _, key := range keys {
		delete(s.entries, key)
	}
	s.mu.Unlock()
}

// Has returns whether the key exists
func (s *SafeKV[K, V]) Has(key K) bool {
	s.mu.RLock()
	_, ok := s.entries[key]
	s.mu.RUnlock()
	return ok
}

// Len returns the length of the map
func (s *SafeKV[K, V]) Len() int {
	s.mu.RLock()
	l := len(s.entries)
	s.mu.RUnlock()
	return l
}

// Keys returns all the keys
func (s *SafeKV[K, V]) Keys() []K {
	keys := make([]K, 0, len(s.entries))
	s.mu.RLock()
	for k := range s.entries {
		keys = append(keys, k)
	}
	s.mu.RUnlock()
	return keys
}

// Values returns all the values
func (s *SafeKV[K, V]) Values() []V {
	values := make([]V, 0, len(s.entries))
	s.mu.RLock()
	for _, v := range s.entries {
		values = append(values, v)
	}
	s.mu.RUnlock()
	return values
}

// Range calls fn sequentially for each key and value present in the map.
func (s *SafeKV[K, V]) Range(fn func(key K, value V) bool) {
	s.mu.RLock()
	for k, v := range s.entries {
		if !fn(k, v) {
			break
		}
	}
	s.mu.RUnlock()
}

// Clear clears the map
func (s *SafeKV[K, V]) Clear() {
	s.mu.Lock()
	s.entries = make(map[K]V, len(s.entries))
	s.mu.Unlock()
}

// Map calls fn with the map
func (s *SafeKV[K, V]) Map(fn func(KV[K, V])) {
	s.mu.Lock()
	fn(s.entries)
	s.mu.Unlock()
}
