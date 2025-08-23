/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"fmt"
)

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
		if node.XY.Equal(coord) {
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

type TreeJSONType int

const (
	CurrentOnly TreeJSONType = iota
	CurrentAndPreferred
	PartialNodes
	Full
)

type NodeJSON struct {
	Color Color `json:"color"`
	Down  []int `json:"down"`
	Depth int   `json:"depth"`
}

type TreeJSON struct {
	Nodes     map[int]*NodeJSON `json:"nodes"`
	Current   int               `json:"current"`
	Preferred []int             `json:"preferred"`
	Depth     int               `json:"depth"`
	Up        int               `json:"up"`
	Root      int               `json:"root"`
}

func (s *State) CreateTreeJSON(t TreeJSONType) *TreeJSON {
	// only really used when we have a partial tree
	up := 0
	root := 0

	// nodes
	var nodes map[int]*NodeJSON

	if t >= PartialNodes {
		nodes = make(map[int]*NodeJSON)
		var start *TreeNode
		// we can choose to send the full or just a partial tree
		// based on which node we start on

		if t == PartialNodes {
			start = s.Current
			up = start.Up.Index
			root = start.Index
		} else if t == Full {
			start = s.Root
		}
		Fmap(func(n *TreeNode) {
			down := []int{}
			for _, c := range n.Down {
				down = append(down, c.Index)
			}
			nodes[n.Index] = &NodeJSON{
				Color: n.Color,
				Down:  down,
				Depth: n.Depth,
			}

		}, start)
	}

	// preferred
	var preferred []int = nil
	if t >= CurrentAndPreferred {
		node := s.Root
		preferred = []int{node.Index}
		for len(node.Down) != 0 {
			node = node.Down[node.PreferredChild]
			preferred = append(preferred, node.Index)
		}
	}

	return &TreeJSON{
		Nodes:     nodes,
		Current:   s.Current.Index,
		Preferred: preferred,
		Depth:     s.Root.MaxDepth(),
		Up:        up,
		Root:      root,
	}
}

/*
func (n *TreeNode) FillGrid(currentIndex int) *Explorer {
	stack := []interface{}{n}
	x := 0
	y := 0
	gridLen := 1
	grid := make(map[[2]int]int)
	loc := make(map[int][2]int)
	colors := make(map[int]Color)
	parents := make(map[int]int)
	prefs := make(map[int]int)

	var currentCoord *Coord
	var currentColor Color

	nodes := []*GridNode{}
	nodeMap := make(map[int]*GridNode)

	for len(stack) > 0 {
		// pop off the stack
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if _, ok := cur.(string); ok {
			x--
			continue
		}

		node := cur.(*TreeNode)
		colors[node.Index] = node.Color
		if node.Up != nil {
			parents[node.Index] = node.Up.Index
		}
		if len(node.Down) > 0 {
			prefs[node.Index] = node.Down[node.PreferredChild].Index
		}

		y = gridLen - 1

		if grid[[2]int{y, x}] != 0 {
			// if there's something in the last row (in the x coord)
			// add a new row
			gridLen++
			y++
		} else {
			for y != 0 {

				// look at the parent
				p := node.Up
				if p != nil {
					a := loc[p.Index]
					x1 := a[0]
					y1 := a[1]
					// actually don't go any farther than the
					// diagonal connecting the parent
					if x-y >= x1-y1 {
						break
					}

					// don't go any farther than the parent row
					if y == y1 {
						break
					}
				}

				// i want to find the earliest row
				// (before going past the parent)
				// that is empty
				if grid[[2]int{y, x}] == 0 && grid[[2]int{y - 1, x}] != 0 {
					break
				}
				y--
			}
		}
		grid[[2]int{y, x}] = node.Index
		loc[node.Index] = [2]int{x, y}

		gridNode := &GridNode{&Coord{x, y}, node.Color, node.Index, nil, nil}
		nodes = append(nodes, gridNode)
		nodeMap[node.Index] = gridNode

		if node.Index == currentIndex {
			currentCoord = &Coord{x, y}
			currentColor = node.Color
		}

		// if the parent is a diagonal away, we have to take up
		// another node
		// (this is for all the "angled" edges")
		p := node.Up
		if p != nil {
			a := loc[p.Index]
			y1 := a[1]
			if y-y1 > 1 {
				if grid[[2]int{y - 1, x - 1}] == 0 {
					grid[[2]int{y - 1, x - 1}] = -1
				}
			}
		}
		x++

		// push on children in reverse order
		for i := len(node.Down) - 1; i >= 0; i-- {
			stack = append(stack, "")
			stack = append(stack, node.Down[i])
		}
	}

	edges := []*GridEdge{}
	for i, l := range loc {
		// gather all the nodes with their color attached
		x := l[0]
		y := l[1]

		// gather all the edges
		p, ok := parents[i]
		if !ok {
			continue
		}

		pCoord := loc[p]
		start := &Coord{pCoord[0], pCoord[1]}
		end := &Coord{x, y}
		edge := &GridEdge{start, end}
		edges = append(edges, edge)
	}


	var currentNode *GridNode
	preferredNodes := []*GridNode{}
	index := 0
	for {
		if l, ok := loc[index]; ok {
			x := l[0]
			y := l[1]
			gridNode := &GridNode{&Coord{x, y}, colors[index], index, nil, nil}

			if len(preferredNodes) > 0 {
				left := preferredNodes[len(preferredNodes)-1]
				left.Right = gridNode
				gridNode.Left = left
			}

			preferredNodes = append(preferredNodes, gridNode)

			if index == currentIndex {
				currentNode = gridNode
			}

			if index, ok = prefs[index]; !ok {
				break
			}
		} else {
			break
		}
	}

	return &Explorer{
		nodes,
		edges,
		preferredNodes,
		currentCoord,
		currentColor,
		currentNode,
	}
}
*/

type GridNode struct {
	Coord *Coord `json:"coord"`
	Color `json:"color"`
	Index int       `json:"index"`
	Left  *GridNode `json:"-"`
	Right *GridNode `json:"-"`
}

func (n *GridNode) String() string {
	if n.Left == nil && n.Right == nil {
		return fmt.Sprintf("%v %v %v", n.Coord, n.Color, n.Index)
	} else if n.Left == nil {
		return fmt.Sprintf("%v %v %v right=%d", n.Coord, n.Color, n.Index, n.Right.Index)
	} else if n.Right == nil {
		return fmt.Sprintf("%v %v %v left=%d", n.Coord, n.Color, n.Index, n.Left.Index)
	}
	return fmt.Sprintf("%v %v %v left=%d right=%d", n.Coord, n.Color, n.Index, n.Left.Index, n.Right.Index)
}

type GridEdge struct {
	Start *Coord `json:"start"`
	End   *Coord `json:"end"`
}

type Explorer struct {
	Nodes          []*GridNode `json:"nodes"`
	Edges          []*GridEdge `json:"edges"`
	PreferredNodes []*GridNode `json:"preferred_nodes"`
	Current        *Coord      `json:"current"`
	CurrentColor   Color       `json:"current_color"`
	CurrentNode    *GridNode   `json:"-"`
}

func (e *Explorer) Left() *Coord {
	// just go left on the saved explorer

	if e == nil {
		return nil
	}

	if e.CurrentNode == nil {
		return nil
	}

	if e.CurrentNode.Left == nil {
		return nil
	}

	e.CurrentNode = e.CurrentNode.Left
	return e.CurrentNode.Coord
}

func (e *Explorer) Right() *Coord {
	// just go left on the saved explorer

	if e == nil {
		return nil
	}

	if e.CurrentNode == nil {
		return nil
	}

	if e.CurrentNode.Right == nil {
		return nil
	}

	e.CurrentNode = e.CurrentNode.Right
	return e.CurrentNode.Coord
}

func (e *Explorer) Rewind() {
}

func (e *Explorer) FastForward() {
}
