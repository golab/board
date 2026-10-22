/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package tree

import (
	"github.com/golab/board/pkg/core/color"
	"github.com/golab/board/pkg/core/coord"
	"github.com/golab/board/pkg/core/fields"
)

// Fmap applies a function f to every node under (and including) root
// even thought State contains Nodes and we can range over all the nodes,
// Fmap is useful for two reasons:
//   - it applies the function to the nodes in a consistent and hierarchical order
//   - it can be applied to a branch only, just supply it with a starting node
//     that isn't the root
func Fmap(f func(*TreeNode), root *TreeNode) {
	stack := []*TreeNode{root}
	for len(stack) > 0 {
		// pop
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]

		// apply step
		f(cur)

		// push on stack in reverse order
		for i := len(cur.Down) - 1; i >= 0; i-- {
			stack = append(stack, cur.Down[i])
		}
	}
}

type TreeNode struct {
	XY             *coord.Coord
	Color          color.Color
	Down           []*TreeNode
	Up             *TreeNode
	Index          int
	PreferredChild int
	fields.Fields
	Diff      *coord.Diff
	Depth     int
	BlackCaps int
	WhiteCaps int
}

func NewTreeNode(crd *coord.Coord, col color.Color, index int, up *TreeNode, flds fields.Fields) *TreeNode {
	down := []*TreeNode{}
	node := &TreeNode{
		XY:             crd,
		Color:          col,
		Down:           down,
		Up:             nil,
		Index:          index,
		PreferredChild: 0,
		Fields:         flds,
		Diff:           nil,
		Depth:          0,
		BlackCaps:      0,
		WhiteCaps:      0,
	}
	if up != nil {
		node.SetParent(up)
	}
	return node
}

func (n *TreeNode) ShallowEqual(m *TreeNode) bool {
	return n.XY.Equal(m.XY) &&
		n.Color == m.Color &&
		len(n.Down) == len(m.Down) &&
		n.Index == m.Index &&
		n.PreferredChild == m.PreferredChild &&
		n.Depth == m.Depth
}

func (n *TreeNode) SetDiff(diff *coord.Diff) {
	n.Diff = diff
	b := 0
	w := 0
	if diff != nil {
		for _, ss := range diff.Remove {
			switch ss.Color {
			case color.White:
				b += len(ss.Coords)
			case color.Black:
				w += len(ss.Coords)
			}
		}
	}
	baseB := 0
	baseW := 0
	if n.Up != nil {
		baseB = n.Up.BlackCaps
		baseW = n.Up.WhiteCaps
	}
	n.BlackCaps = baseB + b
	n.WhiteCaps = baseW + w
}

func (n *TreeNode) TrunkNum(i int) int {
	cur := n
	j := 0
	for j < i {
		if len(cur.Down) == 0 {
			return -1
		}
		cur = cur.Down[0]
		j++
	}
	return cur.Index
}

// SetParent exists to add the depth attribute
func (n *TreeNode) SetParent(up *TreeNode) {
	n.Up = up
	n.Depth = up.Depth + 1
}

// RecomputeDepth is used when grafting, to make sure the depth is set
// correctly for all lower nodes
func (n *TreeNode) RecomputeDepth() {
	Fmap(func(m *TreeNode) {
		if m.Up != nil {
			m.Depth = m.Up.Depth + 1
		}
	}, n)
}

func (n *TreeNode) HasChild(crd *coord.Coord, col color.Color) (int, bool) {
	for _, node := range n.Down {
		if node.XY.Equal(crd) && node.Color == col {
			return node.Index, true
		}
	}
	return 0, false
}

// Copy copies the subtree rooted at n (iteratively, since trees can be
// deep enough to make recursion a problem)
func (n *TreeNode) Copy() *TreeNode {
	copyOne := func(src *TreeNode) *TreeNode {
		// parent will get assigned later
		m := NewTreeNode(
			src.XY.Copy(),
			src.Color,
			0,
			nil,
			src.Fields)
		m.Down = []*TreeNode{}
		m.Diff = src.Diff.Copy()
		return m
	}

	root := copyOne(n)
	type pair struct{ src, dst *TreeNode }
	stack := []pair{{n, root}}
	for len(stack) > 0 {
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]
		for _, d := range cur.src.Down {
			e := copyOne(d)
			e.SetParent(cur.dst)
			cur.dst.Down = append(cur.dst.Down, e)
			stack = append(stack, pair{d, e})
		}
	}
	return root
}

func (n *TreeNode) MaxDepth() int {
	depth := 0
	Fmap(func(m *TreeNode) {
		if m.Depth > depth {
			depth = m.Depth
		}
	}, n)
	return depth
}
