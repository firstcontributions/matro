package generators

import (
	"github.com/firstcontributions/matro/internal/generators/graphql/gocode"
	"github.com/firstcontributions/matro/internal/generators/graphql/schema"
	"github.com/firstcontributions/matro/internal/generators/grpc/proto"
	"github.com/firstcontributions/matro/internal/generators/models/mongo"
	"github.com/firstcontributions/matro/internal/parser"
)

type Type string
type IGenerator interface {
	Generate() error
}

const (
	TypeGQLSchema  Type = "graphql-schema"
	TypeGQLServer  Type = "graphql-server"
	TypeGRPCProto  Type = "grpc-proto"
	TypeMongoModel Type = "mongo-model"
)

func GetGenerator(t Type, path string, d *parser.Definition) IGenerator {
	switch t {
	case TypeGQLSchema:
		return schema.NewGenerator(path, d)
	case TypeGQLServer:
		return gocode.NewGenerator(path, d)
	case TypeGRPCProto:
		return proto.NewGenerator(path, d)
	case TypeMongoModel:
		return mongo.NewGenerator(path, d)
	}

	return nil
}
