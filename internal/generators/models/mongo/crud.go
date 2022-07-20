package mongo

const crudTpl = `
package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gokultp/go-mongoqb"
	"{{- .Repo -}}/internal/models/{{- .Module.Store -}}"
)
func {{ .Name -}}FiltersToQuery(filters *{{ .Module.Store -}}. {{- title .Name -}}Filters) *mongoqb.QueryBuilder {
	qb := mongoqb.NewQueryBuilder()
	if len(filters.Ids) > 0 {
		qb.In("_id", filters.Ids)
	}
	{{- range .Filters }}
	if filters.{{- title .}} != nil {
		qb.Eq("{{- . -}}", filters.{{- title .}})
	}
	{{- end }}

	{{- range .ReferedTypes }}
		if filters.{{- title .Name -}} != nil {
			qb.Eq("{{- .Name -}}_id", filters.{{- title .Name -}}.Id)
		}
	{{- end }}
	{{- if not (empty .SearchFields)}}
	if filters.Search != nil {
		qb.Search(*filters.Search)
	}
	{{- end}}
	return qb
}
func (s *{{- title .Module.Name -}}Store) Create{{- title .Name -}} (ctx context.Context, {{.Name}} *{{-  .Module.Store -}}. {{- title .Name}}) (* {{ .Module.Store -}}. {{- title .Name}}, error) {
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

func (s *{{- title .Module.Name -}}Store) Get{{- title .Name -}}ByID (ctx context.Context, id string) (* {{ .Module.Store -}}. {{- title .Name}}, error) {
	qb := mongoqb.NewQueryBuilder().
			Eq("_id", id)
	var {{.Name}} {{ .Module.Store -}}. {{- title .Name}}
	if err := s.getCollection(Collection{{title (plural .Name)}}).FindOne(ctx, qb.Build()).Decode(&{{- .Name -}}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &{{- .Name}}, nil
}

func (s *{{- title .Module.Name -}}Store) GetOne{{- title .Name -}} (ctx context.Context, filters *{{ .Module.Store -}}. {{- title .Name -}}Filters) (* {{ .Module.Store -}}. {{- title .Name}}, error) {
	qb := {{ .Name -}}FiltersToQuery(filters)
	var {{.Name}} {{ .Module.Store -}}. {{- title .Name}}
	if err := s.getCollection(Collection{{title (plural .Name)}}).FindOne(ctx, qb.Build()).Decode(&{{- .Name -}}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &{{- .Name}}, nil
}

func (s *{{- title .Module.Name -}}Store) Count{{- title (plural .Name) -}} (ctx context.Context, filters *{{ .Module.Store -}}. {{- title .Name -}}Filters) (
	int64,
	error,
) {
	qb := {{ .Name -}}FiltersToQuery(filters)
	
	count, err := s.getCollection(Collection{{title (plural .Name)}}).CountDocuments(ctx, qb.Build())
	if  err != nil {
		return 0, err
	}
	return count, nil
}

func (s *{{- title .Module.Name -}}Store) Get{{- title (plural .Name) -}} (
	ctx context.Context,
	filters *{{ .Module.Store -}}. {{- title .Name -}}Filters,
	after *string,
	before *string,
	first  *int64, 
	last  *int64,
	sortBy {{ .Module.Store -}}.{{title .Name -}}SortBy, 
	sortOrder *string,
) (
	[]*{{ .Module.Store -}}. {{- title .Name}}, 
	bool,
	bool,
	[]string,
	error,
) {
	qb := {{ .Name -}}FiltersToQuery(filters)	
	limit, order, cursorStr := utils.GetLimitAndSortOrderAndCursor(first, last, after, before)
	var c *cursor.Cursor
	if cursorStr != nil {
		c = cursor.FromString(*cursorStr)
		if c != nil {
			if order == 1 {
				qb.Or(
					
					mongoqb.NewQueryBuilder().
						Eq({{ .Module.Store -}}.{{title .Name -}}SortBy(c.SortBy).String(), c.OffsetValue).
						Gt("_id", c.ID),
					mongoqb.NewQueryBuilder().
						Gt({{ .Module.Store -}}.{{title .Name -}}SortBy(c.SortBy).String(), c.OffsetValue),
				)
			} else {
				qb.Or(
					mongoqb.NewQueryBuilder().
						Eq({{ .Module.Store -}}.{{title .Name -}}SortBy(c.SortBy).String(), c.OffsetValue).
						Lt("_id", c.ID), 
					mongoqb.NewQueryBuilder().
						Lt({{ .Module.Store -}}.{{title .Name -}}SortBy(c.SortBy).String(), c.OffsetValue),
				)
			}
		}
	}
	// incrementing limit by 2 to check if next, prev elements are present
	limit += 2
	options := &options.FindOptions{
		Limit: &limit,
		Sort:  utils.GetSortOrder(sortBy.String(), sortOrder, order),
	}

	var hasNextPage, hasPreviousPage bool

	var {{plural .Name}} []*{{ .Module.Store -}}. {{- title .Name}}
	mongoCursor, err := s.getCollection(Collection{{title (plural .Name)}}).Find(ctx, qb.Build(), options)
	if  err != nil {
		return nil, hasNextPage, hasPreviousPage, nil, err
	}
	err = mongoCursor.All(ctx, &{{- plural .Name}})
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, nil, err
	}
	count := len({{ plural .Name}})
	if count == 0 {
		return {{ plural .Name}}, hasNextPage, hasPreviousPage, nil, nil
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

	cursors := make([]string, count)
	for i, {{.Name}} := range {{ plural .Name }} {
		cursors[i] = cursor.NewCursor({{.Name}}.Id, uint8(sortBy), {{.Name}}.Get(sortBy.String()), sortBy.CursorType()).String()
	}

	if order < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		{{ plural .Name}} = utils.ReverseList({{ plural .Name}})
	}
	return {{ plural .Name}}, hasNextPage, hasPreviousPage, cursors, nil
}

func (s *{{- title .Module.Name -}}Store) Update{{- title .Name -}} (ctx context.Context, id string, {{.Name -}}Update *{{-  .Module.Store -}}. {{- title .Name -}}Update) (error) {
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

func (s *{{- title .Module.Name -}}Store) Delete{{- title .Name -}}ByID (ctx context.Context, id string) (error) {
	qb := mongoqb.NewQueryBuilder().
			Eq("_id", id)
	if _,  err := s.getCollection(Collection{{title (plural .Name)}}).DeleteOne(ctx, qb.Build()); err != nil {
		return  err
	}
	return nil
}
`
