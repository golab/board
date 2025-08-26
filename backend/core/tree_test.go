/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core_test

import (
	"github.com/jarednogo/board/backend/core"
	"github.com/jarednogo/board/backend/state"
	"testing"
)

func TestFmap(t *testing.T) {
	s, err := state.FromSGF("(;PW[White]RU[Japanese]KM[6.5]GM[1]FF[4]CA[UTF-8]SZ[19]PB[Black](;B[pd];W[dd];B[pp];W[dp])(;B[dd];W[ee]))")
	if err != nil {
		t.Error(err)
	}

	core.Fmap(func(n *core.TreeNode) {
		n.Color = core.Opposite(n.Color)
	}, s.Root)

	node := s.Nodes[1]
	if node.Color != core.White {
		t.Errorf("fmap failed")
	}
}

func TestField1(t *testing.T) {

	s, err := state.FromSGF("(;PW[White]RU[Japanese]KM[6.5]GM[1]FF[4]CA[UTF-8]SZ[19]PB[Black]C[comment1](;B[pd];W[dd];B[pp];W[dp])(;B[dd];W[ee]))")
	if err != nil {
		t.Error(err)
	}

	root := s.Root
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
	s, err := state.FromSGF("(;PW[White]RU[Japanese]KM[6.5]GM[1]FF[4]CA[UTF-8]SZ[19]PB[Black]C[comment1](;B[pd];W[dd];B[pp];W[dp])(;B[dd];W[ee]))")
	if err != nil {
		t.Error(err)
	}

	m := s.Root.MaxDepth()
	if m != 4 {
		t.Errorf("max depth failed, expected 4, got: %d", m)
	}
}

func TestChild(t *testing.T) {
	s, err := state.FromSGF("(;PW[White]RU[Japanese]KM[6.5]GM[1]FF[4]CA[UTF-8]SZ[19]PB[Black]C[comment1](;B[pd];W[dd];B[pp];W[dp])(;B[dd];W[ee]))")
	if err != nil {
		t.Error(err)
	}
	coord := &core.Coord{3, 3}
	ind, has := s.Root.HasChild(coord, core.Black)
	if !has {
		t.Errorf("failed to find child")
	}
	if ind != 5 {
		t.Errorf("found child at the wrong index %d", ind)
	}
}

func TestGraft(t *testing.T) {
	s, err := state.FromSGF("(;PW[White]RU[Japanese]KM[6.5]GM[1]FF[4]CA[UTF-8]SZ[19]PB[Black]C[comment1](;B[pd];W[dd];B[pp];W[dp])(;B[dd];W[ee]))")
	if err != nil {
		t.Error(err)
	}

	branch := s.Nodes[3].Copy()

	// didn't reindex, but that's not what this is testing
	s.Root.Graft(branch)

	coord := &core.Coord{15, 15}
	_, has := s.Root.HasChild(coord, core.Black)
	if !has {
		t.Errorf("failed to find child")
	}
}
