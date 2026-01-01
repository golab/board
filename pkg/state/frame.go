/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"strconv"
	"strings"

	"github.com/golab/board/pkg/core/color"
	"github.com/golab/board/pkg/core/coord"
	"github.com/golab/board/pkg/core/fields"
)

// FrameType can be either DiffFrame or FullFrame
type FrameType int

const (
	DiffFrame = iota
	FullFrame
)

// Frame provides the data for when the board needs to be updated (not the explorer)
type Frame struct {
	Type      FrameType      `json:"type"`
	Diff      *coord.Diff    `json:"diff"`
	Marks     *Marks         `json:"marks"`
	Comments  []string       `json:"comments"`
	Metadata  *Metadata      `json:"metadata"`
	TreeJSON  *TreeJSON      `json:"tree"`
	BlackCaps int            `json:"black_caps"`
	WhiteCaps int            `json:"white_caps"`
	BlackArea []*coord.Coord `json:"black_area"`
	WhiteArea []*coord.Coord `json:"white_area"`
	Dame      []*coord.Coord `json:"dame"`
}

// Marks provides data for any marks on the board
type Marks struct {
	Current   *coord.Coord   `json:"current"`
	Squares   []*coord.Coord `json:"squares"`
	Triangles []*coord.Coord `json:"triangles"`
	Labels    []*Label       `json:"labels"`
	Pens      []*Pen         `json:"pens"`
}

// Label can be any text, but typically single digits or letters
type Label struct {
	Coord *coord.Coord `json:"coord"`
	Text  string       `json:"text"`
}

// Pen contains a start and end coordinate plus a color
type Pen struct {
	X0    float64 `json:"x0"`
	Y0    float64 `json:"y0"`
	X1    float64 `json:"x1"`
	Y1    float64 `json:"y1"`
	Color string  `json:"color"`
}

// Metadata provides the size of the board plus any fields (usually from the root node)
type Metadata struct {
	Size   int            `json:"size"`
	Fields []fields.Field `json:"fields"`
}

// TreeJSONType defines some options for how much data to send in a TreeJSON
type TreeJSONType int

const (
	CurrentOnly TreeJSONType = iota
	CurrentAndPreferred
	PartialNodes
	Full
)

// NodeJSON is a key component of TreeJSON
type NodeJSON struct {
	Color   color.Color `json:"color"`
	Down    []int       `json:"down"`
	Depth   int         `json:"depth"`
	Comment bool        `json:"comment"`
}

// TreeJSON is the basic struct to encode information about the explorer
// this makes up one part of a Frame
type TreeJSON struct {
	Nodes     map[int]*NodeJSON `json:"nodes"`
	Current   int               `json:"current"`
	Preferred []int             `json:"preferred"`
	Depth     int               `json:"depth"`
	Up        int               `json:"up"`
	Root      int               `json:"root"`
}

func (s *State) generateMarks() *Marks {
	marks := &Marks{}
	if s.current.XY != nil {
		marks.Current = s.current.XY
	}
	if trs := s.current.GetField("TR"); len(trs) > 0 {
		cs := coord.NewCoordSet()
		for _, tr := range trs {
			c := coord.FromLetters(tr)
			cs.Add(c)
		}
		marks.Triangles = cs.List()
	}

	if sqs := s.current.GetField("SQ"); len(sqs) > 0 {
		cs := coord.NewCoordSet()
		for _, sq := range sqs {
			c := coord.FromLetters(sq)
			cs.Add(c)
		}
		marks.Squares = cs.List()
	}
	if lbs := s.current.GetField("LB"); len(lbs) > 0 {
		labels := []*Label{}
		for _, lb := range lbs {
			spl := strings.SplitN(lb, ":", 2)
			c := coord.FromLetters(spl[0])
			text := spl[1]
			label := &Label{Coord: c, Text: text}
			labels = append(labels, label)
		}
		marks.Labels = labels
	}

	if pxs := s.current.GetField("PX"); len(pxs) > 0 {
		pens := []*Pen{}
		for _, px := range pxs {
			spl := strings.Split(px, ":")
			if len(spl) != 5 {
				continue
			}
			x0, err0 := strconv.ParseFloat(spl[0], 64)
			y0, err1 := strconv.ParseFloat(spl[1], 64)
			x1, err2 := strconv.ParseFloat(spl[2], 64)
			y1, err3 := strconv.ParseFloat(spl[3], 64)
			hasErr := err0 != nil || err1 != nil || err2 != nil || err3 != nil
			if hasErr {
				continue
			}
			pen := &Pen{X0: x0, Y0: y0, X1: x1, Y1: y1, Color: spl[4]}
			pens = append(pens, pen)
		}
		marks.Pens = pens
	}
	return marks
}

func (s *State) generateMetadata() *Metadata {
	m := &Metadata{
		Size:   s.size,
		Fields: s.root.AllFields(),
	}
	return m
}

func (s *State) generateComments() []string {
	cmts := []string{}
	if c := s.current.GetField("C"); len(c) > 0 {
		cmts = c
	}
	return cmts
}

func (s *State) GenerateFullFrame(t TreeJSONType) *Frame {
	frame := &Frame{}
	frame.Type = FullFrame
	frame.Diff = s.board.CurrentDiff()
	frame.Marks = s.generateMarks()
	frame.Metadata = s.generateMetadata()
	frame.Comments = s.generateComments()
	frame.TreeJSON = s.saveTree(t)
	frame.BlackCaps = s.current.BlackCaps
	frame.WhiteCaps = s.current.WhiteCaps
	return frame
}

func (s *State) GenerateTreeOnly(t TreeJSONType) *Frame {
	frame := &Frame{}
	frame.Type = DiffFrame
	frame.Diff = nil
	frame.Marks = s.generateMarks()
	frame.Comments = nil
	frame.Metadata = nil
	frame.TreeJSON = s.saveTree(t)
	frame.BlackCaps = s.current.BlackCaps
	frame.WhiteCaps = s.current.WhiteCaps
	return frame
}
