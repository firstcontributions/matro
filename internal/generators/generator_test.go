package generators

import (
	"reflect"
	"testing"

	"github.com/firstcontributions/matro/internal/generators/graphql/schema"
	"github.com/firstcontributions/matro/internal/parser"
)

// TestGetGenerator should test all possibilies of generator factory
func TestGetGenerator(t *testing.T) {
	d := &parser.Definition{
		Modules: []parser.Module{},
	}
	path := "./"
	type args struct {
		t    Type
		path string
		d    *parser.Definition
	}
	tests := []struct {
		name string
		args args
		want IGenerator
	}{
		{
			name: "should return a gql schema generator",
			args: args{
				path: path,
				d:    d,
				t:    TypeGQLSchema,
			},
			want: schema.NewGenerator(path, d),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGenerator(tt.args.t, tt.args.path, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}
