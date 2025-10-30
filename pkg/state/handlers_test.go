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

func TestHandleAddStone(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	evt := &core.EventJSON{
		Value: map[string]interface{}{
			"coords": []interface{}{9.0, 9.0},
			"color":  1.0,
		},
	}
	_, err = s.HandleAddStone(evt)
	assert.NoError(t, err, "s.HandleAddStone")

}

func TestHandlePass(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	evt := &core.EventJSON{
		Value: 1.0,
	}

	_, err = s.HandlePass(evt)
	assert.NoError(t, err, "s.HandlePass")
}

func TestHandleRemoveStone(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	evt := &core.EventJSON{
		Value: []interface{}{9.0, 9.0},
	}

	_, err = s.HandleRemoveStone(evt)
	assert.NoError(t, err, "s.HandleRemoveStone")
}

func TestHandleRemoveMark(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	evt := &core.EventJSON{
		Value: []interface{}{9.0, 9.0},
	}

	_, err = s.HandleRemoveMark(evt)
	assert.NoError(t, err, "s.HandleRemoveMark")
}

func TestHandleLeft(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	_, err = s.HandleLeft()
	assert.NoError(t, err, "s.HandleLeft")
}

func TestHandleRight(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	_, err = s.HandleRight()
	assert.NoError(t, err, "s.HandleRight")
}

func TestHandleUp(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	_, err = s.HandleUp()
	assert.NoError(t, err, "s.HandleUp")
}

func TestHandleDown(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	_, err = s.HandleDown()
	assert.NoError(t, err, "s.HandleDown")
}

func TestHandleRewind(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	_, err = s.HandleRewind()
	assert.NoError(t, err, "s.HandleRewind")
}

func TestHandleFastForward(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	_, err = s.HandleFastForward()
	assert.NoError(t, err, "s.HandleFastForward")
}

func TestHandleGotoGrid(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	evt := &core.EventJSON{
		Value: 1.0,
	}

	_, err = s.HandleGotoGrid(evt)
	assert.NoError(t, err, "s.HandleGotoGrid")
}

func TestHandleGotoCoord(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	evt := &core.EventJSON{
		Value: []interface{}{9.0, 9.0},
	}

	_, err = s.HandleGotoCoord(evt)
	assert.NoError(t, err, "s.HandleGotoCoord")
}

func TestHandleComment(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	evt := &core.EventJSON{
		Value: "some comment",
	}

	_, err = s.HandleComment(evt)
	assert.NoError(t, err, "s.HandleComment")
}

func TestHandleDraw(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	evt := &core.EventJSON{
		Value: []interface{}{0.0, 0.0, 5.0, 5.0, "#000000"},
	}

	_, err = s.HandleDraw(evt)
	assert.NoError(t, err, "s.HandleDraw")
}

func TestHandleErasePen(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	_, err = s.HandleErasePen()
	assert.NoError(t, err, "s.HandleErasePen")
}

func TestHandleCut(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	_, err = s.HandleCut()
	assert.NoError(t, err, "s.HandleCut")
}

func TestHandleCopy(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	_, err = s.HandleCopy()
	assert.NoError(t, err, "s.HandleCopy")
}

func TestHandleClipboard(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	_, err = s.HandleClipboard()
	assert.NoError(t, err, "s.HandleClipboard")
}

func TestHandleGraft(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	evt := &core.EventJSON{
		Value: "a1 b2 c3 d4",
	}

	_, err = s.HandleGraft(evt)
	assert.NoError(t, err, "s.HandleGraft")
}

func TestHandleScore(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	_, err = s.HandleScore()
	assert.NoError(t, err, "s.HandleScore")
}

func TestHandleMarkDead(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NoError(t, err, "state.FromSGF")

	evt := &core.EventJSON{
		Value: []interface{}{9.0, 9.0},
	}

	_, err = s.HandleMarkDead(evt)
	assert.NoError(t, err, "s.HandleMarkDead")
}
