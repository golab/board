/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package board_test

import (
	"fmt"
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/pkg/core/board"
	"github.com/golab/board/pkg/core/color"
	"github.com/golab/board/pkg/core/coord"
)

func TestString(t *testing.T) {
	b := board.NewBoard(9)
	b.Set(coord.NewCoord(1, 0), color.Black)
	b.Set(coord.NewCoord(2, 0), color.Black)
	b.Set(coord.NewCoord(3, 0), color.Black)
	g := b.FindGroup(coord.NewCoord(3, 0))
	assert.Equal(t, g.Color, color.Black)
	assert.Equal(t, len(g.String()), 16)
}

func TestEmptyGroup(t *testing.T) {
	g := board.NewGroup(nil, nil, color.Empty)
	assert.Equal(t, g.Color, color.Empty)
	assert.NotNil(t, g.Coords)
	assert.NotNil(t, g.Libs)
}

func TestGetEmpty(t *testing.T) {
	b := board.NewBoard(9)
	assert.Equal(t, b.Get(coord.NewCoord(9, 0)), color.Empty)
	assert.Equal(t, b.Get(coord.NewCoord(0, 9)), color.Empty)
	assert.Equal(t, b.Get(coord.NewCoord(-1, 0)), color.Empty)
	assert.Equal(t, b.Get(coord.NewCoord(0, -1)), color.Empty)
}

func TestFindGroupEmpty(t *testing.T) {
	b := board.NewBoard(9)
	g := b.FindGroup(coord.NewCoord(0, 0))
	assert.Equal(t, g.Color, color.Empty)
	assert.NotNil(t, g.Coords)
	assert.NotNil(t, g.Libs)
}

func TestExistingLegal(t *testing.T) {
	b := board.NewBoard(9)
	b.Set(coord.NewCoord(1, 0), color.Black)
	legal := b.Legal(coord.NewCoord(1, 0), color.Black)
	assert.False(t, legal)
}

func TestLegalMove(t *testing.T) {
	b := board.NewBoard(9)
	b.Set(coord.NewCoord(1, 0), color.Black)
	diff := b.Move(coord.NewCoord(1, 0), color.Black)
	assert.Zero(t, diff)
}

func TestFindAreaNonempty(t *testing.T) {
	b := board.NewBoard(9)
	b.Set(coord.NewCoord(1, 0), color.Black)
	_, typ := b.FindArea(coord.NewCoord(1, 0), nil)
	assert.Equal(t, typ, board.NotCovered)
}

var coordTests = []struct {
	input string
	x     int
	y     int
}{
	{"ah", 0, 7},
	{"js", 9, 18},
}

func TestCoord(t *testing.T) {
	for i, tt := range coordTests {
		t.Run(fmt.Sprintf("coord%d", i), func(t *testing.T) {
			coord := coord.FromLetters(tt.input)
			if coord.X != tt.x || coord.Y != tt.y {
				t.Errorf("input coord %s mapped to %v, expected: (%d, %d)", tt.input, coord, tt.x, tt.y)
			}
		})
	}
}

func TestCoordFail(t *testing.T) {
	if coord.FromLetters("") != nil {
		t.Errorf("empty string should not convert to coords")
	}

	if coord.FromLetters("a") != nil {
		t.Errorf("length one string should not convert to coords")
	}

	if coord.FromLetters("abc") != nil {
		t.Errorf("length three string should not convert to coords")
	}
}

func TestNeighbors(t *testing.T) {
	b := board.NewBoard(19)
	start := coord.NewCoord(9, 9)
	nbs := b.Neighbors(start)
	assert.Equal(t, len(nbs.List()), 4)
}

func TestFindGroup(t *testing.T) {
	b := board.NewBoard(19)
	b.Move(coord.NewCoord(10, 10), color.Black)
	b.Move(coord.NewCoord(11, 11), color.Black)
	b.Move(coord.NewCoord(9, 11), color.Black)
	b.Move(coord.NewCoord(10, 11), color.Black)
	b.Move(coord.NewCoord(10, 12), color.Black)
	g := b.FindGroup(coord.NewCoord(10, 12))
	assert.Equal(t, g.Color, color.Black)
	assert.Equal(t, len(g.Coords.List()), 5)
}

func TestLibs(t *testing.T) {
	b := board.NewBoard(19)
	b.Set(coord.NewCoord(10, 10), color.Black)
	b.Set(coord.NewCoord(11, 11), color.Black)
	b.Set(coord.NewCoord(9, 11), color.Black)
	b.Set(coord.NewCoord(10, 12), color.Black)
	b.Set(coord.NewCoord(10, 11), color.White)
	g := b.FindGroup(coord.NewCoord(10, 11))
	assert.Equal(t, len(g.Libs.List()), 0)
}

func TestWouldKill(t *testing.T) {
	b := board.NewBoard(19)
	b.Move(coord.NewCoord(10, 10), color.Black)
	b.Move(coord.NewCoord(11, 11), color.Black)
	b.Move(coord.NewCoord(9, 11), color.Black)
	b.Move(coord.NewCoord(10, 11), color.White)
	start := coord.NewCoord(10, 12)
	s := b.WouldKill(start, color.White)
	require.Equal(t, len(s.Coords), 1)
	assert.Equal(t, s.Color, color.White)
	assert.True(t, s.Coords[0].Equal(coord.NewCoord(10, 11)))
}

func TestLegal(t *testing.T) {
	b := board.NewBoard(19)
	b.Set(coord.NewCoord(10, 10), color.Black)
	b.Set(coord.NewCoord(11, 11), color.Black)
	b.Set(coord.NewCoord(9, 11), color.Black)
	legal := b.Legal(coord.NewCoord(10, 11), color.White)
	assert.True(t, legal)
	b.Set(coord.NewCoord(10, 12), color.Black)
	legal = b.Legal(coord.NewCoord(10, 11), color.White)
	assert.False(t, legal)
}

func TestBoard1(t *testing.T) {
	b := board.NewBoard(19)
	c := coord.NewCoord(10, 10)
	b.Set(c, color.Black)
	if b.Get(c) != color.Black {
		t.Errorf("error with set and get")
	}
}

func TestBoard2(t *testing.T) {
	b := board.NewBoard(19)
	b.Move(coord.NewCoord(10, 10), color.Black)
	b.Move(coord.NewCoord(11, 11), color.Black)
	b.Move(coord.NewCoord(9, 11), color.Black)
	b.Move(coord.NewCoord(10, 11), color.White)
	b.Move(coord.NewCoord(10, 12), color.Black)
	g := b.Get(coord.NewCoord(10, 11))
	if g != color.Empty {
		t.Errorf("error with capture, expected %v, got: %v", color.Empty, g)
	}
}

func TestBoard3(t *testing.T) {
	b := board.NewBoard(19)
	b.Move(coord.NewCoord(0, 1), color.Black)
	b.Move(coord.NewCoord(1, 0), color.Black)
	b.Move(coord.NewCoord(2, 0), color.White)
	b.Move(coord.NewCoord(1, 1), color.White)
	b.Move(coord.NewCoord(0, 0), color.White)
	g := b.Get(coord.NewCoord(0, 0))
	if g != color.White {
		t.Errorf("error with capture, expected %v, got: %v", color.White, g)
	}
}

func TestBoard4(t *testing.T) {
	b := board.NewBoard(19)
	cs := coord.NewCoordSet()
	for i := 0; i < 5; i++ {
		c := coord.NewCoord(10, i)
		cs.Add(c)
	}
	b.SetMany(cs.List(), color.Black)

	gps := b.Groups()

	if len(gps) != 1 {
		t.Errorf("error in groups")
	}

	if !gps[0].Coords.Equal(cs) {
		t.Errorf("error in group coords, expected %v, got: %v", cs, gps[0].Coords)
	}
}

func TestCurrentDiff(t *testing.T) {
	b := board.NewBoard(19)
	c1 := coord.NewCoord(0, 1)
	c2 := coord.NewCoord(1, 0)
	c3 := coord.NewCoord(2, 0)
	c4 := coord.NewCoord(1, 1)
	c5 := coord.NewCoord(0, 0)
	b.Move(c1, color.Black)
	b.Move(c2, color.Black)
	b.Move(c3, color.White)
	b.Move(c4, color.White)
	b.Move(c5, color.White)

	diff := b.CurrentDiff()

	c := coord.NewCoordSet()
	c.Add(c1)
	addBlack := coord.NewStoneSet(c, color.Black)

	d := coord.NewCoordSet()
	d.Add(c3)
	d.Add(c4)
	d.Add(c5)
	addWhite := coord.NewStoneSet(d, color.White)

	if !diff.Add[0].Equal(addBlack) {
		t.Errorf("addBlack wrong")
	}

	if !diff.Add[1].Equal(addWhite) {
		t.Errorf("addWhite wrong")
	}

	if diff.Remove != nil {
		t.Errorf("diff.Remove wrong")
	}
}

func TestFindArea(t *testing.T) {
	b := board.NewBoard(19)
	for x := 0; x < 4; x++ {
		b.Move(coord.NewCoord(x, 3-x), color.Black)
	}
	cs, tp := b.FindArea(coord.NewCoord(0, 0), coord.NewCoordSet())
	if len(cs.List()) != 6 {
		t.Errorf("error finding empty area: got %d (expected 6)", len(cs))
	}
	if tp != board.BlackPoint {
		t.Errorf("%v", tp)
	}
}

func TestScore(t *testing.T) {
	b := board.NewBoard(9)
	for x := 0; x < 9; x++ {
		b.Move(coord.NewCoord(x, 4), color.Black)
		b.Move(coord.NewCoord(x, 5), color.White)
	}

	b.Move(coord.NewCoord(7, 3), color.Black)
	b.Move(coord.NewCoord(8, 3), color.Black)

	b.Move(coord.NewCoord(7, 6), color.White)
	b.Move(coord.NewCoord(8, 6), color.White)

	// dead stones
	b.Move(coord.NewCoord(0, 0), color.White)
	b.Move(coord.NewCoord(8, 8), color.Black)

	markedDead := coord.NewCoordSet()
	markedDead.Add(coord.NewCoord(0, 0))
	markedDead.Add(coord.NewCoord(8, 8))

	markedDame := coord.NewCoordSet()
	markedDame.Add(coord.NewCoord(8, 4))
	markedDame.Add(coord.NewCoord(8, 5))

	scoreResult := b.Score(markedDead, markedDame)
	blackArea := scoreResult.BlackArea
	whiteArea := scoreResult.WhiteArea
	blackDead := scoreResult.BlackDead
	whiteDead := scoreResult.WhiteDead
	dame := scoreResult.Dame

	if len(blackArea) != 34 {
		t.Errorf("blackArea wrong, expected 34 (got %d)", len(blackArea))
	}
	if len(whiteArea) != 25 {
		t.Errorf("whiteArea wrong, expected 25 (got %d)", len(whiteArea))
	}
	if len(blackDead) != 1 {
		t.Errorf("blackDead wrong, expected 0 (got %d)", len(blackDead))
	}
	if len(whiteDead) != 1 {
		t.Errorf("whiteDead wrong, expected 0 (got %d)", len(whiteDead))
	}
	if len(dame) != 2 {
		t.Errorf("dame wrong, expected 0 (got %d)", len(dame))
	}
}

func TestScore2(t *testing.T) {
	b := board.NewBoard(9)
	for x := 0; x < 8; x++ {
		b.Move(coord.NewCoord(x, 4), color.Black)
		b.Move(coord.NewCoord(x, 5), color.White)
	}

	b.Move(coord.NewCoord(7, 3), color.Black)
	b.Move(coord.NewCoord(7, 2), color.Black)
	b.Move(coord.NewCoord(8, 2), color.Black)

	b.Move(coord.NewCoord(8, 4), color.White)

	scoreResult := b.Score(coord.NewCoordSet(), coord.NewCoordSet())
	blackArea := scoreResult.BlackArea
	whiteArea := scoreResult.WhiteArea
	blackDead := scoreResult.BlackDead
	whiteDead := scoreResult.WhiteDead
	dame := scoreResult.Dame

	assert.Equal(t, len(blackArea), 32)
	assert.Equal(t, len(whiteArea), 27)
	assert.Equal(t, len(blackDead), 0)
	assert.Equal(t, len(whiteDead), 0)
	assert.Equal(t, len(dame), 2)
}

func TestBoardCopy(t *testing.T) {
	b := board.NewBoard(9)
	for i := 0; i < 9; i++ {
		b.Move(coord.NewCoord(i, i), color.Black)
	}

	c := b.Copy()

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			coord := coord.NewCoord(i, j)
			assert.Equal(t, c.Get(coord), b.Get(coord))
		}
	}
}

func TestBoardString(t *testing.T) {
	b := board.NewBoard(9)
	for i := 0; i < 9; i++ {
		b.Move(coord.NewCoord(i, i), color.Black)
	}

	assert.Equal(t, len(b.String()), 171)
}

func TestDetectAtariDame1(t *testing.T) {
	b := board.NewBoard(9)
	b.Set(coord.NewCoord(1, 0), color.Black)
	b.Set(coord.NewCoord(3, 0), color.Black)
	b.Set(coord.NewCoord(5, 0), color.Black)
	b.Set(coord.NewCoord(7, 0), color.Black)
	b.Set(coord.NewCoord(2, 1), color.Black)
	b.Set(coord.NewCoord(3, 1), color.Black)
	b.Set(coord.NewCoord(4, 1), color.Black)
	b.Set(coord.NewCoord(5, 1), color.Black)
	b.Set(coord.NewCoord(6, 1), color.Black)
	b.Set(coord.NewCoord(7, 1), color.Black)
	b.Set(coord.NewCoord(0, 1), color.White)
	b.Set(coord.NewCoord(1, 1), color.White)
	b.Set(coord.NewCoord(1, 2), color.White)
	b.Set(coord.NewCoord(2, 2), color.White)
	b.Set(coord.NewCoord(3, 2), color.White)
	b.Set(coord.NewCoord(4, 2), color.White)
	b.Set(coord.NewCoord(5, 2), color.White)
	b.Set(coord.NewCoord(6, 2), color.White)
	b.Set(coord.NewCoord(7, 2), color.White)
	b.Set(coord.NewCoord(8, 2), color.White)
	b.Set(coord.NewCoord(8, 1), color.White)
	b.Set(coord.NewCoord(8, 0), color.White)

	dead := coord.NewCoordSet()
	dame := coord.NewCoordSet()
	dame.Add(coord.NewCoord(0, 0))
	ad := b.DetectAtariDame(dead, dame)

	assert.True(t, ad.Has(coord.NewCoord(2, 0)))
}

func TestDetectAtariDame2(t *testing.T) {
	b := board.NewBoard(9)

	b.Set(coord.NewCoord(4, 0), color.Black)
	b.Set(coord.NewCoord(4, 1), color.Black)
	b.Set(coord.NewCoord(4, 2), color.Black)
	b.Set(coord.NewCoord(4, 3), color.Black)
	b.Set(coord.NewCoord(4, 4), color.Black)
	b.Set(coord.NewCoord(4, 5), color.Black)
	b.Set(coord.NewCoord(4, 6), color.Black)
	b.Set(coord.NewCoord(4, 7), color.Black)
	b.Set(coord.NewCoord(4, 8), color.Black)
	b.Set(coord.NewCoord(2, 8), color.Black)

	b.Set(coord.NewCoord(5, 0), color.White)
	b.Set(coord.NewCoord(5, 1), color.White)
	b.Set(coord.NewCoord(5, 2), color.White)
	b.Set(coord.NewCoord(5, 3), color.White)
	b.Set(coord.NewCoord(5, 4), color.White)
	b.Set(coord.NewCoord(5, 5), color.White)
	b.Set(coord.NewCoord(5, 6), color.White)
	b.Set(coord.NewCoord(5, 7), color.White)
	b.Set(coord.NewCoord(5, 8), color.White)
	b.Set(coord.NewCoord(3, 8), color.White)

	dead := coord.NewCoordSet()
	dead.Add(coord.NewCoord(3, 8))
	dame := coord.NewCoordSet()
	ad := b.DetectAtariDame(dead, dame)
	assert.Equal(t, len(ad), 0)
}

func TestScoreDame(t *testing.T) {
	b := board.NewBoard(5)
	b.Set(coord.NewCoord(1, 0), color.Black)
	b.Set(coord.NewCoord(1, 1), color.Black)
	b.Set(coord.NewCoord(1, 2), color.Black)
	b.Set(coord.NewCoord(1, 3), color.Black)
	b.Set(coord.NewCoord(1, 4), color.Black)

	b.Set(coord.NewCoord(3, 0), color.White)
	b.Set(coord.NewCoord(3, 1), color.White)
	b.Set(coord.NewCoord(3, 2), color.White)
	b.Set(coord.NewCoord(3, 3), color.White)
	b.Set(coord.NewCoord(3, 4), color.White)

	dead := coord.NewCoordSet()
	dame := coord.NewCoordSet()
	dame.Add(coord.NewCoord(2, 2))
	sr := b.Score(dead, dame)
	assert.Equal(t, len(sr.Dame), 5)
}

func TestClear(t *testing.T) {
	b := board.NewBoard(5)
	b.Set(coord.NewCoord(1, 0), color.Black)
	b.Set(coord.NewCoord(1, 1), color.Black)
	b.Set(coord.NewCoord(1, 2), color.Black)
	b.Set(coord.NewCoord(1, 3), color.Black)
	b.Set(coord.NewCoord(1, 4), color.Black)

	b.Clear()
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			c := coord.NewCoord(i, j)
			assert.Equal(t, b.Get(c), color.Empty)
		}
	}
}

func TestApplyDiff(t *testing.T) {
	b := board.NewBoard(5)
	b.Set(coord.NewCoord(1, 0), color.Black)
	b.Set(coord.NewCoord(1, 1), color.Black)
	b.Set(coord.NewCoord(1, 2), color.Black)
	b.Set(coord.NewCoord(1, 3), color.White)
	b.Set(coord.NewCoord(1, 4), color.White)

	c := board.NewBoard(5)
	c.ApplyDiff(nil)
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			crd := coord.NewCoord(i, j)
			assert.Equal(t, c.Get(crd), color.Empty)
		}
	}

	diff := b.CurrentDiff()
	c.ApplyDiff(diff)
	assert.Equal(t, c.Get(coord.NewCoord(1, 0)), color.Black)
	assert.Equal(t, c.Get(coord.NewCoord(1, 1)), color.Black)
	assert.Equal(t, c.Get(coord.NewCoord(1, 2)), color.Black)
	assert.Equal(t, c.Get(coord.NewCoord(1, 3)), color.White)
	assert.Equal(t, c.Get(coord.NewCoord(1, 4)), color.White)
}

func TestSnapback(t *testing.T) {
	b, err := board.FromString(`
.........
BB.......
WB.BB....
WWB.WB...
.WWBWB...
...WBBBBB
...WBWBWB
..W.WWWWW
.........`)
	require.NoError(t, err)

	dead := coord.NewCoordSet()
	dame := coord.NewCoordSet()
	dead.Add(coord.NewCoord(4, 4))
	dead.Add(coord.NewCoord(4, 3))

	sr := b.Score(dead, dame)

	assert.Equal(t, len(sr.BlackArea), 30)
	assert.Equal(t, len(sr.WhiteDead), 2)
	assert.Equal(t, len(sr.Dame), 0)
}
