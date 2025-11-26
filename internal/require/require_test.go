/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package require_test

import (
	"fmt"
	"testing"

	"github.com/jarednogo/board/internal/require"
)

type MockTester struct {
	failed bool
}

func (t *MockTester) Errorf(s string, args ...any) {
	t.failed = true
}

func (t *MockTester) Fatalf(s string, args ...any) {
	t.failed = true
}

func (t *MockTester) Helper() {}

func TestEqual(t *testing.T) {
	h := &MockTester{}
	require.Equal(h, 1, 1)
	if h.failed {
		t.Errorf("should succeed")
	}

	h = &MockTester{}
	require.Equal(h, 1, 2)
	if !h.failed {
		t.Errorf("should fail")
	}
}

func TestNotEqual(t *testing.T) {
	h := &MockTester{}
	require.NotEqual(h, 1, 1)
	if !h.failed {
		t.Errorf("should succeed")
	}

	h = &MockTester{}
	require.NotEqual(h, 1, 2)
	if h.failed {
		t.Errorf("should fail")
	}
}

func TestTrue(t *testing.T) {
	h := &MockTester{}
	require.True(h, true)
	if h.failed {
		t.Errorf("should succeed")
	}

	h = &MockTester{}
	require.True(h, false)
	if !h.failed {
		t.Errorf("should fail")
	}
}

func TestFalse(t *testing.T) {
	h := &MockTester{}
	require.False(h, true)
	if !h.failed {
		t.Errorf("should succeed")
	}

	h = &MockTester{}
	require.False(h, false)
	if h.failed {
		t.Errorf("should fail")
	}
}

func TestZero(t *testing.T) {
	h := &MockTester{}
	require.Zero(h, 0)
	if h.failed {
		t.Errorf("should succeed")
	}

	h = &MockTester{}
	require.Zero(h, 1)
	if !h.failed {
		t.Errorf("should fail")
	}
}

func TestNoError(t *testing.T) {
	h := &MockTester{}
	require.NoError(h, nil)
	if h.failed {
		t.Errorf("should succeed")
	}

	h = &MockTester{}
	require.NoError(h, fmt.Errorf("example error"))
	if !h.failed {
		t.Errorf("should fail")
	}
}

func TestErrorIs(t *testing.T) {
	h := &MockTester{}
	err1 := fmt.Errorf("example error")
	err2 := fmt.Errorf("example error")
	require.ErrorIs(h, err1, err2)
	if h.failed {
		t.Errorf("should succeed")
	}

	h = &MockTester{}
	require.ErrorIs(h, nil, err1)
	if !h.failed {
		t.Errorf("should fail")
	}
}

func TestNotNil(t *testing.T) {
	h := &MockTester{}
	require.NotNil(h, 1)
	if h.failed {
		t.Errorf("should succeed")
	}

	h = &MockTester{}
	require.NotNil(h, nil)
	if !h.failed {
		t.Errorf("should fail")
	}
}
