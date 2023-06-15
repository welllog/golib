package mapz

import "sync"

// SafeKV is a thread-safe map
type SafeKV[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

// NewSafeKV creates a new SafeKV
func NewSafeKV[K comparable, V any](cap int) *SafeKV[K, V] {
	return &SafeKV[K, V]{
		m: make(map[K]V, cap),
	}
}

// Get returns the value associated with the key and whether the key existed
func (s *SafeKV[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	value, ok := s.m[key]
	s.mu.RUnlock()

	return value, ok
}

// Set sets the value associated with the key
func (s *SafeKV[K, V]) Set(key K, value V) {
	s.mu.Lock()
	s.m[key] = value
	s.mu.Unlock()
}

// Delete deletes the value associated with the key
func (s *SafeKV[K, V]) Delete(key K) {
	s.mu.Lock()
	delete(s.m, key)
	s.mu.Unlock()
}

// Len returns the length of the map
func (s *SafeKV[K, V]) Len() int {
	s.mu.RLock()
	l := len(s.m)
	s.mu.RUnlock()
	return l
}

// Range calls f sequentially for each key and value present in the map.
func (s *SafeKV[K, V]) Range(f func(key K, value V) bool) {
	s.mu.RLock()
	for k, v := range s.m {
		if !f(k, v) {
			break
		}
	}
	s.mu.RUnlock()
}

// Clear clears the map
func (s *SafeKV[K, V]) Clear() {
	s.mu.Lock()
	s.m = make(map[K]V, len(s.m))
	s.mu.Unlock()
}

// SetNx sets the value associated with the key if the key does not exist
func (s *SafeKV[K, V]) SetNx(key K, value V) bool {
	var isSet bool
	s.mu.Lock()
	if _, ok := s.m[key]; !ok {
		s.m[key] = value
		isSet = true
	}
	s.mu.Unlock()
	return isSet
}

// SetX sets the value associated with the key if the key exists
func (s *SafeKV[K, V]) SetX(key K, value V) bool {
	var isSet bool
	s.mu.Lock()
	if _, ok := s.m[key]; ok {
		s.m[key] = value
		isSet = true
	}
	s.mu.Unlock()
	return isSet
}

// Map calls f with the map
func (s *SafeKV[K, V]) Map(f func(map[K]V)) {
	s.mu.Lock()
	f(s.m)
	s.mu.Unlock()
}
