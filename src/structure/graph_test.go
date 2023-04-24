package structure_test

import (
	"testing"

	"github.com/nitzanhen/crossword/src/structure"
)

func TestGraph(t *testing.T) {
	a, b, c, d := "a", "b", "c", "d"
	nodes := structure.SetFromSlice([]string{a, b, c, d})

	graph := structure.NewGraph(nodes)

	for _, node := range nodes.ToSlice() {
		if !graph.Neighborhood.Has(node) {

			t.Errorf("Expected graph.Neighborhood to have node %s", node)
		} else if neighbors := graph.Neighborhood.Get(node); neighbors.Size() != 0 {
			t.Errorf("Expected node %s to have 0 neighbors, got %d", node, neighbors.Size())
		}
	}

	graph.Connect(a, b)
	if aNeighbors := graph.Neighborhood.Get(a); !aNeighbors.Has(b) {
		t.Errorf("Expected %s to be a neighbor of %s", b, a)
	} else if bNeighbors := graph.Neighborhood.Get(b); !bNeighbors.Has(a) {
		t.Errorf("Expected %s to be a neighbor of %s", a, b)
	}
	if !graph.Connected(a, b) {
		t.Errorf("Expected %s and %s to be connected", a, b)
	}

	graph.Connect(c, d)

	components := graph.Components()
	if num := len(components); num != 2 {
		t.Errorf("Expected graph to have 2 components, got %d", num)
	}

	var comp1, comp2 structure.Set[string]
	if components[0].Has(a) {
		comp1, comp2 = components[0], components[1]
	} else {
		comp1, comp2 = components[1], components[0]
	}

	if comp1.Size() != 2 || !comp1.Has(a) || !comp1.Has(b) {
		t.Errorf("Expected first component to be Set{%q, %q}, got %v", a, b, comp1.ToSlice())
	}
	if comp2.Size() != 2 || !comp2.Has(c) || !comp2.Has(d) {
		t.Errorf("Expected first component to be Set{%q, %q}, got %v", c, d, comp2.ToSlice())
	}
}
