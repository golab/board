package assert

import (
	"testing"
)

func Equal[V comparable](t *testing.T, got, expected V, msg string) {
	t.Helper()
	if got != expected {
		t.Errorf("%s: (expected %v, got %v)", msg, expected, got)
	}
}

func True(t *testing.T, got bool, msg string) {
	Equal(t, got, true, msg)
}

func Zero[V comparable](t *testing.T, got V, msg string) {
	var expected V
	Equal(t, got, expected, msg)
}

func NoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: (expected nil error, got %v)", msg, err)
	}
}
