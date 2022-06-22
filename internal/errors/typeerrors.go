package errors

import "errors"

var (
	// ErrNoModules will the raised when there is no modules found in the config file
	ErrNoModules = errors.New("no modules defined")
)
