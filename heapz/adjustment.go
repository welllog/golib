package heapz

func swap[T any](s []T, i, j int) {
	s[i], s[j] = s[j], s[i]
}

func fix[T any](s []T, cmp func(T, T) bool, swap func([]T, int, int), index, tail int) {
	if !down(s, cmp, swap, index, tail) {
		up(s, cmp, swap, index)
	}
}

func build[T any](s []T, cmp func(T, T) bool, swap func([]T, int, int)) {
	n := len(s)
	for i := n/2 - 1; i >= 0; i-- {
		down(s, cmp, swap, i, n)
	}
}

func up[T any](s []T, cmp func(T, T) bool, swap func([]T, int, int), j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !cmp(s[j], s[i]) {
			break
		}
		swap(s, i, j)
		j = i
	}
}

func down[T any](s []T, cmp func(T, T) bool, swap func([]T, int, int), i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && cmp(s[j2], s[j1]) {
			j = j2 // = 2*i + 2  // right child
		}
		if !cmp(s[j], s[i]) {
			break
		}
		swap(s, i, j)
		i = j
	}
	return i > i0
}
