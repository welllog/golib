//go:build go1.23

package listz

import "iter"

// All returns an iterator that yields all element value in the list.
func (l *List[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := l.Front(); e != nil; e = e.Next() {
			if !yield(e.Value) {
				break
			}
		}
	}
}
