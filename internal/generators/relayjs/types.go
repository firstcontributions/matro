package relayjs

const typesTmpl = `
import { graphql, useFragment } from "react-relay"

const {{title .Name}} = (props) => {
    const data = useFragment(
        graphql` + "`" + `
            fragment {{title .Name}}_node on {{title .Name}} {
			{{- $parent := .}}
			{{- range .Fields}}
			{{- if .IsQuery}}
				...{{- (title (plural .Type)) -}}List_{{- $parent.Name }}
			{{- else}}
				{{- if .IsPrimitive}}
				{{camel .Name}}
				{{- else}}
				{{camel .Name}} {
					id
					...{{title .Type}}_node
				}
				{{- end}}
			{{- end}}
			{{- end}}
            }
        ` + "`" + `, props.{{- .Name }}
    )

    return (
        <>
            {JSON.stringify(data)}
        </>
    )
}

export default {{title .Name}}
`
