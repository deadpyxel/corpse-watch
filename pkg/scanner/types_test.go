package scanner

import (
	"fmt"
	"testing"
)

func TestResultStringNoError(t *testing.T) {
	result := Result{
		URL:    "http://example.com",
		Status: 200,
		Error:  nil,
	}

	expected := "Got status 200 for URL http://example.com with no errors"
	if result.String() != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result.String())
	}
}

func TestResultStringWithError(t *testing.T) {
	err := fmt.Errorf("An error occurred")
	result := Result{
		URL:    "http://example.com",
		Status: 500,
		Error:  err,
	}

	expected := "Got status 500 for URL http://example.com with error [An error occurred]"
	if result.String() != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result.String())
	}
}
