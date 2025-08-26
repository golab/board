/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core_test

import (
	"fmt"
	"github.com/jarednogo/board/backend/core"
	"testing"
)

var oppTests = []struct {
	input  backend.Color
	output backend.Color
}{
	{backend.NoColor, backend.NoColor},
	{backend.Black, backend.White},
	{backend.White, backend.Black},
}

func TestOpposite(t *testing.T) {
	for i, tt := range oppTests {
		t.Run(fmt.Sprintf("opp%d", i), func(t *testing.T) {
			opp := backend.Opposite(tt.input)
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
			coord := backend.LettersToCoord(tt.input)
			if coord.X != tt.x || coord.Y != tt.y {
				t.Errorf("input coord %s mapped to %v, expected: (%d, %d)", tt.input, coord, tt.x, tt.y)
			}
		})
	}
}

func TestCoordFail(t *testing.T) {
	if backend.LettersToCoord("") != nil {
		t.Errorf("empty string should not convert to coords")
	}

	if backend.LettersToCoord("a") != nil {
		t.Errorf("length one string should not convert to coords")
	}

	if backend.LettersToCoord("abc") != nil {
		t.Errorf("length three string should not convert to coords")
	}
}

func TestBoard1(t *testing.T) {
	b := backend.NewBoard(19)
	c := &backend.Coord{10, 10}
	b.Set(c, backend.Black)
	if b.Get(c) != backend.Black {
		t.Errorf("error with set and get")
	}
}

func TestBoard2(t *testing.T) {
	b := backend.NewBoard(19)
	b.Move(&backend.Coord{10, 10}, backend.Black)
	b.Move(&backend.Coord{11, 11}, backend.Black)
	b.Move(&backend.Coord{9, 11}, backend.Black)
	b.Move(&backend.Coord{10, 11}, backend.White)
	b.Move(&backend.Coord{10, 12}, backend.Black)
	g := b.Get(&backend.Coord{10, 11})
	if g != backend.NoColor {
		t.Errorf("error with capture, expected %v, got: %v", backend.NoColor, g)
	}
}

func TestBoard3(t *testing.T) {
	b := backend.NewBoard(19)
	b.Move(&backend.Coord{0, 1}, backend.Black)
	b.Move(&backend.Coord{1, 0}, backend.Black)
	b.Move(&backend.Coord{2, 0}, backend.White)
	b.Move(&backend.Coord{1, 1}, backend.White)
	b.Move(&backend.Coord{0, 0}, backend.White)
	g := b.Get(&backend.Coord{0, 0})
	if g != backend.White {
		t.Errorf("error with capture, expected %v, got: %v", backend.White, g)
	}
}

func TestBoard4(t *testing.T) {
	b := backend.NewBoard(19)
	cs := backend.NewCoordSet()
	for i := 0; i < 5; i++ {
		c := &backend.Coord{10, i}
		cs.Add(c)
	}
	b.SetMany(cs.List(), backend.Black)

	gps := b.Groups()

	if len(gps) != 1 {
		t.Errorf("error in groups")
	}

	if !gps[0].Coords.Equal(cs) {
		t.Errorf("error in group coords, expected %v, got: %v", cs, gps[0].Coords)
	}
}
