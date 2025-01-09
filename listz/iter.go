//go:build go1.23

package listz

import "iter"

// All returns an iterator that yields all element value in the doubly linked list.
func (l *DList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := l.Front(); e != nil; e = e.Next() {
			if !yield(e.Value) {
				break
			}
		}
	}
}

// All returns an iterator that yields all element value in the singly linked list.
func (l *SList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := l.Front(); e != nil; e = e.Next() {
			if !yield(e.Value) {
				break
			}
		}
	}
}

// All returns an iterator that yields all element value in the skip list.
func (s *SkipList[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if s.len == 0 {
			return
		}

		for e := s.head.next[0]; e != nil; e = e.next[0] {
			if !yield(e.key, e.val) {
				break
			}
		}
	}
}

// All returns an iterator that yields all element value in the skip list with custom comparator.
func (s *SkipListWithCmp[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if s.len == 0 {
			return
		}

		for e := s.head.next[0]; e != nil; e = e.next[0] {
			if !yield(e.key, e.val) {
				break
			}
		}
	}
}
