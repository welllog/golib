//go:build go1.23

package slicez

import "iter"

func (f *FlexSlice[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < f.len; i++ {
			if !yield(f.values[f.backwardIndex(f.head+i)]) {
				break
			}
		}
	}
}

func (f *FlexSlice[T]) AllWithIndex() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := 0; i < f.len; i++ {
			if !yield(i, f.values[f.backwardIndex(f.head+i)]) {
				break
			}
		}
	}
}

func (f *FlexSlice[T]) RevAll() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := f.len - 1; i >= 0; i-- {
			if !yield(f.values[f.backwardIndex(f.head+i)]) {
				break
			}
		}
	}
}

func (f *FlexSlice[T]) RevAllWithIndex() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := f.len - 1; i >= 0; i-- {
			if !yield(i, f.values[f.backwardIndex(f.head+i)]) {
				break
			}
		}
	}
}
