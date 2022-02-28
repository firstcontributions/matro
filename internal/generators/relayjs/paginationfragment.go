package relayjs

const paginationFragment = `
import { graphql, usePaginationFragment } from "react-relay"

const {{ (title (plural .Type)) -}}List = (props) => {
    const {data, loadNext, hasNext} = usePaginationFragment(
        graphql` + "`" + `
        fragment {{- (title (plural .Type)) -}}List_{{- .Parent.Name }} on {{- title .Parent.Name }} 
        @refetchable(queryName: "{{- (title (plural .Type)) -}}List_{{- .Parent.Name -}}Query")
        @argumentDefinitions(
            count: {type: "Int", defaultValue: 10}
            cursor: {type: "String"}
        ){
            {{ .Name }}(first:$count, after: $cursor) 
            @connection(key: "{{- (title (plural .Type)) -}}List__{{- plural .Type -}}") {
                edges {
                    node {
                        id
                        ...{{- title  .Type -}}_node
                    }
                }
            }
        }
        ` + "`" + `, props.{{- .Parent.Name -}}
    )

   
    return (
        <>
           {JSON.stringify(data)}
        </>
    )
}

export default {{ (title (plural .Type)) -}}List

`
