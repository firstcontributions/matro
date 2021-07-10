package types

import "testing"

func Test_GetGraphQLType(t *testing.T) {

	tests := []struct {
		name string
		t    string
		want string
	}{
		{
			name: "should find a valid type for a primitive type (id)",
			t:    "id",
			want: "ID",
		},
		{
			name: "should find a valid type for a composite type (id)",
			t:    "user",
			want: "User",
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
