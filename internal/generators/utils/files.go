package utils

import (
	"os"
)

// EnsurePath creates the given path of not existing
func EnsurePath(path string) error {
	return os.MkdirAll(path, 0777)
}
