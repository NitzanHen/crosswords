package structure_test

import (
	"testing"

	"github.com/nitzanhen/crossword/src/structure"
)

func TestSet(t *testing.T) {
	set := structure.NewSet[int](10)

	set.Add(1).Add(2)

	if size := set.Size(); size != 2 {
		t.Errorf("Expected set.Size() = 2, got %d", size)
	}

	if has := set.Has(1); !has {
		t.Errorf("Expected set.Has(1) = true, got %t", has)
	}

	set.Delete(2)

	if has := set.Has(2); has {
		t.Errorf("Expected set.Has(2) = false, got %t", has)
	}

	if slice := set.ToSlice(); len(slice) != 1 || slice[0] != 1 {
		t.Errorf("Expected set.ToSlice() = [1], got %v", slice)
	}

	copy := set.Copy()
	if size, hasOne := copy.Size(), copy.Has(1); size != 1 || !hasOne {
		t.Errorf("Expected copy (%+v) to be same as original (%+v)", copy, set)
	}

	set = structure.NewSet[int](4)
	set.Add(1).Add(2).Add(3).Add(4)

	set.Filter(func(i int) bool {
		return i%2 == 0
	})

	if has := set.Has(3); has {
		t.Errorf("Expected set.Has(3) = false, got %t", has)
	}
	if has := set.Has(2); !has {
		t.Errorf("Expected set.Has(2) = true, got %t", has)
	}

	set.Union(copy)

	if has := set.Has(1); !has {
		t.Errorf("Expected set.Has(1) = true, got %t", has)
	}

	set1 := structure.SetFromSlice([]int{1, 2, 3, 4})
	set2 := structure.SetFromSlice([]int{1, 2})
	set1.Diff(&set2)

	if set1.Has(1) || set1.Has(2) {
		t.Errorf("Expected set1 to not have 1 or 2")
	}
}
