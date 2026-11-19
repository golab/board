/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state_test

import (
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/internal/sgfsamples"
	"github.com/golab/board/pkg/core/color"
	"github.com/golab/board/pkg/core/coord"
	"github.com/golab/board/pkg/state"
)

func TestCommandAddStone1(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	coord := coord.NewCoord(9, 9)
	color := color.Black
	_, err = state.NewAddStoneCommand(coord, color).Execute(s)
	assert.NoError(t, err)
}

func TestCommandAddStone2(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.TwoFalseEyes)
	assert.NoError(t, err)

	coord := coord.NewCoord(10, 10)
	color := color.Black
	frm, err := state.NewAddStoneCommand(coord, color).Execute(s)
	assert.NoError(t, err)
	assert.Zero(t, frm)
}

func TestCommandAddStone3(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	coord := coord.NewCoord(15, 3)
	color := color.Black
	frm, err := state.NewAddStoneCommand(coord, color).Execute(s)
	require.NoError(t, err)
	require.NotNil(t, frm)
	assert.Equal(t, frm.Type, state.FullFrame)
}

func TestCommandPass(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewPassCommand(color.Black).Execute(s)
	assert.NoError(t, err)
}

func TestCommandRemoveStone(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	coord := coord.NewCoord(9, 9)
	_, err = state.NewRemoveStoneCommand(coord).Execute(s)
	assert.NoError(t, err)
}

func TestCommandAddLabel(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	coord := coord.NewCoord(9, 9)
	_, err = state.NewAddLabelCommand(coord, "O_o").Execute(s)
	assert.NoError(t, err)

	root := s.Root()
	lbs := root.GetField("LB")
	require.Equal(t, len(lbs), 1)
	assert.Equal(t, lbs[0], "jj:O_o")
}

func TestCommandRemoveMark(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	coord := coord.NewCoord(9, 9)
	_, err = state.NewRemoveMarkCommand(coord).Execute(s)
	assert.NoError(t, err)
}

func TestCommandLeft(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewLeftCommand().Execute(s)
	assert.NoError(t, err)
}

func TestCommandRight(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewRightCommand().Execute(s)
	assert.NoError(t, err)
}

func TestCommandUp(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewUpCommand().Execute(s)
	assert.NoError(t, err)
}

func TestCommandDown(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewDownCommand().Execute(s)
	assert.NoError(t, err)
}

func TestCommandRewind(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewRewindCommand().Execute(s)
	assert.NoError(t, err)
}

func TestCommandFastForward(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewFastForwardCommand().Execute(s)
	assert.NoError(t, err)
}

func TestCommandGotoGrid(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	index := 1

	_, err = state.NewGotoGridCommand(index).Execute(s)
	assert.NoError(t, err)
}

func TestCommandGotoCoord(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	coord := coord.NewCoord(9, 9)

	_, err = state.NewGotoCoordCommand(coord).Execute(s)
	assert.NoError(t, err)
}

func TestCommandComment(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	comment := "some comment"

	_, err = state.NewCommentCommand(comment).Execute(s)
	assert.NoError(t, err)
}

func TestCommandDraw(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewDrawCommand(0.0, 0.0, 5.0, 5.0, "#000000").Execute(s)
	assert.NoError(t, err)
}

func TestCommandErasePen(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewErasePenCommand().Execute(s)
	assert.NoError(t, err)
}

func TestCommandCut(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewCutCommand().Execute(s)
	assert.NoError(t, err)
}

func TestCommandCopy(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewCopyCommand().Execute(s)
	assert.NoError(t, err)
}

func TestCommandPaste(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewRightCommand().Execute(s)
	assert.NoError(t, err)
	_, err = state.NewRightCommand().Execute(s)
	assert.NoError(t, err)
	_, err = state.NewRightCommand().Execute(s)
	assert.NoError(t, err)
	_, err = state.NewRightCommand().Execute(s)
	assert.NoError(t, err)
	_, err = state.NewCopyCommand().Execute(s)
	assert.NoError(t, err)
	_, err = state.NewLeftCommand().Execute(s)
	assert.NoError(t, err)

	_, err = state.NewPasteCommand().Execute(s)
	assert.NoError(t, err)
}

func TestCommandGraft1(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	text := "a1 b2 c3 d4"

	_, err = state.NewGraftCommand(text).Execute(s)
	assert.NoError(t, err)
}

func TestCommandGraft2(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	text := "4 a1 b2 c3 d4"

	_, err = state.NewGraftCommand(text).Execute(s)
	assert.NoError(t, err)
}

func TestCommandScore(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewScoreCommand().Execute(s)
	assert.NoError(t, err)
}

func TestCommandMarkDead(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	coord := coord.NewCoord(9, 9)

	_, err = state.NewMarkDeadCommand(coord).Execute(s)
	assert.NoError(t, err)
}

func TestCommandAddMarksTriangles(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	// duplicates and off-board coords are skipped
	crds := []*coord.Coord{
		coord.NewCoord(0, 0),
		coord.NewCoord(1, 0),
		coord.NewCoord(1, 0),
		nil,
	}
	frm, err := state.NewAddMarksCommand("triangle", crds, "", nil).Execute(s)
	assert.NoError(t, err)
	require.NotNil(t, frm)
	assert.Equal(t, len(frm.Marks.Triangles), 2)
	assert.Equal(t, len(s.Current().GetField("TR")), 2)
}

func TestCommandAddMarksSkipsTaken(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	c := coord.NewCoord(3, 3)
	_, err = state.NewAddTriangleCommand(c).Execute(s)
	assert.NoError(t, err)

	// the only cell is already marked, so nothing changes
	frm, err := state.NewAddMarksCommand("square", []*coord.Coord{c}, "", nil).Execute(s)
	assert.NoError(t, err)
	assert.Zero(t, frm)
	assert.Equal(t, len(s.Current().GetField("SQ")), 0)

	crds := []*coord.Coord{c, coord.NewCoord(4, 3)}
	frm, err = state.NewAddMarksCommand("square", crds, "", nil).Execute(s)
	assert.NoError(t, err)
	require.NotNil(t, frm)
	sqs := s.Current().GetField("SQ")
	require.Equal(t, len(sqs), 1)
	assert.Equal(t, sqs[0], "ed")
	// the frame carries every mark on the node, not just the new ones
	assert.Equal(t, len(frm.Marks.Triangles), 1)
	assert.Equal(t, len(frm.Marks.Squares), 1)
}

func TestCommandAddMarksLetters(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	// B is already in use, so the server hands out A, C
	_, err = state.NewAddLetterCommand(coord.NewCoord(0, 0), "B").Execute(s)
	assert.NoError(t, err)

	crds := []*coord.Coord{coord.NewCoord(1, 0), coord.NewCoord(2, 0)}
	_, err = state.NewAddMarksCommand("letter", crds, "", nil).Execute(s)
	assert.NoError(t, err)

	lbs := s.Current().GetField("LB")
	require.Equal(t, len(lbs), 3)
	assert.Equal(t, lbs[1], "ba:A")
	assert.Equal(t, lbs[2], "ca:C")
}

func TestCommandAddMarksLettersRunOut(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	crds := []*coord.Coord{}
	for x := 0; x < 19; x++ {
		crds = append(crds, coord.NewCoord(x, 0), coord.NewCoord(x, 1))
	}
	_, err = state.NewAddMarksCommand("letter", crds, "", nil).Execute(s)
	assert.NoError(t, err)
	assert.Equal(t, len(s.Current().GetField("LB")), 26)
}

func TestCommandAddMarksNumbers(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewAddNumberCommand(coord.NewCoord(0, 0), 2).Execute(s)
	assert.NoError(t, err)

	crds := []*coord.Coord{coord.NewCoord(1, 0), coord.NewCoord(2, 0)}
	_, err = state.NewAddMarksCommand("number", crds, "", nil).Execute(s)
	assert.NoError(t, err)

	lbs := s.Current().GetField("LB")
	require.Equal(t, len(lbs), 3)
	assert.Equal(t, lbs[1], "ba:1")
	assert.Equal(t, lbs[2], "ca:3")
}

func TestCommandAddMarksLabel(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	crds := []*coord.Coord{coord.NewCoord(1, 0), coord.NewCoord(2, 0)}
	_, err = state.NewAddMarksCommand("label", crds, "O_o", nil).Execute(s)
	assert.NoError(t, err)

	lbs := s.Current().GetField("LB")
	require.Equal(t, len(lbs), 2)
	assert.Equal(t, lbs[0], "ba:O_o")
	assert.Equal(t, lbs[1], "ca:O_o")
}

func TestCommandAddMarksWrongNode(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	crds := []*coord.Coord{coord.NewCoord(1, 0)}

	// the client was looking at a different node: drop the batch
	stale := s.Current().Index + 1
	frm, err := state.NewAddMarksCommand("triangle", crds, "", &stale).Execute(s)
	assert.NoError(t, err)
	assert.Zero(t, frm)
	assert.Equal(t, len(s.Current().GetField("TR")), 0)

	cur := s.Current().Index
	frm, err = state.NewAddMarksCommand("triangle", crds, "", &cur).Execute(s)
	assert.NoError(t, err)
	require.NotNil(t, frm)
	assert.Equal(t, len(s.Current().GetField("TR")), 1)
}
