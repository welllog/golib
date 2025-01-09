package listz

import (
	"math/bits"
	"math/rand"
	"time"
)

type skipNodeCmp[K any, V any] struct {
	key  K
	val  V
	next []*skipNodeCmp[K, V]
}

type SkipListWithCmp[K any, V any] struct {
	head  skipNodeCmp[K, V]
	len   int
	level int
	cmp   func(K, K) int
	rand  *rand.Rand
}

// NewSkipListWithCmp returns an initialized skip list with custom comparator.
func NewSkipListWithCmp[K any, V any](keyCmp func(K, K) int) *SkipListWithCmp[K, V] {
	var s SkipListWithCmp[K, V]
	s.Init(keyCmp)
	return &s
}

// Init initializes the skip list with custom comparator.
func (s *SkipListWithCmp[K, V]) Init(keyCmp func(K, K) int) {
	s.head.next = make([]*skipNodeCmp[K, V], maxLevel)
	s.len = 0
	s.level = 1
	s.cmp = keyCmp
	s.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Len returns the number of nodes of the skip list.
func (s *SkipListWithCmp[K, V]) Len() int {
	return s.len
}

// Set sets the value associated with the key.
func (s *SkipListWithCmp[K, V]) Set(key K, val V) {
	s.set(key, val, 0)
}

// SetNx sets the value associated with the key if the key does not exist.
func (s *SkipListWithCmp[K, V]) SetNx(key K, val V) bool {
	return s.set(key, val, 2)
}

// SetX sets the value associated with the key if the key exists.
func (s *SkipListWithCmp[K, V]) SetX(key K, val V) bool {
	return s.set(key, val, 1)
}

// Get returns the value associated with the key.
func (s *SkipListWithCmp[K, V]) Get(key K) (V, bool) {
	cur := &s.head
	// if the skip list not initialized, the level is 0, so the loop will not be executed
	for i := s.level - 1; i >= 0; i-- {
		for cur.next[i] != nil {
			next := cur.next[i]
			n := s.cmp(next.key, key)
			if n > 0 {
				break
			}

			if n == 0 {
				return cur.next[i].val, true
			}

			cur = next
		}
	}

	var zero V
	return zero, false
}

// Remove deletes the value associated with the key.
func (s *SkipListWithCmp[K, V]) Remove(key K) (V, bool) {
	update := make([]*skipNodeCmp[K, V], maxLevel)
	cur := &s.head
	var curLevel int
	// if the skip list not initialized, the level is 0, so the loop will not be executed
	for i := s.level - 1; i >= 0; i-- {
		for cur.next[i] != nil {
			next := cur.next[i]
			n := s.cmp(next.key, key)
			if n > 0 {
				break
			}

			if n == 0 {
				if curLevel == 0 {
					curLevel = i + 1
				}
				break
			}

			cur = next
		}

		update[i] = cur
	}

	var val V
	if curLevel == 0 {
		return val, false
	}

	cur = cur.next[0]
	val = cur.val
	for i := 0; i < curLevel; i++ {
		update[i].next[i] = cur.next[i]
	}

	if curLevel >= s.level {
		for s.level > 1 && s.head.next[s.level-1] == nil {
			s.level--
		}
	}

	s.len--
	return val, true
}

// Clear removes all nodes from the skip list.
func (s *SkipListWithCmp[K, V]) Clear() {
	s.head.next = make([]*skipNodeCmp[K, V], maxLevel)
	s.len = 0
	s.level = 1
}

// Range calls f sequentially for each key and value present in the skip list.
func (s *SkipListWithCmp[K, V]) Range(f func(K, V) bool) {
	if s.len == 0 {
		return
	}

	for e := s.head.next[0]; e != nil; e = e.next[0] {
		if !f(e.key, e.val) {
			break
		}
	}
}

// RangeWithStart calls f sequentially for each key and value present in the skip list starting from the key.
// The zone is [start, +∞)
func (s *SkipListWithCmp[K, V]) RangeWithStart(start K, f func(K, V) bool) {
	cur := &s.head
top:
	for i := s.level - 1; i >= 0; i-- {
		for cur.next[i] != nil {
			next := cur.next[i]
			n := s.cmp(next.key, start)
			if n > 0 {
				break
			}

			if n == 0 {
				cur = next
				if !f(next.key, next.val) {
					return
				}
				break top
			}

			cur = next
		}
	}

	for cur.next[0] != nil {
		next := cur.next[0]
		if !f(next.key, next.val) {
			break
		}
		cur = next
	}
}

// RangeWithRange calls f sequentially for each key and value present in the skip list within the range [start, end).
func (s *SkipListWithCmp[K, V]) RangeWithRange(start, end K, f func(K, V) bool) {
	s.RangeWithStart(start, func(key K, val V) bool {
		if s.cmp(key, end) >= 0 {
			return false
		}
		return f(key, val)
	})
}

func (s *SkipListWithCmp[K, V]) randomLevel() int {
	// k is a random number in [0, 2^maxLevel)
	k := s.rand.Uint64() & zoneMask
	return ((maxLevel - bits.Len64(k)) & levelMask) + 1
}

// set sets the value associated with the key.
// mode: 0 set the value don't care if the key exists
//
//	1 set the value if the key exists
//	2 set the value if the key does not exist
func (s *SkipListWithCmp[K, V]) set(key K, val V, mode int) bool {
	update := make([]*skipNodeCmp[K, V], maxLevel)
	cur := &s.head
	for i := s.level - 1; i >= 0; i-- {
		for cur.next[i] != nil {
			next := cur.next[i]
			n := s.cmp(next.key, key)
			if n > 0 {
				break
			}

			if n == 0 {
				if mode == 2 {
					// set the value if the key does not exist
					return false
				}
				next.val = val
				return true
			}

			cur = next
		}

		update[i] = cur
	}

	if mode == 1 {
		// set the value if the key exists
		return false
	}

	level := s.randomLevel()
	if level > s.level {
		level = s.level + 1
		for i := s.level; i < level; i++ {
			update[i] = &s.head
		}
		s.level = level
	}

	node := &skipNodeCmp[K, V]{
		key:  key,
		val:  val,
		next: make([]*skipNodeCmp[K, V], level),
	}

	for i := 0; i < level; i++ {
		node.next[i] = update[i].next[i]
		update[i].next[i] = node
	}

	s.len++
	return true
}
