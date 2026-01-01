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

func TestState(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleFourMoves)

	assert.NoError(t, err)
	assert.Equal(t, s.Size(), 19)
}

func TestState2(t *testing.T) {
	input := sgfsamples.SimpleFourMoves
	s, err := state.FromSGF(input)
	assert.NoError(t, err)

	sgf := s.ToSGF()
	sgfix := s.ToSGFIX()

	assert.Equal(t, len(sgf), len(input))
	assert.Equal(t, len(sgfix), 132)
}

func TestState3(t *testing.T) {
	input := sgfsamples.SimpleTwoBranches
	s, err := state.FromSGF(input)
	assert.NoError(t, err)

	sgf := s.ToSGF()
	assert.Equal(t, len(sgf), len(input))
}

func TestMultifieldSZ(t *testing.T) {
	input := sgfsamples.MultifieldSZ
	_, err := state.FromSGF(input)
	assert.NotNil(t, err)
}

func TestSGFIX(t *testing.T) {
	input := sgfsamples.SGFIX
	s, err := state.FromSGF(input)
	assert.NoError(t, err)
	assert.Equal(t, s.Root().Index, 100)
}

func TestParseTTPAss(t *testing.T) {
	input := sgfsamples.PassWithTT
	s, err := state.FromSGF(input)
	assert.NoError(t, err)
	node := s.Nodes()[358]
	assert.Equal(t, node.XY, nil)
}

func TestNewState(t *testing.T) {
	s1 := state.NewState(19)
	s2 := state.NewEmptyState(19)
	assert.Equal(t, s1.GetNextIndex(), 1)
	assert.Equal(t, s2.GetNextIndex(), 0)
}

func TestHandicap(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.Handicap1)
	assert.NoError(t, err)
	root := s.Root()
	assert.Equal(t, len(root.GetField("AB")), 4)
}

func TestMissingSZ(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.MissingSZ)
	assert.NoError(t, err)
	sz := s.Size()
	require.Equal(t, sz, 19)
}

func TestSuicide(t *testing.T) {
	_, err := state.FromSGF(sgfsamples.Suicide1)
	assert.NotNil(t, err)
}

func TestHeadColor(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	require.NoError(t, err)
	assert.Equal(t, s.HeadColor(), color.White)
}

func TestGetColorAt1(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	require.NoError(t, err)
	assert.Equal(t, s.GetColorAt(3), color.Black)
}

func TestGetColorAt2(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	require.NoError(t, err)
	assert.Equal(t, s.GetColorAt(-1), color.Empty)
}

func TestSetNextIndex(t *testing.T) {
	s := state.NewState(9)
	s.SetNextIndex(100)
	assert.Equal(t, s.GetNextIndex(), 100)
	assert.Equal(t, s.GetNextIndex(), 101)
}

func TestEditPlayerBlack(t *testing.T) {
	s := state.NewState(9)
	s.EditPlayerBlack("pblack")
	pb := s.Root().GetField("PB")
	require.Equal(t, len(pb), 1)
	assert.Equal(t, pb[0], "pblack")
}

func TestEditPlayerWhite(t *testing.T) {
	s := state.NewState(9)
	s.EditPlayerWhite("pwhite")
	pw := s.Root().GetField("PW")
	require.Equal(t, len(pw), 1)
	assert.Equal(t, pw[0], "pwhite")
}

func TestEditKomi(t *testing.T) {
	s := state.NewState(9)
	s.EditKomi("100.5")
	km := s.Root().GetField("KM")
	require.Equal(t, len(km), 1)
	assert.Equal(t, km[0], "100.5")
}

func TestAddNode(t *testing.T) {
	s := state.NewState(19)
	s.AddNode(coord.NewCoord(9, 9), color.Black)
	assert.Equal(t, s.Current().Index, 1)
	s.AddNode(coord.NewCoord(10, 10), color.White)
	assert.Equal(t, s.Current().Index, 2)
	s.AddNode(coord.NewCoord(11, 11), color.Black)
	assert.Equal(t, s.Current().Index, 3)
}

func TestAddStonesToTrunk(t *testing.T) {
	s := state.NewState(19)
	s.PushHead(10, 10, color.Black)
	s.PushHead(11, 11, color.White)
	s.PushHead(12, 12, color.Black)
	s.PushHead(13, 13, color.White)
	stones := []*coord.Stone{}
	stones = append(stones, coord.NewStone(2, 3, color.Black))
	stones = append(stones, coord.NewStone(14, 3, color.White))
	s.AddStonesToTrunk(2, stones)
	assert.Equal(t, s.GetNextIndex(), 7)
}

func TestAddStonesNoChange(t *testing.T) {
	s := state.NewState(19)
	s.PushHead(10, 10, color.Black)
	s.PushHead(11, 11, color.White)
	s.PushHead(12, 12, color.Black)
	s.PushHead(13, 13, color.White)
	stones := []*coord.Stone{}
	stones = append(stones, coord.NewStone(10, 10, color.Black))
	stones = append(stones, coord.NewStone(11, 11, color.White))
	s.AddStones(stones)
	assert.Equal(t, s.GetNextIndex(), 5)
}
