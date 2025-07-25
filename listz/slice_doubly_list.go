package listz

type SliceDList[T any] struct {
	nodes []sdNode[T]
	head  int
	tail  int
	free  []int
}

type sdNode[T any] struct {
	value T
	prev  int
	next  int
}

func (s *SliceDList[T]) Init(cap int) {
	s.nodes = make([]sdNode[T], 0, cap)
	s.head = -1
	s.tail = -1
}

func (s *SliceDList[T]) Len() int {
	if s.head == -1 {
		return 0
	}
	return len(s.nodes)
}

func (s *SliceDList[T]) PushFront() {

}

func (s *SliceDList[T]) PushBack(value T) {
	node := sdNode[T]{value: value, prev: s.tail, next: -1}
	if s.tail != -1 {
		s.nodes[s.tail].next = len(s.nodes)
	}
	s.nodes = append(s.nodes, node)
	s.tail = len(s.nodes) - 1
	if s.head == -1 {
		s.head = 0
	}
}
