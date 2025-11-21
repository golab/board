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
	"github.com/jarednogo/board/pkg/event"
	"github.com/jarednogo/board/pkg/state"
)

func TestDecodeAddStone(t *testing.T) {
	val := make(map[string]any)
	val["color"] = 1.0
	val["coords"] = []any{9.0, 9.0}
	e := event.NewEvent("add_stone", val)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodePass(t *testing.T) {
	e := event.NewEvent("pass", 1.0)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeRemoveStone(t *testing.T) {
	val := []any{9.0, 9.0}
	e := event.NewEvent("remove_stone", val)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeTriangle(t *testing.T) {
	val := []any{9.0, 9.0}
	e := event.NewEvent("triangle", val)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeSquare(t *testing.T) {
	val := []any{9.0, 9.0}
	e := event.NewEvent("square", val)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeLetter(t *testing.T) {
	val := make(map[string]any)
	val["coords"] = []any{9.0, 9.0}
	val["letter"] = "A"
	e := event.NewEvent("letter", val)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeNumber(t *testing.T) {
	val := make(map[string]any)
	val["coords"] = []any{9.0, 9.0}
	val["number"] = 1.0
	e := event.NewEvent("number", val)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeRemoveMark(t *testing.T) {
	val := []any{9.0, 9.0}
	e := event.NewEvent("remove_mark", val)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeCut(t *testing.T) {
	e := event.NewEvent("cut", nil)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeLeft(t *testing.T) {
	e := event.NewEvent("left", nil)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeRight(t *testing.T) {
	e := event.NewEvent("right", nil)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeUp(t *testing.T) {
	e := event.NewEvent("up", nil)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeDown(t *testing.T) {
	e := event.NewEvent("down", nil)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeRewind(t *testing.T) {
	e := event.NewEvent("rewind", nil)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeFastForward(t *testing.T) {
	e := event.NewEvent("fastforward", nil)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeGotoGrid(t *testing.T) {
	e := event.NewEvent("goto_grid", 3.0)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeGotoCoord(t *testing.T) {
	val := []any{9.0, 9.0}
	e := event.NewEvent("goto_coord", val)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeComment(t *testing.T) {
	e := event.NewEvent("comment", "somecomment")

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeDraw(t *testing.T) {
	val := []any{0.0, 0.0, 1.0, 1.0, "#445566"}
	e := event.NewEvent("draw", val)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeErasePen(t *testing.T) {
	e := event.NewEvent("erase_pen", nil)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeCopy(t *testing.T) {
	e := event.NewEvent("copy", nil)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeClipboard(t *testing.T) {
	e := event.NewEvent("clipboard", nil)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeGraft(t *testing.T) {
	e := event.NewEvent("graft", "a1 a2 a3 a4")

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeScore(t *testing.T) {
	e := event.NewEvent("score", nil)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeMarkDead(t *testing.T) {
	val := []any{9.0, 9.0}
	e := event.NewEvent("markdead", val)

	_, err := state.DecodeToCommand(e)
	assert.NoError(t, err)
}

func TestDecodeFail(t *testing.T) {
	e := event.NewEvent("failevent", nil)
	_, err := state.DecodeToCommand(e)
	assert.NotNil(t, err)
}
