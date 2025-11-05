/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state_test

import (
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/internal/sgfsamples"
	"github.com/jarednogo/board/pkg/state"
)

func TestState(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleFourMoves)

	assert.NoError(t, err, "error should be nil")
	assert.Equal(t, s.Size(), 19, "error with state")
}

func TestState2(t *testing.T) {
	input := sgfsamples.SimpleFourMoves
	s, err := state.FromSGF(input)
	assert.NoError(t, err, "error should be nil")

	sgf := s.ToSGF(false)
	sgfix := s.ToSGF(true)

	assert.Equal(t, len(sgf), len(input), "error with state to sgf")

	assert.Equal(t, len(sgfix), 132, "error with state to sgf (indexes)")
}
