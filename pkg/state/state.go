/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golab/board/pkg/core/board"
	"github.com/golab/board/pkg/core/color"
	"github.com/golab/board/pkg/core/coord"
	"github.com/golab/board/pkg/core/fields"
	"github.com/golab/board/pkg/core/parser"
	"github.com/golab/board/pkg/core/tree"
)

const Letters = "ABCDEFGHIJKLNMOPQRSTUVWXYZ"

// as a rule, anything that would need to get sent to new connections
// should be stored here
type State struct {
	root       *tree.TreeNode
	current    *tree.TreeNode
	head       *tree.TreeNode
	nodes      map[int]*tree.TreeNode
	nextIndex  int
	size       int
	board      *board.Board
	clipboard  *tree.TreeNode
	markedDead coord.CoordSet
	markedDame coord.CoordSet
}

func (s *State) HeadColor() color.Color {
	return s.head.Color
}

func (s *State) GetColorAt(t int) color.Color {
	p := s.root.TrunkNum(t)
	if p == -1 {
		return color.Empty
	}
	node := s.nodes[p]
	return Color(node)
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

func (s *State) Board() *board.Board {
	return s.board
}

func (s *State) Current() *tree.TreeNode {
	return s.current
}

func (s *State) Root() *tree.TreeNode {
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

func (s *State) Head() *tree.TreeNode {
	return s.head
}

func (s *State) Nodes() map[int]*tree.TreeNode {
	return s.nodes
}

func (s *State) AnyMove() {
	s.markedDead = coord.NewCoordSet()
	s.markedDame = coord.NewCoordSet()
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
		node := cur.(*tree.TreeNode)
		sb.WriteByte(';')

		// throw in other fields
		node.SortFields()

		for _, field := range node.AllFields() {
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
	p := parser.New(data)
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

	state := NewEmptyState(int(size))
	stack := []any{root}
	for len(stack) > 0 {
		i := len(stack) - 1
		cur := stack[i]
		stack = stack[:i]
		if _, ok := cur.(string); ok {
			state.left()
		} else {
			node := cur.(*parser.SGFNode)

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
			if Coord(node) != nil && !state.board.Legal(Coord(node), Color(node)) {
				return nil, fmt.Errorf("suicide moves are not currently supported")
			}

			if IsPass(node) {
				state.addPassNode(Color(node), node.Fields, index)
			} else if IsMove(node) {
				crd := Coord(node)
				if crd != nil {
					state.addNode(crd, Color(node), node.Fields, index, true)
				}
			} else {
				state.addFieldNode(node.Fields, index)
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

func NewState(size int) *State {
	return newState(size, true)
}

func NewEmptyState(size int) *State {
	return newState(size, false)
}

func newState(size int, initRoot bool) *State {
	nodes := make(map[int]*tree.TreeNode)
	var root *tree.TreeNode
	root = nil
	index := 0
	if initRoot {
		flds := fields.Fields{}
		root = tree.NewTreeNode(nil, color.Empty, 0, nil, flds)
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
	brd := board.NewBoard(size)
	return &State{
		root:       root,
		current:    root,
		head:       root,
		nodes:      nodes,
		nextIndex:  index,
		size:       size,
		board:      brd,
		clipboard:  nil,
		markedDead: coord.NewCoordSet(),
		markedDame: coord.NewCoordSet()}
}
