/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package assert

import (
	"testing"
)

func Equal[V comparable](t *testing.T, got, expected V) {
	t.Helper()
	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func NotEqual[V comparable](t *testing.T, got, expectednot V) {
	t.Helper()
	if got == expectednot {
		t.Errorf("expected something else, got %v", got)
	}
}

func True(t *testing.T, got bool) {
	Equal(t, got, true)
}

func False(t *testing.T, got bool) {
	Equal(t, got, false)
}

func Zero[V comparable](t *testing.T, got V) {
	var expected V
	Equal(t, got, expected)
}

func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func ErrorIs(t *testing.T, got, expected error) {
	Equal(t, got, expected)
}

func NotNil[V comparable](t *testing.T, got V) {
	var vnil V
	NotEqual(t, got, vnil)
}
