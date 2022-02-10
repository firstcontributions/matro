package mongo

const crudTpl = `
package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gokultp/go-mongoqb"
	"{{- .Repo -}}/internal/models/{{- .Module -}}store"
)

func (s *{{- title .Module -}}Store) Create{{- title .Name -}} (ctx context.Context, {{.Name}} *{{-  .Module -}}store. {{- title .Name}}) (* {{ .Module -}}store. {{- title .Name}}, error) {
	now := time.Now()
	{{.Name -}}.TimeCreated = now
	{{.Name -}}.TimeUpdated = now
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	{{.Name -}}.Id = uuid.String()
	if _, err := s.getCollection(Collection{{title (plural .Name)}}).InsertOne(ctx, {{.Name}}); err != nil {
		return nil, err
	}
	return {{ .Name}}, nil
}

func (s *{{- title .Module -}}Store) Get{{- title .Name -}}ByID (ctx context.Context, id string) (* {{ .Module -}}store. {{- title .Name}}, error) {
	qb := mongoqb.NewQueryBuilder().
			Eq("_id", id)
	var {{.Name}} {{ .Module -}}store. {{- title .Name}}
	if err := s.getCollection(Collection{{title (plural .Name)}}).FindOne(ctx, qb.Build()).Decode(&{{- .Name -}}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &{{- .Name}}, nil
}

func (s *{{- title .Module -}}Store) Get{{- title (plural .Name) -}} (
	ctx context.Context,
	ids []string,
	{{- if not (empty .SearchFields)}}
	search *string,
	{{- end}}
	{{- template "getargs" . }}
	after *string,
	before *string,
	first  *int64, 
	last  *int64,
) (
	[]*{{ .Module -}}store. {{- title .Name}}, 
	bool,
	bool,
	string,
	string,
	error,
) {
	qb := mongoqb.NewQueryBuilder()
	if len(ids) > 0 {
		qb.In("_id", ids)
	}
	{{- range .Filters }}
	if {{ camel .}} != nil {
		qb.Eq("{{- . -}}", {{ camel .}})
	}
	{{- end }}

	{{- range .ReferedFields }}
		if {{. -}}ID != nil {
			qb.Eq("{{- . -}}_id", {{. -}}ID)
		}
	{{- end }}
	{{- if not (empty .SearchFields)}}
	if search != nil {
		qb.Search(*search)
	}
	{{- end}}
	
	limit, order, cursorStr := utils.GetLimitAndSortOrderAndCursor(first, last, after, before)
	var c *cursor.Cursor
	if cursorStr != nil {
		c = cursor.FromString(*cursorStr)
		if c != nil {
			if order == 1 {
				qb.Lte("time_created", c.TimeStamp)
				qb.Lte("_id", c.ID)
			} else {
				qb.Gte("time_created", c.TimeStamp)
				qb.Gte("_id", c.ID)
			}
		}
	}
	sortOrder := utils.GetSortOrder(order)
	// incrementing limit by 2 to check if next, prev elements are present
	limit += 2
	options := &options.FindOptions{
		Limit: &limit,
		Sort:  sortOrder,
	}

	var firstCursor, lastCursor string
	var hasNextPage, hasPreviousPage bool

	var {{plural .Name}} []*{{ .Module -}}store. {{- title .Name}}
	mongoCursor, err := s.getCollection(Collection{{title (plural .Name)}}).Find(ctx, qb.Build(), options)
	if  err != nil {
		return nil, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err
	}
	err = mongoCursor.All(ctx, &{{- plural .Name}})
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err
	}
	count := len({{ plural .Name}})
	if count == 0 {
		return {{ plural .Name}}, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
	}

	// check if the cursor element present, if yes that can be a prev elem
	if c != nil && {{ plural .Name -}}[0].Id == c.ID {
		hasPreviousPage = true
		{{ plural .Name }} = {{ plural .Name -}}[1:]
		count--
	}

	// check if actual limit +1 elements are there, if yes trim it to limit
	if count >= int(limit)-1 {
		hasNextPage = true
		{{ plural .Name }} = {{ plural .Name -}}[:limit-2]
		count = len({{ plural .Name }})
	}

	if  count > 0 {
		firstCursor = cursor.NewCursor({{ plural .Name -}}[0].Id, {{ plural .Name -}}[0].TimeCreated).String()
		lastCursor = cursor.NewCursor({{ plural .Name -}}[count-1].Id, {{ plural .Name -}}[count-1].TimeCreated).String()
	}
	if order < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		firstCursor, lastCursor = lastCursor, firstCursor
	}
	return {{ plural .Name}}, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
}
{{- if .Mutatable}}
func (s *{{- title .Module -}}Store) Update{{- title .Name -}} (ctx context.Context, id string, {{.Name -}}Update *{{-  .Module -}}store. {{- title .Name -}}Update) (error) {
	qb := mongoqb.NewQueryBuilder().
			Eq("_id", id)

	now := time.Now()
	{{.Name -}}Update.TimeUpdated = &now

	u := mongoqb.NewUpdateMap().
	SetFields({{- .Name -}}Update)

	um, err := u.BuildUpdate()
	if err != nil {
		return err
	}
	if _, err := s.getCollection(Collection{{title (plural .Name)}}).UpdateOne(ctx, qb.Build(), um); err != nil {
		return err
	}
	return nil
}
{{- end}}
func (s *{{- title .Module -}}Store) Delete{{- title .Name -}}ByID (ctx context.Context, id string) (error) {
	qb := mongoqb.NewQueryBuilder().
			Eq("_id", id)
	if _,  err := s.getCollection(Collection{{title (plural .Name)}}).DeleteOne(ctx, qb.Build()); err != nil {
		return  err
	}
	return nil
}

{{- define "getargs"}}
{{- $t := .}}
{{- range .ReferedFields }}
	{{. -}}ID *string,
{{- end }}
{{- range .Filters}}
	{{ camel .}} *{{$t.FieldType .}},
{{- end}}
{{- end}}
`
