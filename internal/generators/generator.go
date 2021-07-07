package generators

import (
	"github.com/firstcontributions/matro/internal/generators/graphql/gocode"
	"github.com/firstcontributions/matro/internal/generators/graphql/schema"
	"github.com/firstcontributions/matro/internal/generators/grpc/proto"
	"github.com/firstcontributions/matro/internal/generators/models/mongo"
	"github.com/firstcontributions/matro/internal/parser"
)

// Type defines a generator type
type Type string

// IGenerator is an interface with a generate function
type IGenerator interface {
	Generate() error
}

const (
	// TypeGQLSchema defines the type for graphql schema generator
	TypeGQLSchema Type = "graphql-schema"
	// TypeGQLServer defines the type for graphql server code generator
	TypeGQLServer Type = "graphql-server"
	// TypeGRPCProto defines the type for gRPC protocol buffer generator
	TypeGRPCProto Type = "grpc-proto"
	// TypeMongoModel defines the type for mongo model code generator
	TypeMongoModel Type = "mongo-model"
)

// GetGenerator is a factory method to get generator by given type
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
