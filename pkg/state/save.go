/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"encoding/base64"

	"github.com/golab/board/pkg/core/tree"
)

type StateJSON struct {
	SGF       string         `json:"sgf"`
	Location  string         `json:"loc"`
	Prefs     map[string]int `json:"prefs"`
	NextIndex int            `json:"next_index"`
}

func (s *State) Save() *StateJSON {
	sgf := s.ToSGFIX()
	encoded := base64.StdEncoding.EncodeToString([]byte(sgf))
	loc := s.locate()
	prefs := s.prefs()
	stateStruct := &StateJSON{
		SGF:       encoded,
		Location:  loc,
		Prefs:     prefs,
		NextIndex: s.nextIndex,
	}
	return stateStruct
}

func (s *State) saveTree(t TreeJSONType) *TreeJSON {
	// only really used when we have a partial tree
	up := 0
	root := 0

	// nodes
	var nodes map[int]*NodeJSON

	if t >= PartialNodes {
		nodes = make(map[int]*NodeJSON)
		var start *tree.TreeNode
		// we can choose to send the full or just a partial tree
		// based on which node we start on

		switch t {
		case PartialNodes:
			start = s.current
			up = start.Up.Index
			root = start.Index
		case Full:
			start = s.root
		}
		tree.Fmap(func(n *tree.TreeNode) {
			down := []int{}
			for _, c := range n.Down {
				down = append(down, c.Index)
			}
			nodes[n.Index] = &NodeJSON{
				Color:   n.Color,
				Down:    down,
				Depth:   n.Depth,
				Comment: HasComment(n),
			}

		}, start)
	}

	// preferred
	var preferred []int
	if t >= CurrentAndPreferred {
		node := s.root
		preferred = []int{node.Index}
		for len(node.Down) != 0 {
			node = node.Down[node.PreferredChild]
			preferred = append(preferred, node.Index)
		}
	}

	return &TreeJSON{
		Nodes:     nodes,
		Current:   s.current.Index,
		Preferred: preferred,
		Depth:     s.root.MaxDepth(),
		Up:        up,
		Root:      root,
	}
}
