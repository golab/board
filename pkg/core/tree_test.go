/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core_test

import (
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/internal/sgfsamples"
	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/state"
)

func TestFmap(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleTwoBranches)
	if err != nil {
		t.Error(err)
	}

	core.Fmap(func(n *core.TreeNode) {
		n.Color = core.Opposite(n.Color)
	}, s.Root())

	node := s.Nodes()[1]
	if node.Color != core.White {
		t.Errorf("fmap failed")
	}
}

func TestField1(t *testing.T) {

	s, err := state.FromSGF(sgfsamples.SimpleWithComment)
	if err != nil {
		t.Error(err)
	}

	root := s.Root()
	root.AddField("C", "some comment")

	found := false
	for _, comment := range root.Fields["C"] {
		if comment == "some comment" {
			found = true
		}
	}
	if !found {
		t.Errorf("failed to add comment")
	}

	root.RemoveField("C", "comment1")

	for _, comment := range root.Fields["C"] {
		if comment == "comment1" {
			t.Errorf("failed to remove comment")
		}
	}
}

func TestDepth(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleWithComment)
	if err != nil {
		t.Error(err)
	}

	m := s.Root().MaxDepth()
	if m != 4 {
		t.Errorf("max depth failed, expected 4, got: %d", m)
	}
}

func TestChild(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleWithComment)
	if err != nil {
		t.Error(err)
	}
	coord := &core.Coord{3, 3}
	ind, has := s.Root().HasChild(coord, core.Black)
	if !has {
		t.Errorf("failed to find child")
	}
	if ind != 5 {
		t.Errorf("found child at the wrong index %d", ind)
	}
}

func TestTrunkNum(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleTwoBranches)
	if err != nil {
		t.Error(err)
	}

	if tn := s.Root().TrunkNum(5); tn != -1 {
		t.Errorf("error: expected -1 (got %d)", tn)
	}

	if tn := s.Root().TrunkNum(4); tn != 4 {
		t.Errorf("error: expected 4 (got %d)", tn)
	}
}

func TestNodeCopy(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleTwoBranches)
	if err != nil {
		t.Error(err)
	}

	root := s.Root()
	c := root.Copy()
	if !root.ShallowEqual(c) {
		t.Errorf("error copying: %v %v", root, c)
	}
}

func TestIsMove(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.Resignation1)
	assert.NoError(t, err)

	root := s.Root()
	assert.False(t, root.IsMove())

	assert.True(t, len(root.Down) > 0)
	current := root.Down[0]
	assert.True(t, current.IsMove())

	assert.True(t, len(current.Down) > 0)
	current = current.Down[0]
	assert.True(t, current.IsMove())
}

func TestRecomputeDepth(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleFourMoves)
	assert.NoError(t, err)

	root := s.Root()

	root.Depth = 42
	root.RecomputeDepth()
	assert.True(t, len(root.Down) > 0)
	node := root.Down[0]
	assert.Equal(t, node.Depth, 43)
}

func TestOverwriteField(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleFourMoves)
	assert.NoError(t, err)

	root := s.Root()

	assert.Equal(t, len(root.Fields["PB"]), 1)
	assert.Equal(t, root.Fields["PB"][0], "Black")

	root.OverwriteField("PB", "foobar")

	assert.Equal(t, len(root.Fields["PB"]), 1)
	assert.Equal(t, root.Fields["PB"][0], "foobar")
}
