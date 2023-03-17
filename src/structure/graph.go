package structure

type Graph[V comparable] struct {
	Nodes        Set[V]
	Neighborhood OrderedMap[V, Set[V]]
}

func NewGraph[V comparable](nodes Set[V]) Graph[V] {
	graph := Graph[V]{
		Nodes:        nodes,
		Neighborhood: *NewOrderedMap[V, Set[V]](nodes.Size()),
	}

	for _, node := range nodes.ToSlice() {
		graph.Neighborhood.Set(node, NewSet[V](nodes.Size()))
	}

	return graph
}

func (g *Graph[V]) Connect(v1, v2 V) {
	e1 := g.Neighborhood.Get(v1)
	e2 := g.Neighborhood.Get(v2)

	e1.Add(v2)
	e2.Add(v1)
}

func (g *Graph[V]) Connected(v1, v2 V) bool {
	e1 := g.Neighborhood.Get(v1)
	return e1.Has(v2)
}

// Performs a depth-first search on the graph, starting at the given node.
func (g *Graph[V]) dfs(comp *Set[V], node V) {
	if comp.Has(node) {
		return
	}

	comp.Add(node)
	neighbors := g.Neighborhood.Get(node)
	for _, v := range neighbors.ToSlice() {
		g.dfs(comp, v)
	}
}

func (g *Graph[V]) Components() []Set[V] {

	components := List[Set[V]]{}

	nodes := g.Nodes.Copy()

	var node V
	var component Set[V]

	for nodes.Size() > 0 {

		node = nodes.ToSlice()[0]
		component = NewSet[V](nodes.Size())

		g.dfs(&component, node)

		components.Add(component)
		nodes.Diff(&component)
	}

	return components.ToSlice()
}
