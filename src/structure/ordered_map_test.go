package structure_test

import (
	"nitzanhen/crossword/src/structure"
	"testing"
)

type Pair struct {
	a, b int
}

func TestOrderedMap(t *testing.T) {
	m := structure.NewOrderedMap[Pair, int](2)

	p1 := Pair{1, 2}

	m.Set(p1, 3)
	size := m.Size()
	if size != 1 {
		t.Errorf("Expected m.Size() = 1, got %d", size)
	}

	p2 := Pair{1, 2}
	m.Set(p2, 2)

	if size := m.Size(); size != 1 {
		t.Errorf("Expected m.Size() = 1, got %d", size)
	}

	if value := m.Get(p2); value != 2 {
		t.Errorf("Expected m.Get(p2) = 2, got %d", value)
	}

	p3 := Pair{0, 3}
	m.Set(p3, 1)

	if keys := m.Keys(); len(keys) != 2 || keys[0] != p1 || keys[1] != p3 {
		t.Errorf("Expected keys=[p1, p3], got %v", keys)
	}

	if values := m.Values(); len(values) != 2 || values[0] != 2 || values[1] != 1 {
		t.Errorf("Expected values=[2, 1], got %v", values)
	}
}
