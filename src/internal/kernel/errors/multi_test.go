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

// TestNewMultiError tests the NewMultiError function.
func TestNewMultiError(t *testing.T) {
	// Test with no errors
	m := NewMultiError()
	assert.NotNil(t, m, "NewMultiError should return a non-nil MultiError")
	assert.Empty(t, m.Errors, "NewMultiError should return a MultiError with no errors")

	// Test with one error
	err := errors.New("error 1")
	m = NewMultiError(err)
	assert.NotNil(t, m, "NewMultiError should return a non-nil MultiError")
	assert.Equal(t, []error{err}, m.Errors, "NewMultiError should return a MultiError with one error")

	// Test with multiple errors
	errs := []error{
		errors.New("error 1"),
		errors.New("error 2"),
		errors.New("error 3"),
	}
	m = NewMultiError(errs...)
	assert.NotNil(t, m, "NewMultiError should return a non-nil MultiError")
	assert.Equal(t, errs, m.Errors, "NewMultiError should return a MultiError with multiple errors")
}

// TestNew tests the New function.
func TestNew(t *testing.T) {
	err := New("test error")
	assert.Equal(t, "test error", err.Error(), "New should return an error with the specified text")
}
func TestMultiErrorIsError(t *testing.T) {
	// Test with an empty MultiError
	m := &MultiError{}
	assert.Nil(t, m.IsError(), "Empty MultiError should return nil")

	// Test with a MultiError containing one error
	m = &MultiError{
		Errors: []error{errors.New("error 1")},
	}
	assert.Equal(t, m, m.IsError(), "MultiError with one error should return the MultiError itself")

	// Test with a MultiError containing multiple errors
	m = &MultiError{
		Errors: []error{
			errors.New("error 1"),
			errors.New("error 2"),
			errors.New("error 3"),
		},
	}
	assert.Equal(t, m, m.IsError(), "MultiError with multiple errors should return the MultiError itself")
}

// TestMultiError_Count tests the Count method of the MultiError type.
func TestMultiErrorCount(t *testing.T) {
	// Test with an empty MultiError
	m := &MultiError{}
	assert.Equal(t, 0, m.Count(), "Empty MultiError should have a count of 0")

	// Test with a MultiError containing one error
	m = &MultiError{
		Errors: []error{errors.New("error 1")},
	}
	assert.Equal(t, 1, m.Count(), "MultiError with one error should have a count of 1")

	// Test with a MultiError containing multiple errors
	m = &MultiError{
		Errors: []error{
			errors.New("error 1"),
			errors.New("error 2"),
			errors.New("error 3"),
		},
	}
	assert.Equal(t, 3, m.Count(), "MultiError with multiple errors should have a count of 3")
}

// TestMultiError_Del tests the Del method of the MultiError type.
func TestMultiErrorDel(t *testing.T) {
	// Test with an empty MultiError
	m := &MultiError{}
	m.Del(0)
	assert.Empty(t, m.Errors, "Del should not modify an empty MultiError")

	// Test with a MultiError containing one error
	m = &MultiError{
		Errors: []error{errors.New("error 1")},
	}
	m.Del(0)
	assert.Empty(t, m.Errors, "Del should remove the only error from a MultiError with one error")

	// Test with a MultiError containing multiple errors
	m = &MultiError{
		Errors: []error{
			errors.New("error 1"),
			errors.New("error 2"),
			errors.New("error 3"),
		},
	}
	m.Del(1)
	assert.Equal(t, []error{errors.New("error 1"), errors.New("error 3")}, m.Errors, "Del should remove the error at the specified index from a MultiError with multiple errors")

	// Test with an invalid index
	m.Del(10)
	assert.Equal(t, []error{errors.New("error 1"), errors.New("error 3")}, m.Errors, "Del should not modify a MultiError when the index is out of range")
}

// TestMultiError_Add tests the Add method of the MultiError type.
func TestMultiErrorAdd(t *testing.T) {
	// Test with a nil error
	m := &MultiError{}
	m.Add(nil)
	assert.Empty(t, m.Errors, "Add should not add a nil error to the MultiError")

	// Test with a non-nil error
	m = &MultiError{}
	err := errors.New("test error")
	m.Add(err)
	assert.Equal(t, []error{err}, m.Errors, "Add should add a non-nil error to the MultiError")
}
