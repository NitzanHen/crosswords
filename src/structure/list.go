package structure

type node[T any] struct {
	value T
	next  *node[T]
}

type List[T any] struct {
	size int
	head *node[T]
	tail *node[T]
}

func ListFromSlice[T any](slice []T) List[T] {
	list := List[T]{}
	for _, value := range slice {
		list.Add(value)
	}

	return list
}

func (list *List[T]) Size() int {
	return list.size
}

// Adds a value to the end of the list
func (list *List[T]) Add(value T) {
	node := node[T]{value, nil}

	if list.head == nil {
		list.head = &node
	}
	if tail := list.tail; tail != nil {
		tail.next = &node
	}
	list.tail = &node

	list.size++
}

func (list *List[T]) ToSlice() []T {
	slice := make([]T, list.size)

	node := list.head
	for i := 0; i < list.size; i++ {
		slice[i] = node.value
		node = node.next
	}

	return slice
}
