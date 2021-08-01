package schema

const schemaTmpl = `
schema {
	query Query
}
type Query {
	viewer: User
  
	# Fetches an object given its ID
	node(
	  # The ID of an object
	  id: ID!
	): Node

	{{- range .Queries }}

	{{.GraphQLFormattedName}}: {{.GraphQLFortmattedType}}
	{{- end}}
}

interface Node {
	id: ID!
}

type PageInfo {
	startCursor: String
	endCursor: String
	hasNextPage: Bool!
	hasPreviousPage: Bool!
}

{{- range .Types}}
	{{- template "typeDef" .}}
	{{- if .IsEdge}}
	{{- template "connectionDef" .}}
	{{- template "edgeDef" .}}
	{{- end}}
{{- end}}



{{- define "field"}}
	{{.GraphQLFormattedName}}: {{.GraphQLFortmattedType}}
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


`
