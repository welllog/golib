//go:build go1.23

package heapz

import "iter"

// PopAll returns an iterator that yields all elements in the heap.
func (h *Heap[T]) PopAll() iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			e := h.Pop()
			if e == nil {
				break
			}

			if !yield(e.Value) {
				break
			}
		}
	}
}

// PopAll returns an iterator that yields all elements in the heap.
func (s *Slice[T]) PopAll() iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			e, ok := s.Pop()
			if !ok {
				break
			}

			if !yield(e) {
				break
			}
		}
	}
}
