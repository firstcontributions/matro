package writer

import (
	"context"
	"fmt"
	"go/format"
	"os"
	"os/exec"
)

// GoWriter formats and writes go code
// It overrides Compile and Write functions of TextWriter
type GoWriter struct {
	*TextWriter
}

// NewGoWriter is the constructor for GoWriter
func NewGoWriter(path, file string) *GoWriter {
	return &GoWriter{
		TextWriter: NewTextWriter(path, file),
	}
}

// Format formats gocode using gofmt and goimports
func (w *GoWriter) Format(ctx context.Context) error {
	d, err := format.Source(w.data)
	if err != nil {
		return err
	}
	w.data = d
	return w.runGoImports(ctx)
}

func (w *GoWriter) runGoImports(ctx context.Context) error {
	return goImport(fmt.Sprintf("%s/%s", w.path, w.filename))
}

func goImport(file string) error {
	cmd := exec.Command("goimports", "-w", file)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
