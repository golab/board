/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package state_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/pkg/core/color"
	"github.com/golab/board/pkg/core/coord"
	"github.com/golab/board/pkg/state"
)

// marks come from uploaded files, so malformed ones must never panic
// and only on-board points should be sent to clients
func TestMalformedMarks(t *testing.T) {
	tests := []struct {
		sgf       string
		triangles int
		squares   int
		labels    int
		pens      int
	}{
		{"(;SZ[19]TR[aa:cc])", 9, 0, 0, 0},
		{"(;SZ[19]SQ[aa:bb][dd])", 0, 5, 0, 0},
		{"(;SZ[19]TR[x][][zzz])", 0, 0, 0, 0},
		{"(;SZ[9]TR[ss][aa])", 1, 0, 0, 0},
		{"(;SZ[19]LB[aa][bb:ok][zz:x][:y])", 0, 0, 1, 0},
		{"(;SZ[19]PX[NaN:0:0:0:#000000][0:Inf:0:0:#000000][0.1:0.1:0.2:0.2:#000000])", 0, 0, 0, 1},
	}
	for _, tt := range tests {
		t.Run(tt.sgf, func(t *testing.T) {
			s, err := state.FromSGF(tt.sgf)
			require.NoError(t, err)

			frame := s.GenerateFullFrame(state.Full)
			assert.Equal(t, len(frame.Marks.Triangles), tt.triangles)
			assert.Equal(t, len(frame.Marks.Squares), tt.squares)
			assert.Equal(t, len(frame.Marks.Labels), tt.labels)
			assert.Equal(t, len(frame.Marks.Pens), tt.pens)

			// the frame has to be sendable
			_, err = json.Marshal(frame)
			assert.NoError(t, err)
		})
	}
}

func TestCompressedSetupStones(t *testing.T) {
	s, err := state.FromSGF("(;SZ[19]AB[aa:cc]AW[ss];AE[bb])")
	require.NoError(t, err)
	assert.Equal(t, s.Board().Get(coord.NewCoord(1, 1)), color.Black)
	assert.Equal(t, s.Board().Get(coord.NewCoord(2, 2)), color.Black)
	assert.Equal(t, s.Board().Get(coord.NewCoord(18, 18)), color.White)

	// AE on the next node removes one of them
	execute(t, s, state.NewRightCommand())
	assert.Equal(t, s.Board().Get(coord.NewCoord(1, 1)), color.Empty)
	assert.Equal(t, s.Board().Get(coord.NewCoord(0, 0)), color.Black)
}

func TestSetupStonesOffBoard(t *testing.T) {
	s, err := state.FromSGF("(;SZ[9]AB[ss][aa])")
	require.NoError(t, err)
	assert.Equal(t, s.Board().Get(coord.NewCoord(0, 0)), color.Black)
}

func TestRemoveMarkShortLabel(t *testing.T) {
	s, err := state.FromSGF("(;SZ[19]LB[a][aa:x][ab:y])")
	require.NoError(t, err)
	execute(t, s, state.NewRemoveMarkCommand(coord.NewCoord(0, 0)))
	assert.Equal(t, labels(s), "a ab:y")
}

func TestCommentBackslashReloads(t *testing.T) {
	s := state.NewState(19)
	execute(t, s, state.NewCommentCommand(`see C:\`))
	execute(t, s, state.NewCommentCommand(`a]b\]`))

	reloaded, err := state.FromSGF(s.ToSGF())
	require.NoError(t, err)
	assert.Equal(t, reloaded.ToSGF(), s.ToSGF())
}

// IX values come from uploaded files: duplicates, negatives and large
// values must not corrupt the index -> node map
func TestUntrustedIndexes(t *testing.T) {
	for _, sgf := range []string{
		"(;IX[0]SZ[19];B[aa]IX[5];W[bb]IX[5];B[cc]IX[-7])",
		"(;IX[3]SZ[19];B[aa]IX[3];W[bb]IX[1000])",
	} {
		t.Run(sgf, func(t *testing.T) {
			s, err := state.FromSGF(sgf)
			require.NoError(t, err)

			// every node in the tree is in the map under its own index
			count := 0
			execute(t, s, state.NewRewindCommand())
			for node := s.Root(); node != nil; {
				assert.True(t, s.Nodes()[node.Index] == node)
				count++
				if len(node.Down) == 0 {
					break
				}
				node = node.Down[0]
			}
			assert.Equal(t, len(s.Nodes()), count)

			// new nodes don't collide with loaded ones
			execute(t, s, state.NewFastForwardCommand())
			execute(t, s, state.NewAddStoneCommand(coord.NewCoord(9, 9), color.Black))
			assert.Equal(t, len(s.Nodes()), count+1)
		})
	}
}

func passes(n int) string {
	var sb strings.Builder
	sb.WriteString("(;SZ[19]")
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			sb.WriteString(";B[]")
		} else {
			sb.WriteString(";W[]")
		}
	}
	sb.WriteString(")")
	return sb.String()
}

func TestNodeCap(t *testing.T) {
	// uploads over the cap are refused
	_, err := state.FromSGF(passes(state.MaxNodes))
	assert.NotNil(t, err)

	// copy and paste doubles the tree: refused once it would pass the cap
	s, err := state.FromSGF(passes(state.MaxNodes / 3))
	require.NoError(t, err)
	before := len(s.Nodes())
	execute(t, s, state.NewCopyCommand())
	execute(t, s, state.NewPasteCommand())
	_, err = state.NewPasteCommand().Execute(s)
	assert.True(t, errors.Is(err, state.ErrTooManyNodes))
	assert.Equal(t, len(s.Nodes()), 2*before)

	// as are single moves once it's full
	full, err := state.FromSGF(passes(state.MaxNodes - 1))
	require.NoError(t, err)
	_, err = state.NewAddStoneCommand(coord.NewCoord(3, 3), color.Black).Execute(full)
	assert.True(t, errors.Is(err, state.ErrTooManyNodes))
	_, err = state.NewPassCommand(color.Black).Execute(full)
	assert.True(t, errors.Is(err, state.ErrTooManyNodes))
}

// each setup node used to replay the game from the root: 10k of them
// took ~24s to load
func TestLoadManySetupNodes(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("(;SZ[19]")
	for i := 0; i < 20_000; i++ {
		sb.WriteString(";AB[aa]")
	}
	sb.WriteString(";AE[aa];AW[aa])")

	start := time.Now()
	s, err := state.FromSGF(sb.String())
	require.NoError(t, err)
	assert.True(t, time.Since(start) < 5*time.Second)

	execute(t, s, state.NewFastForwardCommand())
	assert.Equal(t, s.Board().Get(coord.NewCoord(0, 0)), color.White)
}

func TestGraftLimit(t *testing.T) {
	s := state.NewState(19)
	moves := strings.Repeat("a1 ", state.MaxGraftMoves+2)
	_, err := state.NewGraftCommand(moves).Execute(s)
	assert.NotNil(t, err)

	_, err = state.NewGraftCommand("   ").Execute(s)
	assert.NotNil(t, err)

	_, err = state.NewGraftCommand("0 d4 q16").Execute(s)
	require.NoError(t, err)
}

func TestBoardSizeRange(t *testing.T) {
	for _, sgf := range []string{"(;SZ[0])", "(;SZ[-3])", "(;SZ[20])"} {
		_, err := state.FromSGF(sgf)
		assert.NotNil(t, err)
	}
	for _, sgf := range []string{"(;SZ[1])", "(;SZ[2];B[aa])", "(;SZ[7];B[aa])"} {
		_, err := state.FromSGF(sgf)
		require.NoError(t, err)
	}
}
