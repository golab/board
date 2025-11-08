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

func Equal[V comparable](t *testing.T, got, expected V, msg string) {
	t.Helper()
	if got != expected {
		t.Errorf("%s: (expected %v, got %v)", msg, expected, got)
	}
}

func True(t *testing.T, got bool, msg string) {
	Equal(t, got, true, msg)
}

func False(t *testing.T, got bool, msg string) {
	Equal(t, got, false, msg)
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

func ErrorIs(t *testing.T, got, expected error, msg string) {
	Equal(t, got, expected, msg)
}
