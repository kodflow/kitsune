package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMultiError_Error tests the Error method of the MultiError type.
// This test verifies that the Error method correctly concatenates the error messages.
// It tests three scenarios:
// 1. An empty MultiError, expecting an empty string.
// 2. A MultiError with one error, expecting the error message followed by a semicolon.
// 3. A MultiError with multiple errors, expecting a concatenation of error messages separated by semicolons.
func TestMultiErrors(t *testing.T) {
	// Test with an empty MultiError
	m := &MultiError{}
	assert.Equal(t, "", m.Error(), "Empty MultiError should return an empty string")

	// Test with a MultiError containing one error
	m = &MultiError{
		Errors: []error{errors.New("error 1")},
	}
	assert.Equal(t, "error 1; ", m.Error(), "MultiError with one error should return the error message followed by a semicolon")

	// Test with a MultiError containing multiple errors
	m = &MultiError{
		Errors: []error{
			errors.New("error 1"),
			errors.New("error 2"),
			errors.New("error 3"),
		},
	}
	assert.Equal(t, "error 1; error 2; error 3; ", m.Error(), "MultiError with multiple errors should return concatenated error messages separated by semicolons")
}
