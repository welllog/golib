package listz

import (
	"math/bits"
	"math/rand"
	"time"

	"github.com/welllog/golib/typez"
)

const (
	maxLevel  = 1 << 5
	zoneMask  = (uint64(1) << maxLevel) - 1
	levelMask = maxLevel - 1
)

type SkipNode[K typez.Ordered, V any] struct {
	key  K
	val  V
	next []*SkipNode[K, V]
}

func (n *SkipNode[K, V]) Key() K {
	return n.key
}

func (n *SkipNode[K, V]) Value() V {
	return n.val
}

func (n *SkipNode[K, V]) SetValue(val V) {
	n.val = val
}

func (n *SkipNode[K, V]) Next() *SkipNode[K, V] {
	return n.next[0]
}

type SkipList[K typez.Ordered, V any] struct {
	head  SkipNode[K, V]
	len   int
	level int
	rand  *rand.Rand
}

// NewSkipList returns an initialized skip list.
func NewSkipList[K typez.Ordered, V any]() *SkipList[K, V] {
	var s SkipList[K, V]
	s.Init()
	return &s
}

// Init initializes the skip list.
func (s *SkipList[K, V]) Init() {
	s.head.next = make([]*SkipNode[K, V], maxLevel)
	s.len = 0
	s.level = 1
	s.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Len returns the number of nodes of the skip list.
func (s *SkipList[K, V]) Len() int {
	return s.len
}

// Set sets the value associated with the key.
func (s *SkipList[K, V]) Set(key K, val V) {
	s.set(key, val, 0)
}

// SetNx sets the value associated with the key if the key does not exist.
func (s *SkipList[K, V]) SetNx(key K, val V) bool {
	return s.set(key, val, 2)
}

// SetX sets the value associated with the key if the key exists.
func (s *SkipList[K, V]) SetX(key K, val V) bool {
	return s.set(key, val, 1)
}

// Get returns the value associated with the key.
func (s *SkipList[K, V]) Get(key K) (V, bool) {
	node := s.GetNode(key)
	if node != nil {
		return node.val, true
	}

	var zero V
	return zero, false
}

// GetNode returns the node associated with the key.
func (s *SkipList[K, V]) GetNode(key K) *SkipNode[K, V] {
	cur := &s.head
	// if the skip list not initialized, the level is 0, so the loop will not be executed
	for i := s.level - 1; i >= 0; i-- {
		for cur.next[i] != nil {
			next := cur.next[i]
			if next.key > key {
				break
			}

			if next.key == key {
				return next
			}

			cur = next
		}
	}

	return nil
}

// Head returns the first node of the skip list.
func (s *SkipList[K, V]) Head() *SkipNode[K, V] {
	if s.len == 0 {
		return nil
	}
	return s.head.next[0]
}

// Remove deletes the value associated with the key.
func (s *SkipList[K, V]) Remove(key K) (V, bool) {
	update := make([]*SkipNode[K, V], maxLevel)
	cur := &s.head
	var curLevel int
	// if the skip list not initialized, the level is 0, so the loop will not be executed
	for i := s.level - 1; i >= 0; i-- {
		for cur.next[i] != nil {
			next := cur.next[i]
			if next.key > key {
				break
			}

			if next.key == key {
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
	cur.next = nil

	if curLevel >= s.level {
		for s.level > 1 && s.head.next[s.level-1] == nil {
			s.level--
		}
	}

	s.len--
	return val, true
}

// Clear removes all nodes from the skip list.
func (s *SkipList[K, V]) Clear() {
	s.head.next = make([]*SkipNode[K, V], maxLevel)
	s.len = 0
	s.level = 1
}

// Range traverses the skip list in ascending order.
func (s *SkipList[K, V]) Range(f func(key K, val V) bool) {
	if s.len == 0 {
		return
	}

	cur := &s.head
	for cur.next[0] != nil {
		next := cur.next[0]
		if !f(next.key, next.val) {
			break
		}
		cur = next
	}
}

// RangeWithStart traverses the skip list in ascending order starting from the start key.
// The zone is [start, +∞)
func (s *SkipList[K, V]) RangeWithStart(start K, f func(key K, val V) bool) {
	cur := &s.head
top:
	for i := s.level - 1; i >= 0; i-- {
		for cur.next[i] != nil {
			next := cur.next[i]
			if next.key > start {
				break
			}

			if next.key == start {
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

// RangeWithRange traverses the skip list in ascending order starting from the start key and ending before the end key.
// The zone is [start, end)
func (s *SkipList[K, V]) RangeWithRange(start, end K, f func(key K, val V) bool) {
	s.RangeWithStart(start, func(key K, val V) bool {
		if key >= end {
			return false
		}
		return f(key, val)
	})
}

// Keys returns all keys in the skip list.
func (s *SkipList[K, V]) Keys() []K {
	if s.len == 0 {
		return nil
	}

	keys := make([]K, s.len)
	var i int
	for e := s.head.next[0]; e != nil; e = e.next[0] {
		keys[i] = e.key
		i++
	}
	return keys
}

// Values returns all values in the skip list.
func (s *SkipList[K, V]) Values() []V {
	if s.len == 0 {
		return nil
	}

	vals := make([]V, s.len)
	var i int
	for e := s.head.next[0]; e != nil; e = e.next[0] {
		vals[i] = e.val
		i++
	}
	return vals
}

func (s *SkipList[K, V]) lazyInit() {
	if s.head.next == nil {
		s.Init()
	}
}

// set the value associated with the key
// mode:
// 0: set the value don't care if the key exists
// 1: set the value if the key exists
// 2: set the value if the key does not exist
func (s *SkipList[K, V]) set(key K, val V, mode int) bool {
	s.lazyInit()
	update := make([]*SkipNode[K, V], maxLevel)
	cur := &s.head
	// find the previous node of the target node
	for i := s.level - 1; i >= 0; i-- {
		for cur.next[i] != nil {
			next := cur.next[i]
			if next.key > key {
				break
			}

			if next.key == key {
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

	level := randomLevel(s.rand)
	if level > s.level {
		level = s.level + 1
		for i := s.level; i < level; i++ {
			update[i] = &s.head
		}
		s.level = level
	}

	node := &SkipNode[K, V]{
		key:  key,
		val:  val,
		next: make([]*SkipNode[K, V], level),
	}

	for i := 0; i < level; i++ {
		node.next[i] = update[i].next[i]
		update[i].next[i] = node
	}

	s.len++
	return true
}

func randomLevel(r *rand.Rand) int {
	// k is a random number in [0, 2^maxLevel)
	k := r.Uint64() & zoneMask
	return ((maxLevel - bits.Len64(k)) & levelMask) + 1
}
