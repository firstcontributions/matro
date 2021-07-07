package utils

import (
	"fmt"
	"go/format"
	"os"
	"os/exec"
)

// EnsurePath creates the given path of not existing
func EnsurePath(path string) error {
	return os.MkdirAll(path, 0777)
}

// GetFileWriter create and opens a file in the given path
// and returns a writer interface to the file pointer
func GetFileWriter(path, filename string) (*os.File, error) {
	if err := EnsurePath(path); err != nil {
		return nil, err
	}
	return os.Create(fmt.Sprintf("%s/%s", path, filename))
}

// writeCode writes code to given file
func writeCode(path, file string, code []byte) error {
	fw, err := GetFileWriter(path, file)
	if err != nil {
		return fmt.Errorf("error on opening file :%w", err)
	}
	_, err = fw.Write(code)
	if err != nil {
		return fmt.Errorf("error on writing to file  :%w", err)
	}
	return fw.Close()
}

// FormatAndWriteGoCode will execute gofmt, goimports and write the results
//  to the given file
func FormatAndWriteGoCode(path, file string, code []byte) error {
	code, err := format.Source(code)
	if err != nil {
		return fmt.Errorf("error on formatting code :%w", err)
	}
	err = writeCode(path, file, code)
	if err != nil {
		return fmt.Errorf("error on formatting code :%w", err)
	}
	return FixImports(fmt.Sprintf("%s/%s", path, file))

}

// FixImports runs goimports on generated code
func FixImports(file string) error {
	_, err := exec.Command("goimports", "-w", file).Output()
	return err
}
