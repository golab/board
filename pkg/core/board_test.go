/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core_test

import (
	"fmt"
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/pkg/core"
)

var oppTests = []struct {
	input  core.Color
	output core.Color
}{
	{core.NoColor, core.NoColor},
	{core.Black, core.White},
	{core.White, core.Black},
}

func TestOpposite(t *testing.T) {
	for i, tt := range oppTests {
		t.Run(fmt.Sprintf("opp%d", i), func(t *testing.T) {
			opp := core.Opposite(tt.input)
			if opp != tt.output {
				t.Errorf("error in TestOpposite")
			}
		})
	}
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
			coord := core.LettersToCoord(tt.input)
			if coord.X != tt.x || coord.Y != tt.y {
				t.Errorf("input coord %s mapped to %v, expected: (%d, %d)", tt.input, coord, tt.x, tt.y)
			}
		})
	}
}

func TestCoordFail(t *testing.T) {
	if core.LettersToCoord("") != nil {
		t.Errorf("empty string should not convert to coords")
	}

	if core.LettersToCoord("a") != nil {
		t.Errorf("length one string should not convert to coords")
	}

	if core.LettersToCoord("abc") != nil {
		t.Errorf("length three string should not convert to coords")
	}
}

func TestNeighbors(t *testing.T) {
	b := core.NewBoard(19)
	start := core.NewCoord(9, 9)
	nbs := b.Neighbors(start)
	assert.Equal(t, len(nbs.List()), 4)
}

func TestFindGroup(t *testing.T) {
	b := core.NewBoard(19)
	b.Move(core.NewCoord(10, 10), core.Black)
	b.Move(core.NewCoord(11, 11), core.Black)
	b.Move(core.NewCoord(9, 11), core.Black)
	b.Move(core.NewCoord(10, 11), core.Black)
	b.Move(core.NewCoord(10, 12), core.Black)
	g := b.FindGroup(core.NewCoord(10, 12))
	assert.Equal(t, g.Color, core.Black)
	assert.Equal(t, len(g.Coords.List()), 5)
}

func TestLibs(t *testing.T) {
	b := core.NewBoard(19)
	b.Set(core.NewCoord(10, 10), core.Black)
	b.Set(core.NewCoord(11, 11), core.Black)
	b.Set(core.NewCoord(9, 11), core.Black)
	b.Set(core.NewCoord(10, 12), core.Black)
	b.Set(core.NewCoord(10, 11), core.White)
	g := b.FindGroup(core.NewCoord(10, 11))
	assert.Equal(t, len(g.Libs.List()), 0)
}

func TestWouldKill(t *testing.T) {
	b := core.NewBoard(19)
	b.Move(core.NewCoord(10, 10), core.Black)
	b.Move(core.NewCoord(11, 11), core.Black)
	b.Move(core.NewCoord(9, 11), core.Black)
	b.Move(core.NewCoord(10, 11), core.White)
	start := core.NewCoord(10, 12)
	s := b.WouldKill(start, core.White)
	t.Logf("%v", s)
}

func TestLegal(t *testing.T) {
	b := core.NewBoard(19)
	b.Set(core.NewCoord(10, 10), core.Black)
	b.Set(core.NewCoord(11, 11), core.Black)
	b.Set(core.NewCoord(9, 11), core.Black)
	legal := b.Legal(core.NewCoord(10, 11), core.White)
	assert.True(t, legal)
	b.Set(core.NewCoord(10, 12), core.Black)
	legal = b.Legal(core.NewCoord(10, 11), core.White)
	assert.False(t, legal)
}

func TestBoard1(t *testing.T) {
	b := core.NewBoard(19)
	c := core.NewCoord(10, 10)
	b.Set(c, core.Black)
	if b.Get(c) != core.Black {
		t.Errorf("error with set and get")
	}
}

func TestBoard2(t *testing.T) {
	b := core.NewBoard(19)
	b.Move(core.NewCoord(10, 10), core.Black)
	b.Move(core.NewCoord(11, 11), core.Black)
	b.Move(core.NewCoord(9, 11), core.Black)
	b.Move(core.NewCoord(10, 11), core.White)
	b.Move(core.NewCoord(10, 12), core.Black)
	g := b.Get(core.NewCoord(10, 11))
	if g != core.NoColor {
		t.Errorf("error with capture, expected %v, got: %v", core.NoColor, g)
	}
}

func TestBoard3(t *testing.T) {
	b := core.NewBoard(19)
	b.Move(core.NewCoord(0, 1), core.Black)
	b.Move(core.NewCoord(1, 0), core.Black)
	b.Move(core.NewCoord(2, 0), core.White)
	b.Move(core.NewCoord(1, 1), core.White)
	b.Move(core.NewCoord(0, 0), core.White)
	g := b.Get(core.NewCoord(0, 0))
	if g != core.White {
		t.Errorf("error with capture, expected %v, got: %v", core.White, g)
	}
}

func TestBoard4(t *testing.T) {
	b := core.NewBoard(19)
	cs := core.NewCoordSet()
	for i := 0; i < 5; i++ {
		c := core.NewCoord(10, i)
		cs.Add(c)
	}
	b.SetMany(cs.List(), core.Black)

	gps := b.Groups()

	if len(gps) != 1 {
		t.Errorf("error in groups")
	}

	if !gps[0].Coords.Equal(cs) {
		t.Errorf("error in group coords, expected %v, got: %v", cs, gps[0].Coords)
	}
}

func TestCurrentDiff(t *testing.T) {
	b := core.NewBoard(19)
	c1 := core.NewCoord(0, 1)
	c2 := core.NewCoord(1, 0)
	c3 := core.NewCoord(2, 0)
	c4 := core.NewCoord(1, 1)
	c5 := core.NewCoord(0, 0)
	b.Move(c1, core.Black)
	b.Move(c2, core.Black)
	b.Move(c3, core.White)
	b.Move(c4, core.White)
	b.Move(c5, core.White)

	diff := b.CurrentDiff()

	c := core.NewCoordSet()
	c.Add(c1)
	addBlack := core.NewStoneSet(c, core.Black)

	d := core.NewCoordSet()
	d.Add(c3)
	d.Add(c4)
	d.Add(c5)
	addWhite := core.NewStoneSet(d, core.White)

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
	b := core.NewBoard(19)
	for x := 0; x < 4; x++ {
		b.Move(core.NewCoord(x, 3-x), core.Black)
	}
	cs, tp := b.FindArea(core.NewCoord(0, 0), core.NewCoordSet())
	if len(cs.List()) != 6 {
		t.Errorf("error finding empty area: got %d (expected 6)", len(cs))
	}
	if tp != core.BlackPoint {
		t.Errorf("%v", tp)
	}
}

func TestScore(t *testing.T) {
	b := core.NewBoard(9)
	for x := 0; x < 9; x++ {
		b.Move(core.NewCoord(x, 4), core.Black)
		b.Move(core.NewCoord(x, 5), core.White)
	}

	b.Move(core.NewCoord(7, 3), core.Black)
	b.Move(core.NewCoord(8, 3), core.Black)

	b.Move(core.NewCoord(7, 6), core.White)
	b.Move(core.NewCoord(8, 6), core.White)

	// dead stones
	b.Move(core.NewCoord(0, 0), core.White)
	b.Move(core.NewCoord(8, 8), core.Black)

	markedDead := core.NewCoordSet()
	markedDead.Add(core.NewCoord(0, 0))
	markedDead.Add(core.NewCoord(8, 8))

	markedDame := core.NewCoordSet()
	markedDame.Add(core.NewCoord(8, 4))
	markedDame.Add(core.NewCoord(8, 5))

	blackArea, whiteArea, blackDead, whiteDead, dame := b.Score(markedDead, markedDame)
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

func TestBoardCopy(t *testing.T) {
	b := core.NewBoard(9)
	for i := 0; i < 9; i++ {
		b.Move(core.NewCoord(i, i), core.Black)
	}

	c := b.Copy()

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			coord := core.NewCoord(i, j)
			assert.Equal(t, c.Get(coord), b.Get(coord))
		}
	}
}

func TestBoardString(t *testing.T) {
	b := core.NewBoard(9)
	for i := 0; i < 9; i++ {
		b.Move(core.NewCoord(i, i), core.Black)
	}

	assert.Equal(t, len(b.String()), 171)
}
