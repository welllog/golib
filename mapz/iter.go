//go:build go1.23

package mapz

import "iter"

// All returns an iterator that yields all key-value pairs in the SafeKV.
func (s *SafeKV[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()
		for k, v := range s.entries {
			if !yield(k, v) {
				break
			}
		}
	}
}

// All returns an iterator that yields all key-value pairs in the Body.
func (b Body) All() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for k, v := range b {
			if k == hiddenKey {
				continue
			}

			if !yield(k, v) {
				break
			}
		}
	}
}
