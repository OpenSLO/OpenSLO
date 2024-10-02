package assert

import (
	"fmt"
	"reflect"
	"testing"
)

// Require fails the test if the provided boolean is false.
// It should be used in conjunction with assert functions.
// Example:
//
//	assert.Require(t, assert.AssertError(t, err))
func Require(t *testing.T, isPassing bool) {
	t.Helper()
	if !isPassing {
		t.FailNow()
	}
}

// Equal fails the test if the expected and actual values are not equal.
func Equal(t *testing.T, expected, actual interface{}) bool {
	t.Helper()
	if !areEqual(expected, actual) {
		return fail(t, "Expected: %v, actual: %v", expected, actual)
	}
	return true
}

// True fails the test if the actual value is not true.
func True(t *testing.T, actual bool) bool {
	t.Helper()
	if !actual {
		return fail(t, "Should be true")
	}
	return true
}

// False fails the test if the actual value is not false.
func False(t *testing.T, actual bool) bool {
	t.Helper()
	if actual {
		return fail(t, "Should be false")
	}
	return true
}

// Len fails the test if the value is not of the expected length.
func Len(t *testing.T, v interface{}, length int) bool {
	t.Helper()
	actual, err := getLen(v)
	if err != nil {
		return fail(t, "Error getting length: %v", err)
	}
	if actual != length {
		return fail(t, "Expected length: %d, actual: %d", length, actual)
	}
	return true
}

// Error fails the test if the error is nil.
func Error(t *testing.T, err error) bool {
	t.Helper()
	if err == nil {
		return fail(t, "An error is expected but actual nil.")
	}
	return true
}

// NoError fails the test if the error is not nil.
func NoError(t *testing.T, err error) bool {
	t.Helper()
	if err != nil {
		return fail(t, "Unexpected error:\n%+v", err)
	}
	return true
}

// NotEmpty fails the test if the value is empty.
func NotEmpty(t *testing.T, v any) bool {
	t.Helper()
	if isEmpty(v) {
		return fail(t, "Value should not be empty.")
	}
	return true
}

func areEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}
	if !reflect.DeepEqual(expected, actual) {
		return false
	}
	return true
}

func getLen(v interface{}) (int, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Map, reflect.String:
		return rv.Len(), nil
	default:
		return -1, fmt.Errorf("invalid type: %v", rv.Kind())
	}
}

func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map, reflect.Slice:
		return rv.Len() == 0
	case reflect.Ptr:
		if rv.IsNil() {
			return true
		}
		deref := rv.Elem().Interface()
		return isEmpty(deref)
	default:
		zero := reflect.Zero(rv.Type())
		return reflect.DeepEqual(v, zero.Interface())
	}
}

func fail(t *testing.T, msg string, a ...interface{}) bool {
	t.Helper()
	t.Errorf(msg, a...)
	return false
}
