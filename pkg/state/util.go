/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"fmt"

	"github.com/jarednogo/board/pkg/core/color"
	"github.com/jarednogo/board/pkg/core/coord"
)

func (s *State) prefs() map[string]int {
	prefs := make(map[string]int)
	for index, node := range s.nodes {
		key := fmt.Sprintf("%d", index)
		prefs[key] = node.PreferredChild
	}
	return prefs
}

func (s *State) setPreferred(index int) error {
	n := s.nodes[index]
	cur := n
	for {
		if cur == nil {
			return fmt.Errorf("error in indexing")
		}
		if cur.Up == nil {
			break
		}
		oldIndex := cur.Index
		cur = cur.Up
		for i, d := range cur.Down {
			if d.Index == oldIndex {
				cur.PreferredChild = i
			}
		}
	}
	return nil
}

// ResetPrefs sets all prefs to 0
func (s *State) resetPrefs() {
	for _, n := range s.nodes {
		n.PreferredChild = 0
	}
}

// SetPrefs takes a map and sets
func (s *State) SetPrefs(prefs map[string]int) {
	for _, n := range s.nodes {
		key := fmt.Sprintf("%d", n.Index)
		p := prefs[key]
		n.PreferredChild = p
	}
}

func (s *State) locate() string {
	dirs := []int{}
	c := s.current
	for {
		myIndex := c.Index
		if c.Up == nil {
			break
		}
		u := c.Up
		for i := 0; i < len(u.Down); i++ {
			if u.Down[i].Index == myIndex {
				dirs = append(dirs, i)
			}
		}
		c = u
	}
	result := ""
	firstComma := false
	for i := len(dirs) - 1; i >= 0; i-- {
		d := dirs[i]
		if !firstComma {
			result = fmt.Sprintf("%d", d)
			firstComma = true
		} else {
			result = fmt.Sprintf("%s,%d", result, d)
		}
	}
	return result
}

func (s *State) computeDiffMove(i int) *coord.Diff {
	save := s.current.Index
	n, ok := s.nodes[i]
	if !ok {
		return nil
	}
	// only for the purposes of checking the board (for removal of stones)
	if n.Up != nil {
		s.gotoIndex(n.Up.Index) //nolint: errcheck
	} else {
		s.board.Clear()
	}

	// at the end, go back to the saved index
	defer s.gotoIndex(save) //nolint: errcheck
	if n.XY == nil {
		return nil
	}

	diff := s.board.Move(n.XY, n.Color)
	return diff
}

func (s *State) computeDiffSetup(i int) *coord.Diff {
	save := s.current.Index
	n, ok := s.nodes[i]
	if !ok {
		return nil
	}

	// only for the purposes of checking the board (for removal of stones)
	if n.Up != nil {
		s.gotoIndex(n.Up.Index) //nolint: errcheck
	} else {
		s.board.Clear()
	}

	// at the end, go back to the saved index
	defer s.gotoIndex(save) //nolint: errcheck

	// find the black stones to add
	diffAdd := []*coord.StoneSet{}
	if val := n.GetField("AB"); len(val) > 0 {
		add := coord.NewCoordSet()
		for _, v := range val {
			add.Add(coord.FromLetters(v))
		}
		stoneSet := coord.NewStoneSet(add, color.Black)
		diffAdd = append(diffAdd, stoneSet)
	}

	// find the white stones to add
	if val := n.GetField("AW"); len(val) > 0 {
		add := coord.NewCoordSet()
		for _, v := range val {
			add.Add(coord.FromLetters(v))
		}
		stoneSet := coord.NewStoneSet(add, color.White)
		diffAdd = append(diffAdd, stoneSet)
	}

	// find the stones to remove
	diffRemove := []*coord.StoneSet{}
	if val := n.GetField("AE"); len(val) > 0 {
		csBlack := coord.NewCoordSet()
		csWhite := coord.NewCoordSet()
		for _, v := range val {
			coord := coord.FromLetters(v)
			col := s.board.Get(coord)
			switch col {
			case color.Black:
				csBlack.Add(coord)
			case color.White:
				csWhite.Add(coord)
			}
		}
		removeBlack := coord.NewStoneSet(csBlack, color.Black)
		removeWhite := coord.NewStoneSet(csWhite, color.White)
		diffRemove = append(diffRemove, removeBlack)
		diffRemove = append(diffRemove, removeWhite)
	}
	diff := coord.NewDiff(diffAdd, diffRemove)
	return diff
}
