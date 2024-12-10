package algz

import (
	"unicode/utf8"

	"github.com/welllog/golib/ringz"
)

type Trie struct {
	root trieNode
}

type trieNode struct {
	name     string
	children []childNode
	fail     *trieNode
	size     int
	isEnd    bool
}

type childNode struct {
	val  rune
	node *trieNode
}

type scope struct {
	start int
	stop  int
}

// Insert inserts a pattern into the trie.
func (t *Trie) Insert(pattern string) {
	if len(pattern) == 0 {
		return
	}

	node := &t.root
	var (
		r    rune
		size int
	)
	for i := 0; i < len(pattern); {
		if b := pattern[i]; b < utf8.RuneSelf {
			r = rune(b)
			i++
		} else {
			r, size = utf8.DecodeRuneInString(pattern[i:])
			i += size
		}

		idx := t.findChildIndex(node.children, r)
		if idx >= len(node.children) {
			child := childNode{val: r, node: &trieNode{size: i}}
			child.node.name = node.name + "/" + string(r)
			node.children = append(node.children, child)
			node = child.node
		} else if node.children[idx].val != r {
			child := childNode{val: r, node: &trieNode{size: i}}
			child.node.name = node.name + "/" + string(r)
			node.children = append(node.children, child)
			copy(node.children[idx+1:], node.children[idx:])
			node.children[idx] = child
			node = child.node
		} else {
			node = node.children[idx].node
		}
	}
	node.isEnd = true
}

// BuildFailureLinks builds failure links for the trie.
func (t *Trie) BuildFailureLinks() {
	var queue ringz.Ring[*trieNode]
	queue.Init(10)

	for i := range t.root.children {
		t.root.children[i].node.fail = &t.root
		queue.PushWithExpand(t.root.children[i].node)
	}

	for !queue.IsEmpty() {
		curr, _ := queue.Pop()
		for _, child := range curr.children {
			failNode := curr.fail
			var idx int
			for failNode != nil {
				idx = t.index(failNode.children, child.val)
				if idx >= 0 {
					break
				}
				failNode = failNode.fail
			}

			if failNode == nil {
				child.node.fail = &t.root
			} else {
				child.node.fail = failNode.children[idx].node
			}

			queue.PushWithExpand(child.node)
		}
	}
}

// Match returns true if the text contains any of the patterns in the trie.
func (t *Trie) Match(text string) bool {
	node := &t.root
	for _, v := range text {
		for {
			idx := t.index(node.children, v)
			if idx >= 0 {
				node = node.children[idx].node
				break
			}

			if node == &t.root {
				break
			}

			node = node.fail
		}

		if node.isEnd {
			return true
		}
	}
	return false
}

func (t *Trie) FindAll(text string) []string {
	var scopes []scope
	t.find(text, &scopes)

	keywords := make([]string, len(scopes))
	for i, v := range scopes {
		keywords[i] = text[v.start:v.stop]
	}

	return keywords
}

func (t *Trie) find(text string, scopes *[]scope) {
	node := &t.root
	var (
		r    rune
		size int
	)

	for i := 0; i < len(text); {
		if b := text[i]; b < utf8.RuneSelf {
			r = rune(b)
			i++
		} else {
			r, size = utf8.DecodeRuneInString(text[i:])
			i += size
		}

		for {
			idx := t.index(node.children, r)
			if idx >= 0 {
				node = node.children[idx].node
				break
			}

			if node == &t.root {
				break
			}

			node = node.fail
		}

		if node.isEnd {
			*scopes = append(*scopes, scope{i - node.size, i})
		}
	}
}

func (t *Trie) index(children []childNode, val rune) int {
	low, high := 0, len(children)
	if high == 0 || val < children[0].val || val > children[high-1].val {
		return -1
	}

	for low < high {
		mid := (low + high) / 2
		if children[mid].val == val {
			return mid
		} else if children[mid].val < val {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return -1
}

func (t *Trie) findChildIndex(children []childNode, val rune) int {
	low, high := 0, len(children)
	for low < high {
		mid := (low + high) / 2
		if children[mid].val < val {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return low
}
