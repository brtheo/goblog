package octogo

import (
	"sort"
)

type By[T any] func(p1, p2 *T) bool
type Sorter[T any] struct {
	slice []*T
	by    func(p1, p2 *T) bool // Closure used in the Less method.
}

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By[T]) Sort(_slice []*T) {
	ps := &Sorter[T]{
		slice: _slice,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// Len is part of sort.Interface.
func (s *Sorter[T]) Len() int {
	return len(s.slice)
}

// Swap is part of sort.Interface.
func (s *Sorter[T]) Swap(i, j int) {
	s.slice[i], s.slice[j] = s.slice[j], s.slice[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *Sorter[T]) Less(i, j int) bool {
	return s.by(s.slice[i], s.slice[j])
}
