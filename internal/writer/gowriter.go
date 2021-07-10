package writer

import (
	"context"
	"go/format"
	"io/ioutil"
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
	if err := w.write("/tmp/matro", "tmp.go"); err != nil {
		return err
	}
	filepath := "/tmp/matro/tmp.go"
	if err := exec.Command("goimports", "-w", filepath).Run(); err != nil {
		return err
	}
	d, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	w.data = d
	return os.Remove(filepath)
}
