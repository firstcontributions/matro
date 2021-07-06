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
{{- template "bulk_get_request" .}}
{{- template "bulk_get_response" .}}
{{- end }}

service {{title .Name -}}Service {
	{{- range .Types}}
	{{- template "rpcs" .}}
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
	bool first_cursor= 3;
	bool last_cursor= 4;
	repeated {{ title .Name}} data=5;
}
{{- end}}

{{- define "bulk_get_request" }} 
message Get{{- title (plural .Name)}}Request {
	string before= 1;
	string after= 2;
	int32 limit= 3;
	string search= 4;
	repeated string ids= 5;
	{{- $t := .}}
	{{- counter 6}}
	{{- range .Filters }}
	{{ $t.FieldType .}} {{.}}= {{counter}};
	{{- end}}
}
{{- end}}
`
