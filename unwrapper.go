package multierror

import "errors"

type unwrapper []error

func (uw unwrapper) Error() string {
	if len(uw) == 0 {
		return ""
	}

	return uw[0].Error()
}

func (uw unwrapper) Unwrap() error {
	if len(uw) == 1 {
		return nil
	}

	return uw[1:]
}

func (uw unwrapper) As(target any) bool {
	if len(uw) == 0 {
		return false
	}

	return errors.As(uw[0], target)
}

func (uw unwrapper) Is(target error) bool {
	if len(uw) == 0 {
		return false
	}

	return errors.Is(uw[0], target)
}
