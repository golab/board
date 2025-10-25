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

func TestState1(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleFourMoves)

	assert.NotNil(t, err, "error should be nil")
	assert.Equal(t, s.Size, 19, "error with state")
}

func TestState2(t *testing.T) {
	input := sgfsamples.SimpleFourMoves
	s, err := state.FromSGF(input)
	assert.NotNil(t, err, "error should be nil")

	sgf := s.ToSGF(false)
	sgfix := s.ToSGF(true)

	assert.Equal(t, len(sgf), len(input), "error with state to sgf")

	assert.Equal(t, len(sgfix), 132, "error with state to sgf (indexes)")
}

func TestState3(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleTwoBranches)
	assert.NotNil(t, err, "error should be nil")

	// index 5 is the start of the second branch from the root
	s.GotoIndex(5) //nolint:errcheck

	prefs := s.Prefs()
	assert.Equal(t, prefs["0"], 1, "prefs[0]")

	err = s.SetPreferred(4)
	assert.NotNil(t, err, "error should not be nil")

	prefs = s.Prefs()
	assert.Equal(t, prefs["0"], 0, "prefs[0]")

	s.ResetPrefs()
	prefs = s.Prefs()
	for _, value := range prefs {
		assert.Equal(t, value, 0, "all prefs should be 0")
	}

	s.GotoIndex(6) //nolint:errcheck
	assert.Equal(t, s.Locate(), "1,0", "s.Locate")

}

func TestState4(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NotNil(t, err, "error should be nil")

	assert.Equal(t, s.Head.Index, 8, "head index")

	// push a node to the head
	s.PushHead(0, 0, core.Black)
	assert.Equal(t, s.Current.Index, 0, "current index")

	// go to the most recently pushed node
	s.GotoIndex(9) //nolint:errcheck

	// since we're at the head, we will "track" along
	s.PushHead(0, 1, core.White)
	assert.Equal(t, s.Current.Index, 10, "current index")
}

func TestAddPatternNodes(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.Empty)
	assert.NotNil(t, err, "error should be nil")

	// add three new moves
	moves := []*core.PatternMove{
		core.NewPatternMove(9, 9, core.Black),
		core.NewPatternMove(10, 10, core.White),
		core.NewPatternMove(11, 11, core.Black),
	}

	s.AddPatternNodes(moves)

	assert.Equal(t, len(s.Current.Down), 1, "len(down)")

	s.FastForward()
	assert.Equal(t, s.Current.Index, 3, "current index")
}

func TestCut(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NotNil(t, err, "error should be nil")

	s.GotoIndex(5) //nolint:errcheck
	s.Cut()
	assert.Equal(t, s.Current.Index, 4, "current index")
	assert.Equal(t, len(s.Current.Down), 0, "failed cut")
}

func TestGraft(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NotNil(t, err, "error should be nil")

	moves := []*core.PatternMove{
		core.NewPatternMove(9, 9, core.Black),
		core.NewPatternMove(10, 10, core.White),
		core.NewPatternMove(11, 11, core.Black),
		core.NewPatternMove(12, 12, core.Black),
	}

	// dumb graft
	s.Graft(4, moves[:3])
	s.GotoIndex(4) //nolint:errcheck
	s.Down()
	s.FastForward()
	assert.Equal(t, s.Current.Index, 11, "current index")

	// smart graft
	s.SmartGraft(4, moves)
	s.FastForward()
	assert.Equal(t, s.Current.Index, 12, "current index")

	// double-checking Up works
	s.GotoIndex(4) //nolint:errcheck
	s.Up()
	s.FastForward()
	assert.Equal(t, s.Current.Index, 8, "current index")

	// check for gotocoord
	s.GotoCoord(3, 3)
	assert.Equal(t, s.Current.Index, 2, "current index")
}

func TestStateJSON(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.SimpleEightMoves)
	assert.NotNil(t, err, "error should be nil")

	s.GotoIndex(4) //nolint:errcheck

	j := s.CreateStateJSON()
	assert.Equal(t, j.NextIndex, 9, "next index")
	assert.Equal(t, j.Location, "0,0,0,0", "location")

	tr := s.CreateTreeJSON(core.FullFrame)
	assert.Equal(t, tr.Depth, 8, "depth")
	assert.Equal(t, tr.Current, 4, "current")
	assert.Equal(t, tr.Up, 0, "up")
	assert.Equal(t, tr.Root, 0, "root")
}
