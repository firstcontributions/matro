package generators

import (
	"github.com/firstcontributions/matro/internal/generators/graphql/gocode"
	"github.com/firstcontributions/matro/internal/generators/graphql/schema"
	"github.com/firstcontributions/matro/internal/generators/grpc/proto"
	"github.com/firstcontributions/matro/internal/generators/models/mongo"
	"github.com/firstcontributions/matro/internal/parser"
)

type IGenerator interface {
	Generate() error
}

func GetGenerator(path, t string, d *parser.Definition) IGenerator {
	switch t {
	case "schema":
		return schema.NewGenerator(path, d)
	case "gocode":
		return gocode.NewGenerator(path, d)
	case "proto":
		return proto.NewGenerator(path, d)
	case "mongo":
		return mongo.NewGenerator(path, d)
	}

	return nil
}
