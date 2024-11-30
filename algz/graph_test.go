package algz

import (
	"sort"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestGraph_GetMaximalCliques(t *testing.T) {
	var g Graph[string]
	g.Init(12)

	g.AddUndirectedEdge("a", "b")
	g.AddUndirectedEdge("b", "c")
	g.AddUndirectedEdge("b", "d")
	g.AddUndirectedEdge("c", "e")
	g.AddUndirectedEdge("d", "e")
	g.AddUndirectedEdge("e", "f")
	g.AddUndirectedEdge("e", "g")
	g.AddUndirectedEdge("f", "g")
	g.AddUndirectedEdge("f", "g")
	g.AddNode("h")
	g.AddNode("i")
	g.AddUndirectedEdge("j", "k")
	g.AddUndirectedEdge("j", "l")
	g.AddEdge("l", "m")

	cliques := make([][]string, 0, 10)
	R := make([]string, 0, 4)
	P := make([]string, 0, len(g.Nodes))

	for v := range g.Nodes {
		P = append(P, v)
	}
	sort.Strings(P)

	g.BronKerbosch(R, P, P[:0], &cliques)

	expected := [][]string{
		{"a", "b"},
		{"b", "c"},
		{"b", "d"},
		{"c", "e"},
		{"d", "e"},
		{"e", "f", "g"},
		{"h"},
		{"i"},
		{"j", "k"},
		{"j", "l"},
	}

	for i, v := range cliques {
		testz.Equal(t, expected[i], v, "unexpected clique")
	}
}

func TestGraph_GetPaths(t *testing.T) {
	var g Graph[string]
	g.Init(12)

	g.AddUndirectedEdge("a", "b")
	g.AddUndirectedEdge("b", "c")
	g.AddUndirectedEdge("b", "d")
	g.AddUndirectedEdge("c", "e")
	g.AddUndirectedEdge("d", "e")
	g.AddUndirectedEdge("e", "f")
	g.AddUndirectedEdge("e", "g")
	g.AddUndirectedEdge("f", "g")
	g.AddUndirectedEdge("f", "g")
	g.AddNode("h")
	g.AddNode("i")
	g.AddUndirectedEdge("j", "k")
	g.AddUndirectedEdge("j", "l")
	g.AddEdge("l", "m")

	paths := make([][]string, 0, 20)
	R := make([]string, 0, len(g.Nodes))
	P := make([]string, 0, len(g.Nodes))
	X := make(map[string]struct{}, len(g.Nodes))

	for v := range g.Nodes {
		P = append(P, v)
	}
	sort.Strings(P)

	policy := func(path []string, s string, r Relationship[string]) bool {
		for _, v := range path {
			if !r.GetNeighbors(v).Contains(s) {
				return false
			}
		}
		return true
	}
	g.Backtrack(R, P, X, &paths, policy)

	expected := [][]string{
		{"a", "b"},
		{"b", "a"},
		{"b", "c"},
		{"b", "d"},
		{"c", "b"},
		{"c", "e"},
		{"d", "b"},
		{"d", "e"},
		{"e", "c"},
		{"e", "d"},
		{"e", "f", "g"},
		{"e", "g", "f"},
		{"f", "e", "g"},
		{"f", "g", "e"},
		{"g", "e", "f"},
		{"g", "f", "e"},
		{"h"},
		{"i"},
		{"j", "k"},
		{"j", "l"},
		{"k", "j"},
		{"l", "j"},
	}

	for i, v := range expected {
		testz.Equal(t, v, paths[i], "unexpected path")
	}

	R = append(R, "e")
	X["e"] = struct{}{}
	paths = paths[:0]
	g.Backtrack(R, P, X, &paths, policy)
	expected = [][]string{
		{"e", "c"},
		{"e", "d"},
		{"e", "f", "g"},
		{"e", "g", "f"},
	}
	for i, v := range expected {
		testz.Equal(t, v, paths[i], "unexpected path")
	}
}
