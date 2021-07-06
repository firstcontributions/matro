package utils

import (
	"fmt"
	"go/format"
	"os"
)

func EnsurePath(path string) error {
	return os.MkdirAll(path, 0777)
}

func GetFileWriter(path, filename string) (*os.File, error) {
	if err := EnsurePath(path); err != nil {
		return nil, err
	}
	return os.Create(fmt.Sprintf("%s/%s", path, filename))
}

func WriteCodeToGoFile(path, file string, code []byte) error {
	fw, err := GetFileWriter(path, file)
	if err != nil {
		return fmt.Errorf("error on opening file :%w", err)
	}
	code, err = format.Source(code)
	if err != nil {
		return fmt.Errorf("error on formatting code :%w", err)
	}
	_, err = fw.Write(code)
	if err != nil {
		return fmt.Errorf("error on writing to file  :%w", err)
	}
	return fw.Close()
}

// func FixImports(file string) error {
// 	if _, err := exec.Command("goimports ")
// }
