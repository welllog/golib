package sortz

import (
	"sort"

	"github.com/welllog/golib/typez"
)

type OrderAsc[T typez.Ordered] []T

func (a OrderAsc[T]) Len() int {
	return len(a)
}

func (a OrderAsc[T]) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a OrderAsc[T]) Less(i, j int) bool {
	return a[i] < a[j]
}

type OrderDesc[T typez.Ordered] []T

func (d OrderDesc[T]) Len() int {
	return len(d)
}

func (d OrderDesc[T]) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d OrderDesc[T]) Less(i, j int) bool {
	return d[i] > d[j]
}

func Asc[T typez.Ordered](a []T) {
	sort.Sort(OrderAsc[T](a))
}

func Desc[T typez.Ordered](d []T) {
	sort.Sort(OrderDesc[T](d))
}
