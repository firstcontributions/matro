package schema

const schemaTmpl = `
schema {
	query: Query
	mutation: Mutation
}
scalar Time
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
	{{- if (and .IsNode (or .GraphqlOps.Create .GraphqlOps.Update))}}
	{{- template "inputType" .}}
	{{- end}}
	{{- if .IsEdge}}
	{{- template "connectionDef" .}}
	{{- template "edgeDef" .}}
	{{- end}}
	{{- end}}
{{- end}}

{{- range .QueryTypes}}
	{{- template "typeDef" .}}
	{{- if .IsEdge}}
	{{- template "connectionDef" .}}
	{{- template "edgeDef" .}}
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

input Create{{- title .Name -}}Input {
	{{- range .Fields}}
	{{- if (not (or (isAditField .Name) .IsQuery .NoGraphql))}}
	{{.GraphQLFormattedName}}: {{.GraphQLFortmattedType}}
	{{- end}}
	{{- end}}
}

input Update{{- title .Name -}}Input {
	{{- range .Fields}}
	{{- if (and .IsMutatable (not (or (isAditField .Name) .IsQuery .NoGraphql)))}}
	{{.GraphQLFormattedName}}: {{.GraphQLFortmattedType}}
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


{{- define "connectionDef" }} 

type {{.ConnectionName}} {
	edges: [{{ .EdgeName}}]!
	pageInfo: PageInfo!
}
{{- end}}


{{- define "mutation" }}
{{- if (.GraphqlOps.Create) }} 
	create{{- title .Name}}({{.Name}}: Create{{- title .Name -}}Input!): {{title .Name}}!
{{- end}}
{{- end}}
`
