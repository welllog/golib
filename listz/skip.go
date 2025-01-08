package listz

import (
	"math/bits"
	"math/rand"
	"time"

	"github.com/welllog/golib/typez"
)

const (
	maxLevel  = 1 << 5
	levelMask = uint64(1)<<maxLevel - 1
)

type skipNode[K typez.Ordered, V any] struct {
	key  K
	val  V
	next []*skipNode[K, V]
}

type SkipList[K typez.Ordered, V any] struct {
	head  *skipNode[K, V]
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

// Init initializes or clears the skip list.
func (s *SkipList[K, V]) Init() {
	s.head = &skipNode[K, V]{
		next: make([]*skipNode[K, V], maxLevel),
	}
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
	cur := s.head
	for i := s.level - 1; i >= 0; i-- {
		for cur.next[i] != nil {
			if cur.next[i].key > key {
				break
			}

			if cur.next[i].key == key {
				return cur.next[i].val, true
			}

			cur = cur.next[i]
		}
	}

	var zero V
	return zero, false
}

// Remove deletes the value associated with the key.
func (s *SkipList[K, V]) Remove(key K) bool {
	update := make([]*skipNode[K, V], maxLevel)
	cur := s.head
	var curLevel int
	for i := s.level - 1; i >= 0; i-- {
		for cur.next[i] != nil {
			if cur.next[i].key > key {
				break
			}

			if cur.next[i].key == key {
				if curLevel == 0 {
					curLevel = i + 1
				}
				break
			}

			cur = cur.next[i]
		}

		update[i] = cur
	}

	if curLevel == 0 {
		return false
	}

	cur = cur.next[0]
	for i := 0; i < curLevel-1; i++ {
		update[i].next[i] = cur.next[i]
	}

	if curLevel >= s.level {
		for s.level > 1 && s.head.next[s.level-1] == nil {
			s.level--
		}
	}

	s.len--
	return true
}

// set the value associated with the key
// flag:
// 0: set the value don't care if the key exists
// 1: set the value if the key exists
// 2: set the value if the key does not exist
func (s *SkipList[K, V]) set(key K, val V, flag int) bool {
	s.lazyInit()
	update := make([]*skipNode[K, V], maxLevel)
	cur := s.head
	// find the previous node of the target node
	for i := s.level; i >= 0; i-- {
		for cur.next[i] != nil {
			if cur.next[i].key > key {
				break
			}

			if cur.next[i].key == key {
				if flag == 2 {
					// set the value if the key does not exist
					return false
				}
				cur.next[i].val = val
				return true
			}

			cur = cur.next[i]
		}

		update[i] = cur
	}

	if flag == 1 {
		// set the value if the key exists
		return false
	}

	level := s.randomLevel()
	if level > s.level {
		for i := s.level; i < level; i++ {
			update[i] = s.head
		}
		s.level = level
	}

	node := &skipNode[K, V]{
		key:  key,
		val:  val,
		next: make([]*skipNode[K, V], level),
	}

	for i := 0; i < level; i++ {
		node.next[i] = update[i].next[i]
		update[i].next[i] = node
	}

	s.len++
	return true
}

func (s *SkipList[K, V]) lazyInit() {
	if s.head == nil {
		s.Init()
	}
}

func (s *SkipList[K, V]) randomLevel() int {
	k := s.rand.Uint64() & levelMask
	return (maxLevel - bits.Len64(k)) & (maxLevel - 1)
}
