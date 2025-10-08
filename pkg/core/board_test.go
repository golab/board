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

func TestBoard1(t *testing.T) {
	b := core.NewBoard(19)
	c := &core.Coord{10, 10}
	b.Set(c, core.Black)
	if b.Get(c) != core.Black {
		t.Errorf("error with set and get")
	}
}

func TestBoard2(t *testing.T) {
	b := core.NewBoard(19)
	b.Move(&core.Coord{10, 10}, core.Black)
	b.Move(&core.Coord{11, 11}, core.Black)
	b.Move(&core.Coord{9, 11}, core.Black)
	b.Move(&core.Coord{10, 11}, core.White)
	b.Move(&core.Coord{10, 12}, core.Black)
	g := b.Get(&core.Coord{10, 11})
	if g != core.NoColor {
		t.Errorf("error with capture, expected %v, got: %v", core.NoColor, g)
	}
}

func TestBoard3(t *testing.T) {
	b := core.NewBoard(19)
	b.Move(&core.Coord{0, 1}, core.Black)
	b.Move(&core.Coord{1, 0}, core.Black)
	b.Move(&core.Coord{2, 0}, core.White)
	b.Move(&core.Coord{1, 1}, core.White)
	b.Move(&core.Coord{0, 0}, core.White)
	g := b.Get(&core.Coord{0, 0})
	if g != core.White {
		t.Errorf("error with capture, expected %v, got: %v", core.White, g)
	}
}

func TestBoard4(t *testing.T) {
	b := core.NewBoard(19)
	cs := core.NewCoordSet()
	for i := 0; i < 5; i++ {
		c := &core.Coord{10, i}
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
