package errors

import "errors"

// MultiError represents a collection of errors.
// It aggregates multiple error values into a single error entity.
// This is useful for handling multiple errors as a single error.
type MultiError struct {
	Errors []error // Errors is a slice storing individual error instances.
}

// Add adds an error to the MultiError.
// If the provided error is not nil, it is appended to the MultiError's Errors slice.
//
// Parameters:
// - err: error The error to be added to the MultiError.
func (m *MultiError) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

// Del deletes an error at the specified index from the MultiError.
// It safely removes an error by index without causing a panic, even if the index is out of range.
//
// Parameters:
// - index: int The index of the error to be deleted in the MultiError's Errors slice.
func (m *MultiError) Del(index int) {
	if index < 0 || index >= len(m.Errors) {
		return
	}
	m.Errors = append(m.Errors[:index], m.Errors[index+1:]...)
}

// Count returns the number of errors in the MultiError.
// It provides a quick way to check the number of errors aggregated in the MultiError.
//
// Returns:
// - int: The count of non-nil errors in the MultiError's Errors slice.
func (m *MultiError) Count() int {
	return len(m.Errors)
}

// IsError checks if the MultiError contains any errors.
// This method allows for easy checking of the presence of non-nil errors in MultiError.
// It returns the MultiError itself if it contains any non-nil errors, otherwise nil.
//
// Returns:
// - error: The MultiError itself if it contains non-nil errors, otherwise nil.
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
// It concatenates the string representations of all contained errors, separated by semicolons.
// This method implements the error interface for MultiError.
//
// Returns:
// - string: A combined string representation of all errors in the MultiError.
func (m *MultiError) Error() string {
	var combinedError string
	for _, err := range m.Errors {
		combinedError += err.Error() + "; "
	}
	return combinedError
}

// NewMultiError creates a new MultiError with the provided errors.
// It initializes a MultiError with a slice of errors passed as arguments.
//
// Parameters:
// - errs: ...error A variadic list of errors to be included in the new MultiError.
//
// Returns:
// - *MultiError: A pointer to the newly created MultiError initialized with the given errors.
func NewMultiError(errs ...error) *MultiError {
	return &MultiError{Errors: errs}
}

// New creates a new error with the provided text.
// It is a convenience function for creating standard error objects.
//
// Parameters:
// - text: string The text for the error to be created.
//
// Returns:
// - error: A new error created from the provided text.
func New(text string) error {
	return errors.New(text)
}
