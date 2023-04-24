package structure_test

import (
	"testing"

	"github.com/nitzanhen/crossword/src/structure"
)

func TestList(t *testing.T) {
	list := structure.List[int]{}

	list.Add(1)
	list.Add(2)
	list.Add(3)
	got := list.Size()
	if got != 3 {
		t.Errorf("list.Size() = %d; want 3", got)
	}

	slice := list.ToSlice()
	if len(slice) != 3 || slice[0] != 1 || slice[1] != 2 || slice[2] != 3 {
		t.Errorf("list.ToSlice() = %v, expected [1, 2, 3]", slice)
	}
}
