/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/internal/sgfsamples"
	"github.com/jarednogo/board/pkg/core"
)

func TestScore1(t *testing.T) {
	s, err := FromSGF(sgfsamples.Scoring1)
	assert.NoError(t, err)

	s.fastForward()

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
		coord := core.NewCoord(d[0], d[1])
		gp := s.Board().FindGroup(coord)
		markedDead.AddAll(gp.Coords)
	}

	scoreResult := s.board.Score(markedDead, core.NewCoordSet())
	blackArea := scoreResult.BlackArea
	whiteArea := scoreResult.WhiteArea
	blackDead := scoreResult.BlackDead
	whiteDead := scoreResult.WhiteDead
	dame := scoreResult.Dame

	current := s.Current()

	assert.Equal(t, len(blackArea), 56)
	assert.Equal(t, len(whiteArea), 40)

	assert.Equal(t, len(blackDead), 9)
	assert.Equal(t, len(whiteDead), 9)

	assert.Equal(t, current.BlackCaps, 27)
	assert.Equal(t, current.WhiteCaps, 25)

	assert.Equal(t, len(dame), 7)
}

func TestScore2(t *testing.T) {
	s, err := FromSGF(sgfsamples.Scoring2)
	assert.NoError(t, err)

	s.fastForward()

	dead := [][2]int{
		{3, 14},
		{15, 12},
	}

	markedDead := core.NewCoordSet()
	for _, d := range dead {
		coord := core.NewCoord(d[0], d[1])
		gp := s.Board().FindGroup(coord)
		markedDead.AddAll(gp.Coords)
	}

	scoreResult := s.board.Score(markedDead, core.NewCoordSet())
	blackArea := scoreResult.BlackArea
	whiteArea := scoreResult.WhiteArea
	blackDead := scoreResult.BlackDead
	whiteDead := scoreResult.WhiteDead
	dame := scoreResult.Dame

	current := s.Current()

	assert.Equal(t, len(blackArea), 72)
	assert.Equal(t, len(whiteArea), 62)

	assert.Equal(t, len(blackDead), 2)
	assert.Equal(t, len(whiteDead), 2)

	assert.Equal(t, current.BlackCaps, 2)
	assert.Equal(t, current.WhiteCaps, 11)

	assert.Equal(t, len(dame), 8)
}

func TestScore3(t *testing.T) {
	s, err := FromSGF(sgfsamples.TwoFalseEyes)
	assert.NoError(t, err)

	s.fastForward()

	markedDead := core.NewCoordSet()

	scoreResult := s.board.Score(markedDead, core.NewCoordSet())
	blackArea := scoreResult.BlackArea
	whiteArea := scoreResult.WhiteArea
	blackDead := scoreResult.BlackDead
	whiteDead := scoreResult.WhiteDead
	dame := scoreResult.Dame

	current := s.Current()

	assert.Equal(t, len(blackArea), 2)
	assert.Equal(t, len(whiteArea), 38)

	assert.Equal(t, len(blackDead), 0)
	assert.Equal(t, len(whiteDead), 0)

	assert.Equal(t, current.BlackCaps, 0)
	assert.Equal(t, current.WhiteCaps, 0)

	assert.Equal(t, len(dame), 2)
}

func TestScore4(t *testing.T) {
	s, err := FromSGF(sgfsamples.AtariDame1)
	assert.NoError(t, err)

	s.fastForward()

	dead := [][2]int{
		{5, 6},
	}

	markedDead := core.NewCoordSet()
	for _, d := range dead {
		coord := core.NewCoord(d[0], d[1])
		gp := s.Board().FindGroup(coord)
		markedDead.AddAll(gp.Coords)
	}

	scoreResult := s.board.Score(markedDead, core.NewCoordSet())
	blackArea := scoreResult.BlackArea
	whiteArea := scoreResult.WhiteArea
	blackDead := scoreResult.BlackDead
	whiteDead := scoreResult.WhiteDead
	dame := scoreResult.Dame

	current := s.Current()

	assert.Equal(t, len(blackArea), 3)
	assert.Equal(t, len(whiteArea), 44)

	assert.Equal(t, len(blackDead), 0)
	assert.Equal(t, len(whiteDead), 1)

	assert.Equal(t, current.BlackCaps, 0)
	assert.Equal(t, current.WhiteCaps, 0)

	assert.Equal(t, len(dame), 1)
}

func TestScore5(t *testing.T) {
	s, err := FromSGF(sgfsamples.AtariDame2)
	assert.NoError(t, err)

	s.fastForward()

	scoreResult := s.board.Score(core.NewCoordSet(), core.NewCoordSet())
	blackArea := scoreResult.BlackArea
	whiteArea := scoreResult.WhiteArea
	blackDead := scoreResult.BlackDead
	whiteDead := scoreResult.WhiteDead
	dame := scoreResult.Dame

	current := s.Current()

	assert.Equal(t, len(blackArea), 10)
	assert.Equal(t, len(whiteArea), 41)

	assert.Equal(t, len(blackDead), 0)
	assert.Equal(t, len(whiteDead), 0)

	assert.Equal(t, current.BlackCaps, 0)
	assert.Equal(t, current.WhiteCaps, 0)

	assert.Equal(t, len(dame), 1)
}

func TestScore6(t *testing.T) {
	s, err := FromSGF(sgfsamples.AtariDame3)
	assert.NoError(t, err)

	s.fastForward()

	scoreResult := s.board.Score(core.NewCoordSet(), core.NewCoordSet())
	blackArea := scoreResult.BlackArea
	whiteArea := scoreResult.WhiteArea
	blackDead := scoreResult.BlackDead
	whiteDead := scoreResult.WhiteDead
	dame := scoreResult.Dame

	current := s.Current()

	assert.Equal(t, len(blackArea), 4)
	assert.Equal(t, len(whiteArea), 52)

	assert.Equal(t, len(blackDead), 0)
	assert.Equal(t, len(whiteDead), 0)

	assert.Equal(t, current.BlackCaps, 0)
	assert.Equal(t, current.WhiteCaps, 0)

	assert.Equal(t, len(dame), 2)
}

func TestScore7(t *testing.T) {
	s, err := FromSGF(sgfsamples.AtariDame4)
	assert.NoError(t, err)

	s.fastForward()

	scoreResult := s.board.Score(core.NewCoordSet(), core.NewCoordSet())
	blackArea := scoreResult.BlackArea
	whiteArea := scoreResult.WhiteArea
	blackDead := scoreResult.BlackDead
	whiteDead := scoreResult.WhiteDead
	dame := scoreResult.Dame

	current := s.Current()

	assert.Equal(t, len(blackArea), 26)
	assert.Equal(t, len(whiteArea), 5)

	assert.Equal(t, len(blackDead), 0)
	assert.Equal(t, len(whiteDead), 0)

	assert.Equal(t, current.BlackCaps, 0)
	assert.Equal(t, current.WhiteCaps, 0)

	assert.Equal(t, len(dame), 1)
}

func TestScore8(t *testing.T) {
	s, err := FromSGF(sgfsamples.AtariDame5)
	assert.NoError(t, err)

	s.fastForward()

	scoreResult := s.board.Score(core.NewCoordSet(), core.NewCoordSet())
	blackArea := scoreResult.BlackArea
	whiteArea := scoreResult.WhiteArea
	blackDead := scoreResult.BlackDead
	whiteDead := scoreResult.WhiteDead
	dame := scoreResult.Dame

	current := s.Current()

	assert.Equal(t, len(blackArea), 2)
	assert.Equal(t, len(whiteArea), 36)

	assert.Equal(t, len(blackDead), 0)
	assert.Equal(t, len(whiteDead), 0)

	assert.Equal(t, current.BlackCaps, 0)
	assert.Equal(t, current.WhiteCaps, 0)

	assert.Equal(t, len(dame), 6)
}
