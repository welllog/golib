package sortz

import (
	"sort"

	"github.com/welllog/golib/typez"
)

// orderAsc is a type that implements the sort.Interface to sort in ascending order.
type orderAsc[T typez.Ordered] []T

func (a orderAsc[T]) Len() int {
	return len(a)
}

func (a orderAsc[T]) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a orderAsc[T]) Less(i, j int) bool {
	return a[i] < a[j]
}

// orderDesc is a type that implements the sort.Interface to sort in descending order.
type orderDesc[T typez.Ordered] []T

func (d orderDesc[T]) Len() int {
	return len(d)
}

func (d orderDesc[T]) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d orderDesc[T]) Less(i, j int) bool {
	return d[i] > d[j]
}

type orderAscObj[T any, K typez.Ordered] struct {
	s  []T
	fn func(T) K
}

func (o orderAscObj[T, K]) Len() int {
	return len(o.s)
}

func (o orderAscObj[T, K]) Swap(i, j int) {
	o.s[i], o.s[j] = o.s[j], o.s[i]
}

func (o orderAscObj[T, K]) Less(i, j int) bool {
	return o.fn(o.s[i]) < o.fn(o.s[j])
}

type orderDescObj[T any, K typez.Ordered] struct {
	s  []T
	fn func(T) K
}

func (o orderDescObj[T, K]) Len() int {
	return len(o.s)
}

func (o orderDescObj[T, K]) Swap(i, j int) {
	o.s[i], o.s[j] = o.s[j], o.s[i]
}

func (o orderDescObj[T, K]) Less(i, j int) bool {
	return o.fn(o.s[i]) > o.fn(o.s[j])
}

// Asc sorts the given slice in ascending order.
func Asc[T typez.Ordered](a []T) {
	sort.Sort(orderAsc[T](a))
}

// Desc sorts the given slice in descending order.
func Desc[T typez.Ordered](d []T) {
	sort.Sort(orderDesc[T](d))
}

// AscByKey sorts the given slice in ascending order by the key returned by keyFn.
func AscByKey[T any, K typez.Ordered](a []T, keyFn func(T) K) {
	sort.Sort(orderAscObj[T, K]{s: a, fn: keyFn})
}

// DescByKey sorts the given slice in descending order by the key returned by keyFn.
func DescByKey[T any, K typez.Ordered](d []T, keyFn func(T) K) {
	sort.Sort(orderDescObj[T, K]{s: d, fn: keyFn})
}

// AscStableByKey sorts the given slice in ascending order by the key returned by keyFn.
func AscStableByKey[T any, K typez.Ordered](a []T, keyFn func(T) K) {
	sort.Stable(orderAscObj[T, K]{s: a, fn: keyFn})
}

// DescStableByKey sorts the given slice in descending order by the key returned by keyFn.
func DescStableByKey[T any, K typez.Ordered](d []T, keyFn func(T) K) {
	sort.Stable(orderDescObj[T, K]{s: d, fn: keyFn})
}
