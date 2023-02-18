package multierror

// Join joins one error with one or more other errors.
// The resulting value is a multiError containing all non-nil errors provided.
// If the first error is a multiError, the rest of the errors will be appended to the existing multiError.
// If the first error is not a multiError, a new multiError will be created for it and all the rest of the errors as well.
func Join(errors ...error) *multiError {
	var errs []error
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
