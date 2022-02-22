package writer

import (
	"bytes"
	"context"
	"log"
	"text/template"

	"github.com/firstcontributions/matro/internal/generators/utils"
)

const (
	autoGenText = "// Code generated by github.com/firstcontributions/matro. DO NOT EDIT.\n"
)

// TextWriter writes any text files
type TextWriter struct {
	data     []byte
	path     string
	filename string
}

// NewTextWriter returns an instance of TextWriter
func NewTextWriter(path, filename string) *TextWriter {
	return &TextWriter{
		path:     path,
		filename: filename,
	}
}

// Compile should compile the template and get the code as bytes
func (w *TextWriter) Compile(ctx context.Context, tmpl string, data interface{}) error {
	t, err := template.New("t").
		Funcs(FuncMap()).
		Parse(tmpl)
	if err != nil {
		log.Fatalf("error on parsing remplate %v", err)
		return err
	}

	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		log.Fatalf("error on exc remplate %v", err)

		return err
	}

	w.data = b.Bytes()
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
	fw, err := utils.GetFileWriter(path, filename)
	if err != nil {
		return err
	}
	if _, err := fw.Write(w.data); err != nil {
		return err
	}
	return fw.Close()
}

// Write will write contents to given file
func (w *TextWriter) Write(ctx context.Context) error {
	return w.write(w.path, w.filename)
}
