// pkg/errors/not_found_error.go
package errors

import "fmt"

// NotFoundError is a custom error type for file or directory not found errors
type NotFoundError struct {
	path string
	err  error
}

// NewNotFoundError creates a new instance of NotFoundError
func NewNotFoundError(path string, err error) ScribblerError {
	return &NotFoundError{path: path, err: err}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("File or directory not found: %s", e.path)
}

// Path returns the path associated with the error
func (e *NotFoundError) Path() string {
	return e.path
}

// OriginalError returns the original underlying error
func (e *NotFoundError) OriginalError() error {
	return e.err
}
