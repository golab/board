/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core_test

import (
	"encoding/json"
	"github.com/jarednogo/board/backend/core"
	"testing"
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
	var ifc interface{}
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
