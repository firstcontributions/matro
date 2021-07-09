package writer

import (
	"strings"
)

// ICodeWriter implements a code writer
type IWriter interface {
	Compile(string, interface{}) error
	Format() error
	Write() error
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

func GetWriter(path, filename string) IWriter {
	switch getFileExtension(filename) {
	case "go":
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

func CompileAndWrite(
	path string,
	filename string,
	tmpl string,
	data interface{},
) error {
	w := GetWriter(path, filename)
	if err := w.Compile(tmpl, data); err != nil {
		return err
	}
	if err := w.Format(); err != nil {
		return err
	}
	if err := w.Write(); err != nil {
		return err
	}
	return nil
}
