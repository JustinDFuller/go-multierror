package multierror

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type multiError struct {
	errors []error
}

func (m *multiError) MarshalJSON() ([]byte, error) {
	if m == nil {
		return nil, nil
	}

	var errs []error
	for _, e := range m.errors {
		errs = append(errs, flatten(e)...)
	}

	var errStrings []string
	for _, err := range errs {
		errStrings = append(errStrings, err.Error())
	}

	return json.Marshal(strings.Join(errStrings, ", "))
}

func (m *multiError) GoString() string {
	if m == nil {
		return "[]error{nil}"
	}

	var errs []error
	for _, e := range m.errors {
		errs = append(errs, flatten(e)...)
	}

	var errStrings []string
	for _, err := range errs {
		errStrings = append(errStrings, fmt.Sprintf(`"%s"`, err.Error()))
	}

	return fmt.Sprintf("[%d]error{%s}", len(errStrings), strings.Join(errStrings, ","))
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
		errs = append(errs, flatten(e)...)
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

func flatten(err error) []error {
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
			flattened = append(flattened, flatten(e)...)
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

func (m *multiError) Unwrap() []error {
	return m.errors
}

func (m *multiError) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(m.Error()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
