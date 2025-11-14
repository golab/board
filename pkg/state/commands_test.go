/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state_test

import (
	"testing"

	"github.com/jarednogo/board/internal/assert"
	"github.com/jarednogo/board/internal/sgfsamples"
	"github.com/jarednogo/board/pkg/core"
	"github.com/jarednogo/board/pkg/state"
)

func TestCommandAddStone(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	coord := &core.Coord{X: 9, Y: 9}
	color := core.Black
	_, err = state.NewAddStoneCommand(coord, color).Execute(s)
	assert.NoError(t, err)

}

func TestCommandPass(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	_, err = state.NewPassCommand(core.Black).Execute(s)
	assert.NoError(t, err)
}

func TestCommandRemoveStone(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	coord := &core.Coord{X: 9, Y: 9}
	_, err = state.NewRemoveStoneCommand(coord).Execute(s)
	assert.NoError(t, err)
}

func TestCommandRemoveMark(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	coord := &core.Coord{X: 9, Y: 9}
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

	coord := &core.Coord{X: 9, Y: 9}

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

func TestCommandGraft(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err)

	text := "a1 b2 c3 d4"

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

	coord := &core.Coord{X: 9, Y: 9}

	_, err = state.NewMarkDeadCommand(coord).Execute(s)
	assert.NoError(t, err)
}
