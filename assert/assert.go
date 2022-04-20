package assert

import (
	"reflect"
	"testing"
)

func AssertTrue(t *testing.T, expected bool) {
	t.Helper()

	if !expected {
		t.Fatalf("expected true")
	}
}

func AssertEqual[T comparable](t *testing.T, expected T, actual T) {
	t.Helper()

	if expected != actual {
		t.Fatalf("expected: %v, actual: %v", expected, actual)
	}
}

func AssertSliceEqual[T comparable](t *testing.T, expected []T, actual []T) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Fatalf("slices length are different, expected: %d, actual: %d", len(expected), len(actual))
	}

	for i := range expected {
		if expected[i] != actual[i] {
			t.Fatalf("index: %d, expected: %v, actual: %v", i, expected[i], actual[i])
		}
	}
}

func AssertListEqual(t *testing.T, expected []interface{}, actual []interface{}) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Fatalf("slices length are different, expected: %d, actual: %d", len(expected), len(actual))
	}

	for i := range expected {
		if expected[i] != actual[i] {
			t.Fatalf("index: %d, expected: %v (%T), actual: %v (%T)", i, expected[i], expected[i], actual[i], actual[i])
		}
	}
}

func AssertNil(t *testing.T, expected interface{}) {
	t.Helper()

	if expected == nil ||
		(reflect.ValueOf(expected).Kind() == reflect.Ptr &&
			reflect.ValueOf(expected).IsNil()) {
		// pass
	} else {
		t.Fatal("expected nil")
	}
}
