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

func (s *State) HandleAddStone(evt *core.EventJSON) (*core.Frame, error) {
	c, err := core.InterfaceToCoord(evt.Value)
	if err != nil {
		return nil, err
	}
	x := c.X
	y := c.Y
	if x >= s.Size || y >= s.Size || x < 0 || y < 0 {
		return nil, nil
	}

	col := core.Color(evt.Color)

	// if a child already exists with that coord and col, then actually
	// this is just a gotoindex operation
	if child, ok := s.Current.HasChild(c, col); ok {
		s.GotoIndex(child) //nolint: errcheck
		return s.GenerateFullFrame(core.CurrentAndPreferred), nil
	}

	// do nothing on a suicide move
	if !s.Board.Legal(c, col) {
		return nil, nil
	}

	fields := make(map[string][]string)
	key := "B"
	if col == core.White {
		key = "W"
	}
	fields[key] = []string{c.ToLetters()}

	diff := s.AddNode(c, col, fields, -1, false)

	marks := s.GenerateMarks()

	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      diff,
		Marks:     marks,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.CreateTreeJSON(core.PartialNodes),
		BlackCaps: s.Current.BlackCaps,
		WhiteCaps: s.Current.WhiteCaps,
	}, nil
}

func (s *State) HandlePass(evt *core.EventJSON) (*core.Frame, error) {
	fields := make(map[string][]string)
	col := core.Color(evt.Color)
	key := "B"
	if col == core.White {
		key = "W"
	}
	fields[key] = []string{""}
	s.AddPassNode(core.Color(evt.Color), fields, -1)

	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      nil,
		Marks:     nil,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.CreateTreeJSON(core.PartialNodes),
		BlackCaps: s.Current.BlackCaps,
		WhiteCaps: s.Current.WhiteCaps,
	}, nil
}

func (s *State) HandleRemoveStone(evt *core.EventJSON) (*core.Frame, error) {
	c, err := core.InterfaceToCoord(evt.Value)
	if err != nil {
		return nil, err
	}

	x := c.X
	y := c.Y
	if x >= s.Size || y >= s.Size || x < 0 || y < 0 {
		return nil, nil
	}

	fields := make(map[string][]string)
	fields["AE"] = []string{c.ToLetters()}
	diff := s.AddFieldNode(fields, -1)

	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      diff,
		Marks:     nil,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.CreateTreeJSON(core.PartialNodes),
		BlackCaps: s.Current.BlackCaps,
		WhiteCaps: s.Current.WhiteCaps,
	}, nil
}

func (s *State) HandleAddTriangle(evt *core.EventJSON) (*core.Frame, error) {
	c, err := core.InterfaceToCoord(evt.Value)
	if err != nil {
		return nil, err
	}

	x := c.X
	y := c.Y
	if x >= s.Size || y >= s.Size || x < 0 || y < 0 {
		return nil, nil
	}
	l := c.ToLetters()
	s.Current.AddField("TR", l)
	return nil, nil
}

func (s *State) HandleAddSquare(evt *core.EventJSON) (*core.Frame, error) {
	c, err := core.InterfaceToCoord(evt.Value)
	if err != nil {
		return nil, err
	}

	x := c.X
	y := c.Y
	if x >= s.Size || y >= s.Size || x < 0 || y < 0 {
		return nil, nil
	}
	l := c.ToLetters()
	s.Current.AddField("SQ", l)
	return nil, nil
}

func (s *State) HandleAddLetter(evt *core.EventJSON) (*core.Frame, error) {
	val := evt.Value.(map[string]interface{})
	c, err := core.InterfaceToCoord(val["coords"])
	if err != nil {
		return nil, err
	}

	x := c.X
	y := c.Y
	if x >= s.Size || y >= s.Size || x < 0 || y < 0 {
		return nil, nil
	}

	l := c.ToLetters()
	letter := val["letter"].(string)
	lb := fmt.Sprintf("%s:%s", l, letter)
	s.Current.AddField("LB", lb)
	return nil, nil
}

func (s *State) HandleAddNumber(evt *core.EventJSON) (*core.Frame, error) {
	val := evt.Value.(map[string]interface{})
	c, err := core.InterfaceToCoord(val["coords"])
	if err != nil {
		return nil, err
	}

	x := c.X
	y := c.Y
	if x >= s.Size || y >= s.Size || x < 0 || y < 0 {
		return nil, nil
	}

	l := c.ToLetters()
	number := int(val["number"].(float64))
	lb := fmt.Sprintf("%s:%d", l, number)
	s.Current.AddField("LB", lb)
	return nil, nil
}

func (s *State) HandleRemoveMark(evt *core.EventJSON) (*core.Frame, error) {
	c, err := core.InterfaceToCoord(evt.Value)
	if err != nil {
		return nil, err
	}

	l := c.ToLetters()
	for key, values := range s.Current.Fields {
		for _, value := range values {
			if key == "LB" && value[:2] == l {
				s.Current.RemoveField("LB", value)
			} else if key == "SQ" && value == l {
				s.Current.RemoveField("SQ", l)
			} else if key == "TR" && value == l {
				s.Current.RemoveField("TR", l)
			}
		}
	}
	return nil, nil
}

func (s *State) HandleLeft() (*core.Frame, error) {
	diff := s.Left()
	marks := s.GenerateMarks()
	comments := s.GenerateComments()
	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      diff,
		Marks:     marks,
		Comments:  comments,
		Metadata:  nil,
		TreeJSON:  s.CreateTreeJSON(core.CurrentOnly),
		BlackCaps: s.Current.BlackCaps,
		WhiteCaps: s.Current.WhiteCaps,
	}, nil
}

func (s *State) HandleRight() (*core.Frame, error) {
	diff := s.Right()
	marks := s.GenerateMarks()
	comments := s.GenerateComments()

	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      diff,
		Marks:     marks,
		Comments:  comments,
		Metadata:  nil,
		TreeJSON:  s.CreateTreeJSON(core.CurrentOnly),
		BlackCaps: s.Current.BlackCaps,
		WhiteCaps: s.Current.WhiteCaps,
	}, nil
}

func (s *State) HandleUp() (*core.Frame, error) {
	s.Up()
	// for the current mark
	marks := s.GenerateMarks()

	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      nil,
		Marks:     marks,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.CreateTreeJSON(core.CurrentAndPreferred),
		BlackCaps: s.Current.BlackCaps,
		WhiteCaps: s.Current.WhiteCaps,
	}, nil
}

func (s *State) HandleDown() (*core.Frame, error) {
	s.Down()

	// for the current mark
	marks := s.GenerateMarks()

	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      nil,
		Marks:     marks,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.CreateTreeJSON(core.CurrentAndPreferred),
		BlackCaps: s.Current.BlackCaps,
		WhiteCaps: s.Current.WhiteCaps,
	}, nil
}

func (s *State) HandleRewind() (*core.Frame, error) {
	s.Rewind()
	return s.GenerateFullFrame(core.CurrentOnly), nil
}

func (s *State) HandleFastForward() (*core.Frame, error) {
	s.FastForward()
	return s.GenerateFullFrame(core.CurrentOnly), nil
}

func (s *State) HandleGotoGrid(evt *core.EventJSON) (*core.Frame, error) {
	index := int(evt.Value.(float64))
	s.GotoIndex(index) //nolint: errcheck
	return s.GenerateFullFrame(core.CurrentAndPreferred), nil
}

func (s *State) HandleGotoCoord(evt *core.EventJSON) (*core.Frame, error) {
	coords := make([]int, 0)
	// coerce the value to an array
	val := evt.Value.([]interface{})
	for _, v := range val {
		i := int(v.(float64))
		coords = append(coords, i)
	}
	x := coords[0]
	y := coords[1]
	s.GotoCoord(x, y)
	return s.GenerateFullFrame(core.CurrentAndPreferred), nil

}

func (s *State) HandleComment(evt *core.EventJSON) (*core.Frame, error) {
	val := evt.Value.(string)
	s.Current.AddField("C", val+"\n")
	return nil, nil
}

func (s *State) HandleDraw(evt *core.EventJSON) (*core.Frame, error) {
	vals := evt.Value.([]interface{})
	var x0 float64
	var y0 float64
	if vals[0] == nil {
		x0 = -1.0
	} else {
		x0 = vals[0].(float64)
	}

	if vals[1] == nil {
		y0 = -1.0
	} else {
		y0 = vals[1].(float64)
	}

	x1 := vals[2].(float64)
	y1 := vals[3].(float64)
	color := vals[4].(string)

	value := fmt.Sprintf("%.4f:%.4f:%.4f:%.4f:%s", x0, y0, x1, y1, color)
	s.Current.AddField("PX", value)
	return nil, nil

}

func (s *State) HandleErasePen() (*core.Frame, error) {
	delete(s.Current.Fields, "PX")
	return nil, nil
}

func (s *State) HandleCut() (*core.Frame, error) {
	diff := s.Cut()
	marks := s.GenerateMarks()
	comments := s.GenerateComments()
	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      diff,
		Marks:     marks,
		Comments:  comments,
		Metadata:  nil,
		TreeJSON:  s.CreateTreeJSON(core.Full),
		BlackCaps: s.Current.BlackCaps,
		WhiteCaps: s.Current.WhiteCaps,
	}, nil
}

func (s *State) HandleCopy() (*core.Frame, error) {
	s.Clipboard = s.Current.Copy()
	return nil, nil
}

func (s *State) HandleClipboard() (*core.Frame, error) {
	if s.Clipboard == nil {
		return nil, nil
	}

	// keep a copy of the clipboard unaltered
	branch := s.Clipboard.Copy()

	// first give the copy indexes
	// only possible with state context because of GetNextIndex
	// consider other ways of reindexing, or maybe this should be its
	// own function
	core.Fmap(func(n *core.TreeNode) {
		i := s.GetNextIndex()
		n.Index = i
		s.Nodes[i] = n
	}, branch)

	// set parent and child relationships
	branch.SetParent(s.Current)
	s.Current.Down = append(s.Current.Down, branch)

	// save the parent pref
	savedPref := s.Current.PreferredChild

	// recompute depth
	branch.RecomputeDepth()

	// recompute diffs
	core.Fmap(func(n *core.TreeNode) {
		if n.IsMove() {
			n.SetDiff(s.ComputeDiffMove(n.Index))
		} else {
			n.SetDiff(s.ComputeDiffSetup(n.Index))
		}
	}, branch)

	// restore savedpref
	s.Current.PreferredChild = savedPref

	marks := s.GenerateMarks()
	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      nil,
		Marks:     marks,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.CreateTreeJSON(core.Full),
		BlackCaps: s.Current.BlackCaps,
		WhiteCaps: s.Current.WhiteCaps,
	}, nil
}

func (s *State) HandleGraft(evt *core.EventJSON) (*core.Frame, error) {
	// convert the event value to a string and split into tokens
	v := evt.Value.(string)
	tokens := strings.Split(v, " ")

	/*
		if len(tokens) < 2 {
			return nil, nil
		}
	*/

	// currently accepting the first arg as a move number or not
	var parentIndex int
	start := 0

	// convert the move number to an int
	mv64, err := strconv.ParseInt(tokens[0], 10, 64)

	if err != nil {
		parentIndex = s.Current.Index

	} else {
		start = 1

		mv := int(mv64)

		// interpret the move number as a trunk number, and find the
		// corresponding index
		parentIndex = s.Root.TrunkNum(mv)

		if parentIndex == -1 {
			return nil, fmt.Errorf("trunk too short")
		}
	}

	/*
		if parentIndex == 0 {
			return nil, fmt.Errorf("won't graft onto root")
		}
	*/

	// setup the moves array and initial color
	moves := []*core.PatternMove{}
	col := s.Nodes[parentIndex].Color
	if col == core.NoColor {
		col = core.White
	}

	// go through each token
	for _, tok := range tokens[start:] {

		// convert to a Coord
		coord, err := core.AlphanumericToCoord(tok)
		if err != nil {
			return nil, err
		}

		// get new color
		col = core.Opposite(col)

		// create a new move
		move := &core.PatternMove{Coord: coord, Color: col}
		moves = append(moves, move)
	}

	// call the state's graft function
	s.SmartGraft(parentIndex, moves)

	return nil, nil
}

func (s *State) HandleScore() (*core.Frame, error) {
	blackArea, whiteArea, blackDead, whiteDead, dame := s.Board.Score(s.MarkedDead, s.MarkedDame)
	//fmt.Println()
	//fmt.Println(s.Board)
	frame := &core.Frame{
		BlackCaps: s.Current.BlackCaps + len(blackArea) + len(whiteDead),
		WhiteCaps: s.Current.WhiteCaps + len(whiteArea) + len(blackDead),
		BlackArea: blackArea,
		WhiteArea: whiteArea,
		Dame:      dame,
	}

	return frame, nil
}

func (s *State) HandleMarkDead(evt *core.EventJSON) (*core.Frame, error) {
	c, err := core.InterfaceToCoord(evt.Value)
	if err != nil {
		return nil, err
	}
	x := c.X
	y := c.Y
	if x >= s.Size || y >= s.Size || x < 0 || y < 0 {
		return nil, nil
	}

	if s.Board.Get(c) == core.NoColor {
		dame, _ := s.Board.FindArea(c, core.NewCoordSet())
		if s.MarkedDame.Has(c) {
			s.MarkedDame.RemoveAll(dame)
		} else {
			s.MarkedDame.AddAll(dame)
		}
	} else {
		gp := s.Board.FindGroup(c)
		if s.MarkedDead.Has(c) {
			s.MarkedDead.RemoveAll(gp.Coords)
		} else {
			s.MarkedDead.AddAll(gp.Coords)
		}
	}
	return s.HandleScore()
}
