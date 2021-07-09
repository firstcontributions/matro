package writer

import (
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

func (w *GoWriter) Format() error {
	if d, err := format.Source(w.data); err != nil {
		return err
	} else {
		w.data = d
	}
	return w.runGoImports()
}

func (w *GoWriter) runGoImports() error {
	if err := w.write("/tmp/matro", "tmp.go"); err != nil {
		return err
	}
	filepath := "/tmp/matro/tmp.go"
	if err := exec.Command("goimports", "-w", filepath).Run(); err != nil {
		return err
	}
	if d, err := ioutil.ReadFile(filepath); err != nil {
		return err
	} else {
		w.data = d
	}
	return os.Remove(filepath)
}
