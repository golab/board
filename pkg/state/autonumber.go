/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"fmt"

	"github.com/golab/board/pkg/core/tree"
)

func (s *State) AutoNumber() bool {
	return s.autoNumber
}

func (s *State) SetAutoNumber(b bool) {
	s.autoNumber = b
}

// branchNodes returns the nodes from the first node off the main trunk
// down to (and including) the current node
// if the current node is on the trunk, it returns nil
func (s *State) branchNodes() []*tree.TreeNode {
	// path from current up to root
	path := []*tree.TreeNode{}
	for n := s.current; n != nil; n = n.Up {
		path = append(path, n)
	}

	// walk back down from the root, looking for the first step that
	// doesn't follow the trunk (i.e. isn't the first child)
	for i := len(path) - 2; i >= 0; i-- {
		parent := path[i+1]
		if parent.Down[0] != path[i] {
			branch := []*tree.TreeNode{}
			for j := i; j >= 0; j-- {
				branch = append(branch, path[j])
			}
			return branch
		}
	}
	return nil
}

// numberBranch labels every stone of the current branch that is still
// on the board with its move number, counting from 1 at the first move
// off the trunk
// the labels are added to the current node so they stay visible
// as the branch progresses
func (s *State) numberBranch() {
	branch := s.branchNodes()

	// later moves on the same point (after a capture) take precedence
	numbers := make(map[string]int)
	for i, node := range branch {
		if node.XY == nil || s.board.Get(node.XY) != node.Color {
			continue
		}
		numbers[node.XY.ToLetters()] = i + 1
	}

	// add labels in move order so the output is deterministic
	for i, node := range branch {
		if node.XY == nil {
			continue
		}
		l := node.XY.ToLetters()
		if numbers[l] != i+1 {
			continue
		}
		s.current.AddField("LB", fmt.Sprintf("%s:%d", l, i+1))
	}
}
