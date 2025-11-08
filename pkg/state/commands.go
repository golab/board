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

type Command interface {
	Execute(*State) (*core.Frame, error)
}

type addStoneCommand struct {
	coord *core.Coord
	color core.Color
}

func NewAddStoneCommand(coord *core.Coord, color core.Color) Command {
	return &addStoneCommand{coord, color}
}

func (cmd *addStoneCommand) Execute(s *State) (*core.Frame, error) {
	x := cmd.coord.X
	y := cmd.coord.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}

	// if a child already exists with that coord and col, then actually
	// this is just a gotoindex operation
	if child, ok := s.current.HasChild(cmd.coord, cmd.color); ok {
		s.gotoIndex(child) //nolint: errcheck
		return s.GenerateFullFrame(core.CurrentAndPreferred), nil
	}

	// do nothing on a suicide move
	if !s.board.Legal(cmd.coord, cmd.color) {
		return nil, nil
	}

	fields := make(map[string][]string)
	key := "B"
	if cmd.color == core.White {
		key = "W"
	}
	fields[key] = []string{cmd.coord.ToLetters()}

	diff := s.addNode(cmd.coord, cmd.color, fields, -1, false)

	marks := s.generateMarks()

	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      diff,
		Marks:     marks,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.saveTree(core.PartialNodes),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type passCommand struct {
	color core.Color
}

func NewPassCommand(color core.Color) Command {
	return &passCommand{color}
}

func (cmd *passCommand) Execute(s *State) (*core.Frame, error) {
	fields := make(map[string][]string)
	key := "B"
	if cmd.color == core.White {
		key = "W"
	}
	fields[key] = []string{""}
	s.addPassNode(cmd.color, fields, -1)

	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      nil,
		Marks:     nil,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.saveTree(core.PartialNodes),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type removeStoneCommand struct {
	coord *core.Coord
}

func NewRemoveStoneCommand(coord *core.Coord) Command {
	return &removeStoneCommand{coord}
}

func (cmd *removeStoneCommand) Execute(s *State) (*core.Frame, error) {
	x := cmd.coord.X
	y := cmd.coord.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}

	fields := make(map[string][]string)
	fields["AE"] = []string{cmd.coord.ToLetters()}
	diff := s.addFieldNode(fields, -1)

	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      diff,
		Marks:     nil,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.saveTree(core.PartialNodes),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type addTriangleCommand struct {
	coord *core.Coord
}

func NewAddTriangleCommand(coord *core.Coord) Command {
	return &addTriangleCommand{coord}
}

func (cmd *addTriangleCommand) Execute(s *State) (*core.Frame, error) {
	x := cmd.coord.X
	y := cmd.coord.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}
	l := cmd.coord.ToLetters()
	s.current.AddField("TR", l)
	return nil, nil
}

type addSquareCommand struct {
	coord *core.Coord
}

func NewAddSquareCommand(coord *core.Coord) Command {
	return &addSquareCommand{coord}
}

func (cmd *addSquareCommand) Execute(s *State) (*core.Frame, error) {
	x := cmd.coord.X
	y := cmd.coord.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}
	l := cmd.coord.ToLetters()
	s.current.AddField("SQ", l)
	return nil, nil
}

type addLetterCommand struct {
	coord  *core.Coord
	letter string
}

func NewAddLetterCommand(coord *core.Coord, letter string) Command {
	return &addLetterCommand{coord, letter}
}

func (cmd *addLetterCommand) Execute(s *State) (*core.Frame, error) {
	x := cmd.coord.X
	y := cmd.coord.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}

	l := cmd.coord.ToLetters()
	lb := fmt.Sprintf("%s:%s", l, cmd.letter)
	s.current.AddField("LB", lb)
	return nil, nil
}

type addNumberCommand struct {
	coord  *core.Coord
	number int
}

func NewAddNumberCommand(coord *core.Coord, number int) Command {
	return &addNumberCommand{coord, number}
}

func (cmd *addNumberCommand) Execute(s *State) (*core.Frame, error) {
	x := cmd.coord.X
	y := cmd.coord.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}

	l := cmd.coord.ToLetters()
	lb := fmt.Sprintf("%s:%d", l, cmd.number)
	s.current.AddField("LB", lb)
	return nil, nil
}

type removeMarkCommand struct {
	coord *core.Coord
}

func NewRemoveMarkCommand(coord *core.Coord) Command {
	return &removeMarkCommand{coord}
}

func (cmd *removeMarkCommand) Execute(s *State) (*core.Frame, error) {
	l := cmd.coord.ToLetters()
	for key, values := range s.current.Fields {
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

func (cmd *leftCommand) Execute(s *State) (*core.Frame, error) {
	diff := s.left()
	marks := s.generateMarks()
	comments := s.generateComments()
	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      diff,
		Marks:     marks,
		Comments:  comments,
		Metadata:  nil,
		TreeJSON:  s.saveTree(core.CurrentOnly),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type rightCommand struct{}

func NewRightCommand() Command {
	return &rightCommand{}
}

func (cmd *rightCommand) Execute(s *State) (*core.Frame, error) {
	diff := s.right()
	marks := s.generateMarks()
	comments := s.generateComments()

	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      diff,
		Marks:     marks,
		Comments:  comments,
		Metadata:  nil,
		TreeJSON:  s.saveTree(core.CurrentOnly),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type upCommand struct{}

func NewUpCommand() Command {
	return &upCommand{}
}

func (cmd *upCommand) Execute(s *State) (*core.Frame, error) {
	s.up()
	// for the current mark
	marks := s.generateMarks()

	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      nil,
		Marks:     marks,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.saveTree(core.CurrentAndPreferred),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type downCommand struct{}

func NewDownCommand() Command {
	return &downCommand{}
}

func (cmd *downCommand) Execute(s *State) (*core.Frame, error) {
	s.down()

	// for the current mark
	marks := s.generateMarks()

	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      nil,
		Marks:     marks,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.saveTree(core.CurrentAndPreferred),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type rewindCommand struct{}

func NewRewindCommand() Command {
	return &rewindCommand{}
}

func (cmd *rewindCommand) Execute(s *State) (*core.Frame, error) {
	s.rewind()
	return s.GenerateFullFrame(core.CurrentOnly), nil
}

type fastForwardCommand struct{}

func NewFastForwardCommand() Command {
	return &fastForwardCommand{}
}

func (cmd *fastForwardCommand) Execute(s *State) (*core.Frame, error) {
	s.fastForward()
	return s.GenerateFullFrame(core.CurrentOnly), nil
}

type gotoGridCommand struct {
	index int
}

func NewGotoGridCommand(index int) Command {
	return &gotoGridCommand{index}
}

func (cmd *gotoGridCommand) Execute(s *State) (*core.Frame, error) {
	s.gotoIndex(cmd.index) //nolint: errcheck
	return s.GenerateFullFrame(core.CurrentAndPreferred), nil
}

type gotoCoordCommand struct {
	coord *core.Coord
}

func NewGotoCoordCommand(coord *core.Coord) Command {
	return &gotoCoordCommand{coord}
}

func (cmd *gotoCoordCommand) Execute(s *State) (*core.Frame, error) {
	x := cmd.coord.X
	y := cmd.coord.Y
	s.gotoCoord(x, y)
	return s.GenerateFullFrame(core.CurrentAndPreferred), nil
}

type commentCommand struct {
	text string
}

func NewCommentCommand(text string) Command {
	return &commentCommand{text}
}

func (cmd *commentCommand) Execute(s *State) (*core.Frame, error) {
	s.current.AddField("C", cmd.text+"\n")
	return nil, nil
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

func (cmd *drawCommand) Execute(s *State) (*core.Frame, error) {
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

func (cmd *erasePenCommand) Execute(s *State) (*core.Frame, error) {
	delete(s.current.Fields, "PX")
	return nil, nil
}

type cutCommand struct{}

func NewCutCommand() Command {
	return &cutCommand{}
}

func (cmd *cutCommand) Execute(s *State) (*core.Frame, error) {
	diff := s.cut()
	marks := s.generateMarks()
	comments := s.generateComments()
	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      diff,
		Marks:     marks,
		Comments:  comments,
		Metadata:  nil,
		TreeJSON:  s.saveTree(core.Full),
		BlackCaps: s.current.BlackCaps,
		WhiteCaps: s.current.WhiteCaps,
	}, nil
}

type copyCommand struct{}

func NewCopyCommand() Command {
	return &copyCommand{}
}

func (cmd *copyCommand) Execute(s *State) (*core.Frame, error) {
	s.clipboard = s.current.Copy()
	return nil, nil
}

type pasteCommand struct{}

func NewPasteCommand() Command {
	return &pasteCommand{}
}

func (cmd *pasteCommand) Execute(s *State) (*core.Frame, error) {
	if s.clipboard == nil {
		return nil, nil
	}

	s.paste()

	marks := s.generateMarks()
	return &core.Frame{
		Type:      core.DiffFrame,
		Diff:      nil,
		Marks:     marks,
		Comments:  nil,
		Metadata:  nil,
		TreeJSON:  s.saveTree(core.Full),
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

func (cmd *graftCommand) Execute(s *State) (*core.Frame, error) {
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
	moves := []*core.Stone{}
	col := s.nodes[parentIndex].Color
	if col == core.NoColor {
		col = core.White
	}

	// go through each token
	for _, tok := range tokens[start:] {

		// convert to a Coord
		coord, err := core.AlphanumericToCoord(tok, s.size)
		if err != nil {
			return nil, err
		}

		// get new color
		col = core.Opposite(col)

		// create a new move
		move := &core.Stone{Coord: coord, Color: col}
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

func (cmd *scoreCommand) Execute(s *State) (*core.Frame, error) {
	blackArea, whiteArea, blackDead, whiteDead, dame := s.board.Score(s.markedDead, s.markedDame)
	frame := &core.Frame{
		BlackCaps: s.current.BlackCaps + len(blackArea) + len(whiteDead),
		WhiteCaps: s.current.WhiteCaps + len(whiteArea) + len(blackDead),
		BlackArea: blackArea,
		WhiteArea: whiteArea,
		Dame:      dame,
	}

	return frame, nil
}

type markDeadCommand struct {
	coord *core.Coord
}

func NewMarkDeadCommand(coord *core.Coord) Command {
	return &markDeadCommand{coord}
}

func (cmd *markDeadCommand) Execute(s *State) (*core.Frame, error) {
	x := cmd.coord.X
	y := cmd.coord.Y
	if x >= s.size || y >= s.size || x < 0 || y < 0 {
		return nil, nil
	}

	if s.board.Get(cmd.coord) == core.NoColor {
		dame, _ := s.board.FindArea(cmd.coord, core.NewCoordSet())
		if s.markedDame.Has(cmd.coord) {
			s.markedDame.RemoveAll(dame)
		} else {
			s.markedDame.AddAll(dame)
		}
	} else {
		gp := s.board.FindGroup(cmd.coord)
		if s.markedDead.Has(cmd.coord) {
			s.markedDead.RemoveAll(gp.Coords)
		} else {
			s.markedDead.AddAll(gp.Coords)
		}
	}
	return (&scoreCommand{}).Execute(s)
}
