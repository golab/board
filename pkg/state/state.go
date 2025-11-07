/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jarednogo/board/pkg/core"
)

const Letters = "ABCDEFGHIJKLNMOPQRSTUVWXYZ"

// as a rule, anything that would need to get sent to new connections
// should be stored here
type State struct {
	root       *core.TreeNode
	current    *core.TreeNode
	head       *core.TreeNode
	nodes      map[int]*core.TreeNode
	nextIndex  int
	size       int
	board      *core.Board
	clipboard  *core.TreeNode
	markedDead core.CoordSet
	markedDame core.CoordSet
}

func (s *State) HeadColor() core.Color {
	return s.head.Color
}

func (s *State) SetNextIndex(i int) {
	s.nextIndex = i
}

func (s *State) GetNextIndex() int {
	i := s.nextIndex
	s.nextIndex++
	return i
}

func (s *State) Size() int {
	return s.size
}

func (s *State) Board() *core.Board {
	return s.board
}

func (s *State) Current() *core.TreeNode {
	return s.current
}

func (s *State) Root() *core.TreeNode {
	return s.root
}

func (s *State) Head() *core.TreeNode {
	return s.head
}

func (s *State) Nodes() map[int]*core.TreeNode {
	return s.nodes
}

func (s *State) AnyMove() {
	s.markedDead = core.NewCoordSet()
	s.markedDame = core.NewCoordSet()
}

func (s *State) ToSGF(indexes bool) string {
	result := "("
	stack := []interface{}{s.root}
	for len(stack) > 0 {
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]
		if str, ok := cur.(string); ok {
			result += str
			continue
		}
		node := cur.(*core.TreeNode)
		result += ";"
		// throw in other fields
		for key, multifield := range node.Fields {
			if key == "IX" {
				continue
			}
			result += key
			for _, fieldValue := range multifield {
				m := strings.ReplaceAll(fieldValue, "]", "\\]")
				result += fmt.Sprintf("[%s]", m)
			}

		}

		if indexes {
			result += fmt.Sprintf("IX[%d]", node.Index)
		}

		if len(node.Down) == 1 {
			stack = append(stack, node.Down[0])
		} else if len(node.Down) > 1 {
			// go backward through array
			for i := len(node.Down) - 1; i >= 0; i-- {
				n := node.Down[i]
				stack = append(stack, ")")
				stack = append(stack, n)
				stack = append(stack, "(")
			}
		}
	}

	result += ")"
	return result
}

func FromSGF(data string) (*State, error) {
	p := core.NewParser(data)
	root, err := p.Parse()
	if err != nil {
		return nil, err
	}

	var size int64 = 19
	if _, ok := root.Fields["SZ"]; ok {
		sizeField := root.Fields["SZ"]
		if len(sizeField) != 1 {
			return nil, fmt.Errorf("SZ cannot be a multifield")
		}
		size, err = strconv.ParseInt(sizeField[0], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	state := NewState(int(size), false)
	stack := []interface{}{root}
	for len(stack) > 0 {
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]
		if _, ok := cur.(string); ok {
			state.left()
		} else {
			node := cur.(*core.SGFNode)

			index := -1
			if indexes, ok := node.Fields["IX"]; ok {
				if len(indexes) > 0 {
					_index, err := strconv.ParseInt(indexes[0], 10, 64)
					index = int(_index)
					if err != nil {
						index = -1
					}
				}
			}

			// refuse to process sgfs with a suicide move
			if node.Coord() != nil && !state.board.Legal(node.Coord(), node.Color()) {
				return nil, fmt.Errorf("suicide moves are not currently supported")
			}

			if node.IsPass() {
				state.addPassNode(node.Color(), node.Fields, index)
			} else if node.IsMove() {
				state.addNode(node.Coord(), node.Color(), node.Fields, index, true)
			} else {
				state.addFieldNode(node.Fields, index)
			}
			for i := len(node.Down) - 1; i >= 0; i-- {
				stack = append(stack, "<")
				stack = append(stack, node.Down[i])
			}
			// TODO: this might be wrong in some cases
			state.head = state.current
		}
	}
	state.rewind()
	state.resetPrefs()
	return state, nil
}

func NewState(size int, initRoot bool) *State {
	nodes := make(map[int]*core.TreeNode)
	var root *core.TreeNode
	root = nil
	index := 0
	if initRoot {
		fields := map[string][]string{}
		fields["GM"] = []string{"1"}
		fields["FF"] = []string{"4"}
		fields["CA"] = []string{"UTF-8"}
		fields["SZ"] = []string{fmt.Sprintf("%d", size)}
		fields["PB"] = []string{"Black"}
		fields["PW"] = []string{"White"}
		fields["RU"] = []string{"Japanese"}
		fields["KM"] = []string{"6.5"}

		// coord, color, index, up, fields
		root = core.NewTreeNode(nil, core.NoColor, 0, nil, fields)
		nodes[0] = root
		index = 1
	}
	board := core.NewBoard(size)
	return &State{
		root:       root,
		current:    root,
		head:       root,
		nodes:      nodes,
		nextIndex:  index,
		size:       size,
		board:      board,
		clipboard:  nil,
		markedDead: core.NewCoordSet(),
		markedDame: core.NewCoordSet()}
}
