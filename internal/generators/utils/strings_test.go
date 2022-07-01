package utils

import (
	"testing"
)

func TestToTitleCase(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "simple lowercase word",
			s:    "simple",
			want: "Simple",
		},
		{
			name: "simple uppercase word",
			s:    "SIMPLE",
			want: "SIMPLE",
		},
		{
			name: "camelcase word",
			s:    "simpleTestCase",
			want: "SimpleTestCase",
		},
		{
			name: "snake_case word",
			s:    "snakes_are_venemous_not_poisonous",
			want: "SnakesAreVenemousNotPoisonous",
		},
		{
			name: "kebab-case-string",
			s:    "chicken-kebab-is-yum",
			want: "ChickenKebabIsYum",
		},
		{
			name: "alpha numeric strings",
			s:    "cOmPlicat3d_-3ord5",
			want: "COmPlicat3d3ord5",
		},
		{
			name: "reverse titlecase",
			s:    "sIMPLE",
			want: "SIMPLE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToTitleCase(tt.s); got != tt.want {
				t.Errorf("ToTitleCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToTitle(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "simple word",
			s:    "simple",
			want: "Simple",
		},
		{
			name: "number",
			s:    "12143",
			want: "12143",
		},
		{
			name: "non alpha numeric",
			s:    "!$##--",
			want: "!$##--",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToTitle(tt.s); got != tt.want {
				t.Errorf("ToTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name string
		s string
		want string
	}{
		{
			name: "simple lowercase word",
			s:    "simple",
			want: "simple",
		},
		{
			name: "simple uppercase word",
			s:    "SIMPLE",
			want: "sIMPLE",
		},
		{
			name: "camelcase word",
			s:    "simpleTestCase",
			want: "simpleTestCase",
		},
		{
			name: "snake_case word",
			s:    "snakes_are_venemous_not_poisonous",
			want: "snakesAreVenemousNotPoisonous",
		},
		{
			name: "kebab-case-string",
			s:    "chicken-kebab-is-yum",
			want: "chickenKebabIsYum",
		},
		{
			name: "alpha numeric strings",
			s:    "cOmPlicat3d_-3ord5",
			want: "cOmPlicat3d3ord5",
		},
		{
			name: "reverse titlecase",
			s:    "sIMPLE",
			want: "sIMPLE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCamelCase(tt.s); got != tt.want {
				t.Errorf("ToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
