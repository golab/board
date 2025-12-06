/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package coord_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/internal/require"
	"github.com/jarednogo/board/pkg/core/board"
	"github.com/jarednogo/board/pkg/core/color"
	"github.com/jarednogo/board/pkg/core/coord"
)

func TestCoordSetHas(t *testing.T) {
	c := coord.NewCoordSet()
	c.Add(coord.NewCoord(9, 9))
	assert.True(t, c.Has(coord.NewCoord(9, 9)))
}

func TestCoordSetAdd(t *testing.T) {
	c := coord.NewCoordSet()
	c.Add(coord.NewCoord(9, 9))
	assert.Equal(t, len(c.List()), 1)
}

func TestDiff1(t *testing.T) {
	b := board.NewBoard(19)
	b.Move(coord.NewCoord(10, 10), color.Black)
	b.Move(coord.NewCoord(11, 11), color.Black)
	b.Move(coord.NewCoord(9, 11), color.Black)

	mv1 := coord.NewCoord(10, 11)
	b.Move(mv1, color.White)

	// capture
	mv2 := coord.NewCoord(10, 12)
	diff := b.Move(mv2, color.Black)

	e := diff.Copy()

	if len(e.Add) != 1 {
		t.Errorf("incorrect addition: %v", e.Add)
	}

	if len(e.Remove) != 1 {
		t.Errorf("incorrect addition: %v", e.Remove)
	}

	add := e.Add[0]
	if add.Color != color.Black {
		t.Errorf("wrong color: %v", add.Color)
	}
	found := false
	for _, c := range add.Coords {
		if c.Equal(mv2) {
			found = true
		}
	}

	if !found {
		t.Errorf("didn't find: %v", mv1)
	}

	remove := e.Remove[0]
	if remove.Color != color.White {
		t.Errorf("wrong color: %v", add.Color)
	}
	found = false
	for _, c := range remove.Coords {
		if c.Equal(mv1) {
			found = true
		}
	}

	if !found {
		t.Errorf("didn't find: %v", mv2)
	}
}

func TestInterface1(t *testing.T) {
	data := []byte("[3, 17]")
	var ifc any
	err := json.Unmarshal(data, &ifc)
	if err != nil {
		t.Error(err)
	}

	c, err := coord.FromInterface(ifc)
	if err != nil {
		t.Error(err)
	}
	if c.X != 3 || c.Y != 17 {
		t.Errorf("wrong coord, expected (3, 17), got %v", c)
	}
}

func TestCoordSetRemove(t *testing.T) {
	cs := coord.NewCoordSet()
	cs.Add(coord.NewCoord(0, 1))
	cs.Add(coord.NewCoord(0, 2))
	cs.Add(coord.NewCoord(0, 3))
	cs.Remove(coord.NewCoord(0, 2))
	if cs.Has(coord.NewCoord(0, 2)) {
		t.Errorf("error removing coord")
	}
}

func TestCoordSetRemoveAll(t *testing.T) {
	cs := coord.NewCoordSet()
	ds := coord.NewCoordSet()
	for x := 0; x < 10; x++ {
		cs.Add(coord.NewCoord(0, x))
		if x < 5 {
			ds.Add(coord.NewCoord(0, x))
		}
	}

	cs.RemoveAll(ds)
	for x := 0; x < 5; x++ {
		if cs.Has(coord.NewCoord(0, x)) {
			t.Errorf("error removing multiple coords")
		}
	}
}

var alphaTests = []struct {
	input    string
	output   *coord.Coord
	hasError bool
}{
	{"a1", coord.NewCoord(0, 18), false},
	{"j1", coord.NewCoord(8, 18), false},
	{"i1", nil, true},
}

func TestAlphanumericToCoord(t *testing.T) {
	for i, tt := range alphaTests {
		t.Run(fmt.Sprintf("alpha%d", i), func(t *testing.T) {
			coord, err := coord.FromAlphanumeric(tt.input, 19)
			if tt.output != nil && !coord.Equal(tt.output) {
				t.Errorf(
					"wrong output in TestAlphanumeric: %v (expected %v)",
					coord,
					tt.output)
			}
			if err != nil && !tt.hasError {
				t.Errorf("unexpected error")
			}

			if err == nil && tt.hasError {
				t.Errorf("expected error")
			}
		})
	}
}

func TestAddAll(t *testing.T) {
	cs := coord.NewCoordSet()
	ds := coord.NewCoordSet()
	ds.Add(coord.NewCoord(0, 0))
	ds.Add(coord.NewCoord(1, 1))
	ds.Add(coord.NewCoord(2, 2))
	cs.AddAll(ds)
	assert.True(t, cs.Has(coord.NewCoord(0, 0)))
	assert.True(t, cs.Has(coord.NewCoord(1, 1)))
	assert.True(t, cs.Has(coord.NewCoord(2, 2)))
}

func TestIntersect(t *testing.T) {
	cs := coord.NewCoordSet()
	ds := coord.NewCoordSet()
	cs.Add(coord.NewCoord(0, 0))
	cs.Add(coord.NewCoord(2, 2))
	ds.Add(coord.NewCoord(0, 0))
	ds.Add(coord.NewCoord(1, 1))

	i := cs.Intersect(ds)
	assert.Equal(t, len(i.List()), 1)
	assert.True(t, i.Has(coord.NewCoord(0, 0)))
}

func TestString(t *testing.T) {
	cs := coord.NewCoordSet()
	cs.Add(coord.NewCoord(0, 0))
	cs.Add(coord.NewCoord(2, 2))
	assert.Equal(t, len(cs.String()), 8)
}

func TestSubset(t *testing.T) {
	cs := coord.NewCoordSet()
	ds := coord.NewCoordSet()
	ds.Add(coord.NewCoord(0, 0))
	ds.Add(coord.NewCoord(1, 1))
	ds.Add(coord.NewCoord(2, 2))
	cs.Add(coord.NewCoord(0, 0))
	assert.True(t, cs.IsSubsetOf(ds))
	assert.False(t, ds.IsSubsetOf(cs))
}

func TestEqual(t *testing.T) {
	cs := coord.NewCoordSet()
	ds := coord.NewCoordSet()
	ds.Add(coord.NewCoord(0, 0))
	ds.Add(coord.NewCoord(1, 1))
	ds.Add(coord.NewCoord(2, 2))
	cs.AddAll(ds)
	assert.True(t, cs.Equal(ds))
}

func TestCopyStoneSet(t *testing.T) {
	cs := coord.NewCoordSet()
	cs.Add(coord.NewCoord(0, 0))
	s := coord.NewStoneSet(cs, color.Black)
	r := s.Copy()
	assert.Equal(t, len(r.Coords), 1)
}

func TestCopyStoneSetNil(t *testing.T) {
	var s *coord.StoneSet
	r := s.Copy()
	assert.Zero(t, r)
}

func TestEqualStoneSet(t *testing.T) {
	cs := coord.NewCoordSet()
	ds := coord.NewCoordSet()

	cs.Add(coord.NewCoord(0, 0))
	ds.Add(coord.NewCoord(0, 0))

	s := coord.NewStoneSet(cs, color.Black)
	r := coord.NewStoneSet(ds, color.Black)
	assert.True(t, r.Equal(s))
}

func TestStoneSetString(t *testing.T) {
	cs := coord.NewCoordSet()
	cs.Add(coord.NewCoord(0, 0))
	s := coord.NewStoneSet(cs, color.Black)
	assert.Equal(t, len(s.String()), 12)
}

func TestNewCoord(t *testing.T) {
	assert.Zero(t, coord.NewCoord(0, 1000))
	assert.Zero(t, coord.NewCoord(1000, 0))
	assert.Zero(t, coord.NewCoord(-1, 0))
	assert.Zero(t, coord.NewCoord(0, -1))
}

func TestCoordString(t *testing.T) {
	c := coord.NewCoord(3, 3)
	assert.Equal(t, c.String(), "(3, 3)")
}

func TestToLetters(t *testing.T) {
	c := coord.NewCoord(3, 3)
	assert.Equal(t, c.ToLetters(), "dd")
}

func TestCoordEqual1(t *testing.T) {
	var c *coord.Coord
	assert.True(t, c.Equal(c))
}

func TestCoordEqual2(t *testing.T) {
	var c *coord.Coord
	d := coord.NewCoord(1, 1)
	assert.False(t, c.Equal(d))
}

func TestCoordCopy(t *testing.T) {
	var c *coord.Coord
	assert.Zero(t, c.Copy())
}

func TestFromLetters1(t *testing.T) {
	c := coord.FromLetters("abc")
	assert.Zero(t, c)
}

func TestFromLetters2(t *testing.T) {
	c := coord.FromLetters("cc")
	assert.True(t, c.Equal(coord.NewCoord(2, 2)))
}

func TestFromInterface(t *testing.T) {
	_, err := coord.FromInterface(0)
	assert.NotNil(t, err)
}

func TestNewStone(t *testing.T) {
	s := coord.NewStone(3, 3, color.Black)
	assert.True(t, s.Coord.Equal(coord.NewCoord(3, 3)))
	assert.Equal(t, s.Color, color.Black)
}

func TestFromAlphanumeric1(t *testing.T) {
	_, err := coord.FromAlphanumeric("x", 19)
	assert.NotNil(t, err)
}

func TestFromAlphanumeric2(t *testing.T) {
	_, err := coord.FromAlphanumeric("cy", 19)
	assert.NotNil(t, err)
}

func TestFromAlphanumeric3(t *testing.T) {
	_, err := coord.FromAlphanumeric("c0", 19)
	assert.NotNil(t, err)
}

func TestFromAlphanumeric4(t *testing.T) {
	_, err := coord.FromAlphanumeric("c20", 19)
	assert.NotNil(t, err)
}

func TestDiff(t *testing.T) {
	cs := coord.NewCoordSet()
	ds := coord.NewCoordSet()

	cs.Add(coord.NewCoord(0, 0))
	ds.Add(coord.NewCoord(1, 1))

	s1 := []*coord.StoneSet{coord.NewStoneSet(cs, color.Black)}
	s2 := []*coord.StoneSet{coord.NewStoneSet(ds, color.White)}

	diff := coord.NewDiff(s1, s2)
	idiff := diff.Invert()

	require.Equal(t, len(idiff.Add), 1)
	require.Equal(t, len(idiff.Remove), 1)
	add := idiff.Add[0]
	remove := idiff.Remove[0]

	require.Equal(t, len(add.Coords), 1)
	require.Equal(t, len(remove.Coords), 1)

	a := add.Coords[0]
	r := remove.Coords[0]

	assert.True(t, a.Equal(coord.NewCoord(1, 1)))
	assert.True(t, r.Equal(coord.NewCoord(0, 0)))
}

func TestInvertNil(t *testing.T) {
	var d *coord.Diff
	assert.Zero(t, d.Invert())
}

func TestDiffNil(t *testing.T) {
	var d *coord.Diff
	assert.Zero(t, d.Copy())
}
