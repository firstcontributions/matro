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
{{- template "update_request" .}}
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
	{{- range .ReferedTypes}}
	string {{.Name -}}_id= {{counter}};
	{{- end}}
}
{{- end}}

{{- define "rpcs" }}

	// {{ plural .Name}} crud operations
	rpc Create{{- title .Name}} ({{ title .Name}}) returns ({{ title .Name}}) {}
	rpc Get{{- title .Name}}ByID (RefByIDRequest) returns ({{ title .Name}}) {}
	rpc Get{{- title (plural .Name) }} (Get{{- title (plural .Name)}}Request) returns (Get{{- title (plural .Name)}}Response){};
	rpc Update{{- title .Name}} (Update{{- title .Name -}}Request) returns (StatusResponse) {}
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

{{- define "update_request" }} 
message Update{{- title .Name -}}Request {
	{{- counter 0}}
	string id= {{counter}};
	{{- range .Fields}}
	{{- if (and .IsMutatable (not .IsQuery))}}
	{{- if not .IsJoinedData}}
	{{- if .IsList}}
	repeated {{ grpcType .Type}} {{.Name}}= {{counter}};
	{{- else}} 
	optional {{ grpcType .Type}} {{.Name}}= {{counter}};
	{{-  end}}
	{{- else}}
	{{- if not .IsList }}
	optional string {{.Name}}= {{counter}};
	{{- end }}
	{{- end}}
	{{- end}}
	{{- end}}
}

{{- end}}

{{- define "bulk_get_request" }} 
message Get{{- title (plural .Name)}}Request {
	optional string before= 1;
	optional string after= 2;
	optional int64 first= 3;
	optional int64 last= 4;
	repeated string ids= 5;
	{{- if not (empty .SearchFields)}}
	optional string search= 6;
	{{- counter 6}}
	{{- else }}
	{{- counter 5}}
	{{- end}}
	{{- $t := .}}
	{{- range .ReferedTypes }}
	optional string {{.Name -}}_id= {{counter}};
	{{- end }}
	{{- range .Filters }}
	optional {{ $t.FieldType .}} {{.}}= {{counter}};
	{{- end}}
}
{{- end}}
`
