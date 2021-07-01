package utils

import (
	"fmt"
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
