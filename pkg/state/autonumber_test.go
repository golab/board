/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state_test

import (
	"strings"
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/internal/sgfsamples"
	"github.com/golab/board/pkg/core/color"
	"github.com/golab/board/pkg/core/coord"
	"github.com/golab/board/pkg/state"
)

func execute(t *testing.T, s *state.State, cmd state.Command) *state.Frame {
	t.Helper()
	frm, err := cmd.Execute(s)
	require.NoError(t, err)
	return frm
}

func labels(s *state.State) string {
	return strings.Join(s.Current().GetField("LB"), " ")
}

func TestAutoNumberTrunk(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	require.NoError(t, err)
	s.SetAutoNumber(true)
	execute(t, s, state.NewFastForwardCommand())

	// adding to the end of the trunk extends the trunk
	execute(t, s, state.NewAddStoneCommand(coord.NewCoord(9, 9), color.Black))
	assert.Equal(t, labels(s), "")
}

func TestAutoNumberDisabled(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	require.NoError(t, err)
	execute(t, s, state.NewFastForwardCommand())

	execute(t, s, state.NewLeftCommand())
	execute(t, s, state.NewAddStoneCommand(coord.NewCoord(9, 9), color.White))
	assert.Equal(t, labels(s), "")
}

func TestAutoNumberBranch(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	require.NoError(t, err)
	s.SetAutoNumber(true)
	execute(t, s, state.NewFastForwardCommand())

	// step back one move and play something different
	execute(t, s, state.NewLeftCommand())
	frm := execute(t, s, state.NewAddStoneCommand(coord.NewCoord(9, 9), color.White))
	assert.Equal(t, labels(s), "jj:1")
	require.NotNil(t, frm.Marks)
	assert.Equal(t, len(frm.Marks.Labels), 1)

	execute(t, s, state.NewAddStoneCommand(coord.NewCoord(2, 9), color.Black))
	assert.Equal(t, labels(s), "jj:1 cj:2")

	// passes count as moves and keep the numbers visible
	frm = execute(t, s, state.NewPassCommand(color.White))
	assert.Equal(t, labels(s), "jj:1 cj:2")
	require.NotNil(t, frm.Marks)
	assert.Equal(t, len(frm.Marks.Labels), 2)

	execute(t, s, state.NewAddStoneCommand(coord.NewCoord(15, 9), color.Black))
	assert.Equal(t, labels(s), "jj:1 cj:2 pj:4")

	// earlier nodes keep their own labels
	execute(t, s, state.NewLeftCommand())
	execute(t, s, state.NewLeftCommand())
	assert.Equal(t, labels(s), "jj:1 cj:2")

	// a branch off the branch keeps counting from the trunk
	execute(t, s, state.NewAddStoneCommand(coord.NewCoord(12, 12), color.White))
	assert.Equal(t, labels(s), "jj:1 cj:2 mm:3")
}

func TestAutoNumberCapture(t *testing.T) {
	s := state.NewState(9)
	s.SetAutoNumber(true)

	// trunk move, then go back and start a branch
	execute(t, s, state.NewAddStoneCommand(coord.NewCoord(0, 0), color.Black))
	execute(t, s, state.NewLeftCommand())

	execute(t, s, state.NewAddStoneCommand(coord.NewCoord(8, 8), color.Black))
	execute(t, s, state.NewAddStoneCommand(coord.NewCoord(7, 8), color.White))
	execute(t, s, state.NewAddStoneCommand(coord.NewCoord(4, 4), color.Black))
	assert.Equal(t, labels(s), "ii:1 hi:2 ee:3")

	// captures the stone numbered 1, so its number goes away
	execute(t, s, state.NewAddStoneCommand(coord.NewCoord(8, 7), color.White))
	assert.Equal(t, labels(s), "hi:2 ee:3 ih:4")
}

func TestAutoNumberMetadata(t *testing.T) {
	s := state.NewState(19)
	assert.False(t, s.GenerateFullFrame(state.Full).Metadata.AutoNumber)
	s.SetAutoNumber(true)
	assert.True(t, s.GenerateFullFrame(state.Full).Metadata.AutoNumber)
}
