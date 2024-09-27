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
