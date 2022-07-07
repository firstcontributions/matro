package schema

const schemaTmpl = `
schema {
	query: Query
	mutation: Mutation
}
scalar Time

enum SortOrder {
	asc,
	desc,
}

type Query {
	viewer: User
  
	# Fetches an object given its ID
	node(
	  # The ID of an object
	  id: ID!
	): Node

	{{- range .Queries }}
	{{- if (not .Parent)}}
	{{.GraphQLFormattedName}}: {{.GraphQLFortmattedType}}
	{{- end}}
	{{- end}}
}

interface Node {
	id: ID!
}

type PageInfo {
	startCursor: String
	endCursor: String
	hasNextPage: Boolean!
	hasPreviousPage: Boolean!
}

{{- range .Types}}
	{{- if (not .NoGraphql)}}
	{{- template "typeDef" .}}
	{{- if .GraphqlOps.Create}}
	{{- template "inputType" .}}
	{{- end}}
	{{- if (and .IsNode .GraphqlOps.Update)}}
	{{- template "inputTypeUpdate" .}}
	{{- end}}
	{{- if .IsEdge}}
	{{- template "connectionDef" .}}
	{{- template "edgeDef" .}}
	{{- template "sortby" .}}
	{{- end}}
	{{- end}}
{{- end}}

{{- range .QueryTypes}}
	{{- template "typeDef" .}}
	{{- if .IsEdge}}
	{{- template "connectionDef" .}}
	{{- template "edgeDef" .}}
	{{- template "sortby" .}}
	{{- end}}
{{- end}}

type Mutation {
	{{- range .Types }} 
		{{- template "mutation" .}}
	{{- end}}
}


{{- define "field"}}
	{{- if (not .NoGraphql)}}
	{{.GraphQLFormattedName}}: {{.GraphQLFortmattedType}}
	{{- end}}
{{- end}}

{{- define "inputType"}}

input {{ title .Name -}}Input {
	{{- range .Fields}}
	{{- if (not (or (isAditField .Name) .IsQuery .NoGraphql .ViewerRefence))}}
	{{- if .IsPrimitive}}
	{{.GraphQLFormattedName}}: {{.GraphQLFortmattedType}}
	{{- else}} 
	{{.GraphQLFormattedName}}: {{.GraphQLFortmattedInputType}}
	{{- end}}
	{{- end}}
	{{- end}}

	{{- range .ReferedTypes}}
	{{- if (not .IsViewerType)}}
	{{.Name -}}ID: ID!
	{{- end}}
	{{- end}}
}
{{- end}}

{{- define "inputTypeUpdate"}}
input Update{{- title .Name -}}Input {
	id: ID!
	{{- range .Fields}}
	{{- if (and .IsMutatable (not (or (isAditField .Name) .IsQuery .NoGraphql .ViewerRefence)))}}
	{{- if .IsPrimitive}}
	{{.GraphQLFormattedName}}: {{.GraphQLFortmattedType true}}
	{{- else}} 
	{{.GraphQLFormattedName}}: {{.GraphQLFortmattedInputType true}}
	{{- end}}
	{{- end}}
	{{- end}}
}
{{- end}}

{{- define "typeDef"}}
{{- if .IsNode}}

type {{title .Name}} implements Node {
{{- else}}

type {{title .Name}} {
{{- end}}

	{{- range .Fields}}
	{{- template "field" .}}
	{{- end}}
}
{{- end}}


{{- define "edgeDef" }} 

type {{.EdgeName}} {
	node: {{title .Name}}!
	cursor: String!
}
{{- end}}

{{- define "sortby" }} 

enum {{ title .Name -}}SortBy {
	{{- range .SortBy }}
	{{ . }},
	{{- end}}
}
{{- end}}


{{- define "connectionDef" }} 

type {{.ConnectionName}} {
	edges: [{{ .EdgeName}}]!
	pageInfo: PageInfo!
	totalCount: Int!
	{{- if (ne .ViewerRefenceField "") }}
	hasViewerAssociation: Boolean!
	{{- end}}
}
{{- end}}


{{- define "mutation" }}
{{- if (and .IsNode .GraphqlOps.Create) }} 
	create{{- title .Name}}({{.Name}}: {{- title .Name -}}Input!): {{title .Name}}!
{{- end}}
{{- if (and .IsNode .GraphqlOps.Update) }} 
	update{{- title .Name}}({{.Name}}: Update{{- title .Name -}}Input!): {{title .Name}}!
{{- end}}
{{- end}}
`
