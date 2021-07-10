package utils

import "testing"

// TestIsPrimitiveType should test IsPrimitiveType
// for primities like int, string etc and composites
func TestIsPrimitiveType(t *testing.T) {

	tests := []struct {
		name string
		t    string
		want bool
	}{
		{
			name: "should return true for string",
			t:    "string",
			want: true,
		},
		{
			name: "should return true for int",
			t:    "int",
			want: true,
		},
		{
			name: "should return true for bool",
			t:    "bool",
			want: true,
		},
		{
			name: "should return true for time",
			t:    "time",
			want: true,
		},
		{
			name: "should return true for float",
			t:    "float",
			want: true,
		},
		{
			name: "should return true for float",
			t:    "float",
			want: true,
		},
		{
			name: "should return false for user",
			t:    "user",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPrimitiveType(tt.t); got != tt.want {
				t.Errorf("IsPrimitiveType() = %v, want %v", got, tt.want)
			}
		})
	}
}
