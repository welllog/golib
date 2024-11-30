package algz

type Graph[T comparable] struct {
	Nodes map[T]map[T]struct{}
}

// Init initializes the graph with the given capacity.
func (g *Graph[T]) Init(cap int) {
	g.Nodes = make(map[T]map[T]struct{}, cap)
}

// AddNode adds a node to the graph.
func (g *Graph[T]) AddNode(node T) {
	g.lazyInit()
	if _, ok := g.Nodes[node]; !ok {
		g.Nodes[node] = make(map[T]struct{})
	}
}

// AddEdge adds an edge from one node to another.
func (g *Graph[T]) AddEdge(from, to T) {
	g.lazyInit()
	if _, ok := g.Nodes[from]; !ok {
		g.Nodes[from] = make(map[T]struct{}, 2)
	}
	g.Nodes[from][to] = struct{}{}
}

// AddUndirectedEdge adds an undirected edge between two nodes.
func (g *Graph[T]) AddUndirectedEdge(from, to T) {
	g.AddEdge(from, to)
	g.AddEdge(to, from)
}

// Len returns the number of nodes in the graph.
func (g *Graph[T]) Len() int {
	return len(g.Nodes)
}

// GetMaximalCliques finds all maximal cliques in the graph.
func (g *Graph[T]) GetMaximalCliques() [][]T {
	cliques := make([][]T, 0, 2)
	R := make([]T, 0, len(g.Nodes))
	P := make([]T, 0, len(g.Nodes))

	for v := range g.Nodes {
		P = append(P, v)
	}

	g.BronKerbosch(R, P, P[:0], &cliques)
	return cliques
}

// GetPaths finds all paths in the graph.
func (g *Graph[T]) GetPaths(pathPolicy func([]T, T, Relationship[T]) bool) [][]T {
	paths := make([][]T, 0, 2)
	R := make([]T, 0, len(g.Nodes))
	P := make([]T, 0, len(g.Nodes))
	X := make(map[T]struct{}, len(g.Nodes))

	for v := range g.Nodes {
		P = append(P, v)
	}

	g.Backtrack(R, P, X, &paths, pathPolicy)
	return paths
}

// BronKerbosch finds all maximal cliques in the graph.
// R is the current clique, P is the potential nodes to add, and X is the nodes that have been excluded.
func (g *Graph[T]) BronKerbosch(R, P, X []T, cliques *[][]T) {
	if len(P) == 0 && len(X) == 0 {
		clique := append(([]T)(nil), R...)
		*cliques = append(*cliques, clique)
		return
	}

	for _, v := range P {
		neighbors := g.Nodes[v]
		g.BronKerbosch(append(R, v), intersect(P, neighbors), intersect(X, neighbors), cliques)

		P = P[1:]
		X = append(X, v)
	}
}

// Backtrack finds all paths in the graph.
// R is the current path, first call should be with R = nil.
// P is the potential nodes to add, and X is the nodes that have been excluded.
// pathPolicy is a function that determines whether to add a node to the path.
func (g *Graph[T]) Backtrack(R, P []T, X map[T]struct{}, paths *[][]T,
	pathPolicy func([]T, T, Relationship[T]) bool) {
	for _, v := range P {
		if _, ok := X[v]; ok {
			continue
		}

		if pathPolicy(R, v, g.Nodes) {
			X[v] = struct{}{}
			R = append(R, v)

			c := *paths
			if len(R) > 1 && len(c) > 0 && isSub(c[len(c)-1], R) {
				(*paths)[len(c)-1] = append(c[len(c)-1], v)
			} else {
				*paths = append(c, append([]T(nil), R...))
			}

			g.Backtrack(R, P, X, paths, pathPolicy)

			delete(X, v)
			R = R[:len(R)-1]
		}

	}
}

func (g *Graph[T]) lazyInit() {
	if g.Nodes == nil {
		g.Nodes = make(map[T]map[T]struct{}, 4)
	}
}

func intersect[T comparable](a []T, b map[T]struct{}) []T {
	var ret []T
	for _, v := range a {
		if _, ok := b[v]; ok {
			ret = append(ret, v)
		}
	}
	return ret
}

func isSub[T comparable](sub, super []T) bool {
	if len(sub)+1 != len(super) {
		return false
	}

	super = super[:len(sub)]
	for i := range sub {
		if super[i] != sub[i] {
			return false
		}
	}

	return true
}

type Neighbors[T comparable] map[T]struct{}

func (n Neighbors[T]) Contains(node T) bool {
	if n == nil {
		return false
	}

	_, ok := n[node]
	return ok
}

type Relationship[T comparable] map[T]map[T]struct{}

func (r Relationship[T]) GetNeighbors(node T) Neighbors[T] {
	if r == nil {
		return nil
	}

	return r[node]
}
