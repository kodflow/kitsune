package errors

import "errors"

type MultiError struct {
	Errors []error
}

func (m *MultiError) Add(err error) {
	m.Errors = append(m.Errors, err)
}

func (m *MultiError) Del(index int) {
	if index < 0 || index >= len(m.Errors) {
		return
	}
	m.Errors = append(m.Errors[:index], m.Errors[index+1:]...)
}

func (m *MultiError) Count() int {
	return len(m.Errors)
}

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

func (m *MultiError) Error() string {
	var combinedError string
	for _, err := range m.Errors {
		combinedError += err.Error() + "; "
	}
	return combinedError
}

func NewMultiError(errs ...error) *MultiError {
	return &MultiError{Errors: errs}
}

func New(text string) error {
	return errors.New(text)
}
