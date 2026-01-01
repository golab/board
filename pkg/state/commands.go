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

	"github.com/golab/board/pkg/core/color"
	"github.com/golab/board/pkg/core/coord"
	"github.com/golab/board/pkg/core/fields"
)

type Command interface {
	Execute(*State) (*Frame, error)
}

type addStoneCommand struct {
	crd   *coord.Coord
	color color.Color
}

func NewAddStoneCommand(crd *coord.Coord, color color.Color) Command {
	return &addStoneCommand{crd, color}
}

func (cmd *addStoneCommand) Execute(s *State) (*Frame, error) {
	x := cmd.crd.X
	y := cmd.crd.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}

	// if a child already exists with that coord and col, then actually
	// this is just a gotoindex operation
	if child, ok := s.current.HasChild(cmd.crd, cmd.color); ok {
		s.gotoIndex(child) //nolint: errcheck
		return s.GenerateFullFrame(CurrentAndPreferred), nil
	}

	// do nothing on a suicide move
	if !s.board.Legal(cmd.crd, cmd.color) {
		return nil, nil
	}

	flds := fields.Fields{}
	key := "B"
	if cmd.color == color.White {
		key = "W"
	}
	flds.AddField(key, cmd.crd.ToLetters())

	diff := s.addNode(cmd.crd, cmd.color, flds, -1, false)

	marks := s.generateMarks()

	return &Frame{
		Type:      DiffFrame,
		Diff:      diff,
		Marks:     marks,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.saveTree(PartialNodes),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type passCommand struct {
	color color.Color
}

func NewPassCommand(color color.Color) Command {
	return &passCommand{color}
}

func (cmd *passCommand) Execute(s *State) (*Frame, error) {
	flds := fields.Fields{}
	key := "B"
	if cmd.color == color.White {
		key = "W"
	}
	flds.AddField(key, "")
	s.addPassNode(cmd.color, flds, -1)

	return &Frame{
		Type:      DiffFrame,
		Diff:      nil,
		Marks:     nil,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.saveTree(PartialNodes),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type removeStoneCommand struct {
	crd *coord.Coord
}

func NewRemoveStoneCommand(crd *coord.Coord) Command {
	return &removeStoneCommand{crd}
}

func (cmd *removeStoneCommand) Execute(s *State) (*Frame, error) {
	x := cmd.crd.X
	y := cmd.crd.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}

	flds := fields.Fields{}
	flds.AddField("AE", cmd.crd.ToLetters())
	diff := s.addFieldNode(flds, -1)

	return &Frame{
		Type:      DiffFrame,
		Diff:      diff,
		Marks:     nil,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.saveTree(PartialNodes),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type addTriangleCommand struct {
	crd *coord.Coord
}

func NewAddTriangleCommand(crd *coord.Coord) Command {
	return &addTriangleCommand{crd}
}

func (cmd *addTriangleCommand) Execute(s *State) (*Frame, error) {
	x := cmd.crd.X
	y := cmd.crd.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}
	l := cmd.crd.ToLetters()
	s.current.AddField("TR", l)
	return nil, nil
}

type addSquareCommand struct {
	crd *coord.Coord
}

func NewAddSquareCommand(crd *coord.Coord) Command {
	return &addSquareCommand{crd}
}

func (cmd *addSquareCommand) Execute(s *State) (*Frame, error) {
	x := cmd.crd.X
	y := cmd.crd.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}
	l := cmd.crd.ToLetters()
	s.current.AddField("SQ", l)
	return nil, nil
}

type addLetterCommand struct {
	crd    *coord.Coord
	letter string
}

func NewAddLetterCommand(crd *coord.Coord, letter string) Command {
	return &addLetterCommand{crd, letter}
}

func (cmd *addLetterCommand) Execute(s *State) (*Frame, error) {
	x := cmd.crd.X
	y := cmd.crd.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}

	l := cmd.crd.ToLetters()
	lb := fmt.Sprintf("%s:%s", l, cmd.letter)
	s.current.AddField("LB", lb)
	return nil, nil
}

type addNumberCommand struct {
	crd    *coord.Coord
	number int
}

func NewAddNumberCommand(crd *coord.Coord, number int) Command {
	return &addNumberCommand{crd, number}
}

func (cmd *addNumberCommand) Execute(s *State) (*Frame, error) {
	x := cmd.crd.X
	y := cmd.crd.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}

	l := cmd.crd.ToLetters()
	lb := fmt.Sprintf("%s:%d", l, cmd.number)
	s.current.AddField("LB", lb)
	return nil, nil
}

type addLabelCommand struct {
	crd   *coord.Coord
	label string
}

func NewAddLabelCommand(crd *coord.Coord, label string) Command {
	return &addLabelCommand{crd, label}
}

func (cmd *addLabelCommand) Execute(s *State) (*Frame, error) {
	x := cmd.crd.X
	y := cmd.crd.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}

	l := cmd.crd.ToLetters()
	lb := fmt.Sprintf("%s:%s", l, cmd.label)
	s.current.AddField("LB", lb)
	return nil, nil
}

type removeMarkCommand struct {
	crd *coord.Coord
}

func NewRemoveMarkCommand(crd *coord.Coord) Command {
	return &removeMarkCommand{crd}
}

func (cmd *removeMarkCommand) Execute(s *State) (*Frame, error) {
	l := cmd.crd.ToLetters()
	for _, field := range s.current.AllFields() {
		key := field.Key
		values := field.Values
		for _, value := range values {
			if key == "LB" && value[:2] == l {
				s.current.RemoveField("LB", value)
			} else if key == "SQ" && value == l {
				s.current.RemoveField("SQ", l)
			} else if key == "TR" && value == l {
				s.current.RemoveField("TR", l)
			}
		}
	}
	return nil, nil
}

type leftCommand struct{}

func NewLeftCommand() Command {
	return &leftCommand{}
}

func (cmd *leftCommand) Execute(s *State) (*Frame, error) {
	diff := s.left()
	marks := s.generateMarks()
	comments := s.generateComments()
	return &Frame{
		Type:      DiffFrame,
		Diff:      diff,
		Marks:     marks,
		Comments:  comments,
		Metadata:  nil,
		TreeJSON:  s.saveTree(CurrentOnly),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type rightCommand struct{}

func NewRightCommand() Command {
	return &rightCommand{}
}

func (cmd *rightCommand) Execute(s *State) (*Frame, error) {
	diff := s.right()
	marks := s.generateMarks()
	comments := s.generateComments()

	return &Frame{
		Type:      DiffFrame,
		Diff:      diff,
		Marks:     marks,
		Comments:  comments,
		Metadata:  nil,
		TreeJSON:  s.saveTree(CurrentOnly),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type upCommand struct{}

func NewUpCommand() Command {
	return &upCommand{}
}

func (cmd *upCommand) Execute(s *State) (*Frame, error) {
	s.up()
	// for the current mark
	marks := s.generateMarks()

	return &Frame{
		Type:      DiffFrame,
		Diff:      nil,
		Marks:     marks,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.saveTree(CurrentAndPreferred),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type downCommand struct{}

func NewDownCommand() Command {
	return &downCommand{}
}

func (cmd *downCommand) Execute(s *State) (*Frame, error) {
	s.down()

	// for the current mark
	marks := s.generateMarks()

	return &Frame{
		Type:      DiffFrame,
		Diff:      nil,
		Marks:     marks,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.saveTree(CurrentAndPreferred),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type rewindCommand struct{}

func NewRewindCommand() Command {
	return &rewindCommand{}
}

func (cmd *rewindCommand) Execute(s *State) (*Frame, error) {
	s.rewind()
	return s.GenerateFullFrame(CurrentOnly), nil
}

type fastForwardCommand struct{}

func NewFastForwardCommand() Command {
	return &fastForwardCommand{}
}

func (cmd *fastForwardCommand) Execute(s *State) (*Frame, error) {
	s.fastForward()
	return s.GenerateFullFrame(CurrentOnly), nil
}

type gotoGridCommand struct {
	index int
}

func NewGotoGridCommand(index int) Command {
	return &gotoGridCommand{index}
}

func (cmd *gotoGridCommand) Execute(s *State) (*Frame, error) {
	s.gotoIndex(cmd.index) //nolint: errcheck
	return s.GenerateFullFrame(CurrentAndPreferred), nil
}

type gotoCoordCommand struct {
	crd *coord.Coord
}

func NewGotoCoordCommand(crd *coord.Coord) Command {
	return &gotoCoordCommand{crd}
}

func (cmd *gotoCoordCommand) Execute(s *State) (*Frame, error) {
	x := cmd.crd.X
	y := cmd.crd.Y
	s.gotoCoord(x, y)
	return s.GenerateFullFrame(CurrentAndPreferred), nil
}

type commentCommand struct {
	text string
}

func NewCommentCommand(text string) Command {
	return &commentCommand{text}
}

func (cmd *commentCommand) Execute(s *State) (*Frame, error) {
	s.current.AddField("C", cmd.text+"\n")
	return s.GenerateFullFrame(Full), nil
}

type drawCommand struct {
	x0    float64
	y0    float64
	x1    float64
	y1    float64
	color string
}

func NewDrawCommand(x0, y0, x1, y1 float64, color string) Command {
	return &drawCommand{x0, y0, x1, y1, color}
}

func (cmd *drawCommand) Execute(s *State) (*Frame, error) {
	value := fmt.Sprintf(
		"%.4f:%.4f:%.4f:%.4f:%s",
		cmd.x0,
		cmd.y0,
		cmd.x1,
		cmd.y1,
		cmd.color)
	s.current.AddField("PX", value)
	return nil, nil
}

type erasePenCommand struct{}

func NewErasePenCommand() Command {
	return &erasePenCommand{}
}

func (cmd *erasePenCommand) Execute(s *State) (*Frame, error) {
	s.current.DeleteField("PX")
	return nil, nil
}

type cutCommand struct{}

func NewCutCommand() Command {
	return &cutCommand{}
}

func (cmd *cutCommand) Execute(s *State) (*Frame, error) {
	diff := s.cut()
	marks := s.generateMarks()
	comments := s.generateComments()
	return &Frame{
		Type:      DiffFrame,
		Diff:      diff,
		Marks:     marks,
		Comments:  comments,
		Metadata:  nil,
		TreeJSON:  s.saveTree(Full),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type copyCommand struct{}

func NewCopyCommand() Command {
	return &copyCommand{}
}

func (cmd *copyCommand) Execute(s *State) (*Frame, error) {
	s.clipboard = s.current.Copy()
	return nil, nil
}

type pasteCommand struct{}

func NewPasteCommand() Command {
	return &pasteCommand{}
}

func (cmd *pasteCommand) Execute(s *State) (*Frame, error) {
	if s.clipboard == nil {
		return nil, nil
	}

	s.paste()

	marks := s.generateMarks()
	return &Frame{
		Type:      DiffFrame,
		Diff:      nil,
		Marks:     marks,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.saveTree(Full),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type graftCommand struct {
	text string
}

func NewGraftCommand(text string) Command {
	return &graftCommand{text}
}

func (cmd *graftCommand) Execute(s *State) (*Frame, error) {
	tokens := strings.Split(cmd.text, " ")

	// currently accepting the first arg as a move number or not
	var parentIndex int
	start := 0

	// convert the move number to an int
	mv64, err := strconv.ParseInt(tokens[0], 10, 64)

	if err != nil {
		parentIndex = s.current.Index

	} else {
		start = 1

		mv := int(mv64)

		// interpret the move number as a trunk number, and find the
		// corresponding index
		parentIndex = s.root.TrunkNum(mv)

		if parentIndex == -1 {
			return nil, fmt.Errorf("trunk too short")
		}
	}

	// setup the moves array and initial color
	moves := []*coord.Stone{}
	col := s.nodes[parentIndex].Color
	if col == color.Empty {
		col = color.White
	}

	// go through each token
	for _, tok := range tokens[start:] {

		// convert to a Coord
		crd, err := coord.FromAlphanumeric(tok, s.size)
		if err != nil {
			return nil, err
		}

		// get new color
		col = col.Opposite()

		// create a new move
		move := &coord.Stone{Coord: crd, Color: col}
		moves = append(moves, move)
	}

	// call the state's graft function
	s.smartGraft(parentIndex, moves)

	return nil, nil
}

type scoreCommand struct{}

func NewScoreCommand() Command {
	return &scoreCommand{}
}

func (cmd *scoreCommand) Execute(s *State) (*Frame, error) {
	scoreResult := s.board.Score(s.markedDead, s.markedDame)
	blackArea := scoreResult.BlackArea
	whiteArea := scoreResult.WhiteArea
	blackDead := scoreResult.BlackDead
	whiteDead := scoreResult.WhiteDead
	dame := scoreResult.Dame
	frame := &Frame{
		BlackCaps: s.current.BlackCaps + len(blackArea) + len(whiteDead),
		WhiteCaps: s.current.WhiteCaps + len(whiteArea) + len(blackDead),
		BlackArea: blackArea,
		WhiteArea: whiteArea,
		Dame:      dame,
	}

	return frame, nil
}

type markDeadCommand struct {
	crd *coord.Coord
}

func NewMarkDeadCommand(crd *coord.Coord) Command {
	return &markDeadCommand{crd}
}

func (cmd *markDeadCommand) Execute(s *State) (*Frame, error) {
	x := cmd.crd.X
	y := cmd.crd.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}

	if s.board.Get(cmd.crd) == color.Empty {
		dame, _ := s.board.FindArea(cmd.crd, coord.NewCoordSet())
		if s.markedDame.Has(cmd.crd) {
			s.markedDame.RemoveAll(dame)
		} else {
			s.markedDame.AddAll(dame)
		}
	} else {
		gp := s.board.FindGroup(cmd.crd)
		if s.markedDead.Has(cmd.crd) {
			s.markedDead.RemoveAll(gp.Coords)
		} else {
			s.markedDead.AddAll(gp.Coords)
		}
	}
	return (&scoreCommand{}).Execute(s)
}
