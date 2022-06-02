package service

const crudTmpl = `
package service

func (s *Service) Create{{- title .Name -}} (ctx context.Context, in *proto.{{- title .Name }}) (*proto.{{- title .Name }}, error) {

	{{.Name}} := {{.Module.Store}}.New{{- title .Name}}()
	{{.Name}}.FromProto(in)
	res, err := s.Store.Create{{- title .Name -}}(ctx, {{.Name}})
	if err != nil {
		return nil, err
	}
	return res.ToProto(), nil
}

func (s *Service) Get{{- title .Name -}}ByID(ctx context.Context, in *proto.RefByIDRequest) (*proto.{{- title .Name }}, error) {
	{{.Name}}, err := s.Store.Get{{- title .Name -}}ByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return {{.Name}}.ToProto(), nil
}

func (s *Service) Get{{- (plural (title .Name))}}(ctx context.Context, in *proto.Get{{- (plural (title .Name))}}Request) (*proto.Get{{- (plural (title .Name))}}Response, error) {
	{{- range .ReferedTypes -}}
	var {{.Name}} *{{.Module.Store}}.{{- title .Name}}
	if in.{{- title .Name -}}Id != nil {
		{{.Name}} = &{{.Module.Store}}.{{- title .Name}}{Id: *in.{{- title .Name -}}Id}
	}
	{{- end }}
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err :=  s.Store.Get{{- plural (title .Name)}} (
		ctx,
		in.Ids,
		{{- if not (empty .SearchFields)}}
		in.Search,
		{{- end}}
		
		{{- range .Filters -}}
			in.{{ title . }},
		{{- end }}
		{{- range .ReferedTypes }}
			{{.Name}},
		{{- end }}
		in.After, 
		in.Before, 
		in.First ,
		in.Last,
	)
	if err != nil {
		return nil, err
	}
	res := []*proto.{{- title .Name}}{}
	for _, d := range data {
		res = append(res, d.ToProto())
	}
	return &proto.Get{{- (plural (title .Name))}}Response{
		HasNext: hasNextPage,
		HasPrevious: hasPreviousPage,
		FirstCursor: firstCursor,
		LastCursor: lastCursor,
		Data: res,
	}, nil
}
func (s *Service) Update{{- title .Name -}}(ctx context.Context, in *proto.Update{{- title .Name -}}Request) (*proto.StatusResponse, error) {
	update{{- title .Name -}} := &{{.Module.Store -}}.{{- title .Name -}}Update{}
	update{{- title .Name -}}.FromProto(in)
	
	err := s.Store.Update{{- title .Name -}}(ctx, in.Id, update{{- title .Name -}})
	if err != nil {
		return nil, err
	}
	return &proto.StatusResponse{Status: true}, nil
}
func (s *Service) Delete{{- title .Name -}}(ctx context.Context, in *proto.RefByIDRequest) (*proto.StatusResponse, error) {

	err := s.Store.Delete{{- title .Name -}}ByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &proto.StatusResponse{Status: true}, nil
}
`
