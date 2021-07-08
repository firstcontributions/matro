package parser

import "testing"

func Test_validateTypeAndGetFirstNonEmptyIdx(t *testing.T) {

	tests := []struct {
		name    string
		b       []byte
		want    int
		wantErr error
	}{
		{
			name:    "should parse a primitive type (int here)",
			b:       []byte("int"),
			want:    0,
			wantErr: nil,
		},
		{
			name:    "should parse a primitive type with some spaces as prefix (int here)",
			b:       []byte("    int"),
			want:    4,
			wantErr: nil,
		},
		{
			name:    "should parse a composite type with some spaces as prefix",
			b:       []byte("    {}"),
			want:    4,
			wantErr: nil,
		},
		{
			name:    "should parse a composite type without any spaces as prefix",
			b:       []byte("{}"),
			want:    0,
			wantErr: nil,
		},
		{
			name:    "should throw error if empty string",
			b:       []byte(""),
			want:    0,
			wantErr: errEmptyType,
		},
		{
			name:    "should throw error for string with only spaces",
			b:       []byte("         "),
			want:    0,
			wantErr: errEmptyType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateTypeAndGetFirstNonEmptyIdx(tt.b)
			if err != tt.wantErr {
				t.Errorf("validateTypeAndGetFirstNonEmptyIdx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateTypeAndGetFirstNonEmptyIdx() = %v, want %v", got, tt.want)
			}
		})
	}
}
