package multierror

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
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

	errStrings := make([]string, 0, len(errs))

	for _, err := range errs {
		errStrings = append(errStrings, err.Error())
	}

	b, err := json.Marshal(strings.Join(errStrings, ", "))
	if err != nil {
		return nil, fmt.Errorf("error Marshaling multiError to json: %w", err)
	}

	return b, nil
}

func (m *multiError) GoString() string {
	if m == nil {
		return "[]error{nil}"
	}

	var errs []error
	for _, e := range m.errors {
		errs = append(errs, flatten(e)...)
	}

	errStrings := make([]string, 0, len(errs))

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

	var str string

	var errs []error
	for _, e := range m.errors {
		errs = append(errs, flatten(e)...)
	}

	if len(errs) == 1 {
		str += "Found one error:\n"
	} else {
		str += fmt.Sprintf("Found %d errors:\n", len(errs))
	}

	for _, err := range errs {
		str += fmt.Sprintf("\t%s\n", err.Error())
	}

	return str
}

func flatten(err error) []error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*multiError); ok {
		var flattened []error

		if len(e.errors) == 0 {
			return nil
		}

		for _, e := range e.errors {
			flattened = append(flattened, flatten(e)...)
		}

		return flattened
	}

	return []error{err}
}

func (m *multiError) Unwrap() error {
	if m == nil || len(m.errors) == 0 {
		return nil
	}

	if len(m.errors) == 1 {
		return m.errors[0]
	}

	errs := make([]error, len(m.errors))
	copy(errs, m.errors)

	return unwrapper(errs)
}

func (m *multiError) GobEncode() ([]byte, error) {
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(m.Error()); err != nil {
		return nil, fmt.Errorf("error encoding multiError with gob: %w", err)
	}

	return buf.Bytes(), nil
}

func (m *multiError) MarshalText() ([]byte, error) {
	return []byte(m.Error()), nil
}

func (m *multiError) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(m); err != nil {
		return nil, fmt.Errorf("error encoding multiError with gob: %w", err)
	}

	return buf.Bytes(), nil
}
