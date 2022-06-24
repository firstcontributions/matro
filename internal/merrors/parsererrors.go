package merrors

import "errors"

var (
	// ErrEmptyType will be raised when type defenition is given as ""
	ErrEmptyType = errors.New("type can't be empty")

	// ErrInvalidTypeDefinition will be raised when the given type definition is not matching with the schema
	ErrInvalidTypeDefinition = errors.New("invalid type definition")
)
