package utils

import (
	"fmt"
	"go/format"
	"os"
	"os/exec"
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

	if err := fw.Close(); err != nil {
		return err
	}
	return FixImports(fmt.Sprintf("%s/%s", path, file))

}

func FixImports(file string) error {
	_, err := exec.Command("goimports", "-w", file).Output()
	return err
}
