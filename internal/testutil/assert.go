package testutil

import (
	"reflect"
	"testing"
)

// AssertEqual verifica se dois valores são iguais
func AssertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v (type %T), got %v (type %T)", expected, expected, actual, actual)
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
	}
}

// AssertNotEqual verifica se dois valores são diferentes
func AssertNotEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected values to be different, but both are %v", expected)
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
	}
}

// AssertNil verifica se um valor é nil
func AssertNil(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if value != nil && !reflect.ValueOf(value).IsNil() {
		t.Errorf("Expected nil, got %v (type %T)", value, value)
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
	}
}

// AssertNotNil verifica se um valor não é nil
func AssertNotNil(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if value == nil || reflect.ValueOf(value).IsNil() {
		t.Errorf("Expected non-nil value, got nil")
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
	}
}

// AssertTrue verifica se uma condição é verdadeira
func AssertTrue(t *testing.T, condition bool, msgAndArgs ...interface{}) {
	t.Helper()
	if !condition {
		t.Error("Expected condition to be true, got false")
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
	}
}

// AssertFalse verifica se uma condição é falsa
func AssertFalse(t *testing.T, condition bool, msgAndArgs ...interface{}) {
	t.Helper()
	if condition {
		t.Error("Expected condition to be false, got true")
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
	}
}

// AssertNoError verifica se um erro é nil
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
	}
}

// AssertError verifica se um erro não é nil
func AssertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err == nil {
		t.Error("Expected an error, got nil")
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
	}
}

// AssertErrorContains verifica se um erro contém uma mensagem específica
func AssertErrorContains(t *testing.T, err error, expectedMsg string, msgAndArgs ...interface{}) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error containing '%s', got nil", expectedMsg)
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
		return
	}

	if !contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error containing '%s', got '%s'", expectedMsg, err.Error())
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
	}
}

// AssertLen verifica o comprimento de uma slice/array/map
func AssertLen(t *testing.T, value interface{}, expected int, msgAndArgs ...interface{}) {
	t.Helper()

	v := reflect.ValueOf(value)
	kind := v.Kind()

	if kind != reflect.Slice && kind != reflect.Array && kind != reflect.Map && kind != reflect.String {
		t.Errorf("Cannot check length of type %T", value)
		return
	}

	actual := v.Len()
	if actual != expected {
		t.Errorf("Expected length %d, got %d", expected, actual)
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
	}
}

// AssertContains verifica se uma string contém uma substring
func AssertContains(t *testing.T, haystack, needle string, msgAndArgs ...interface{}) {
	t.Helper()
	if !contains(haystack, needle) {
		t.Errorf("Expected '%s' to contain '%s'", haystack, needle)
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
	}
}

// AssertNotContains verifica se uma string não contém uma substring
func AssertNotContains(t *testing.T, haystack, needle string, msgAndArgs ...interface{}) {
	t.Helper()
	if contains(haystack, needle) {
		t.Errorf("Expected '%s' to not contain '%s'", haystack, needle)
		if len(msgAndArgs) > 0 {
			t.Log(msgAndArgs...)
		}
	}
}

// AssertGreaterThan verifica se um valor é maior que outro
func AssertGreaterThan(t *testing.T, actual, expected interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	av := reflect.ValueOf(actual)
	ev := reflect.ValueOf(expected)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if av.Int() <= ev.Int() {
			t.Errorf("Expected %v to be greater than %v", actual, expected)
		}
	case reflect.Float32, reflect.Float64:
		if av.Float() <= ev.Float() {
			t.Errorf("Expected %v to be greater than %v", actual, expected)
		}
	default:
		t.Errorf("Cannot compare type %T", actual)
	}

	if len(msgAndArgs) > 0 {
		t.Log(msgAndArgs...)
	}
}

// contains verifica se uma string contém outra
func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (haystack == needle || len(needle) == 0 ||
		func() bool {
			for i := 0; i <= len(haystack)-len(needle); i++ {
				if haystack[i:i+len(needle)] == needle {
					return true
				}
			}
			return false
		}())
}
