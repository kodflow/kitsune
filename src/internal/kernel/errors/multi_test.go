package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiError_Error(t *testing.T) {
	// Test case 1: Empty errors slice
	m := &MultiError{}
	assert.Equal(t, "", m.Error())

	// Test case 2: Single error
	m = &MultiError{
		Errors: []error{errors.New("error 1")},
	}
	assert.Equal(t, "error 1; ", m.Error())

	// Test case 3: Multiple errors
	m = &MultiError{
		Errors: []error{
			errors.New("error 1"),
			errors.New("error 2"),
			errors.New("error 3"),
		},
	}
	assert.Equal(t, "error 1; error 2; error 3; ", m.Error())
}
