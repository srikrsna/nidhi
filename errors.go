package nidhi

const (
	// ErrInvalidCursor is returned when the pagination cursor is invalid.
	ErrInvalidCursor Error = "invalid-cursor"
	// ErrNotFound is returned when a document is not found.
	ErrNotFound Error = "not-found"
)

// Error is the underlying error type returned by functions/methods of this package.
//
// Should always be checked using `errors.Is`. They are almost always wrapped by another error.
type Error string

// Error implements the error interface
func (e Error) Error() string {
	return string(e)
}
