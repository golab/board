/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state_test

import (
	"testing"

	"github.com/jarednogo/board/internal/sgfsamples"
	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/state"
)

func TestScore(t *testing.T) {

	s, err := state.FromSGF(sgfsamples.Scoring1)
	if err != nil {
		t.Error(err)
		return
	}

	s.FastForward()

	dead := [][2]int{
		{4, 3},
		{5, 2},
		{6, 3},
		{13, 0},
		{14, 1},
		{16, 0},
		{11, 5},
		{18, 5},
		{1, 15},
		{3, 13},
	}

	markedDead := core.NewCoordSet()
	for _, d := range dead {
		coord := &core.Coord{X: d[0], Y: d[1]}
		gp := s.Board().FindGroup(coord)
		markedDead.AddAll(gp.Coords)
	}

	blackArea, whiteArea, blackDead, whiteDead, dame := s.Board().Score(
		markedDead, core.NewCoordSet())

	current := s.Current()

	if len(blackArea) != 56 {
		t.Errorf("expected len(blackArea) == 56 (got %d)", len(blackArea))
	}

	if len(whiteArea) != 40 {
		t.Errorf("expected len(whiteArea) == 40 (got %d)", len(whiteArea))
	}

	if len(blackDead) != 9 {
		t.Errorf("expected len(blackDead) == 9 (got %d)", len(blackDead))
	}

	if len(whiteDead) != 9 {
		t.Errorf("expected len(whiteDead) == 9 (got %d)", len(whiteDead))
	}

	if current.BlackCaps != 27 {
		t.Errorf("expected BlackCaps == 27 (got %d)", current.BlackCaps)
	}

	if current.WhiteCaps != 25 {
		t.Errorf("expected WhiteCaps == 25 (got %d)", current.WhiteCaps)
	}

	if len(dame) != 7 {
		t.Errorf("expected len(dame) == 7 (got %d)", len(dame))
	}

	blackScore := len(blackArea) + len(whiteDead) + current.BlackCaps
	whiteScore := len(whiteArea) + len(blackDead) + current.WhiteCaps

	if blackScore != 92 {
		t.Errorf("expected blackScore == 92 (got %d)", blackScore)
	}

	if whiteScore != 74 {
		t.Errorf("expected whiteScore == 74 (got %d)", whiteScore)
	}
}
