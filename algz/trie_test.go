package algz

import (
	"fmt"
	"strings"
	"testing"
)

func TestTrie_FindAll(t *testing.T) {
	trie := Trie{trieNode{name: "root"}}
	// trie := Trie{trieNode{}}
	trie.Insert("he")
	trie.Insert("she")
	trie.Insert("his")
	trie.Insert("hers")
	trie.BuildFailureLinks()

	queue := []*trieNode{&trie.root}
	for i := 0; i < len(queue); i++ {
		node := queue[i]
		var name string
		if node.fail != nil {
			name = node.fail.name
		}
		fmt.Println(node.name, "-------", name)
		for _, child := range node.children {
			queue = append(queue, child.node)
		}
	}

	text := "ahishers"
	fmt.Println(trie.FindAll(text))
}

func TestTrie2(t *testing.T) {
	s := "\\\"level\\\":\\\"error\\\"|fatal|panic:|\\\"level\\\":\\\"alert\\\"|\\\"level\\\":\\\"severe\\\"|\\\"level\\\":\\\"stack\\\""
	ss := strings.Split(s, "|")
	// caseStr := `{"message":"{\"@timestamp\":\"2024-12-09T14:21:31\",\"level\":\"error\",\"caller\":\"userbase/authLogic.go:139\",\"content\":\"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537\",\"trace\":\"8bb377a4c45b1019990808dd939c2b3f\",\"span\":\"d1486113e267fd62\",\"method\":\"/userBase.UserBase/Auth\"}"}`
	// caseStr := `{"message":"{\"level\":\"error\"}"}`
	caseStr := `{"message":{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth""}}`
	// caseStr := `{"message":{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"\"level\":\"error\""}}`

	tt := newAhoCorasick()
	for _, s := range ss {
		tt.Insert(s)
	}
	tt.BuildFailureLinks()

	var tt2 Trie
	for _, s := range ss {
		tt2.Insert(s)
	}
	tt2.BuildFailureLinks()

	fmt.Println(tt.Match(caseStr))
	fmt.Println(tt2.Match(caseStr))
}

func BenchmarkTrie_Search(b *testing.B) {
	s := "\\\"level\\\":\\\"error\\\"|fatal|panic:|\\\"level\\\":\\\"alert\\\"|\\\"level\\\":\\\"severe\\\"|\\\"level\\\":\\\"stack\\\""
	ss := strings.Split(s, "|")
	// caseStr := `{"message":"{\"@timestamp\":\"2024-12-09T14:21:31\",\"level\":\"error\",\"caller\":\"userbase/authLogic.go:139\",\"content\":\"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537\",\"trace\":\"8bb377a4c45b1019990808dd939c2b3f\",\"span\":\"d1486113e267fd62\",\"method\":\"/userBase.UserBase/Auth\"}"}`
	// caseStr := `{"message":"{\"level\":\"error\"}"}`
	caseStr := `{"message":{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth"}{"@timestamp":"2024-12-09T14:21:31","level":"error","caller":"userbase/authLogic.go:139","content":"接口/invoiceAdmin/getDomainManagePageList不在权限列表中,account_id: 537","trace":"8bb377a4c45b1019990808dd939c2b3f","span":"d1486113e267fd62","method":"/userBase.UserBase/Auth","}}`

	b.Run("trie1", func(b *testing.B) {
		b.ReportAllocs()

		t := newAhoCorasick()
		for _, s := range ss {
			t.Insert(s)
		}
		t.BuildFailureLinks()

		for i := 0; i < b.N; i++ {
			t.Match(caseStr)
		}
	})

	b.Run("trie2", func(b *testing.B) {
		b.ReportAllocs()

		var t Trie
		for _, s := range ss {
			t.Insert(s)
		}
		t.BuildFailureLinks()

		for i := 0; i < b.N; i++ {
			t.Match(caseStr)
		}
	})

	b.Run("trie3", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			contains(ss, caseStr)
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

// 插入模式串
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

// 构建失败指针
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
		}

		if node.isEnd {
			return true
		}
	}
	return false
}
