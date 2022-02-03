package proto

const tmpl = `
syntax = "proto3";
package {{ .Name }};

import "google/protobuf/timestamp.proto";
option go_package = "internal/grpc/{{- plural .Name -}}/proto";

message StatusResponse {
	bool status= 1;
}

message RefByIDRequest {
	string id= 1;
}

{{- range .Types}}
{{- template "typeDef" .}}
{{- if .IsNode }}
{{- template "bulk_get_request" .}}
{{- template "bulk_get_response" .}}
{{- end}}
{{- end }}

service {{title .Name -}}Service {
	{{- range .Types}}
	{{- if .IsNode }}
	{{- template "rpcs" .}}
	{{- end }}
	{{- end }}
}


{{- define "typeDef" }}

message {{title .Name}} {
	{{- counter 0}}
	{{- range .Fields}}
	{{- template "fields" .}}
	{{- end}}
	{{- range .ReferedFields}}
	string {{. -}}_id= {{counter}};
	{{- end}}
}
{{- end}}

{{- define "rpcs" }}

	// {{ plural .Name}} crud operations
	rpc Create{{- title .Name}} ({{ title .Name}}) returns ({{ title .Name}}) {}
	rpc Get{{- title .Name}}ByID (RefByIDRequest) returns ({{ title .Name}}) {}
	rpc Get{{- title (plural .Name) }} (Get{{- title (plural .Name)}}Request) returns (Get{{- title (plural .Name)}}Response){};
	rpc Update{{- title .Name}} ({{- title .Name}}) returns ({{ title .Name}}) {}
	rpc Delete{{- title .Name}} (RefByIDRequest) returns  (StatusResponse){}
{{- end}}

{{- define "fields" }}
{{- if not .IsJoinedData}}
	{{- if .IsList}}
	repeated {{ grpcType .Type}} {{.Name}}= {{counter}};
	{{- else}} 
	{{ grpcType .Type}} {{.Name}}= {{counter}};
	{{-  end}}
{{- else}}
	{{- if not .IsList }}
	string {{.Name}}= {{counter}};
	{{- end }}
{{- end}}
{{- end}}

{{- define "bulk_get_response" }} 
message Get{{- title (plural .Name)}}Response {
	bool has_next= 1;
	bool has_previous= 2;
	string first_cursor= 3;
	string last_cursor= 4;
	repeated {{ title .Name}} data=5;
}
{{- end}}

{{- define "bulk_get_request" }} 
message Get{{- title (plural .Name)}}Request {
	string before= 1;
	string after= 2;
	int64 first= 3;
	int64 last= 4;
	repeated string ids= 5;
	{{- if not (empty .SearchFields)}}
	string search= 6;
	{{- counter 6}}
	{{- else }}
	{{- counter 5}}
	{{- end}}
	{{- $t := .}}
	{{- range .ReferedFields }}
	string {{. -}}_id= {{counter}};
	{{- end }}
	{{- range .Filters }}
	{{ $t.FieldType .}} {{.}}= {{counter}};
	{{- end}}
}
{{- end}}
`
