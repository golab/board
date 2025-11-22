/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"fmt"
	"sort"
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

func (s *State) EditPlayerBlack(value string) {
	s.root.OverwriteField("PB", value)
}

func (s *State) EditPlayerWhite(value string) {
	s.root.OverwriteField("PW", value)
}

func (s *State) EditKomi(value string) {
	s.root.OverwriteField("KM", value)
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

func (s *State) ToSGF() string {
	return s.toSGF(false)
}

func (s *State) ToSGFIX() string {
	return s.toSGF(true)
}

func (s *State) toSGF(indexes bool) string {
	sb := strings.Builder{}
	// approximate preallocation, each move being at least 6 characters
	// i.e. ;B[aa]
	sb.Grow(6 * len(s.Nodes()))
	sb.WriteByte('(')
	stack := []any{s.root}
	for len(stack) > 0 {
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]
		if str, ok := cur.(string); ok {
			sb.WriteString(str)
			continue
		}
		node := cur.(*core.TreeNode)
		sb.WriteByte(';')

		// throw in other fields
		sort.Slice(node.Fields, func(i, j int) bool {
			return node.Fields[i].Key < node.Fields[j].Key
		})

		for _, field := range node.Fields {
			key := field.Key
			multifield := field.Values
			if key == "IX" {
				continue
			}
			sb.WriteString(key)
			for _, fieldValue := range multifield {
				m := strings.ReplaceAll(fieldValue, "]", "\\]")
				sb.WriteByte('[')
				sb.WriteString(m)
				sb.WriteByte(']')
			}

		}

		if indexes {
			sb.WriteString("IX[")
			sb.WriteString(strconv.Itoa(node.Index))
			sb.WriteByte(']')
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

	sb.WriteByte(')')
	return sb.String()
}

func FromSGF(data string) (*State, error) {
	p := core.NewParser(data)
	root, err := p.Parse()
	if err != nil {
		return nil, err
	}

	var size int64 = 19
	sizeField := root.GetField("SZ")
	if len(sizeField) > 1 {
		return nil, fmt.Errorf("SZ cannot be a multifield")
	}
	if len(sizeField) == 1 {
		size, err = strconv.ParseInt(sizeField[0], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	state := NewState(int(size), false)
	stack := []any{root}
	for len(stack) > 0 {
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]
		if _, ok := cur.(string); ok {
			state.left()
		} else {
			node := cur.(*core.SGFNode)

			index := -1
			indexes := node.GetField("IX")
			if len(indexes) > 0 {
				_index, err := strconv.ParseInt(indexes[0], 10, 64)
				index = int(_index)
				if err != nil {
					index = -1
				}
			}

			// refuse to process sgfs with a suicide move
			if node.Coord() != nil && !state.board.Legal(node.Coord(), node.Color()) {
				return nil, fmt.Errorf("suicide moves are not currently supported")
			}

			if node.IsPass() {
				state.addPassNode(node.Color(), node.Fields(), index)
			} else if node.IsMove() {
				state.addNode(node.Coord(), node.Color(), node.Fields(), index, true)
			} else {
				state.addFieldNode(node.Fields(), index)
			}
			for i := node.NumChildren() - 1; i >= 0; i-- {
				stack = append(stack, "<")
				stack = append(stack, node.GetChild(i))
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
		fields := []core.Field{}
		root = core.NewTreeNode(nil, core.NoColor, 0, nil, fields)
		root.AddField("GM", "1")
		root.AddField("FF", "4")
		root.AddField("CA", "UTF-8")
		root.AddField("SZ", strconv.Itoa(size))
		root.AddField("PB", "Black")
		root.AddField("PW", "White")
		root.AddField("RU", "Japanese")
		root.AddField("KM", "6.5")

		// coord, color, index, up, fields
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
