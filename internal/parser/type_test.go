package parser

import (
	"errors"
	"testing"

	"github.com/firstcontributions/matro/internal/merrors"
	"github.com/go-test/deep"
)

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
			wantErr: merrors.ErrEmptyType,
		},
		{
			name:    "should throw error for string with only spaces",
			b:       []byte("         "),
			want:    0,
			wantErr: merrors.ErrEmptyType,
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

func TestOps_Create(t *testing.T) {

	tests := []struct {
		name string
		op   *Ops
		want bool
	}{
		{
			name: "should return false if object is nil",
			op:   nil,
			want: false,
		},
		{
			name: "should return false if create field is false",
			op:   &Ops{create: false},
			want: false,
		},
		{
			name: "should return true if create field is true",
			op:   &Ops{create: false},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.op.Create(); got != tt.want {
				t.Errorf("Ops.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOps_Read(t *testing.T) {

	tests := []struct {
		name string
		op   *Ops
		want bool
	}{
		{
			name: "should return false if object is nil",
			op:   nil,
			want: false,
		},
		{
			name: "should return false if read field is false",
			op:   &Ops{read: false},
			want: false,
		},
		{
			name: "should return true if read field is true",
			op:   &Ops{read: false},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.op.Read(); got != tt.want {
				t.Errorf("Ops.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOps_Update(t *testing.T) {

	tests := []struct {
		name string
		op   *Ops
		want bool
	}{
		{
			name: "should return false if object is nil",
			op:   nil,
			want: false,
		},
		{
			name: "should return false if update field is false",
			op:   &Ops{update: false},
			want: false,
		},
		{
			name: "should return true if update field is true",
			op:   &Ops{update: false},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.op.Update(); got != tt.want {
				t.Errorf("Ops.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOps_Delete(t *testing.T) {

	tests := []struct {
		name string
		op   *Ops
		want bool
	}{
		{
			name: "should return false if object is nil",
			op:   nil,
			want: false,
		},
		{
			name: "should return false if delete field is false",
			op:   &Ops{delete: false},
			want: false,
		},
		{
			name: "should return true if delete field is true",
			op:   &Ops{delete: false},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.op.Delete(); got != tt.want {
				t.Errorf("Ops.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOps_Union(t *testing.T) {

	tests := []struct {
		name string
		op1  Ops
		op2  Ops
		want Ops
	}{
		{
			name: "should take union of all operations",
			op1: Ops{
				create: true,
				read:   true,
				update: false,
				delete: false,
			},
			op2: Ops{
				create: false,
				read:   true,
				update: true,
				delete: false,
			},
			want: Ops{
				create: true,
				read:   true,
				update: true,
				delete: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := tt.op1
			op.Union(tt.op2)
			if diff := deep.Equal(op, tt.want); diff != nil {
				t.Errorf("Ops.Union() %v", diff)
			}
		})
	}
}

func TestOps_UnmarshalJSON(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    Ops
	}{
		{
			name:    "should parse create",
			input:   "C",
			wantErr: false,
			want:    Ops{create: true},
		},
		{
			name:    "should parse read",
			input:   "R",
			wantErr: false,
			want:    Ops{read: true},
		},
		{
			name:    "should parse update",
			input:   "U",
			wantErr: false,
			want:    Ops{update: true},
		},
		{
			name:    "should parse delete",
			input:   "D",
			wantErr: false,
			want:    Ops{delete: true},
		},
		{
			name:    "should parse create and update",
			input:   "CU",
			wantErr: false,
			want:    Ops{create: true, update: true},
		},
		{
			name:    "should parse read and update",
			input:   "RU",
			wantErr: false,
			want:    Ops{read: true, update: true},
		},
		{
			name:    "should parse update and delete",
			input:   "U",
			wantErr: false,
			want:    Ops{update: true, delete: true},
		},
		{
			name:    "should parse create and delete",
			input:   "CD",
			wantErr: false,
			want:    Ops{delete: true, create: true},
		},
		{
			name:    "should parse all",
			input:   "CRUD",
			wantErr: false,
			want:    Ops{create: true, read: true, delete: true, update: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := Ops{}
			if err := op.UnmarshalJSON([]byte(tt.input)); (err != nil) != tt.wantErr {
				t.Errorf("Ops.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := deep.Equal(op, tt.want); diff != nil {
				t.Errorf("Ops.UnmarshalJSON(), %v", diff)
			}
		})
	}
}

func TestType_IsPrimitive(t *testing.T) {
	tests := []struct {
		name   string
		fields _Type
		want   bool
	}{
		{
			name: "should check for int",
			fields: _Type{
				Type: "int",
			},
			want: true,
		},
		{
			name: "should check for id",
			fields: _Type{
				Type: "id",
			},
			want: true,
		},
		{
			name: "should check for a composite type user",
			fields: _Type{
				Type: "user",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Type{
				_Type: tt.fields,
			}
			if got := tr.IsPrimitive(); got != tt.want {
				t.Errorf("Type.IsPrimitive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Type
		wantErr error
	}{
		{
			name:    "should return empty type error if definition is empty",
			want:    Type{},
			wantErr: merrors.ErrEmptyType,
		},
		{
			name: "should parse static type names",
			args: "string",
			want: Type{
				_Type: _Type{
					Type: "string",
				},
			},
			wantErr: nil,
		},
		{
			name:    "should raise invalid type definition error if the object definition is not valid",
			args:    `{"Type": asdf`,
			want:    Type{},
			wantErr: merrors.ErrInvalidTypeDefinition,
		},
		{
			name: "should parse valid type definition in object format",
			args: `{"type": "int", "paginated": true}`,
			want: Type{
				_Type: _Type{
					Type:      "int",
					Paginated: true,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Type{}
			if err := tr.UnmarshalJSON([]byte(tt.args)); !errors.Is(err, tt.wantErr) {
				t.Errorf("Type.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := deep.Equal(tr, tt.want); diff != nil {
				t.Errorf("Type.UnmarshalJSON()  %v", diff)
			}
		})
	}
}

func TestType_Validate(t *testing.T) {

	tests := []struct {
		name    string
		t       *Type
		wantErr error
	}{
		{
			name:    "type name cannot be empty",
			t:       &Type{},
			wantErr: merrors.ErrEmptyType,
		},
		{
			name:    "custom type should have a schema or properties",
			t:       &Type{_Type: _Type{Type: "list"}},
			wantErr: merrors.ErrNoSchemaDefinedForCustomType,
		},
		{
			name: "type name cannot be empty",
			t: &Type{_Type: _Type{
				Type: "list",
				Properties: map[string]*Type{
					"id": {_Type: _Type{Type: "id"}},
				},
			}},
			wantErr: merrors.ErrNoNameForInlineCustomTypes,
		},
		{
			name: "should report undefined search fields",
			t: &Type{_Type: _Type{
				Name: "user",
				Type: "list",
				Meta: Meta{
					SearchFields: []string{"name"},
				},
				Properties: map[string]*Type{
					"id": {_Type: _Type{Type: "id"}},
				},
			}},
			wantErr: merrors.ErrUndefinedSearchField,
		},
		{
			name: "should report undefined filter fields",
			t: &Type{_Type: _Type{
				Name: "user",
				Type: "list",
				Meta: Meta{
					Filters: []string{"name"},
				},
				Properties: map[string]*Type{
					"id": {_Type: _Type{Type: "id"}},
				},
			}},
			wantErr: merrors.ErrUndefinedFilter,
		},
		{
			name: "should report undefined mutable fields",
			t: &Type{_Type: _Type{
				Name: "user",
				Type: "list",
				Meta: Meta{
					MutatableFields: []string{"name"},
				},
				Properties: map[string]*Type{
					"id": {_Type: _Type{Type: "id"}},
				},
			}},
			wantErr: merrors.ErrUndefinedMutableField,
		},
		{
			name: "should report undefined sort by fields",
			t: &Type{_Type: _Type{
				Name: "user",
				Type: "list",
				Meta: Meta{
					SortBy: []string{"name"},
				},
				Properties: map[string]*Type{
					"id": {_Type: _Type{Type: "id"}},
				},
			}},
			wantErr: merrors.ErrUndefinedSortField,
		},
		{
			name: "should not raise any errors if all the given fields are defined",
			t: &Type{
				_Type: _Type{
					Name: "user",
					Type: "list",
					Meta: Meta{
						MutatableFields: []string{"id"},
					},
					Properties: map[string]*Type{
						"id": {_Type: _Type{Type: "id"}},
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.t.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("Type.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestType_raiseErrorIfFieldsNotDefined(t *testing.T) {

	type args struct {
		fields []string
		err    error
	}
	tests := []struct {
		name    string
		t       *Type
		args    args
		wantErr error
	}{
		{
			name: "should not raise any errors if all the given fields are defined",
			t: &Type{
				_Type: _Type{
					Properties: map[string]*Type{
						"id": {_Type: _Type{Type: "id"}},
					},
				},
			},
			args: args{
				fields: []string{"id"},
				err:    merrors.ErrUndefinedFilter,
			},
			wantErr: nil,
		},
		{
			name: "should raise given error if field not defined",
			t:    &Type{},
			args: args{
				fields: []string{"name"},
				err:    merrors.ErrUndefinedFilter,
			},
			wantErr: merrors.ErrUndefinedFilter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.raiseErrorIfFieldsNotDefined(tt.args.fields, tt.args.err); !errors.Is(err, tt.wantErr) {
				t.Errorf("Type.raiseErrorIfFieldsNotDefined() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
