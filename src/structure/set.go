package structure

type Set[T comparable] struct {
	members *OrderedMap[T, bool]
}

func NewSet[T comparable](cap int) Set[T] {
	members := NewOrderedMap[T, bool](cap)

	return Set[T]{members}
}

func SetFromSlice[T comparable](slice []T) Set[T] {
	set := NewSet[T](len(slice))
	for _, element := range slice {
		set.Add(element)
	}

	return set
}

func (s *Set[T]) Size() int {
	return s.members.Size()
}

// Adds an element to the set. Returns the set.
func (s *Set[T]) Add(element T) *Set[T] {
	s.members.Set(element, true)

	return s
}

// Deletes an element from the set. Returns the set.
func (s *Set[T]) Delete(element T) *Set[T] {
	s.members.Delete(element)

	return s
}

// Checks and returns whether the set has the given element
func (s *Set[T]) Has(element T) bool {
	return s.members.Has(element)
}

// Adds all of the other set's elements to the current set.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	for _, element := range other.members.Keys() {
		s.Add(element)
	}

	return s
}

// Removes from this set all the elements that for which the predicate returns false
func (s *Set[T]) Filter(predicate func(el T) bool) *Set[T] {
	for _, item := range s.members.Keys() {
		if !predicate(item) {
			s.Delete(item)
		}
	}

	return s
}

// Removes from this set all the elements that are not in the other set.
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	return s.Filter(func(el T) bool { return other.Has(el) })
}

// Removes from this set all the elements that are also in the other set.
func (s *Set[T]) Diff(other *Set[T]) *Set[T] {
	return s.Filter(func(el T) bool { return !other.Has(el) })
}

// Creates and returns a new copy of the current set.
func (s *Set[T]) Copy() *Set[T] {
	copy := NewSet[T](s.members.Size())
	copy.Union(s)

	return &copy
}

func (s *Set[T]) ToSlice() []T {
	list := List[T]{}
	for _, item := range s.members.Keys() {
		list.Add(item)
	}

	return list.ToSlice()
}
