package merrors

import "errors"

var (
	// ErrEmptyType will be raised when type defenition is given as ""
	ErrEmptyType = errors.New("type can't be empty")

	// ErrInvalidTypeDefinition will be raised when the given type definition is not matching with the schema
	ErrInvalidTypeDefinition = errors.New("invalid type definition")

	// ErrNoSchemaDefinedForCustomType raised when there is no schema or properties defined for a custom type
	ErrNoSchemaDefinedForCustomType = errors.New("no schema defined for custom type")

	// ErrNoNameForInlineCustomTypes raised when custom types defined with properties but without a name
	ErrNoNameForInlineCustomTypes = errors.New("no name defined for inline custom type")

	// ErrUndefinedSearchField raised when the given search field is not defined as a property for the type
	ErrUndefinedSearchField = errors.New("undefined search field")

	// ErrUndefinedFilter raised when the given filter field is not defined as a property for the type
	ErrUndefinedFilter = errors.New("undefined filter")

	// ErrUndefinedMutableField raised when the given mutable field is not defined as a property for the type
	ErrUndefinedMutableField = errors.New("undefined mutable field")
)
