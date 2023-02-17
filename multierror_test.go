package multierror_test

import (
	"errors"
	"testing"

	"github.com/justindfuller/go-multierror"
)

var (
	errSentinelOne   = errors.New("sentinel one")
	errSentinelTwo   = errors.New("sentinel two")
	errSentinelThree = errors.New("sentinel three")
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

	if s := multierror.Append(errSentinelOne, errSentinelTwo, errSentinelThree).String(); s != "Found 3 errors:\n\tsentinel one\n\tsentinel two\n\tsentinel three\n" {
		t.Errorf("Unexpected string, got %s", s)
	}
}
