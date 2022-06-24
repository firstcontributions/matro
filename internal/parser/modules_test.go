package parser

import "testing"

func TestModule_Store(t *testing.T) {
	type fields struct {
		Name       string
		DataSource string
		DB         string
		Entities   map[string]*Type
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should generate store name",
			fields: fields{
				Name: "user",
			},
			want: "usersstore",
		},
		{
			name: "should generate store name",
			fields: fields{
				Name: "story",
			},
			want: "storiesstore",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				Name:       tt.fields.Name,
				DataSource: tt.fields.DataSource,
				DB:         tt.fields.DB,
				Entities:   tt.fields.Entities,
			}
			if got := m.Store(); got != tt.want {
				t.Errorf("Module.Store() = %v, want %v", got, tt.want)
			}
		})
	}
}
