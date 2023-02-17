package multierror_test

import (
	"errors"
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
	if err := multierror.Append(nil); err != nil {
		t.Errorf("Expected err to be nil, got %s", err)
	}

	if err := multierror.Append(nil, nil, nil); err != nil {
		t.Errorf("Expected err to be nil, got %s", err)
	}

	if err := multierror.Append(nil, errSentinelOne); !errors.Is(err, errSentinelOne) {
		t.Errorf("Expected err to contain %s, got %s", errSentinelOne, err)
	}

	if err := multierror.Append(errSentinelOne); !errors.Is(err, errSentinelOne) {
		t.Errorf("Expected err to contain %s, got %s", errSentinelOne, err)
	}

	if err := multierror.Append(errSentinelOne, errSentinelTwo, errSentinelThree); !errors.Is(err, errSentinelOne) || !errors.Is(err, errSentinelTwo) || !errors.Is(err, errSentinelThree) {
		t.Errorf("Missing one of the sentinel errors, got %s", err)
	}

	if err := multierror.Append(errSentinelOne, nil, errSentinelThree); !errors.Is(err, errSentinelOne) || errors.Is(err, errSentinelTwo) || !errors.Is(err, errSentinelThree) {
		t.Errorf("Missing one of the sentinel errors, got %s", err)
	}

	if err := multierror.Append(nil, nil, errSentinelThree); errors.Is(err, errSentinelOne) || errors.Is(err, errSentinelTwo) || !errors.Is(err, errSentinelThree) {
		t.Errorf("Missing one of the sentinel errors, got %s", err)
	}

	if s := multierror.Append(errSentinelOne, errSentinelTwo, errSentinelThree).Error(); s != "Found 3 errors:\n\tsentinel one\n\tsentinel two\n\tsentinel three\n" {
		t.Errorf("Unexpected string, got %s", s)
	}

	err1 := multierror.Append(errSentinelOne, errSentinelTwo)
	err2 := multierror.Append(errSentinelThree, errSentinelFour)
	if err := multierror.Append(err1, err2); !errors.Is(err, errSentinelOne) || !errors.Is(err, errSentinelTwo) || !errors.Is(err, errSentinelThree) || !errors.Is(err, errSentinelFour) {
		t.Errorf("(Recursive) Missing one of the sentinel errors, got %s", err)
	}

	if s := multierror.Append(err1, err2).Error(); s != "Found 4 errors:\n\tsentinel one\n\tsentinel two\n\tsentinel three\n\tsentinel four\n" {
		t.Errorf("(Recursive) Unexpected string, got %s", s)
	}

	if _, err := os.Open("non-existing"); err != nil {
		err := multierror.Append(err, errSentinelOne)
		var pathError *fs.PathError
		if !errors.As(err, &pathError) {
			t.Errorf("Expected Append to support errors.As: %s", err)
		}
	}

	if _, err := os.Open("non-existing"); err != nil {
		err := multierror.Append(errSentinelOne, err)
		var pathError *fs.PathError
		if !errors.As(err, &pathError) {
			t.Errorf("Expected Append to support errors.As: %s", err)
		}
	}
}
