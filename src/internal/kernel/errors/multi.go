package errors

import "errors"

// MultiError represents a collection of errors.
type MultiError struct {
	Errors []error
}

// Add adds an error to the MultiError.
func (m *MultiError) Add(err error) {
	m.Errors = append(m.Errors, err)
}

// Del deletes an error at the specified index from the MultiError.
func (m *MultiError) Del(index int) {
	if index < 0 || index >= len(m.Errors) {
		return
	}
	m.Errors = append(m.Errors[:index], m.Errors[index+1:]...)
}

// Count returns the number of errors in the MultiError.
func (m *MultiError) Count() int {
	return len(m.Errors)
}

// IsError checks if the MultiError contains any errors.
// It returns the MultiError itself if there are errors, otherwise it returns nil.
func (m *MultiError) IsError() error {
	if m.Count() > 0 {
		for _, err := range m.Errors {
			if err != nil {
				return m
			}
		}
	}

	return nil
}

// Error returns a string representation of the MultiError.
func (m *MultiError) Error() string {
	var combinedError string
	for _, err := range m.Errors {
		combinedError += err.Error() + "; "
	}
	return combinedError
}

// NewMultiError creates a new MultiError with the provided errors.
func NewMultiError(errs ...error) *MultiError {
	return &MultiError{Errors: errs}
}

// New creates a new error with the provided text.
func New(text string) error {
	return errors.New(text)
}
