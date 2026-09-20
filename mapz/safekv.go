package mapz

import (
	"fmt"
	"sync"
)

type call[V any] struct {
	wg  sync.WaitGroup
	val V
	err error
}

// SafeKV is a thread-safe map
type SafeKV[K comparable, V any] struct {
	entries KV[K, V]
	mu      sync.RWMutex
	calls   map[K]*call[V]
}

// NewSafeKV creates a new SafeKV
func NewSafeKV[K comparable, V any](capacity int) *SafeKV[K, V] {
	return &SafeKV[K, V]{
		entries: make(KV[K, V], capacity),
		calls:   make(map[K]*call[V]),
	}
}

// Get returns the value associated with the key and whether the key existed
func (s *SafeKV[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	value, ok := s.entries[key]
	s.mu.RUnlock()

	return value, ok
}

// GetSet sets the value associated with the key and returns the old value and whether the key existed
func (s *SafeKV[K, V]) GetSet(key K, value V) (V, bool) {
	s.mu.Lock()
	oldValue, ok := s.entries.GetSet(key, value)
	s.mu.Unlock()

	return oldValue, ok
}

// GetDel deletes the value associated with the key and returns the old value and whether the key existed
func (s *SafeKV[K, V]) GetDel(key K) (V, bool) {
	s.mu.Lock()
	oldValue, ok := s.entries.GetDel(key)
	s.mu.Unlock()

	return oldValue, ok
}

// GetOrSet returns the value associated with the key if it exists.
// Otherwise, it sets the value associated with the key to the provided value and returns that value.
// The loaded result is true if the value was loaded, false if stored.
func (s *SafeKV[K, V]) GetOrSet(key K, value V) (actual V, loaded bool) {
	s.mu.RLock()
	actual, loaded = s.entries[key]
	s.mu.RUnlock()

	if !loaded {
		s.mu.Lock()
		actual, loaded = s.entries[key]
		if !loaded {
			s.entries[key] = value
			actual = value
		}
		s.mu.Unlock()
	}
	return actual, loaded
}

// GetOrSetFunc returns the value associated with the key if it exists.
// Otherwise, it sets the value associated with the key to the result of fn and returns that value.
// loaded indicates whether the value was already present.
func (s *SafeKV[K, V]) GetOrSetFunc(key K, fn func() (V, error)) (actual V, loaded bool, err error) {
	s.mu.RLock()
	actual, loaded = s.entries[key]
	s.mu.RUnlock()
	if loaded {
		return actual, true, nil
	}

	s.mu.Lock()
	actual, loaded = s.entries[key]
	if loaded {
		s.mu.Unlock()
		return actual, true, nil
	}

	// Check if another goroutine is already running fn() for this key.
	if c, ok := s.calls[key]; ok {
		s.mu.Unlock()
		c.wg.Wait() // Wait for the other goroutine to finish.

		if c.err != nil {
			return c.val, false, c.err
		}
		return c.val, true, nil
	}

	// No other goroutine is running fn() for this key. Create a new call.
	c := new(call[V])
	c.wg.Add(1)
	s.calls[key] = c
	s.mu.Unlock()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in GetOrSetFunc: %v", r)
			c.err = err
		}

		s.mu.Lock()
		delete(s.calls, key)
		if c.err == nil {
			s.entries[key] = c.val
		}
		s.mu.Unlock()

		c.wg.Done()
	}()
	c.val, c.err = fn()

	return c.val, false, c.err
}

// FillMap fills the provided map with values from the SafeKV for existing keys
func (s *SafeKV[K, V]) FillMap(m map[K]V) {
	s.mu.RLock()
	for k := range m {
		v, ok := s.entries[k]
		if ok {
			m[k] = v
		}
	}
	s.mu.RUnlock()
}

// GetWithMap returns the value associated with the key and whether the key existed
// Deprecated: use FillMap instead
func (s *SafeKV[K, V]) GetWithMap(m map[K]V) {
	s.FillMap(m)
}

// GetWithLock calls fn with the value associated with the key if the key existed
func (s *SafeKV[K, V]) GetWithLock(key K, fn func(V)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.entries[key]
	if ok {
		fn(value)
	}
}

// Set sets the value associated with the key
func (s *SafeKV[K, V]) Set(key K, value V) {
	s.mu.Lock()
	s.entries[key] = value
	s.mu.Unlock()
}

// SetBatch sets multiple key-value pairs at once
func (s *SafeKV[K, V]) SetBatch(kvs map[K]V) {
	s.mu.Lock()
	for k, v := range kvs {
		s.entries[k] = v
	}
	s.mu.Unlock()
}

// SetIf sets the value associated with the key if the key does not exist or if fn returns true for the old value
func (s *SafeKV[K, V]) SetIf(key K, value V, fn func(oldValue V) bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.entries.SetIf(key, value, fn)
}

// SetIfPresent sets the value associated with the key if the key exists
func (s *SafeKV[K, V]) SetIfPresent(key K, value V) bool {
	s.mu.Lock()
	ok := s.entries.SetIfPresent(key, value)
	s.mu.Unlock()
	return ok
}

// SetIfAbsent sets the value associated with the key if the key does not exist
func (s *SafeKV[K, V]) SetIfAbsent(key K, value V) bool {
	s.mu.Lock()
	ok := s.entries.SetIfAbsent(key, value)
	s.mu.Unlock()
	return ok
}

// SetNx sets the value associated with the key if the key does not exist
// alias of SetIfAbsent
func (s *SafeKV[K, V]) SetNx(key K, value V) bool {
	return s.SetIfAbsent(key, value)
}

// SetX sets the value associated with the key if the key exists
// alias of SetIfPresent
func (s *SafeKV[K, V]) SetX(key K, value V) bool {
	return s.SetIfPresent(key, value)
}

// Delete deletes the value associated with the key
// alias of Remove
func (s *SafeKV[K, V]) Delete(keys ...K) {
	s.Remove(keys...)
}

// Remove deletes the value associated with the key
func (s *SafeKV[K, V]) Remove(keys ...K) {
	s.mu.Lock()
	for _, key := range keys {
		delete(s.entries, key)
	}
	s.mu.Unlock()
}

// RemoveIf deletes the value associated with the key if fn returns true for the old value
func (s *SafeKV[K, V]) RemoveIf(key K, fn func(value V) bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.entries.RemoveIf(key, fn)
}

// Has returns whether the key exists
// alias of Contains
func (s *SafeKV[K, V]) Has(key K) bool {
	return s.Contains(key)
}

// Contains returns whether the key exists
func (s *SafeKV[K, V]) Contains(key K) bool {
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
	s.mu.RLock()
	keys := make([]K, 0, len(s.entries))
	for k := range s.entries {
		keys = append(keys, k)
	}
	s.mu.RUnlock()
	return keys
}

// Values returns all the values
func (s *SafeKV[K, V]) Values() []V {
	s.mu.RLock()
	values := make([]V, 0, len(s.entries))
	for _, v := range s.entries {
		values = append(values, v)
	}
	s.mu.RUnlock()
	return values
}

// Range calls fn sequentially for each key and value present in the map.
func (s *SafeKV[K, V]) Range(fn func(key K, value V) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.entries {
		if !fn(k, v) {
			break
		}
	}
}

// Clear clears the map
func (s *SafeKV[K, V]) Clear() {
	s.mu.Lock()
	s.entries = make(map[K]V, len(s.entries))
	s.mu.Unlock()
}

// Map calls fn with the map under write lock
func (s *SafeKV[K, V]) Map(fn func(KV[K, V])) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(s.entries)
}

// ReadMap calls fn with the map under read lock
func (s *SafeKV[K, V]) ReadMap(fn func(KV[K, V])) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn(s.entries)
}
