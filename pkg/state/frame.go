/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state

import (
	"strconv"
	"strings"

	"github.com/jarednogo/board/pkg/core"
)

func (s *State) generateMarks() *core.Marks {
	marks := &core.Marks{}
	if s.current.XY != nil {
		marks.Current = s.current.XY
	}
	if trs := s.current.GetField("TR"); len(trs) > 0 {
		cs := core.NewCoordSet()
		for _, tr := range trs {
			c := core.LettersToCoord(tr)
			cs.Add(c)
		}
		marks.Triangles = cs.List()
	}

	if sqs := s.current.GetField("SQ"); len(sqs) > 0 {
		cs := core.NewCoordSet()
		for _, sq := range sqs {
			c := core.LettersToCoord(sq)
			cs.Add(c)
		}
		marks.Squares = cs.List()
	}
	if lbs := s.current.GetField("LB"); len(lbs) > 0 {
		labels := []*core.Label{}
		for _, lb := range lbs {
			spl := strings.Split(lb, ":")
			c := core.LettersToCoord(spl[0])
			text := spl[1]
			label := &core.Label{Coord: c, Text: text}
			labels = append(labels, label)
		}
		marks.Labels = labels
	}

	if pxs := s.current.GetField("PX"); len(pxs) > 0 {
		pens := []*core.Pen{}
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
			pen := &core.Pen{X0: x0, Y0: y0, X1: x1, Y1: y1, Color: spl[4]}
			pens = append(pens, pen)
		}
		marks.Pens = pens
	}
	return marks
}

func (s *State) generateMetadata() *core.Metadata {
	m := &core.Metadata{
		Size:   s.size,
		Fields: s.root.Fields,
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

func (s *State) GenerateFullFrame(t core.TreeJSONType) *core.Frame {
	frame := &core.Frame{}
	frame.Type = core.FullFrame
	frame.Diff = s.board.CurrentDiff()
	frame.Marks = s.generateMarks()
	frame.Metadata = s.generateMetadata()
	frame.Comments = s.generateComments()
	frame.TreeJSON = s.saveTree(t)
	frame.BlackCaps = s.current.BlackCaps
	frame.WhiteCaps = s.current.WhiteCaps
	return frame
}
