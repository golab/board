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

func TestState3(t *testing.T) {
	s, err := FromSGF(sgfsamples.SimpleTwoBranches)
	assert.NoError(t, err)

	// index 5 is the start of the second branch from the root
	s.gotoIndex(5) //nolint:errcheck

	prefs := s.prefs()
	assert.Equal(t, prefs["0"], 1)

	err = s.setPreferred(4)
	assert.NoError(t, err)

	prefs = s.prefs()
	assert.Equal(t, prefs["0"], 0)

	s.resetPrefs()
	prefs = s.prefs()
	for _, value := range prefs {
		assert.Equal(t, value, 0)
	}

	s.gotoIndex(6) //nolint:errcheck
	assert.Equal(t, s.locate(), "1,0")

}

func TestState4(t *testing.T) {
	s, err := FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	assert.Equal(t, s.Head().Index, 8)

	// push a node to the head
	s.PushHead(0, 0, core.Black)
	assert.Equal(t, s.Current().Index, 0)

	// go to the most recently pushed node
	s.gotoIndex(9) //nolint:errcheck

	// since we're at the head, we will "track" along
	s.PushHead(0, 1, core.White)
	assert.Equal(t, s.Current().Index, 10)
}

func TestAddStones(t *testing.T) {
	s, err := FromSGF(sgfsamples.Empty)
	assert.NoError(t, err)

	// add three new moves
	moves := []*core.Stone{
		core.NewStone(9, 9, core.Black),
		core.NewStone(10, 10, core.White),
		core.NewStone(11, 11, core.Black),
	}

	s.AddStones(moves)

	assert.Equal(t, len(s.Current().Down), 1)

	s.fastForward()
	assert.Equal(t, s.Current().Index, 3)
}

func TestCut(t *testing.T) {
	s, err := FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	s.gotoIndex(5) //nolint:errcheck
	s.cut()
	assert.Equal(t, s.Current().Index, 4)
	assert.Equal(t, len(s.Current().Down), 0)
}

func TestGraft(t *testing.T) {
	s, err := FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	moves := []*core.Stone{
		core.NewStone(9, 9, core.Black),
		core.NewStone(10, 10, core.White),
		core.NewStone(11, 11, core.Black),
		core.NewStone(12, 12, core.Black),
	}

	// dumb graft
	s.graft(4, moves[:3])
	s.gotoIndex(4) //nolint:errcheck
	s.down()
	s.fastForward()
	assert.Equal(t, s.Current().Index, 11)

	// smart graft
	s.smartGraft(4, moves)
	s.fastForward()
	assert.Equal(t, s.Current().Index, 12)

	// double-checking Up works
	s.gotoIndex(4) //nolint:errcheck
	s.up()
	s.fastForward()
	assert.Equal(t, s.Current().Index, 8)

	// check for gotocoord
	s.gotoCoord(3, 3)
	assert.Equal(t, s.Current().Index, 2)
}

func TestStateJSON(t *testing.T) {
	s, err := FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	s.gotoIndex(4) //nolint:errcheck

	j := s.Save()
	assert.Equal(t, j.NextIndex, 9)
	assert.Equal(t, j.Location, "0,0,0,0")

	tr := s.saveTree(FullFrame)
	assert.Equal(t, tr.Depth, 8)
	assert.Equal(t, tr.Current, 4)
	assert.Equal(t, tr.Up, 0)
	assert.Equal(t, tr.Root, 0)
}

func TestGenerateMarks(t *testing.T) {
	s, err := FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	s.right()

	// add a triangle
	coord := core.NewCoord(9, 9)
	_, err = NewAddTriangleCommand(coord).Execute(s)
	assert.NoError(t, err)

	// add a square
	coord = core.NewCoord(10, 10)
	_, err = NewAddSquareCommand(coord).Execute(s)
	assert.NoError(t, err)

	// add a letter
	coord = core.NewCoord(11, 11)
	letter := "A"
	_, err = NewAddLetterCommand(coord, letter).Execute(s)
	assert.NoError(t, err)

	// add a number
	coord = core.NewCoord(12, 12)
	number := 1
	_, err = NewAddNumberCommand(coord, number).Execute(s)
	assert.NoError(t, err)

	// add pen stroke
	_, err = NewDrawCommand(10.0, 10.0, 20.0, 20.0, "#AAAAAA").Execute(s)
	assert.NoError(t, err)

	marks := s.generateMarks()

	assert.True(t, marks.Current.Equal(core.NewCoord(15, 3)))

	assert.Equal(t, len(marks.Triangles), 1)
	assert.True(t, marks.Triangles[0].Equal(core.NewCoord(9, 9)))

	assert.Equal(t, len(marks.Squares), 1)
	assert.True(t, marks.Squares[0].Equal(core.NewCoord(10, 10)))

	assert.Equal(t, len(marks.Labels), 2)

	assert.True(t, marks.Labels[0].Coord.Equal(core.NewCoord(11, 11)))
	assert.True(t, marks.Labels[1].Coord.Equal(core.NewCoord(12, 12)))

	assert.Equal(t, marks.Labels[0].Text, "A")
	assert.Equal(t, marks.Labels[1].Text, "1")

	assert.Equal(t, len(marks.Pens), 1)
	assert.Equal(t, marks.Pens[0].Color, "#AAAAAA")
}
