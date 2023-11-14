package mapz

// Keys returns all keys of the map.
func Keys[K comparable, V any](ms ...map[K]V) []K {
	var n int
	for _, m := range ms {
		n += len(m)
	}

	ret := make([]K, 0, n)
	for _, m := range ms {
		for k := range m {
			ret = append(ret, k)
		}
	}
	return ret
}

// Values returns all values of the map.
func Values[K comparable, V any](ms ...map[K]V) []V {
	var n int
	for _, m := range ms {
		n += len(m)
	}

	ret := make([]V, 0, n)
	for _, m := range ms {
		for _, v := range m {
			ret = append(ret, v)
		}
	}
	return ret
}
