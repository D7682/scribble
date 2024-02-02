// Package errors provides custom error types for scribble package
package errors

import (
	"fmt"
)

// FileIOError represents an error related to file I/O operations
type FileIOError struct {
	path string
	err  error
}

// Error implements the error interface for FileIOError
func (e *FileIOError) Error() string {
	return fmt.Sprintf("file I/O error at path %v: %v", e.path, e.err)
}

// Path returns the path associated with the error
func (e *FileIOError) Path() string {
	return e.path
}

// OriginalError returns the original underlying error
func (e *FileIOError) OriginalError() error {
	return e.err
}

// NewFileIOError creates a new instance of FileIOError
func NewFileIOError(path string, err error) ScribblerError {
	return &FileIOError{
		path: path,
		err:  err,
	}
}
