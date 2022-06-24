package utils

import (
	"reflect"
	"sort"
	"testing"
)

func TestNewSet_Int(t *testing.T) {
	type args struct {
		elems []int
	}
	tests := []struct {
		name string
		args args
		want *Set[int]
	}{
		{
			name: "should create a set with given elements",
			args: args{
				elems: []int{1, 2, 3},
			},
			want: &Set[int]{
				map[int]struct{}{
					1: {},
					2: {},
					3: {},
				},
			},
		},
		{
			name: "should create and remove duplicate elements",
			args: args{
				elems: []int{1, 2, 3, 3, 2, 1},
			},
			want: &Set[int]{
				map[int]struct{}{
					1: {},
					2: {},
					3: {},
				},
			},
		},
		{
			name: "should create even if no elements are povided",
			want: &Set[int]{
				map[int]struct{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSet(tt.args.elems...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSet_String(t *testing.T) {
	type args struct {
		elems []string
	}
	tests := []struct {
		name string
		args args
		want *Set[string]
	}{
		{
			name: "should create a set with given elements",
			args: args{
				elems: []string{"A", "B", "c"},
			},
			want: &Set[string]{
				map[string]struct{}{
					"A": {},
					"B": {},
					"c": {},
				},
			},
		},
		{
			name: "should create and remove duplicate elements",
			args: args{
				elems: []string{"1", "2", "3", "3", "2", "1"},
			},
			want: &Set[string]{
				map[string]struct{}{
					"1": {},
					"2": {},
					"3": {},
				},
			},
		},
		{
			name: "should create even if no elements are povided",
			want: &Set[string]{
				map[string]struct{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSet(tt.args.elems...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_String_Add(t *testing.T) {
	type args struct {
		elems []string
	}
	tests := []struct {
		name string
		set  *Set[string]
		args args
		want *Set[string]
	}{
		{
			name: "should add elements to an empty set",
			args: args{
				elems: []string{"A", "B", "c"},
			},
			set: NewSet[string](),
			want: &Set[string]{
				map[string]struct{}{
					"A": {},
					"B": {},
					"c": {},
				},
			},
		},
		{
			name: "should add elements to a set with content",
			args: args{
				elems: []string{"A", "B", "c"},
			},
			set: NewSet("E", "F"),
			want: &Set[string]{
				map[string]struct{}{
					"A": {},
					"B": {},
					"c": {},
					"E": {},
					"F": {},
				},
			},
		},
		{
			name: "should add elements by eliminating duplicates",
			args: args{
				elems: []string{"A", "B", "c", "A", "B"},
			},
			set: NewSet("E", "F"),
			want: &Set[string]{
				map[string]struct{}{
					"A": {},
					"B": {},
					"c": {},
					"E": {},
					"F": {},
				},
			},
		},
		{
			name: "should not add existing elements again",
			args: args{
				elems: []string{"A", "B", "c", "E", "F"},
			},
			set: NewSet("E", "F"),
			want: &Set[string]{
				map[string]struct{}{
					"A": {},
					"B": {},
					"c": {},
					"E": {},
					"F": {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := tt.set
			set.Add(tt.args.elems...)
			if !reflect.DeepEqual(set, tt.want) {
				t.Errorf("Set = %v, want %v", set, tt.want)
			}
		})
	}
}

func TestSet_Int_Iter(t *testing.T) {
	tests := []struct {
		name string
		set  *Set[int]
		want []int
	}{
		{
			name: "should return an iterable object with all elements in the set",

			set:  NewSet(1, 2, 3, 4, 5),
			want: []int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elems := []int{}
			for item := range tt.set.Iter() {
				elems = append(elems, item)
			}
			sort.Ints(elems)
			sort.Ints(tt.want)

			if !reflect.DeepEqual(elems, tt.want) {
				t.Errorf("set.Iter() expected = %v, got %v", elems, tt.want)
			}
		})
	}
}
