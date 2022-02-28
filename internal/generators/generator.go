package generators

import (
	"context"
)

// IGenerator is an interface with a generate function
type IGenerator interface {
	Generate(ctx context.Context) error
}
