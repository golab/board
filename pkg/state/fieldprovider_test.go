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
	"github.com/jarednogo/board/pkg/core/color"
	"github.com/jarednogo/board/pkg/core/coord"
	"github.com/jarednogo/board/pkg/core/fields"
	"github.com/jarednogo/board/pkg/state"
)

func TestIsMove1(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("B", "aa")
	assert.True(t, state.IsMove(f))
}

func TestIsMove2(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("W", "aa")
	assert.True(t, state.IsMove(f))
}

func TestIsMove3(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("b", "aa")
	assert.False(t, state.IsMove(f))
}

func TestIsMove4(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("w", "aa")
	assert.False(t, state.IsMove(f))
}

func TestIsPass1(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("B", "")
	assert.True(t, state.IsPass(f))
}

func TestIsPass2(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("W", "")
	assert.True(t, state.IsPass(f))
}

func TestColor1(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("B", "aa")
	assert.Equal(t, state.Color(f), color.Black)
}

func TestColor2(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("W", "aa")
	assert.Equal(t, state.Color(f), color.White)
}

func TestColor3(t *testing.T) {
	f := &fields.Fields{}
	assert.Equal(t, state.Color(f), color.Empty)
}

func TestCoord1(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("B", "aa")
	crd := state.Coord(f)
	assert.True(t, crd.Equal(coord.NewCoord(0, 0)))
}

func TestCoord2(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("W", "bb")
	crd := state.Coord(f)
	assert.True(t, crd.Equal(coord.NewCoord(1, 1)))
}

func TestCoordNil(t *testing.T) {
	f := &fields.Fields{}
	crd := state.Coord(f)
	assert.Zero(t, crd)
}

func TestTreeNodeIsMove(t *testing.T) {
	s, err := state.FromSGF(sgfsamples.Resignation1)
	assert.NoError(t, err)

	root := s.Root()
	assert.False(t, state.IsMove(root))

	assert.True(t, len(root.Down) > 0)
	current := root.Down[0]
	assert.True(t, state.IsMove(current))

	assert.True(t, len(current.Down) > 0)
	current = current.Down[0]
	assert.True(t, state.IsMove(current))
}
