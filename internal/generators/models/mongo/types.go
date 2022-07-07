package mongo

const modelTyp = `
package {{ .Module.Name -}}store

type {{title .Name}} struct {
	{{- range .ReferedTypes}}
	{{ title .Name}}ID string ` + "`bson:\"{{- .Name -}}_id\"`" + `
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
		{{- range .ReferedTypes}}
			{{ title .Name}}Id : {{ $t.Name -}}.{{- title .Name}}ID ,
		{{- end}}
		{{- range .Fields}}
			{{- if  (not (and .IsJoinedData  .IsList))}}
				{{- if (and (not .IsJoinedData) (isCompositeType .Type))}}
					{{ .GoName true}} :  {{ $t.Name -}}.{{- .GoName true -}}.ToProto(),
				{{- else}}
					{{- if  (eq .Type "time")}}
						{{ .GoName true}} :  timestamppb.New({{ $t.Name -}}.{{- .GoName true}}),
					{{- else }}
						{{- if (and .IsList (eq .Type "string"))}}
						{{ .GoName true}} :  utils.ToStringArray({{ $t.Name -}}.{{- .GoName true}}),
						{{- else}}
						{{ .GoName true}} :  {{ $t.Name -}}.{{- .GoName true}},
						{{- end}}
					{{- end}}
				{{- end}}
			{{- end}}
		{{- end}}
	}
}

func  ({{ .Name }} *{{title .Name}}) FromProto(proto{{- title .Name }} *proto.{{- title .Name}}) *{{title .Name}}{
	{{- $t := . -}}
	
	{{- range .ReferedTypes}}
		{{ $t.Name -}}.{{- title .Name}}ID = proto{{- title $t.Name -}}.{{- title .Name}}Id 
	{{- end}}
	{{- range .Fields}}
		{{- if  (not (and .IsJoinedData  .IsList))}}
			{{- if (and (not .IsJoinedData) (isCompositeType .Type))}}
				{{ $t.Name -}}.{{- .GoName true}}   = New{{- title .Type -}}().FromProto(proto{{- title $t.Name -}}.{{- .GoName true}})
			{{- else}}
				{{- if  (eq .Type "time")}}
				{{ $t.Name -}}.{{- .GoName true}}  = proto{{- title $t.Name -}}.{{- .GoName true}}.AsTime()
				{{- else }}
				{{- if (and .IsList (eq .Type "string"))}}
				{{ $t.Name -}}.{{- .GoName true}}   = utils.FromStringArray(proto{{- title $t.Name -}}.{{- .GoName true}})
				{{- else}}
				{{ $t.Name -}}.{{- .GoName true}}   = proto{{- title $t.Name -}}.{{- .GoName true}}
				{{- end}}
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
func ({{ .Name }} *{{title .Name -}}Update) ToProto() *proto.Update{{- title .Name -}}Request {
	{{- $t := . -}}
	p := &proto.Update{{- title .Name -}}Request {}
		
	{{- range .Fields}}
	{{- if  .IsMutatable }}
	if {{ $t.Name -}}.{{- .GoName true}} != nil {
		{{- if  (not (and .IsJoinedData  .IsList))}}
			{{- if (and (not .IsJoinedData) (isCompositeType .Type))}}
				p.{{- .GoName true}} =  {{ $t.Name -}}.{{- .GoName true -}}.ToProto()
			{{- else}}
				{{- if  (eq .Type "time")}}
				p.{{ .GoName true}} = timestamppb.New(*{{ $t.Name -}}.{{- .GoName true}})
				{{- else }}
					{{- if (and .IsList (eq .Type "string"))}}
					p.{{ .GoName true}} =  utils.ToStringArray({{ $t.Name -}}.{{- .GoName true}})
					{{- else}}
					p.{{ .GoName true}} =  {{ $t.Name -}}.{{- .GoName true}}
					{{- end}}
				{{- end}}
			{{- end}}
		{{- end}}
	}
	{{- end}}
	{{- end}}
	return p
}

func ({{.Name}} *{{title .Name -}}Update) FromProto(proto{{- title .Name }} *proto.Update{{- title .Name -}}Request) {
	{{- $t := . -}}
		
	{{- range .Fields}}
	{{- if  .IsMutatable }}
	{{- if  (not (and .IsJoinedData  .IsList))}}
		{{- if (and (not .IsJoinedData) (isCompositeType .Type))}}
			{{ $t.Name -}}.{{- .GoName true}}   = New{{- title .Type -}}().FromProto(proto{{- title $t.Name -}}.{{- .GoName true}})
		{{- else}}
			{{- if  (eq .Type "time")}}
			{{.Name}} := proto{{- title $t.Name -}}.{{- .GoName true}}.AsTime()
			{{ $t.Name -}}.{{- .GoName true}}  = &{{- .Name}}
			{{- else }}
			{{- if (and .IsList (eq .Type "string"))}}
			{{ $t.Name -}}.{{- .GoName true}}   = utils.FromStringArray(proto{{- title $t.Name -}}.{{- .GoName true}})
			{{- else}}
			{{ $t.Name -}}.{{- .GoName true}}   = proto{{- title $t.Name -}}.{{- .GoName true}}
			{{- end}}
			{{- end}}
		{{- end}}
		{{- end}}
	{{- end}}
	{{- end}}
}
{{- end}}
{{- end}}



type {{ title .Name -}}Filters struct {
	Ids []string
	{{ if not (empty .SearchFields) -}}
		Search *string
	{{ end }}
	{{- $t := . -}}
	{{- range .Filters -}}
		{{ title . }} *{{$t.FieldType . }}
	{{- end -}}
	{{- range .ReferedTypes -}}
		{{- if (eq $t.Module.Name .Module.Name)}}
		{{ title .Name}} *{{- title .Name}}
		{{- else}}
		{{ title .Name}} *{{.Module.Store}}.{{- title .Name}}
		{{- end}}
	{{- end -}}
}
`
