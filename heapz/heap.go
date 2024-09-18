package heapz

type Element[T any] struct {
	index int
	heap  *Heap[T]
	Value T
}

func (e *Element[T]) Index() int {
	return e.index
}

type Heap[T any] struct {
	values []*Element[T]
	cmp    func(*Element[T], *Element[T]) bool
}

// New returns a new heap with the given capacity and compare function.
func New[T any](cap int, cmp func(T, T) bool) *Heap[T] {
	var h Heap[T]
	values := make([]*Element[T], 0, cap)
	h.init(values, cmp)
	return &h
}

// Init initializes a heap with the given elements and compare function.
func (h *Heap[T]) Init(s []T, cmp func(T, T) bool) {
	values := make([]*Element[T], len(s))
	for i, v := range s {
		values[i] = &Element[T]{Value: v, heap: h, index: i}
	}

	h.init(values, cmp)
}

// Push pushes the element x onto the heap.
func (h *Heap[T]) Push(x T) *Element[T] {
	index := len(h.values)
	e := &Element[T]{Value: x, heap: h, index: index}
	h.values = append(h.values, e)
	up(h.values, h.cmp, swapEle[T], index)
	return e
}

// Pop removes and returns the minimum element (according to compare function) from the heap.
func (h *Heap[T]) Pop() *Element[T] {
	n := len(h.values)
	if n == 0 {
		return nil
	} else if n == 1 {
		return h.pop()
	}

	n--
	swapEle(h.values, 0, n)
	down(h.values, h.cmp, swapEle[T], 0, n)
	return h.pop()
}

// Peek returns the minimum element (according to compare function) from the heap without removing it.
func (h *Heap[T]) Peek() *Element[T] {
	if len(h.values) == 0 {
		return nil
	}

	return h.values[0]
}

// Len returns the number of elements in the heap.
func (h *Heap[T]) Len() int {
	return len(h.values)
}

// Remove removes the element from the heap.
func (h *Heap[T]) Remove(e *Element[T]) {
	if e.heap == nil || e.heap != h {
		return
	}

	if e.index < 0 || e.index >= len(h.values) {
		panic("heap: invalid index")
	}

	n := len(h.values) - 1
	if n != e.index {
		index := e.index
		swapEle(h.values, e.index, n)
		fix(h.values, h.cmp, swapEle[T], index, n)
	}

	_ = h.pop()
}

// Fix re-establishes the heap ordering after the element has changed its value.
func (h *Heap[T]) Fix(e *Element[T]) {
	if e.heap == nil || e.heap != h {
		return
	}

	if e.index < 0 || e.index >= len(h.values) {
		panic("heap: invalid index")
	}

	fix(h.values, h.cmp, swapEle[T], e.index, len(h.values))
}

func (h *Heap[T]) init(values []*Element[T], cmp func(T, T) bool) {
	rcmp := func(a *Element[T], b *Element[T]) bool {
		return cmp(a.Value, b.Value)
	}

	build(values, rcmp, swapEle[T])

	h.values = values
	h.cmp = rcmp
}

func (h *Heap[T]) pop() *Element[T] {
	n := len(h.values) - 1
	e := h.values[n]
	h.values[n] = nil
	h.values = h.values[:n]

	e.heap = nil
	e.index = -1
	return e
}

func swapEle[T any](s []*Element[T], i, j int) {
	s[i], s[j] = s[j], s[i]
	s[i].index = i
	s[j].index = j
}
