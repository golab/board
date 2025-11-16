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

	assert.NoError(t, err)
	assert.Equal(t, s.Size(), 19)
}

func TestState2(t *testing.T) {
	input := sgfsamples.SimpleFourMoves
	s, err := state.FromSGF(input)
	assert.NoError(t, err)

	sgf := s.ToSGF()
	sgfix := s.ToSGFIX()

	assert.Equal(t, len(sgf), len(input))

	assert.Equal(t, len(sgfix), 132)
}

func TestParseTTPAss(t *testing.T) {
	input := sgfsamples.PassWithTT
	s, err := state.FromSGF(input)
	assert.NoError(t, err)
	node := s.Nodes()[358]
	assert.Equal(t, node.XY, nil)
}

func TestNewState(t *testing.T) {
	s1 := state.NewState(19, true)
	s2 := state.NewState(19, false)
	assert.Equal(t, s1.GetNextIndex(), 1)
	assert.Equal(t, s2.GetNextIndex(), 0)
}
