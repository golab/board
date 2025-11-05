/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"strings"

	"github.com/jarednogo/board/pkg/core"
)

func (s *State) SetLocation(loc string) {
	if loc != "" {
		dirs := strings.Split(loc, ",")
		// don't need to assign to a variable if we don't use it
		for range dirs {
			s.right()
		}
	}

}

func (s *State) left() *core.Diff {
	s.AnyMove()
	if s.current.Up != nil {
		d := s.current.Diff.Invert()
		s.board.ApplyDiff(d)
		s.current = s.current.Up
		return d
	}
	return nil
}

func (s *State) right() *core.Diff {
	s.AnyMove()
	if len(s.current.Down) > 0 {
		index := s.current.PreferredChild
		s.current = s.current.Down[index]
		d := s.current.Diff
		s.board.ApplyDiff(d)
		return d
	}
	return nil
}

func (s *State) gotoIndex(index int) error {
	err := s.setPreferred(index)
	if err != nil {
		return err
	}
	s.rewind()
	last := s.current.Index
	for s.current.Index != index {
		s.right()
		// to prevent infinite loops
		if s.current.Index == last {
			break
		}
	}
	//s.Current = s.Nodes[index]
	return nil
}

func (s *State) rewind() {
	s.AnyMove()
	s.current = s.root
	s.board.Clear()
	s.board.ApplyDiff(s.current.Diff)
}

func (s *State) fastForward() {
	for len(s.current.Down) != 0 {
		s.right()
	}
}

func (s *State) up() {
	if len(s.current.Down) == 0 {
		return
	}
	c := s.current.PreferredChild
	mod := len(s.current.Down)
	s.current.PreferredChild = (((c - 1) % mod) + mod) % mod
}

func (s *State) down() {
	if len(s.current.Down) == 0 {
		return
	}
	c := s.current.PreferredChild
	mod := len(s.current.Down)
	s.current.PreferredChild = (((c + 1) % mod) + mod) % mod
}

func (s *State) gotoCoord(x, y int) {
	cur := s.current
	// look forward
	for {
		if cur.XY != nil && cur.XY.X == x && cur.XY.Y == y {
			s.gotoIndex(cur.Index) //nolint: errcheck
			return
		}
		if len(cur.Down) == 0 {
			break
		}
		cur = cur.Down[cur.PreferredChild]
	}

	cur = s.current
	// look backward
	for {
		if cur.XY != nil && cur.XY.X == x && cur.XY.Y == y {
			s.gotoIndex(cur.Index) //nolint: errcheck
			return
		}
		if cur.Up == nil {
			break
		}
		cur = cur.Up
	}

	// didn't find anything
	// do nothing
}
