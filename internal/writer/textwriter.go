package writer

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/sirupsen/logrus"
)

// TextWriter writes any text files
type TextWriter struct {
	header   string
	data     []byte
	path     string
	filename string
}

// NewTextWriter returns an instance of TextWriter
func NewTextWriter(path, filename, header string) *TextWriter {
	return &TextWriter{
		path:     path,
		filename: filename,
		header:   header,
	}
}

// Compile should compile the template and get the code as bytes
func (w *TextWriter) Compile(ctx context.Context, tmpl string, data interface{}) error {
	t, err := template.New("t").
		Funcs(FuncMap()).
		Parse(tmpl)
	if err != nil {
		logrus.Errorf("error on parsing remplate %v", err)
		return err
	}

	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		logrus.Errorf("error on exc remplate %v", err)
		return err
	}

	w.data = append([]byte(w.header), b.Bytes()...)
	return nil
}

// Format is supposed to format the code, as text formatter can be
// any type of code we cant implement any formatter. This leave as
// an empty function to implement Writer interface
func (*TextWriter) Format(ctx context.Context) error {
	return nil
}

// write will write contents to given file
func (w *TextWriter) write(path, filename string) error {
	if err := utils.EnsurePath(path); err != nil {
		return err
	}
	logrus.Infof("generating %s/%s", path, filename)
	return ioutil.WriteFile(fmt.Sprintf("%s/%s", path, filename), w.data, 0644)
}

// Write will write contents to given file
func (w *TextWriter) Write(ctx context.Context) error {
	return w.write(w.path, w.filename)
}
