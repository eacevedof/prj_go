package tests

import (
	"reflect"
	"strings"
	"testing"
)

// las aserciones de PHPUnit que usan los tests, con el mismo nombre

func AssertEquals(t *testing.T, expected any, actual any, message ...string) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("%s expected %#v, got %#v", strings.Join(message, " "), expected, actual)
	}
}

func AssertEmpty(t *testing.T, actual any) {
	t.Helper()
	if !isEmpty(actual) {
		t.Fatalf("expected empty, got %#v", actual)
	}
}

func AssertNotEmpty(t *testing.T, actual any) {
	t.Helper()
	if isEmpty(actual) {
		t.Fatalf("expected not empty, got %#v", actual)
	}
}

func AssertStringContainsString(t *testing.T, needle string, haystack string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected %q in %q", needle, haystack)
	}
}

func isEmpty(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Slice, reflect.Map, reflect.String, reflect.Array:
		return reflected.Len() == 0
	case reflect.Ptr:
		return reflected.IsNil()
	}
	return reflected.IsZero()
}
