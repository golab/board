/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package require

type Tester interface {
	Errorf(string, ...any)
	Fatalf(string, ...any)
	Helper()
}

func Equal[V comparable](t Tester, got, expected V) {
	t.Helper()
	if got != expected {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func NotEqual[V comparable](t Tester, got, expectednot V) {
	t.Helper()
	if got == expectednot {
		t.Fatalf("expected something else, got %v", got)
	}
}

func True(t Tester, got bool) {
	Equal(t, got, true)
}

func False(t Tester, got bool) {
	Equal(t, got, false)
}

func Zero[V comparable](t Tester, got V) {
	var expected V
	Equal(t, got, expected)
}

func NoError(t Tester, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func ErrorIs(t Tester, got, expected error) {
	if got == nil || expected == nil {
		t.Fatalf("unexpected nil error")
		return
	}
	Equal(t, got.Error(), expected.Error())
}

func NotNil(t Tester, got any) {
	NotEqual(t, got, nil)
}
