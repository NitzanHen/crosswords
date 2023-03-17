package crossword

import (
	"nitzanhen/crossword/src/structure"
)

func Filter[T any](items []T, predicate func(T) bool) []T {
	filtered := structure.List[T]{}

	for _, item := range items {
		if predicate(item) {
			filtered.Add(item)
		}
	}

	return filtered.ToSlice()
}

func Map[T any, S any](items []T, transform func(T) S) []S {
	transformed := make([]S, len(items))
	for i, item := range items {
		transformed[i] = transform(item)
	}

	return transformed
}

func MakeMatrix[T any](rows, cols int, factory func(i, j int) T) [][]T {
	matrix := make([][]T, rows)
	for i := range matrix {
		matrix[i] = make([]T, cols)
		for j := range matrix[i] {
			matrix[i][j] = factory(i, j)
		}
	}

	return matrix
}

func Chars(str string) []string {
	runes := []rune(str)

	return Map(runes, func(r rune) string {
		return string(r)
	})
}

func FirstIndex[T any](arr []T, predicate func(el T) bool) int {
	for i := 0; i < len(arr); i++ {
		if predicate(arr[i]) {
			return i
		}
	}
	return -1
}

func LastIndex[T any](arr []T, predicate func(el T) bool) int {
	for i := len(arr) - 1; i >= 0; i-- {
		if predicate(arr[i]) {
			return i
		}
	}
	return -1
}

func IndexArray(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = i
	}

	return arr
}
