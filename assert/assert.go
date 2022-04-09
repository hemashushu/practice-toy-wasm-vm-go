package assert

import "testing"

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
		AssertEqual(t, expected[i], actual[i])
	}
}
