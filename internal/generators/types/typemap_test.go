package types

import (
	"testing"
)

func Test_GetGraphQLType(t *testing.T) {

	tests := []struct {
		name string
		t    *Field
		want string
	}{
		{
			name: "should find a valid type for a primitive type (id)",
			t:    &Field{Type: "id"},
			want: "ID",
		},
		{
			name: "should find a valid type for a composite type",
			t:    &Field{Type: "user"},
			want: "User",
		},
		{
			name: "should find a valid type for a composite type array",
			t:    &Field{Type: "user", IsList: true},
			want: "[User]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGraphQLType(tt.t); got != tt.want {
				t.Errorf("GetGraphQLType() = %v, want %v", got, tt.want)
			}
		})
	}
}
