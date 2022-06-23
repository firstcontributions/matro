package utils

// Set implements a set of comparable types
type Set[T comparable] struct {
	set map[T]struct{}
}

// NewSet returns a set of given elements
func NewSet[T comparable](elems ...T) *Set[T] {
	s := &Set[T]{
		set: map[T]struct{}{},
	}
	s.Add(elems...)
	return s
}

// Add adds an element to the set
func (s *Set[T]) Add(elems ...T) {
	for _, e := range elems {
		s.set[e] = struct{}{}
	}
}

// Iter returns an iteratable map
func (s *Set[T]) Iter() map[T]struct{} {
	return s.set
}

// Union executes a set union operatio with given set
func (s *Set[T]) Union(t *Set[T]) {
	for e := range t.Iter() {
		s.Add(e)
	}
}

// IsElem says if the given string is an element of set or not
func (s *Set[T]) IsElem(e T) bool {
	_, ok := s.set[e]
	return ok
}

// Elems returns the array of elements
func (s *Set[T]) Elems() []T {
	elems := []T{}
	for e := range s.Iter() {
		elems = append(elems, e)
	}
	return elems
}
