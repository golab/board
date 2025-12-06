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
