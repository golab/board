/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package core

// Fmap applies a function f to every node under (and including) root
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
	XY             *Coord
	Color          Color
	Down           []*TreeNode
	Up             *TreeNode
	Index          int
	PreferredChild int
	Fields         map[string][]string
	Diff           *Diff
	Depth          int
}

func NewTreeNode(coord *Coord, col Color, index int, up *TreeNode, fields map[string][]string) *TreeNode {
	if fields == nil {
		fields = make(map[string][]string)
	}
	down := []*TreeNode{}
	node := &TreeNode{coord, col, down, nil, index, 0, fields, nil, 0}
	if up != nil {
		node.SetParent(up)
	}
	return node
}

func (n *TreeNode) IsMove() bool {
	if _, ok := n.Fields["B"]; ok {
		return true
	}
	if _, ok := n.Fields["W"]; ok {
		return true
	}
	return false
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

func (n *TreeNode) HasChild(coord *Coord, col Color) (int, bool) {
	for _, node := range n.Down {
		if node.XY.Equal(coord) && node.Color == col {
			return node.Index, true
		}
	}
	return 0, false
}

func (n *TreeNode) Copy() *TreeNode {
	// copy fields
	fields := make(map[string][]string)
	for key, value := range n.Fields {
		newValue := make([]string, len(value))
		copy(newValue, value)
		fields[key] = newValue
	}

	// parent will get assigned later
	m := NewTreeNode(
		n.XY.Copy(),
		n.Color,
		0,
		nil,
		fields)

	// copy children
	down := []*TreeNode{}
	for _, d := range n.Down {
		e := d.Copy()
		e.SetParent(m)
		down = append(down, e)
	}

	m.Down = down
	m.Diff = n.Diff.Copy()

	return m
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

func (n *TreeNode) AddField(key, value string) {
	if _, ok := n.Fields[key]; !ok {
		n.Fields[key] = []string{}
	}
	n.Fields[key] = append(n.Fields[key], value)
}

func (n *TreeNode) RemoveField(key, value string) {
	if _, ok := n.Fields[key]; !ok {
		return
	}
	index := -1
	for i, v := range n.Fields[key] {
		if v == value {
			index = i
		}
	}
	if index == -1 {
		return
	}
	n.Fields[key] = append(n.Fields[key][:index], n.Fields[key][index+1:]...)
	if len(n.Fields[key]) == 0 {
		delete(n.Fields, key)
	}
}
