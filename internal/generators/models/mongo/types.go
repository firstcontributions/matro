package mongo

const modelTyp = `
package {{ .Module.Name -}}store

type {{title .Name}} struct {
	{{- range .ReferedFields}}
	{{ title .}}ID string ` + "`bson:\"{{- . -}}_id\"`" + `
	{{- end}}
	{{- range .Fields}}
	{{- if  (not (and .IsJoinedData  .IsList))}}
	{{- if eq .Name "id" }}
	{{ .GoName true}}  {{ .GoType }}` + "`bson:\"_id\"`" + `
	{{- else}}
	{{ .GoName true}}  {{ .GoType }}` + "`bson:\"{{- .Name}},omitempty\"`" + `  
	{{- end}}
	{{- end}}
	{{- end}}
	{{- if (eq .Module.DB "")}}
	Cursor string
	{{- end}}
}

func New{{- title .Name }}() *{{- title .Name }} {
	return &{{ title .Name }}{}
}

{{- if  .HaveProto}}
func ({{ .Name }} *{{title .Name}}) ToProto() *proto.{{- title .Name -}} {
	{{- $t := . -}}
	
	return &proto.{{- title .Name -}} {
		{{- range .ReferedFields}}
			{{ title .}}Id : {{ $t.Name -}}.{{- title .}}ID ,
		{{- end}}
		{{- range .Fields}}
			{{- if  (not (and .IsJoinedData  .IsList))}}
				{{- if (and (not .IsJoinedData) (isCompositeType .Type))}}
					{{ .GoName true}} :  {{ $t.Name -}}.{{- .GoName true -}}.ToProto(),
				{{- else}}
					{{- if  (eq .Type "time")}}
						{{ .GoName true}} :  timestamppb.New({{ $t.Name -}}.{{- .GoName true}}),
					{{- else }}
						{{ .GoName true}} :  {{ $t.Name -}}.{{- .GoName true}},
					{{- end}}
				{{- end}}
			{{- end}}
		{{- end}}
	}
}

func  ({{ .Name }} *{{title .Name}}) FromProto(proto{{- title .Name }} *proto.{{- title .Name}}) *{{title .Name}}{
	{{- $t := . -}}
	
	{{- range .ReferedFields}}
		{{ $t.Name -}}.{{- title .}}ID = proto{{- title $t.Name -}}.{{- title .}}Id 
	{{- end}}
	{{- range .Fields}}
		{{- if  (not (and .IsJoinedData  .IsList))}}
			{{- if (and (not .IsJoinedData) (isCompositeType .Type))}}
				{{ $t.Name -}}.{{- .GoName true}}   = New{{- title .Type -}}().FromProto(proto{{- title $t.Name -}}.{{- .GoName true}})
			{{- else}}
				{{- if  (eq .Type "time")}}
				{{ $t.Name -}}.{{- .GoName true}}  = proto{{- title $t.Name -}}.{{- .GoName true}}.AsTime()
				{{- else }}
				{{ $t.Name -}}.{{- .GoName true}}   = proto{{- title $t.Name -}}.{{- .GoName true}}
				{{- end}}
			{{- end}}
		{{- end}}
	{{- end}}
	return {{ .Name }}
}
{{- end}}

{{- if (and (ne .Module.DB "") .IsNode)}}
type {{title .Name -}}Update struct {
	{{- range .Fields}}
	{{- if  .IsMutatable }}
	{{ .GoName true}}  {{ .GoType false true}}` + "`bson:\"{{- .Name}},omitempty\"`" + `  
	{{- end}}
	{{- end}}
}
{{- if  .HaveProto}}
func ({{ .Name }} *{{title .Name -}}Update) ToProto() *proto.{{- title .Name -}} {
	{{- $t := . -}}
	p := &proto.{{- title .Name -}} {}
		
	{{- range .Fields}}
	{{- if  .IsMutatable }}
	if {{ $t.Name -}}.{{- .GoName true}} != nil {
		{{- if  (eq .Type "time")}}
		p.{{- .GoName true}} =  timestamppb.New(*{{ $t.Name -}}.{{- .GoName true}})
		{{- else }}
			p.{{- .GoName true}} =  *{{- $t.Name -}}.{{- .GoName true}}
		{{- end}}
	}
	{{- end}}
	{{- end}}
	return p
}
{{- end}}
{{- end}}
`
