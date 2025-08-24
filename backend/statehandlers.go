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

func (s *State) HandleAddStone(evt *EventJSON) (*Frame, error) {
	c, err := InterfaceToCoord(evt.Value)
	if err != nil {
		return nil, err
	}
	x := c.X
	y := c.Y
	if x >= s.Size || y >= s.Size || x < 0 || y < 0 {
		return nil, nil
	}

	col := Color(evt.Color)

	// if a child already exists with that coord and col, then actually
	// this is just a gotoindex operation
	if child, ok := s.Current.HasChild(c, col); ok {
		s.GotoIndex(child)
		return s.GenerateFullFrame(CurrentAndPreferred), nil
	}

	// do nothing on a suicide move
	if !s.Board.Legal(c, col) {
		return nil, nil
	}

	fields := make(map[string][]string)
	key := "B"
	if col == White {
		key = "W"
	}
	fields[key] = []string{c.ToLetters()}

	diff := s.AddNode(c, col, fields, -1, false)

	marks := s.GenerateMarks()

	return &Frame{DiffFrame, diff, marks, nil, nil, s.CreateTreeJSON(PartialNodes)}, nil
}

func (s *State) HandlePass(evt *EventJSON) (*Frame, error) {
	fields := make(map[string][]string)
	col := Color(evt.Color)
	key := "B"
	if col == White {
		key = "W"
	}
	fields[key] = []string{""}
	s.AddPassNode(Color(evt.Color), fields, -1)

	return &Frame{DiffFrame, nil, nil, nil, nil, s.CreateTreeJSON(PartialNodes)}, nil
}

func (s *State) HandleRemoveStone(evt *EventJSON) (*Frame, error) {
	c, err := InterfaceToCoord(evt.Value)
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

	return &Frame{DiffFrame, diff, nil, nil, nil, s.CreateTreeJSON(PartialNodes)}, nil
}

func (s *State) HandleAddTriangle(evt *EventJSON) (*Frame, error) {
	c, err := InterfaceToCoord(evt.Value)
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

func (s *State) HandleAddSquare(evt *EventJSON) (*Frame, error) {
	c, err := InterfaceToCoord(evt.Value)
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

func (s *State) HandleAddLetter(evt *EventJSON) (*Frame, error) {
	val := evt.Value.(map[string]interface{})
	c, err := InterfaceToCoord(val["coords"])
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

func (s *State) HandleAddNumber(evt *EventJSON) (*Frame, error) {
	val := evt.Value.(map[string]interface{})
	c, err := InterfaceToCoord(val["coords"])
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

func (s *State) HandleRemoveMark(evt *EventJSON) (*Frame, error) {
	c, err := InterfaceToCoord(evt.Value)
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

func (s *State) HandleLeft() (*Frame, error) {
	diff := s.Left()
	marks := s.GenerateMarks()
	comments := s.GenerateComments()
	return &Frame{DiffFrame, diff, marks, comments, nil, s.CreateTreeJSON(CurrentOnly)}, nil
}

func (s *State) HandleRight() (*Frame, error) {
	diff := s.Right()
	marks := s.GenerateMarks()
	comments := s.GenerateComments()

	return &Frame{DiffFrame, diff, marks, comments, nil, s.CreateTreeJSON(CurrentOnly)}, nil
}

func (s *State) HandleUp() (*Frame, error) {
	s.Up()
	// for the current mark
	marks := s.GenerateMarks()

	return &Frame{DiffFrame, nil, marks, nil, nil, s.CreateTreeJSON(CurrentAndPreferred)}, nil
}

func (s *State) HandleDown() (*Frame, error) {
	s.Down()

	// for the current mark
	marks := s.GenerateMarks()

	return &Frame{DiffFrame, nil, marks, nil, nil, s.CreateTreeJSON(CurrentAndPreferred)}, nil
}

func (s *State) HandleRewind() (*Frame, error) {
	s.Rewind()
	return s.GenerateFullFrame(CurrentOnly), nil
}

func (s *State) HandleFastForward() (*Frame, error) {
	s.FastForward()
	return s.GenerateFullFrame(CurrentOnly), nil
}

func (s *State) HandleGotoGrid(evt *EventJSON) (*Frame, error) {
	index := int(evt.Value.(float64))
	s.GotoIndex(index)
	return s.GenerateFullFrame(CurrentAndPreferred), nil
}

func (s *State) HandleGotoCoord(evt *EventJSON) (*Frame, error) {
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
	return s.GenerateFullFrame(CurrentAndPreferred), nil

}

func (s *State) HandleComment(evt *EventJSON) (*Frame, error) {
	val := evt.Value.(string)
	s.Current.AddField("C", val+"\n")
	return nil, nil
}

func (s *State) HandleDraw(evt *EventJSON) (*Frame, error) {
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

func (s *State) HandleErasePen() (*Frame, error) {
	delete(s.Current.Fields, "PX")
	return nil, nil
}

func (s *State) HandleCut() (*Frame, error) {
	diff := s.Cut()
	marks := s.GenerateMarks()
	comments := s.GenerateComments()
	return &Frame{DiffFrame, diff, marks, comments, nil, s.CreateTreeJSON(Full)}, nil
}

func (s *State) HandleCopy() (*Frame, error) {
	s.Clipboard = s.Current.Copy()
	return nil, nil
}

func (s *State) HandleClipboard() (*Frame, error) {
	if s.Clipboard == nil {
		return nil, nil
	}

	// keep a copy of the clipboard unaltered
	clipboard := s.Clipboard.Copy()

	// first give the copy indexes
	Fmap(func(n *TreeNode) {
		i := s.GetNextIndex()
		n.Index = i
		s.Nodes[i] = n
	}, clipboard)

	// add the clipboard branch to the children of the current node
	s.Current.Down = append(s.Current.Down, clipboard)

	// set the current node to be the parent of the clipboard branch
	// this also sets the depth appropriately
	clipboard.SetParent(s.Current)

	// adusts the depth of all the lower nodes
	clipboard.RecomputeDepth()

	marks := s.GenerateMarks()
	return &Frame{DiffFrame, nil, marks, nil, nil, s.CreateTreeJSON(Full)}, nil
}
