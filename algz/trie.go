package algz

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"github.com/welllog/golib/ringz"
)

type Trie struct {
	root trieNode
}

type trieNode struct {
	// name     string
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
		r, size = decodeRune(pattern, i)
		i += size

		idx := t.findChildIndex(node.children, r)
		if idx >= len(node.children) || node.children[idx].val != r {
			child := childNode{val: r, node: &trieNode{size: i}}
			// child.node.name = node.name + "/" + string(r)
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
		idx := t.index(node.children, v)
		for node != &t.root && idx < 0 {
			node = node.fail
			idx = t.index(node.children, v)
		}

		if idx >= 0 {
			node = node.children[idx].node
			tempNode := node
			for tempNode != &t.root {
				if tempNode.isEnd {
					return true
				}
				tempNode = tempNode.fail
			}
		}

	}
	return false
}

// FindAll returns all patterns found in the text.
func (t *Trie) FindAll(text string) []string {
	var scopes []scope
	t.find(text, &scopes)

	keywords := make([]string, len(scopes))
	for i, v := range scopes {
		keywords[i] = text[v.start:v.stop]
	}

	return keywords
}

// ReplaceWithMask replaces all patterns found in the text with the mask.
func (t *Trie) ReplaceWithMask(text string, mask rune) string {
	var scopes []scope
	t.find(text, &scopes)

	t.mergeScopes(&scopes)

	var buf strings.Builder
	buf.Grow(len(text))

	var begin int
	for _, v := range scopes {
		buf.WriteString(text[begin:v.start])
		num := utf8.RuneCountInString(text[v.start:v.stop])
		for i := 0; i < num; i++ {
			buf.WriteRune(mask)
		}
		begin = v.stop
	}
	buf.WriteString(text[begin:])

	return buf.String()
}

// Replace replaces all patterns found in the text with the repl.
func (t *Trie) Replace(text string, repl string) string {
	var scopes []scope
	t.find(text, &scopes)

	t.mergeScopes(&scopes)

	var buf strings.Builder
	buf.Grow(len(text))

	var begin int
	for _, v := range scopes {
		buf.WriteString(text[begin:v.start])
		buf.WriteString(repl)
		begin = v.stop
	}
	buf.WriteString(text[begin:])

	return buf.String()
}

// PrefixSearch returns all patterns that have the key as prefix.
func (t *Trie) PrefixSearch(key string) []string {
	node := &t.root
	for _, v := range key {
		idx := t.index(node.children, v)
		if idx < 0 {
			return nil
		}

		node = node.children[idx].node
	}

	nodeStack := make([]*trieNode, 0, 10)
	nodeStack = append(nodeStack, node)
	runeStack := make([]runeDepth, 0, 10)

	var buf bytes.Buffer
	buf.Grow(len(key) + 64)
	buf.WriteString(key)

	var ret []string
	var depth int
	for len(nodeStack) > 0 {
		cur := nodeStack[len(nodeStack)-1]
		nodeStack = nodeStack[:len(nodeStack)-1]

		if len(runeStack) > 0 {
			buf.WriteRune(runeStack[len(runeStack)-1].r)
			runeStack = runeStack[:len(runeStack)-1]
		}

		if cur.isEnd {
			ret = append(ret, buf.String())
		}

		depth++

		if len(cur.children) == 0 {
			if len(runeStack) == 0 {
				break
			}

			curDepth := runeStack[len(runeStack)-1].depth
			buf.Truncate(buf.Len() - (depth - curDepth))
			depth = curDepth
			continue
		}

		for _, child := range cur.children {
			nodeStack = append(nodeStack, child.node)
			runeStack = append(runeStack, runeDepth{child.val, depth})
		}
	}

	return ret
}

// FuzzySearch returns all patterns that are similar to the key.
func (t *Trie) FuzzySearch(key string) []string {
	node := &t.root
	for _, v := range key {
		idx := t.index(node.children, v)
		for node != &t.root && idx < 0 {
			node = node.fail
			idx = t.index(node.children, v)
		}

		if idx < 0 {
			return nil
		}

		node = node.children[idx].node
	}

	nodeStack := make([]*trieNode, 0, 10)
	runeStack := make([]runeDepth, 0, 10)
	var ret []string
	var buf bytes.Buffer
	buf.Grow(len(key) + 64)

	for node != &t.root {
		nodeStack = append(nodeStack, node)
		buf.WriteString(key[len(key)-node.size:])

		var depth int
		for len(nodeStack) > 0 {
			cur := nodeStack[len(nodeStack)-1]
			nodeStack = nodeStack[:len(nodeStack)-1]

			if len(runeStack) > 0 {
				buf.WriteRune(runeStack[len(runeStack)-1].r)
				runeStack = runeStack[:len(runeStack)-1]
			}

			if cur.isEnd {
				ret = append(ret, buf.String())
			}

			depth++

			if len(cur.children) == 0 {
				if len(runeStack) == 0 {
					break
				}

				curDepth := runeStack[len(runeStack)-1].depth
				buf.Truncate(buf.Len() - (depth - curDepth))
				depth = curDepth
				continue
			}

			for _, child := range cur.children {
				nodeStack = append(nodeStack, child.node)
				runeStack = append(runeStack, runeDepth{child.val, depth})
			}
		}

		buf.Reset()
		node = node.fail
	}

	return ret
}

func (t *Trie) mergeScopes(sp *[]scope) {
	scopes := *sp
	for i := 0; i < len(scopes)-1; {
		if scopes[i].stop > scopes[i+1].start {
			if scopes[i].stop < scopes[i+1].stop {
				scopes[i].stop = scopes[i+1].stop
			}
			if scopes[i].start > scopes[i+1].start {
				scopes[i].start = scopes[i+1].start
			}
			scopes = append(scopes[:i+1], scopes[i+2:]...)
		} else {
			i++
		}
	}
	*sp = scopes
}

func (t *Trie) find(text string, scopes *[]scope) {
	node := &t.root
	var (
		r    rune
		size int
	)

	for i := 0; i < len(text); {
		r, size = decodeRune(text, i)
		i += size

		idx := t.index(node.children, r)
		for node != &t.root && idx < 0 {
			node = node.fail
			idx = t.index(node.children, r)
		}

		if idx >= 0 {
			node = node.children[idx].node
			tempNode := node
			for tempNode != &t.root {
				if tempNode.isEnd {
					*scopes = append(*scopes, scope{i - tempNode.size, i})
				}
				tempNode = tempNode.fail
			}
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

func decodeRune(s string, i int) (rune, int) {
	if b := s[i]; b < utf8.RuneSelf {
		return rune(b), 1
	}

	r, size := utf8.DecodeRuneInString(s[i:])
	return r, size
}

type runeDepth struct {
	r     rune
	depth int
}
