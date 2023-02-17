package multierror

import (
	"errors"
	"fmt"
)

type multiError struct {
	errors []error
}

func (m *multiError) String() string {
	if m == nil {
		return ""
	}

	return m.Error()
}

func (m *multiError) Error() string {
	if m == nil {
		return ""
	}

	if len(m.errors) == 0 {
		return ""
	}

	var s string

	var errs []error
	for _, e := range m.errors {
		errs = append(errs, getRecursiveErrors(e)...)
	}

	if len(errs) == 1 {
		s += "Found one error:\n"
	} else {
		s += fmt.Sprintf("Found %d errors:\n", len(errs))
	}

	for _, err := range errs {
		s += fmt.Sprintf("\t%s\n", err.Error())
	}

	return s
}

func getRecursiveErrors(err error) []error {
	if err == nil {
		return nil
	}

	var flattened []error

	switch e := err.(type) {
	case *multiError:
		if len(e.errors) == 0 {
			return nil
		}

		for _, e := range e.errors {
			flattened = append(flattened, getRecursiveErrors(e)...)
		}
	case error:
		return []error{e}
	}

	return flattened
}

func (m *multiError) Is(target error) bool {
	if m == nil {
		return false
	}

	if target == nil {
		return false
	}

	if m.errors == nil {
		return false
	}

	for _, err := range m.errors {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func (m *multiError) As(target any) bool {
	if m == nil {
		return false
	}

	if target == nil {
		return false
	}

	if m.errors == nil {
		return false
	}

	for _, err := range m.errors {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}

// Append joins one error with one or more other errors.
// The resulting value is a multiError containing all non-nil errors provided.
// If the first error is a multiError, the rest of the errors will be appended to the existing multiError.
// If the first error is not a multiError, a new multiError will be created for it and all the rest of the errors as well.
func Append(err error, errors ...error) *multiError {
	if err == nil && len(errors) == 0 {
		return nil
	}

	switch e := err.(type) {
	case *multiError:
		return &multiError{
			errors: append(e.errors, errors...),
		}
	default:
		var errs []error
		if err != nil {
			errs = append(errs, err)
		}
		for _, e := range errors {
			if e != nil {
				errs = append(errs, e)
			}
		}
		if len(errs) == 0 {
			return nil
		}
		return &multiError{
			errors: errs,
		}
	}
	return nil
}
