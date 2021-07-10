package utils

import (
	"fmt"
	"os"
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
