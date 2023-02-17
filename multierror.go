package multierror

import "errors"

type MultiError struct {
	errors []error
}

func (m MultiError) Error() string {
	return ""
}

func (m MultiError) Is(target error) bool {
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

func Append(err error, errors ...error) error {
	if err == nil && len(errors) == 0 {
		return nil
	}

	switch e := err.(type) {
	case MultiError:
		return MultiError{
			errors: append(e.errors, errors...),
		}
	default:
		var errs []error
		if err != nil {
			errs = append(errs, err)
		}
		if len(errors) != 0 {
			errs = append(errs, errors...)
		}
		return MultiError{
			errors: errs,
		}
	}
	return nil
}
