/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jarednogo/board/pkg/core"
)

func TestDiff1(t *testing.T) {
	b := core.NewBoard(19)
	b.Move(&core.Coord{10, 10}, core.Black)
	b.Move(&core.Coord{11, 11}, core.Black)
	b.Move(&core.Coord{9, 11}, core.Black)

	mv1 := &core.Coord{10, 11}
	b.Move(mv1, core.White)

	// capture
	mv2 := &core.Coord{10, 12}
	diff := b.Move(mv2, core.Black)

	e := diff.Copy()

	if len(e.Add) != 1 {
		t.Errorf("incorrect addition: %v", e.Add)
	}

	if len(e.Remove) != 1 {
		t.Errorf("incorrect addition: %v", e.Remove)
	}

	add := e.Add[0]
	if add.Color != core.Black {
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
	if remove.Color != core.White {
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

	c, err := core.InterfaceToCoord(ifc)
	if err != nil {
		t.Error(err)
	}
	if c.X != 3 || c.Y != 17 {
		t.Errorf("wrong coord, expected (3, 17), got %v", c)
	}
}

func TestCoordSetRemove(t *testing.T) {
	cs := core.NewCoordSet()
	cs.Add(&core.Coord{0, 1})
	cs.Add(&core.Coord{0, 2})
	cs.Add(&core.Coord{0, 3})
	cs.Remove(&core.Coord{0, 2})
	if cs.Has(&core.Coord{0, 2}) {
		t.Errorf("error removing coord")
	}
}

func TestCoordSetRemoveAll(t *testing.T) {
	cs := core.NewCoordSet()
	ds := core.NewCoordSet()
	for x := 0; x < 10; x++ {
		cs.Add(&core.Coord{0, x})
		if x < 5 {
			ds.Add(&core.Coord{0, x})
		}
	}

	cs.RemoveAll(ds)
	for x := 0; x < 5; x++ {
		if cs.Has(&core.Coord{0, x}) {
			t.Errorf("error removing multiple coords")
		}
	}
}

var alphaTests = []struct {
	input    string
	output   *core.Coord
	hasError bool
}{
	{"a1", &core.Coord{0, 18}, false},
	{"j1", &core.Coord{8, 18}, false},
	{"i1", nil, true},
}

func TestAlphanumericToCoord(t *testing.T) {
	for i, tt := range alphaTests {
		t.Run(fmt.Sprintf("alpha%d", i), func(t *testing.T) {
			coord, err := core.AlphanumericToCoord(tt.input, 19)
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

var sanitizeTests = []struct {
	input  string
	output string
}{
	{"  AAA  ", "AAA"},
	{"http://www.example.com?foo=bar&baz=bot", "httpwwwexamplecomfoobarbazbot"},
}

func TestSanitize(t *testing.T) {
	for i, tt := range sanitizeTests {
		t.Run(fmt.Sprintf("sanitize%d", i), func(t *testing.T) {
			output := core.Sanitize(tt.input)
			if tt.output != output {
				t.Errorf("error in sanitize: got %s (expected %s)", output, tt.output)
			}
		})
	}
}
