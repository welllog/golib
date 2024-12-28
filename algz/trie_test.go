package algz

import (
	"strings"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestTrie_Match(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		text     string
		want     bool
	}{
		{
			name:     "case1",
			patterns: []string{"he", "she", "his", "hers"},
			text:     "ahishers",
			want:     true,
		},
		{
			name:     "case2",
			patterns: []string{"sheri", "he"},
			text:     "ahishers",
			want:     true,
		},
		{
			name:     "case3",
			patterns: []string{"sheri", "her"},
			text:     "ahishers",
			want:     true,
		},
		{
			name:     "case4",
			patterns: []string{"sheri", "hei"},
			text:     "ahishers",
			want:     false,
		},
		{
			name:     "case5",
			patterns: []string{"sheri", "😀"},
			text:     "her😀i",
			want:     true,
		},
		{
			name:     "case6",
			patterns: []string{"sheri", "??"},
			text:     "her😀i",
			want:     false,
		},
		{
			name:     "case7",
			patterns: []string{"sheri", "he"},
			text:     "ahishers",
			want:     true,
		},
		{
			name:     "case8",
			patterns: []string{"sheri", "her"},
			text:     "ahishers",
			want:     true,
		},
		{
			name:     "case9",
			patterns: []string{"sheri", "hec"},
			text:     "ahishers",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := Trie{}
			for _, pattern := range tt.patterns {
				trie.Insert(pattern)
			}
			trie.BuildFailureLinks()
			if got := trie.Match(tt.text); got != tt.want {
				t.Errorf("Trie.Match() = %v, want %v", got, tt.want)
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := newAhoCorasick()
			for _, pattern := range tt.patterns {
				trie.Insert(pattern)
			}
			trie.BuildFailureLinks()
			if got := trie.Match(tt.text); got != tt.want {
				t.Errorf("AhoCorasick.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrie_FindAll(t *testing.T) {
	// trie := Trie{trieNode{name: "root"}}
	trie := Trie{trieNode{}}
	trie.Insert("he")
	trie.Insert("she")
	trie.Insert("his")
	trie.Insert("hers")
	trie.BuildFailureLinks()

	// queue := []*trieNode{&trie.root}
	// for i := 0; i < len(queue); i++ {
	// 	node := queue[i]
	// 	var name string
	// 	if node.fail != nil {
	// 		name = node.fail.name
	// 	}
	// 	fmt.Println(node.name, "-------", name)
	// 	for _, child := range node.children {
	// 		queue = append(queue, child.node)
	// 	}
	// }

	text := "ahishers"
	keywords := trie.FindAll(text)
	expected := []string{"his", "she", "he", "hers"}
	for i, v := range keywords {
		if v != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], v)
		}
	}
}

func TestTrie_ReplaceWithMask(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		text     string
		mask     rune
		want     string
	}{
		{
			name:     "case1",
			patterns: []string{"he", "she", "his", "hers"},
			text:     "ahishers",
			mask:     '*',
			want:     "a*******",
		},
		{
			name:     "case2",
			patterns: []string{"he", "llo"},
			text:     "hello world",
			mask:     '#',
			want:     "##### world",
		},
		{
			name:     "case3",
			patterns: []string{"o w", " "},
			text:     "hello world",
			mask:     '#',
			want:     "hell###orld",
		},
		{
			name:     "case4",
			patterns: []string{"hello", "elo", "llo", "o ", "ld"},
			text:     "hello world",
			mask:     '*',
			want:     "******wor**",
		},
		{
			name:     "case5",
			patterns: []string{"o", "r"},
			text:     "hello world",
			mask:     '*',
			want:     "hell* w**ld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := Trie{}
			for _, pattern := range tt.patterns {
				trie.Insert(pattern)
			}
			trie.BuildFailureLinks()
			if got := trie.ReplaceWithMask(tt.text, tt.mask); got != tt.want {
				t.Errorf("Trie.ReplaceWithMask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrie_Replace(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		text     string
		repl     string
		want     string
	}{
		{
			name:     "case1",
			patterns: []string{"he", "she", "his", "hers"},
			text:     "ahishers",
			repl:     "it",
			want:     "ait",
		},
		{
			name:     "case2",
			patterns: []string{"he", "llo"},
			text:     "hello world",
			repl:     "it",
			want:     "itit world",
		},
		{
			name:     "case3",
			patterns: []string{"o w", " "},
			text:     "hello world",
			repl:     "it",
			want:     "hellitorld",
		},
		{
			name:     "case4",
			patterns: []string{"hello", "elo", "llo", "o ", "ld"},
			text:     "hello world",
			repl:     "it",
			want:     "itworit",
		},
		{
			name:     "case5",
			patterns: []string{"o", "r"},
			text:     "hello world",
			repl:     "it",
			want:     "hellit wititld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := Trie{}
			for _, pattern := range tt.patterns {
				trie.Insert(pattern)
			}
			trie.BuildFailureLinks()
			if got := trie.Replace(tt.text, tt.repl); got != tt.want {
				t.Errorf("Trie.Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrie_PrefixSearch(t *testing.T) {
	trie1 := Trie{}
	trie1.Insert("hello")
	trie1.Insert("hey")
	trie1.Insert("hello world")
	trie1.Insert("hello bob")
	trie1.Insert("happy")
	trie1.Insert("happiness")
	trie1.Insert("world")
	trie1.Insert("world peace")
	trie1.Insert("world war")
	trie1.BuildFailureLinks()

	ret1 := trie1.PrefixSearch("he")
	expect1 := []string{"hey", "hello", "hello world", "hello bob"}
	for i, v := range ret1 {
		if v != expect1[i] {
			t.Errorf("expected %s, got %s", expect1[i], v)
		}
	}

	ret2 := trie1.PrefixSearch("hello")
	expect2 := []string{"hello", "hello world", "hello bob"}
	for i, v := range ret2 {
		if v != expect2[i] {
			t.Errorf("expected %s, got %s", expect2[i], v)
		}
	}

	ret3 := trie1.PrefixSearch("world")
	expect3 := []string{"world", "world war", "world peace"}
	for i, v := range ret3 {
		if v != expect3[i] {
			t.Errorf("expected %s, got %s", expect3[i], v)
		}
	}

	ret4 := trie1.PrefixSearch("happy world")
	if len(ret4) != 0 {
		t.Errorf("expected empty, got %v", ret4)
	}

	ret5 := trie1.PrefixSearch("")
	expect5 := []string{"world", "world war", "world peace", "hey", "hello", "hello world", "hello bob", "happy", "happiness"}
	for i, v := range ret5 {
		if v != expect5[i] {
			t.Errorf("expected %s, got %s", expect5[i], v)
		}
	}
}

func TestTrie_FuzzySearch(t *testing.T) {
	trie := Trie{}
	trie.Insert("hello")
	trie.Insert("hello world")
	trie.Insert("hello bob")
	trie.Insert("happy")
	trie.Insert("happiness")
	trie.Insert("world")
	trie.Insert("world peace")
	trie.Insert("world war")
	trie.BuildFailureLinks()

	ret1 := trie.FuzzySearch("hello")
	expect1 := []string{"hello", "hello world", "hello bob"}
	for i, v := range ret1 {
		if v != expect1[i] {
			t.Errorf("expected %s, got %s", expect1[i], v)
		}
	}

	ret2 := trie.FuzzySearch("hello world")
	expect2 := []string{"hello world", "world", "world war", "world peace"}
	for i, v := range ret2 {
		if v != expect2[i] {
			t.Errorf("expected %s, got %s", expect2[i], v)
		}
	}

	ret3 := trie.FuzzySearch("heyha")
	expect3 := []string{"happy", "happiness"}
	for i, v := range ret3 {
		if v != expect3[i] {
			t.Errorf("expected %s, got %s", expect3[i], v)
		}
	}

	ret4 := trie.FuzzySearch("happa")
	if len(ret4) != 0 {
		t.Errorf("expected empty, got %v", ret4)
	}

	ret5 := trie.FuzzySearch("i happa")
	if len(ret5) != 0 {
		t.Errorf("expected empty, got %v", ret5)
	}
}

func TestTrieNodeQueue_Init(t *testing.T) {
	q := trieNodeQueue{}
	q.Init(10)

	for i := 0; i < 10; i++ {
		q.Push(&trieNode{size: i})
	}
	testz.Equal(t, 10, q.Len())
	for i := 0; i < 10; i++ {
		node := q.Pop()
		testz.Equal(t, i, node.size)
		q.Push(&trieNode{size: i})
	}
	testz.Equal(t, 10, q.Len())

	for i := 10; i < 20; i++ {
		q.Push(&trieNode{size: i})
	}
	testz.Equal(t, 20, q.Len())
	for i := 0; i < 20; i++ {
		node := q.Pop()
		testz.Equal(t, i, node.size)
	}
	if q.Pop() != nil {
		t.Error("expected nil")
	}

	for i := 0; i < 100000; i++ {
		for j := 0; j < 10; j++ {
			q.Push(&trieNode{size: i + j})
		}
		for j := 0; j < 10; j++ {
			testz.Equal(t, q.Pop().size, i+j)
		}
	}
}

func BenchmarkTrie_Search(b *testing.B) {
	patterns := []string{
		"Unlimited and uncontrolled growth", "yes", "excited", "happy", "boring",
		"I used to be one", "this phenomenon is named", "when all the failed",
		"Human constructions"}
	str := `Like a lot of people, I used to be a pathological maximalist. A phone with more features is necessarily better, a company with more people is better, a program with more lines of code is better, a house with more stuff is better.
Until the day when reality hit me in the face: there is a direct relationship between "more" and "complexity". The more a system is big and complex, the more we use our (limited) time and energy to do useless, busy work.
The more a house is stuffed, the more we take time to clean it, the more features a phone has, the more buttons we need to click to perform even basic actions...
In thermodynamics, this phenomenon is called Entropy: a property originally introduced to explain the part of the internal energy of a system that is unavailable as a source for useful work.
These are the important words: useful work.
The second law of thermodynamics states that the entropy of an isolated system always increases because isolated systems spontaneously evolve towards thermodynamic equilibrium: the state of maximum entropy, maximum chaos.
Today I want to show you that this law also applies to human organizations and constructions. When all the energy is directed towards useless work: complicated bureaucratic forms, layers and layers of management, websites and apps full buttons, NFTs, there is no more time to do the things that matters actually.
But it's even worse: when all the available energy of a living organism or organization is directed towards useless work, and no energy remains for self-sustainability, it leads to a slow and painful death. It's especially true in our world where, in all likelihood, we are going to have to deal with less energy available in the future than in the past, with the end of fossil fuels and the population plateauing in most western countries.
In one word: Entropy is fatal.
We all know jobs, teams, organizations... that long time ago were created to do useful work, and now only exist to justify their own existence.
This is also why prices inevitably go up despite technological innovation: more and more of the costs are directed towards useless work.
But, unlike in thermodynamics, there are ways to control entropy in organizations and Human constructions.`

	b.Run("innerTrie", func(b *testing.B) {
		b.ReportAllocs()

		t := newAhoCorasick()
		for _, v := range patterns {
			t.Insert(v)
		}
		t.BuildFailureLinks()

		for i := 0; i < b.N; i++ {
			t.Match(str)
		}
	})

	b.Run("Trie", func(b *testing.B) {
		b.ReportAllocs()

		var t Trie
		for _, v := range patterns {
			t.Insert(v)
		}
		t.BuildFailureLinks()

		for i := 0; i < b.N; i++ {
			t.Match(str)
		}
	})

	b.Run("contains", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			contains(patterns, str)
		}
	})

}

func contains(ss []string, s string) bool {
	for _, each := range ss {
		if strings.Contains(s, each) {
			return true
		}
	}

	return false
}

type trieNode1 struct {
	children map[rune]*trieNode1
	fail     *trieNode1
	isEnd    bool
}

type ahoCorasick struct {
	root *trieNode1
}

func newAhoCorasick() *ahoCorasick {
	return &ahoCorasick{
		root: &trieNode1{
			children: make(map[rune]*trieNode1, 1),
		},
	}
}

func (ac *ahoCorasick) Insert(pattern string) {
	if len(pattern) == 0 {
		return
	}

	node := ac.root
	for _, ch := range pattern {
		if node.children[ch] == nil {
			node.children[ch] = &trieNode1{
				children: make(map[rune]*trieNode1, 1),
			}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

func (ac *ahoCorasick) BuildFailureLinks() {
	queue := make([]*trieNode1, 0, 3)

	for _, child := range ac.root.children {
		child.fail = ac.root
		queue = append(queue, child)
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for ch, child := range curr.children {
			failNode := curr.fail
			for failNode != nil && failNode.children[ch] == nil {
				failNode = failNode.fail
			}

			if failNode == nil {
				child.fail = ac.root
			} else {
				child.fail = failNode.children[ch]
			}
			queue = append(queue, child)
		}
	}
}

func (ac *ahoCorasick) Match(text string) bool {
	node := ac.root
	for _, r := range text {
		for node != ac.root && node.children[r] == nil {
			node = node.fail
		}

		if child, ok := node.children[r]; ok {
			node = child

			tempNode := node
			for tempNode != ac.root {
				if tempNode.isEnd {
					return true
				}
				tempNode = tempNode.fail
			}
		}
	}
	return false
}
