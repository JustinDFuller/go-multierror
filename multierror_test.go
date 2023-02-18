package multierror_test

import (
	"encoding"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/justindfuller/go-multierror"
)

var (
	errSentinelOne   = errors.New("sentinel one")
	errSentinelTwo   = errors.New("sentinel two")
	errSentinelThree = errors.New("sentinel three")
	errSentinelFour  = errors.New("sentinel four")
)

func TestMultiError(t *testing.T) {
	if err := multierror.Join(nil); err != nil {
		t.Errorf("Expected err to be nil, got %s", err)
	}

	if err := multierror.Join(nil, nil, nil); err != nil {
		t.Errorf("Expected err to be nil, got %s", err)
	}

	if err := multierror.Join(nil, errSentinelOne); !errors.Is(err, errSentinelOne) {
		t.Errorf("Expected err to contain %s, got %s", errSentinelOne, err)
	}

	if err := multierror.Join(errSentinelOne); !errors.Is(err, errSentinelOne) {
		t.Errorf("Expected err to contain %s, got %s", errSentinelOne, err)
	}

	if err := multierror.Join(errSentinelOne, errSentinelTwo, errSentinelThree); !errors.Is(err, errSentinelOne) || !errors.Is(err, errSentinelTwo) || !errors.Is(err, errSentinelThree) {
		t.Errorf("Missing one of the sentinel errors, got %s", err)
	}

	if err := multierror.Join(errSentinelOne, nil, errSentinelThree); !errors.Is(err, errSentinelOne) || errors.Is(err, errSentinelTwo) || !errors.Is(err, errSentinelThree) {
		t.Errorf("Missing one of the sentinel errors, got %s", err)
	}

	if err := multierror.Join(nil, nil, errSentinelThree); errors.Is(err, errSentinelOne) || errors.Is(err, errSentinelTwo) || !errors.Is(err, errSentinelThree) {
		t.Errorf("Missing one of the sentinel errors, got %s", err)
	}

	if s := multierror.Join(errSentinelOne, errSentinelTwo, errSentinelThree).Error(); s != "Found 3 errors:\n\tsentinel one\n\tsentinel two\n\tsentinel three\n" {
		t.Errorf("Unexpected string, got %s", s)
	}

	stringer, ok := multierror.Join(errSentinelOne, errSentinelTwo, errSentinelThree).(fmt.Stringer)
	if !ok {
		t.Error("Expected the resulting error to implement fmt.Stringer")
	}

	if s := stringer.String(); s != "Found 3 errors:\n\tsentinel one\n\tsentinel two\n\tsentinel three\n" {
		t.Errorf("Unexpected string, got %s", s)
	}

	gostringer, ok := multierror.Join(errSentinelOne, errSentinelTwo, errSentinelThree).(fmt.GoStringer)
	if !ok {
		t.Error("Expected the resulting error to implement GoStringer")
	}
	if s := gostringer.GoString(); s != `[3]error{"sentinel one","sentinel two","sentinel three"}` {
		t.Errorf("Unexpected string, got %s", s)
	}

	err1 := multierror.Join(errSentinelOne, errSentinelTwo)
	err2 := multierror.Join(errSentinelThree, errSentinelFour)
	if err := multierror.Join(err1, err2); !errors.Is(err, errSentinelOne) || !errors.Is(err, errSentinelTwo) || !errors.Is(err, errSentinelThree) || !errors.Is(err, errSentinelFour) {
		t.Errorf("(Recursive) Missing one of the sentinel errors, got %s", err)
	}

	if s := multierror.Join(err1, err2).Error(); s != "Found 4 errors:\n\tsentinel one\n\tsentinel two\n\tsentinel three\n\tsentinel four\n" {
		t.Errorf("(Recursive) Unexpected string, got %s", s)
	}

	stringer, ok = multierror.Join(err1, err2).(fmt.Stringer)
	if !ok {
		t.Error("Expected the resulting error to implement fmt.Stringer")
	}
	if s := stringer.String(); s != "Found 4 errors:\n\tsentinel one\n\tsentinel two\n\tsentinel three\n\tsentinel four\n" {
		t.Errorf("(Recursive) Unexpected string, got %s", s)
	}

	gostringer, ok = multierror.Join(err1, err2).(fmt.GoStringer)
	if !ok {
		t.Error("Expected the resulting error to implement GoStringer")
	}
	if s := gostringer.GoString(); s != `[4]error{"sentinel one","sentinel two","sentinel three","sentinel four"}` {
		t.Errorf("(Recursive) Unexpected GoString, got %s", s)
	}

	if _, err := os.Open("non-existing"); err != nil {
		err := multierror.Join(err, errSentinelOne)
		var pathError *fs.PathError
		if !errors.As(err, &pathError) {
			t.Errorf("Expected Join to support errors.As: %s", err)
		}
	}

	if _, err := os.Open("non-existing"); err != nil {
		err := multierror.Join(errSentinelOne, err)
		var pathError *fs.PathError
		if !errors.As(err, &pathError) {
			t.Errorf("Expected Join to support errors.As: %s", err)
		}
	}

	if b, err := json.Marshal(multierror.Join(err1, err2)); err != nil || string(b) != `"sentinel one, sentinel two, sentinel three, sentinel four"` {
		t.Errorf("Expected Join to support json.Marshal: %s, %s", err, string(b))
	}

	unwrapped := errors.Unwrap(multierror.Join(err1, err2))
	if !errors.Is(unwrapped, errSentinelOne) || !errors.Is(unwrapped, errSentinelTwo) {
		t.Errorf("Expected to unwrap to return wrapped error, got: %s", unwrapped)
	}
	unwrapped = errors.Unwrap(unwrapped)
	if !errors.Is(unwrapped, errSentinelThree) || !errors.Is(unwrapped, errSentinelFour) {
		t.Errorf("Expected unwrap to unwrap second error, got %s", unwrapped)
	}
	unwrapped = errors.Unwrap(unwrapped)
	if unwrapped != nil {
		t.Errorf("Expected unwrap to be nil, got %s", unwrapped)
	}

	var builder strings.Builder
	if err := gob.NewEncoder(&builder).Encode(multierror.Join(err1, err2)); err != nil {
		t.Errorf("Expected Join to support gob.Encode, got err: %s", err)
	}
	if s := builder.String(); !strings.Contains(s, "Found 4 errors:\n\tsentinel one\n\tsentinel two\n\tsentinel three\n\tsentinel four\n") {
		t.Errorf("Expected gob to create string, got: %s", s)
	}
}

func ExampleJoin() {
	err1 := errors.New("something bad happened")
	err2 := errors.New("something is broken")

	err := multierror.Join(err1, err2)

	fmt.Println(err)
	// output: Found 2 errors:
	//	something bad happened
	//	something is broken
}

func ExampleJoin_nested() {
	err1 := errors.New("something bad happened")
	err2 := errors.New("something is broken")
	err3 := errors.New("something is REALLY broken")

	err := multierror.Join(err1, err2)
	err = multierror.Join(err, err3)
	fmt.Println(err)
	// output: Found 3 errors:
	//	something bad happened
	//	something is broken
	//	something is REALLY broken
}

func ExampleJoin_errorsIs() {
	err1 := errors.New("something bad happened")
	err2 := errors.New("something is broken")
	err3 := errors.New("something is REALLY broken")

	err := multierror.Join(err1, err2)
	err = multierror.Join(err, err3)
	fmt.Println(errors.Is(err, err1))
	// output: true
}

func ExampleJoin_errorsAs() {
	_, err := os.Open("non-existing")
	if err == nil {
		fmt.Println("No error")
	}

	err = multierror.Join(err, errSentinelOne)

	var pathError *fs.PathError
	fmt.Println(errors.As(err, &pathError))
	// output: true
}

func ExampleJoin_jsonMarshal() {
	err1 := errors.New("something bad happened")
	err2 := errors.New("something is broken")

	err := multierror.Join(err1, err2)
	b, _ := json.Marshal(err)

	fmt.Println(string(b))
	// output: "something bad happened, something is broken"
}

func ExampleJoin_fmtGoSyntaxRepresentation() {
	err1 := errors.New("something bad happened")
	err2 := errors.New("something is broken")

	err := multierror.Join(err1, err2)

	fmt.Printf("%#v", err)
	// output: [2]error{"something bad happened","something is broken"}
}

func ExampleJoin_fmtDefault() {
	err1 := errors.New("something bad happened")
	err2 := errors.New("something is broken")

	err := multierror.Join(err1, err2)

	fmt.Printf("%v", err)
	// output: Found 2 errors:
	//	something bad happened
	//	something is broken
}

func ExampleJoin_gobEncode() {
	err1 := errors.New("something bad happened")
	err2 := errors.New("something is broken")

	var builder strings.Builder
	if err := gob.NewEncoder(&builder).Encode(multierror.Join(err1, err2)); err != nil {
		fmt.Println(err)
		return
	}

	b, err := json.Marshal(builder.String())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))
	// output: "\n\ufffd\ufffd\u0005\u0001\u0002\ufffd\ufffd\u0000\u0000\u0000E\ufffd\ufffd\u0000A@\u000c\u0000=Found 2 errors:\n\tsomething bad happened\n\tsomething is broken\n"
}

func ExampleJoin_textMarshaler() {
	err1 := errors.New("something bad happened")
	err2 := errors.New("something is broken")

	marshaler, ok := multierror.Join(err1, err2).(encoding.TextMarshaler)
	if !ok {
		fmt.Println("Not an encoding.TextMarshaler")
		return
	}

	b, err := marshaler.MarshalText()
	if err != nil {
		fmt.Printf("Error encoding text: %s\n", err)
		return
	}

	fmt.Println(string(b))
	// output: Found 2 errors:
	//	something bad happened
	//	something is broken
}

func ExampleJoin_binaryMarshaler() {
	err1 := errors.New("something bad happened")
	err2 := errors.New("something is broken")

	marshaler, ok := multierror.Join(err1, err2).(encoding.BinaryMarshaler)
	if !ok {
		fmt.Println("Not an encoding.TextMarshaler")
		return
	}

	b, err := marshaler.MarshalBinary()
	if err != nil {
		fmt.Printf("Error encoding text: %s\n", err)
		return
	}

	b, err = json.Marshal(b)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))
	// output: "Cv+BBQEC/4QAAABF/4IAQUAMAD1Gb3VuZCAyIGVycm9yczoKCXNvbWV0aGluZyBiYWQgaGFwcGVuZWQKCXNvbWV0aGluZyBpcyBicm9rZW4K"
}

func ExampleJoin_nil() {
	err := multierror.Join(nil, nil, nil)
	fmt.Println(err)
	// output: <nil>
}
