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
	sortBy *string, 
	sortOrder *string,
) (
	[]*{{ .Module.Store -}}. {{- title .Name}}, 
	bool,
	bool,
	string,
	string,
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
					qb.And(
						qb.Eq(c.SortBy, c.OffsetValue), 
						qb.Gt("_id", c.ID),
					),  
					qb.Gt(c.SortBy, c.OffsetValue),
				)
			} else {
				qb.Or(
					qb.And(
						qb.Eq(c.SortBy, c.OffsetValue), 
						qb.Lt("_id", c.ID),
					),  
					qb.Lt(c.SortBy, c.OffsetValue),
				)
			}
		}
	}
	// incrementing limit by 2 to check if next, prev elements are present
	limit += 2
	options := &options.FindOptions{
		Limit: &limit,
		Sort:  utils.GetSortOrder(sortBy, sortOrder, order),
	}

	var firstCursor, lastCursor string
	var hasNextPage, hasPreviousPage bool

	var {{plural .Name}} []*{{ .Module.Store -}}. {{- title .Name}}
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
		firstCursor = cursor.NewCursor({{ plural .Name -}}[0].Id, "time_created", {{ plural .Name -}}[0].TimeCreated).String()
		lastCursor = cursor.NewCursor({{ plural .Name -}}[count-1].Id,  "time_created",  {{ plural .Name -}}[count-1].TimeCreated).String()
	}
	if order < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		firstCursor, lastCursor = lastCursor, firstCursor
		{{ plural .Name}} = utils.ReverseList({{ plural .Name}})
	}
	return {{ plural .Name}}, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
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
