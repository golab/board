/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state_test

import (
	//"fmt"
	"github.com/jarednogo/board/backend/state"
	"testing"
)

func TestState1(t *testing.T) {
	s, err := state.FromSGF("(;PW[White]RU[Japanese]KM[6.5]GM[1]FF[4]CA[UTF-8]SZ[19]PB[Black];B[pd];W[dd];B[pp];W[dp];B[];W[])")
	if err != nil {
		t.Error(err)
	}
	if s.Size != 19 {
		t.Errorf("error with state")
	}
}

func TestState2(t *testing.T) {
	input := "(;PW[White]RU[Japanese]KM[6.5]GM[1]FF[4]CA[UTF-8]SZ[19]PB[Black];B[pd];W[dd];B[pp];W[dp];B[];W[])"
	s, err := state.FromSGF(input)
	if err != nil {
		t.Error(err)
	}

	sgf := s.ToSGF(false)
	sgfix := s.ToSGF(true)

	if len(sgf) != len(input) {
		t.Errorf("error with state to sgf, expected %d, got: %d", len(input), len(sgf))
	}

	if len(sgfix) != 132 {
		t.Errorf("error with state to sgf (indexes), expected %d, got: %d", 132, len(sgfix))
	}
}
