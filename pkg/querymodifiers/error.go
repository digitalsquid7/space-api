package querymodifiers

import "fmt"

type InvalidInputError struct {
	field   string
	message string
}

func NewInvalidInputError(field, message string) *InvalidInputError {
	return &InvalidInputError{field: field, message: message}
}

func (e *InvalidInputError) Error() string {
	return fmt.Sprintf(`"%s": %s`, e.field, e.message)
}
