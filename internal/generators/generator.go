package generators

import (
	"context"

	"github.com/firstcontributions/matro/internal/generators/gomod"
	"github.com/firstcontributions/matro/internal/generators/graphql/gocode"
	"github.com/firstcontributions/matro/internal/generators/graphql/schema"
	"github.com/firstcontributions/matro/internal/generators/grpc/proto"
	"github.com/firstcontributions/matro/internal/generators/models/mongo"
	"github.com/firstcontributions/matro/internal/parser"
)

// Type defines a generator type
type Type int

// IGenerator is an interface with a generate function
type IGenerator interface {
	Generate(ctx context.Context) error
}

const (
	// TypeGQLSchema defines the type for graphql schema generator
	TypeGQLSchema Type = iota
	// TypeGQLServer defines the type for graphql server code generator
	TypeGQLServer
	// TypeGRPCProto defines the type for gRPC protocol buffer generator
	TypeGRPCProto
	// TypeMongoModel defines the type for mongo model code generator
	TypeMongoModel
	// TypeGoMod inits go.mod and get all dependencies
	TypeGoMod
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
	case TypeGoMod:
		return gomod.NewGenerator(path, d)
	}

	return nil
}
