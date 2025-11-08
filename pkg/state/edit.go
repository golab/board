/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"github.com/jarednogo/board/pkg/core"
)

func (s *State) addFieldNode(fields map[string][]string, index int) *core.Diff {
	s.AnyMove()
	tmp := s.GetNextIndex()
	if index == -1 {
		index = tmp
	}
	n := core.NewTreeNode(nil, core.NoColor, index, s.current, fields)
	s.nodes[index] = n
	if s.root == nil {
		s.root = n
	} else {
		s.current.Down = append(s.current.Down, n)
		s.current.PreferredChild = len(s.current.Down) - 1
	}
	s.current = n

	// compute diff
	diff := s.computeDiffSetup(index)
	s.board.ApplyDiff(diff)
	s.current.SetDiff(diff)
	return diff
}

func (s *State) addPassNode(col core.Color, fields map[string][]string, index int) {
	s.AnyMove()
	tmp := s.GetNextIndex()
	if index == -1 {
		index = tmp
	}
	n := core.NewTreeNode(nil, col, index, s.current, fields)
	s.nodes[index] = n
	if s.root == nil {
		s.root = n
	} else {
		s.current.Down = append(s.current.Down, n)
		s.current.PreferredChild = len(s.current.Down) - 1
	}
	s.current = n
	// no need to add a diff
	// but actually, SetDiff also sets the score
	s.current.SetDiff(nil)
}

func (s *State) PushHead(x, y int, col core.Color) {
	coord := &core.Coord{X: x, Y: y}
	if x == -1 || y == -1 {
		coord = nil
	}
	index := s.GetNextIndex()
	fields := make(map[string][]string)
	var key string
	if col == core.Black {
		key = "B"
	} else {
		key = "W"
	}
	value := ""
	if x != -1 {
		value = coord.ToLetters()
	}
	fields[key] = []string{value}
	n := core.NewTreeNode(coord, col, index, s.head, fields)
	s.nodes[index] = n
	if len(s.head.Down) > 0 {
		s.head.PreferredChild++
	}
	s.head.Down = append([]*core.TreeNode{n}, s.head.Down...)

	// tracking the head or not
	tracking := s.current == s.head

	var diff *core.Diff

	// if we're not tracking the head
	if !tracking {
		// save where we currently are
		save := s.current.Index

		// goto head
		s.gotoIndex(s.head.Index) //nolint: errcheck

		// compute diff
		diff = s.board.Move(coord, col)

		// go back to saved index
		s.gotoIndex(save) //nolint: errcheck
	} else {
		// do nothing if it's a pass

		// otherwise
		if x != -1 {
			// if we are tracking, just compute the diff
			diff = s.board.Move(coord, col)
		}

		// and follow along
		s.current = n
	}

	// set new head
	s.head = n

	// set diff
	s.head.SetDiff(diff)
}

func (s *State) addNode(coord *core.Coord, col core.Color, fields map[string][]string, index int, force bool) *core.Diff {
	s.AnyMove()
	if fields == nil {
		fields = make(map[string][]string)
	}

	if !force {
		// check to see if it's already there
		for i, node := range s.current.Down {
			coordOld := node.XY
			if coordOld != nil &&
				coord != nil &&
				coordOld.X == coord.X &&
				coordOld.Y == coord.Y &&
				node.Color == core.Color(col) {
				s.current.PreferredChild = i
				s.right()
				return s.current.Diff
			}
		}
	}

	tmp := s.GetNextIndex()
	if index == -1 {
		index = tmp
	}
	n := core.NewTreeNode(coord, core.Color(col), index, s.current, fields)

	s.nodes[index] = n
	if s.root == nil {
		s.root = n
	} else {
		s.current.Down = append(s.current.Down, n)
		s.current.PreferredChild = len(s.current.Down) - 1
	}
	s.current = n
	diff := s.board.Move(coord, core.Color(col))
	s.current.SetDiff(diff)
	return diff
}

func (s *State) AddStones(moves []*core.Stone) {
	node := s.root
	locationSave := s.current.Index

	for _, move := range moves {
		found := false

		for _, child := range node.Down {
			if (child.XY == nil && move.Coord == nil) ||
				(child.XY != nil && move.Coord != nil &&
					child.XY.X == move.Coord.X && child.XY.Y == move.Coord.Y &&
					child.Color == move.Color) {
				node = child
				found = true
				break
			}
		}

		if !found {
			s.gotoIndex(node.Index) //nolint: errcheck

			fields := make(map[string][]string)
			key := "B"
			if move.Color == core.White {
				key = "W"
			}

			if move.Coord == nil {
				fields[key] = []string{""}
				s.addPassNode(move.Color, fields, -1)
			} else {
				fields[key] = []string{move.Coord.ToLetters()}
				s.addNode(move.Coord, move.Color, fields, -1, false)
			}
			node = s.current
		}
	}
	s.gotoIndex(locationSave) //nolint: errcheck
}

// SmartGraft doesn't duplicate existing moves
func (s *State) smartGraft(parentIndex int, moves []*core.Stone) {
	parent := s.nodes[parentIndex]
	savedPrefs := make(map[int]int)
	save := s.current.Index

	var graft *core.TreeNode
	up := parent

	for _, move := range moves {

		// go to the parent
		s.gotoIndex(up.Index) //nolint: errcheck

		// save the preferences on each node that already exists
		savedPrefs[up.Index] = up.PreferredChild

		// if the move exists in a child node, then follow it
		if i, ok := s.current.HasChild(move.Coord, move.Color); ok {
			up = s.nodes[i]
			continue
		}

		// once we get here we are adding new nodes

		// each node needs an index
		index := s.GetNextIndex()

		// each node needs either B[] or W[] field
		fields := make(map[string][]string)
		var key string
		if move.Color == core.Black {
			key = "B"
		} else {
			key = "W"
		}
		fields[key] = []string{move.Coord.ToLetters()}

		// create the node, up is the parent of the new node
		node := core.NewTreeNode(move.Coord, move.Color, index, up, fields)

		// keep track of the first new node
		if graft == nil {
			graft = node
		}

		// add the node to the state's node map
		s.nodes[index] = node

		// follow along so we can set child nodes
		up.Down = append(up.Down, node)

		// calculate the diff
		diff := s.board.Move(move.Coord, move.Color)
		node.SetDiff(diff)

		// set the new parent for the next node
		up = node
	}

	// cleanup

	// (this is only necessary if we added something)
	if graft != nil {
		graft.RecomputeDepth()
	}

	s.gotoIndex(save) //nolint: errcheck
	for index, pref := range savedPrefs {
		s.nodes[index].PreferredChild = pref
	}

}

// Graft may duplicate existing moves
func (s *State) graft(parentIndex int, moves []*core.Stone) {
	parent := s.nodes[parentIndex]
	savedPref := parent.PreferredChild
	save := s.current.Index

	var graft *core.TreeNode
	up := parent

	for _, move := range moves {

		// go to the parent
		s.gotoIndex(up.Index) //nolint: errcheck

		// each node needs an index
		index := s.GetNextIndex()

		// each node needs either B[] or W[] field
		fields := make(map[string][]string)
		var key string
		if move.Color == core.Black {
			key = "B"
		} else {
			key = "W"
		}
		fields[key] = []string{move.Coord.ToLetters()}

		// create the node, up is the parent of the new node
		node := core.NewTreeNode(move.Coord, move.Color, index, up, fields)

		// keep track of the first node
		if graft == nil {
			graft = node
		}

		// add the node to the state's node map
		s.nodes[index] = node

		// follow along so we can set child nodes
		up.Down = append(up.Down, node)

		// calculate the diff
		diff := s.board.Move(move.Coord, move.Color)
		node.SetDiff(diff)

		// set the new parent for the next node
		up = node
	}

	// cleanup
	graft.RecomputeDepth()
	s.gotoIndex(save) //nolint: errcheck
	parent.PreferredChild = savedPref
}

func (s *State) cut() *core.Diff {
	s.AnyMove()
	// store the current index
	index := s.current.Index

	// go left
	diff := s.left()

	// find the child that matches the index to cut
	j := -1
	for i := 0; i < len(s.current.Down); i++ {
		node := s.current.Down[i]
		if node.Index == index {
			j = i
			break
		}
	}

	// if we didn't find anything return (shouldn't really happen)
	if j == -1 {
		return nil
	}

	// store the branch (child index j)
	branch := s.current.Down[j]

	// cut the branch out from the children
	s.current.Down = append(s.current.Down[:j], s.current.Down[j+1:]...)

	// delete all the nodes from the nodes map
	core.Fmap(func(n *core.TreeNode) {
		delete(s.nodes, n.Index)
	}, branch)

	// adjust prefs
	if s.current.PreferredChild >= len(s.current.Down) {
		s.current.PreferredChild = 0
	}

	// save the branch to the clipboard
	s.clipboard = branch

	return diff
}

func (s *State) paste() {
	// keep a copy of the clipboard unaltered
	branch := s.clipboard.Copy()

	// first give the copy indexes
	// only possible with state context because of GetNextIndex
	// consider other ways of reindexing, or maybe this should be its
	// own function
	core.Fmap(func(n *core.TreeNode) {
		i := s.GetNextIndex()
		n.Index = i
		s.nodes[i] = n
	}, branch)

	// set parent and child relationships
	branch.SetParent(s.current)
	s.current.Down = append(s.current.Down, branch)

	// save the parent pref
	savedPref := s.current.PreferredChild

	// recompute depth
	branch.RecomputeDepth()

	// recompute diffs
	core.Fmap(func(n *core.TreeNode) {
		if n.IsMove() {
			n.SetDiff(s.computeDiffMove(n.Index))
		} else {
			n.SetDiff(s.computeDiffSetup(n.Index))
		}
	}, branch)

	// restore savedpref
	s.current.PreferredChild = savedPref
}
