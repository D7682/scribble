// pkg/errors/errors.go
package errors

import "errors"

var (
	// ErrMissingCollection is the error for missing collection
	ErrMissingCollection = errors.New("missing collection - no place to save record")

	// ErrResourceNotFound is the error for missing resource
	ErrResourceNotFound = errors.New("missing resource - unable to save record")
)

// ScribblerError is a custom error interface for the scribbler package with enhanced error handling methods.
type ScribblerError interface {
	error
	Path() string         // Path returns the path associated with the error
	OriginalError() error // OriginalError returns the original underlying error
}
