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
	testname := "scoring1"
	s, err := FromSGF(sgfsamples.Scoring1)
	assert.NoError(t, err, testname)

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
		coord := &core.Coord{X: d[0], Y: d[1]}
		gp := s.Board().FindGroup(coord)
		markedDead.AddAll(gp.Coords)
	}

	blackArea, whiteArea, blackDead, whiteDead, dame := s.Board().Score(
		markedDead, core.NewCoordSet())

	current := s.Current()

	assert.Equal(t, len(blackArea), 56, testname)
	assert.Equal(t, len(whiteArea), 40, testname)

	assert.Equal(t, len(blackDead), 9, testname)
	assert.Equal(t, len(whiteDead), 9, testname)

	assert.Equal(t, current.BlackCaps, 27, testname)
	assert.Equal(t, current.WhiteCaps, 25, testname)

	assert.Equal(t, len(dame), 7, testname)
}

func TestScore2(t *testing.T) {
	testname := "scoring2"
	s, err := FromSGF(sgfsamples.Scoring2)
	assert.NoError(t, err, testname)

	s.fastForward()

	dead := [][2]int{
		{3, 14},
		{15, 12},
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

	assert.Equal(t, len(blackArea), 72, testname)
	assert.Equal(t, len(whiteArea), 62, testname)

	assert.Equal(t, len(blackDead), 2, testname)
	assert.Equal(t, len(whiteDead), 2, testname)

	assert.Equal(t, current.BlackCaps, 2, testname)
	assert.Equal(t, current.WhiteCaps, 11, testname)

	assert.Equal(t, len(dame), 8, testname)
}

func TestState3(t *testing.T) {
	s, err := FromSGF(sgfsamples.SimpleTwoBranches)
	assert.NoError(t, err, "error should be nil")

	// index 5 is the start of the second branch from the root
	s.gotoIndex(5) //nolint:errcheck

	prefs := s.prefs()
	assert.Equal(t, prefs["0"], 1, "prefs[0]")

	err = s.setPreferred(4)
	assert.NoError(t, err, "error should not be nil")

	prefs = s.prefs()
	assert.Equal(t, prefs["0"], 0, "prefs[0]")

	s.resetPrefs()
	prefs = s.prefs()
	for _, value := range prefs {
		assert.Equal(t, value, 0, "all prefs should be 0")
	}

	s.gotoIndex(6) //nolint:errcheck
	assert.Equal(t, s.locate(), "1,0", "s.Locate")

}

func TestState4(t *testing.T) {
	s, err := FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "error should be nil")

	assert.Equal(t, s.Head().Index, 8, "head index")

	// push a node to the head
	s.PushHead(0, 0, core.Black)
	assert.Equal(t, s.Current().Index, 0, "current index")

	// go to the most recently pushed node
	s.gotoIndex(9) //nolint:errcheck

	// since we're at the head, we will "track" along
	s.PushHead(0, 1, core.White)
	assert.Equal(t, s.Current().Index, 10, "current index")
}

func TestAddStones(t *testing.T) {
	s, err := FromSGF(sgfsamples.Empty)
	assert.NoError(t, err, "error should be nil")

	// add three new moves
	moves := []*core.Stone{
		core.NewStone(9, 9, core.Black),
		core.NewStone(10, 10, core.White),
		core.NewStone(11, 11, core.Black),
	}

	s.AddStones(moves)

	assert.Equal(t, len(s.Current().Down), 1, "len(down)")

	s.fastForward()
	assert.Equal(t, s.Current().Index, 3, "current index")
}

func TestCut(t *testing.T) {
	s, err := FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "error should be nil")

	s.gotoIndex(5) //nolint:errcheck
	s.cut()
	assert.Equal(t, s.Current().Index, 4, "current index")
	assert.Equal(t, len(s.Current().Down), 0, "failed cut")
}

func TestGraft(t *testing.T) {
	s, err := FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "error should be nil")

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
	assert.Equal(t, s.Current().Index, 11, "current index")

	// smart graft
	s.smartGraft(4, moves)
	s.fastForward()
	assert.Equal(t, s.Current().Index, 12, "current index")

	// double-checking Up works
	s.gotoIndex(4) //nolint:errcheck
	s.up()
	s.fastForward()
	assert.Equal(t, s.Current().Index, 8, "current index")

	// check for gotocoord
	s.gotoCoord(3, 3)
	assert.Equal(t, s.Current().Index, 2, "current index")
}

func TestStateJSON(t *testing.T) {
	s, err := FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "error should be nil")

	s.gotoIndex(4) //nolint:errcheck

	j := s.Save()
	assert.Equal(t, j.NextIndex, 9, "next index")
	assert.Equal(t, j.Location, "0,0,0,0", "location")

	tr := s.saveTree(core.FullFrame)
	assert.Equal(t, tr.Depth, 8, "depth")
	assert.Equal(t, tr.Current, 4, "current")
	assert.Equal(t, tr.Up, 0, "up")
	assert.Equal(t, tr.Root, 0, "root")
}

func TestGenerateMarks(t *testing.T) {
	s, err := FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "error should be nil")

	s.right()

	// add a triangle
	coord := &core.Coord{X: 9, Y: 9}
	_, err = NewAddTriangleCommand(coord).Execute(s)
	assert.NoError(t, err, "err should be nil")

	// add a square
	coord = &core.Coord{X: 10, Y: 10}
	_, err = NewAddSquareCommand(coord).Execute(s)
	assert.NoError(t, err, "err should be nil")

	// add a letter
	coord = &core.Coord{X: 11, Y: 11}
	letter := "A"
	_, err = NewAddLetterCommand(coord, letter).Execute(s)
	assert.NoError(t, err, "err should be nil")

	coord = &core.Coord{X: 12, Y: 12}
	number := 1
	_, err = NewAddNumberCommand(coord, number).Execute(s)
	assert.NoError(t, err, "err should be nil")

	marks := s.generateMarks()

	assert.True(t, marks.Current.Equal(&core.Coord{X: 15, Y: 3}), "current")

	assert.Equal(t, len(marks.Triangles), 1, "len(triangles)")
	assert.True(t, marks.Triangles[0].Equal(&core.Coord{X: 9, Y: 9}), "triangles")

	assert.Equal(t, len(marks.Squares), 1, "len(squares)")
	assert.True(t, marks.Squares[0].Equal(&core.Coord{X: 10, Y: 10}), "squares")

	assert.Equal(t, len(marks.Labels), 2, "len(labels)")

	assert.True(t, marks.Labels[0].Coord.Equal(&core.Coord{X: 11, Y: 11}), "letter coord")
	assert.True(t, marks.Labels[1].Coord.Equal(&core.Coord{X: 12, Y: 12}), "number coord")

	assert.Equal(t, marks.Labels[0].Text, "A", "letter")
	assert.Equal(t, marks.Labels[1].Text, "1", "number")

}
