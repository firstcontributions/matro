package writer

import (
	"context"
	"strings"

	"github.com/firstcontributions/matro/internal/ctxkeys"
	"github.com/firstcontributions/matro/pkg/spinner"
)

// IWriter implements a code writer
type IWriter interface {
	Compile(context.Context, string, interface{}) error
	Format(context.Context) error
	Write(context.Context) error
}

// Type will be used to define enums for code types
type Type int

const (
	// TypeGoCode writes go source code
	TypeGoCode Type = iota
	// TypeText writes any text file, this will be used as a common code
	// writers for all types which does not have any formatter available
	TypeText
)

// GetWriter is a factory method for writers
func GetWriter(path, filename string) IWriter {
	if getFileExtension(filename) == "go" {
		return NewGoWriter(path, filename)
	}
	return NewTextWriter(path, filename)
}

func getFileExtension(file string) string {
	tmp := strings.Split(file, ".")
	if len(tmp) < 2 {
		return ""
	}
	return tmp[1]
}

// CompileAndWrite will compile the given template, autoformat
// and write to file
func CompileAndWrite(
	ctx context.Context,
	path string,
	filename string,
	tmpl string,
	data interface{},
) error {
	w := GetWriter(path, filename)
	// getSpinner(ctx).Update(fmt.Sprintf("%s/%s", path, filename))
	if err := w.Compile(ctx, tmpl, data); err != nil {
		return err
	}
	if err := w.Write(ctx); err != nil {
		return err
	}
	return w.Format(ctx)
}
func getSpinner(ctx context.Context) *spinner.Spinner {
	return ctx.Value(ctxkeys.Spinner).(*spinner.Spinner)
}
