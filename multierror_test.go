package multierror_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
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

	if s := multierror.Join(errSentinelOne, errSentinelTwo, errSentinelThree).String(); s != "Found 3 errors:\n\tsentinel one\n\tsentinel two\n\tsentinel three\n" {
		t.Errorf("Unexpected string, got %s", s)
	}

	if s := multierror.Join(errSentinelOne, errSentinelTwo, errSentinelThree).GoString(); s != "[3]error{sentinel one,sentinel two,sentinel three}" {
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

	if s := multierror.Join(err1, err2).String(); s != "Found 4 errors:\n\tsentinel one\n\tsentinel two\n\tsentinel three\n\tsentinel four\n" {
		t.Errorf("(Recursive) Unexpected string, got %s", s)
	}

	if s := multierror.Join(err1, err2).GoString(); s != "[4]error{sentinel one,sentinel two,sentinel three,sentinel four}" {
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

func ExampleJoin_jsonmarshal() {
	err1 := errors.New("something bad happened")
	err2 := errors.New("something is broken")

	err := multierror.Join(err1, err2)
	b, _ := json.Marshal(err)

	fmt.Println(string(b))
	// output: "something bad happened, something is broken"
}
