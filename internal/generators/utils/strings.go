package utils

import "strings"

func ToTitleCase(s string) string {
	parts := strings.Split(s, "_")
	for i, v := range parts {
		parts[i] = strings.Title(v)
	}
	return strings.Join(parts, "")
}

type Set struct {
	set map[string]struct{}
}

func NewSet(elems ...string) *Set {
	s := &Set{
		set: map[string]struct{}{},
	}
	s.Add(elems...)
	return s
}

func (s *Set) Add(elems ...string) {
	for _, e := range elems {
		s.set[e] = struct{}{}
	}
}
func (s *Set) Iter() map[string]struct{} {
	return s.set
}

func (s *Set) Union(t *Set) {
	for e := range t.Iter() {
		s.Add(e)
	}
}
func (s *Set) IsElem(e string) bool {
	_, ok := s.set[e]
	return ok
}

func (s *Set) Elems() []string {
	elems := []string{}
	for e := range s.Iter() {
		elems = append(elems, e)
	}
	return elems
}
